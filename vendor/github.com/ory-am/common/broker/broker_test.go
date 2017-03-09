package broker

import (
	"fmt"
	"os"
	"gopkg.in/ory-am/dockertest.v3"
	"log"
	"testing"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var nc *nats.Conn

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("nats", "0.9.4", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		var err error
		nc, err = nats.Connect(fmt.Sprintf("nats://localhost:%s", resource.GetPort("4222/tcp")))
		if err != nil {
			return err
		}
		if nc.Status() != nats.CONNECTED {
			return errors.New("Not connected yet")
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	s := m.Run()
	pool.Purge(resource)
	os.Exit(s)
}

func TestBroker(t *testing.T) {
	b := New(nc, "0.0.1")

	type req struct{ Foo string `json:"foo"` }
	type res struct{ Bar string `json:"bar"` }

	t.Run("case=1/description=should encode and decode messages", func(t *testing.T) {
		sub, err := nc.Subscribe("foo", func(m *nats.Msg) {
			var f req
			c, err := b.Parse(m, &f)
			require.Nil(t, err)

			b.Reply(m, c.RequestID, &res{Bar: f.Foo})
		})
		require.Nil(t, err)
		defer sub.Unsubscribe()

		var r res
		c, err := b.Request("foo", "request-id", &req{Foo: "bar"}, &r)
		require.Nil(t, err)
		assert.Equal(t, r.Bar, "bar")
		assert.Equal(t, c.RequestID, "request-id")
	})

	t.Run("case=2/description=should be able to write errors", func(t *testing.T) {
		sub, err := nc.Subscribe("foo", func(m *nats.Msg) {
			var f req
			c, err := b.Parse(m, &f)
			require.Nil(t, err)
			b.WriteErrorCode(m.Reply, c.RequestID, 404, errors.New("some error"))
		})
		require.Nil(t, err)
		defer sub.Unsubscribe()

		var r res
		c, err := b.Request("foo", "request-id", &req{Foo: "bar"}, &r)
		assert.NotNil(t, err)
		require.NotNil(t, c)
		assert.Equal(t, c.Status, 404)
	})

	t.Run("case=3/description=parse should always return a container", func(t *testing.T) {
		c, err := b.Parse(&nats.Msg{}, nil)
		require.NotNil(t, err)
		require.NotNil(t, c)

	})
}
