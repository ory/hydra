package metrics_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/ory/hydra/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/negroni"
	//"time"
	"math/rand"
	"time"

	"encoding/json"

	"github.com/ory/herodot"
	"github.com/ory/hydra/health"
	"github.com/sirupsen/logrus"
)

func TestMiddleware(t *testing.T) {
	rand.Seed(time.Now().Unix())
	mw := metrics.NewMetricsManager(logrus.StandardLogger())
	n := negroni.New()
	r := httprouter.New()

	time.Sleep(time.Second)

	n.Use(mw)
	r.GET("/", handle)
	r.GET("/oauth2/introspect", handle)
	r.POST("/", handle)
	n.UseHandler(r)

	s := httptest.NewServer(n)
	defer s.Close()

	for i := 1; i <= 100; i++ {
		t.Run(fmt.Sprintf("case=%d", i), func(t *testing.T) {
			res, err := http.Get(s.URL)
			require.NoError(t, err)
			res.Body.Close()

			mw.RLock()
			require.Equal(t, http.StatusOK, res.StatusCode)
			mw.RUnlock()
		})
	}

	i := 100
	assert.EqualValues(t, i, mw.Snapshot.Requests)
	assert.EqualValues(t, i, mw.Snapshot.Requests)
	assert.EqualValues(t, i, mw.Snapshot.Responses)

	mw.Snapshot.Update()
	assert.True(t, mw.Snapshot.UpTime > 0)
	assert.True(t, mw.Snapshot.GetUpTime() > 0)

	assert.EqualValues(t, 0, mw.Snapshot.Status[http.StatusOK].Requests)
	assert.EqualValues(t, i, mw.Snapshot.Status[http.StatusOK].Responses)

	assert.EqualValues(t, i, mw.Snapshot.Methods["GET"].Requests)
	assert.EqualValues(t, i, mw.Snapshot.Methods["GET"].Responses)

	res, err := http.Get(s.URL + "/oauth2/introspect/1231")
	require.NoError(t, err)
	res.Body.Close()

	assert.EqualValues(t, 1, mw.Snapshot.Path("/oauth2/introspect").Requests)

	mw.Lock()
	mw.Update()
	assert.NotEqual(t, 0, mw.UpTime)
	mw.Unlock()

	out, _ := json.MarshalIndent(mw, "\t", "  ")
	t.Logf("%s", out)
}

func TestRacyMiddleware(t *testing.T) {
	rand.Seed(time.Now().Unix())
	mw := metrics.NewMetricsManager(logrus.StandardLogger())
	n := negroni.New()
	r := httprouter.New()

	h := health.Handler{
		H:       herodot.NewJSONWriter(nil),
		Metrics: mw,
	}

	n.Use(mw)
	h.SetRoutes(r)
	n.UseHandler(r)

	s := httptest.NewServer(n)
	defer s.Close()

	for i := 1; i <= 100; i++ {
		t.Run("type=concurrent", func(t *testing.T) {
			go func() {
				res, err := http.Get(s.URL + "/health")
				require.NoError(t, err)
				res.Body.Close()
			}()

		})
	}

	time.Sleep(time.Second)
}

func handle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	time.Sleep(time.Duration(random(0, 10)) * time.Millisecond)
	w.WriteHeader(http.StatusOK)

	for i := 0; i <= random(1, 100); i++ {
		w.Write([]byte("ok"))
	}
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}
