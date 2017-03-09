package broker

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats"
	"github.com/pkg/errors"
	"github.com/pborman/uuid"
	"time"
	"fmt"
)

type jsonError struct {
	Message string `json:"message"`
}

type Broker struct {
	Logger  *logrus.Logger
	N       *nats.Conn
	Version string
	Timeout time.Duration
}

func New(n *nats.Conn, version string) *Broker {
	return &Broker{
		Logger: logrus.New(),
		Version: version,
		N: n,
		Timeout: time.Second * 5,
	}
}

type Container struct {
	ID        string `json:"i"`
	Version   string `json:"v"`
	RequestID string `json:"r"`
	Status    int `json:"s"`
	Payload   interface{} `json:"p"`
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func (h *Broker) RID(r *http.Request) string {
	return r.Header.Get("X-REQUEST-ID")
}

func (h *Broker) GetTimeout() time.Duration {
	if h.Timeout == 0 {
		return time.Second * 5
	}
	return h.Timeout
}

func (h *Broker) GetVersion() string {
	if h.Version == "" {
		return "0.0.0"
	}
	return h.Version
}

func (h *Broker) Reply(m *nats.Msg, rid string, e interface{}) {
	h.WriteCode(m.Reply, rid, http.StatusOK, e)
}

func (h *Broker) WriteCode(message string, rid string, code int, e interface{}) {
	p, err := json.Marshal(&Container{
		ID: uuid.New(),
		Version: h.Version,
		Status: code,
		Payload: e,
		RequestID: rid,
	})
	if err != nil {
		h.WriteErrorCode(message, rid, http.StatusInternalServerError, errors.Wrap(err, "Could not marshal container"))
	}

	if err := h.N.Publish(message, p); err != nil {
		h.Logger.WithError(err).WithField("request", rid).Errorln("Message can not be published.")
	}
}

func (h *Broker) Parse(m *nats.Msg, e interface{}) (*Container, error) {
	var c = &Container{}
	if err := json.Unmarshal(m.Data, c); err != nil {
		return c, errors.Wrap(err, "Could not unmarshal message container")
	}

	if c.Status < 200 || c.Status >= 300 {
		var e jsonError
		c = &Container{Payload: &e}
		if err := json.Unmarshal(m.Data, c); err != nil {
			return c, errors.Wrap(err, "Could not unmarshal message error")
		}

		return c, errors.Errorf("An error code (%d) occurred on the other side: %s", c.Status, e.Message)
	}

	c = &Container{Payload: e}
	if err := json.Unmarshal(m.Data, c); err != nil {
		return c, errors.Wrap(err, "Could not unmarshal message container")
	}

	return c, nil
}

func (h *Broker) Request(message string, rid string, in, out interface{}) (*Container, error) {
	p, err := json.Marshal(&Container{
		ID: uuid.New(),
		Version: h.Version,
		Payload: in,
		Status: http.StatusOK,
		RequestID: rid,
	})
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	rep, err := h.N.Request(message, p, h.GetTimeout())
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return h.Parse(rep, out)
}

func (h *Broker) Publish(message string, rid string, in interface{}) (error) {
	p, err := json.Marshal(&Container{
		ID: uuid.New(),
		Version: h.Version,
		Payload: in,
		Status: http.StatusOK,
		RequestID: rid,
	})
	if err != nil {
		return errors.Wrap(err, "")
	}

	if err := h.N.Publish(message, p); err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

func (h *Broker) MessageLogger(f func (m *nats.Msg)) func (m *nats.Msg) {
	return func (m *nats.Msg) {
		c, _ := h.Parse(m, nil)
		logrus.WithField("id", c.ID).WithField("request", c.RequestID).WithField("subject", m.Subject).Info("Received message")
		f(m)
		logrus.WithField("id", c.ID).WithField("request", c.RequestID).WithField("subject", m.Subject).Info("Handled message")
	}
}

func (h *Broker) WriteErrorCode(message string, rid string, code int, err error) {
	if code == 0 {
		code = http.StatusInternalServerError
	}

	var stack = "not available"
	if e, ok := err.(stackTracer); ok {
		stack = fmt.Sprintf("%+v", e.StackTrace())
	} else if e, ok := errors.Cause(err).(stackTracer); ok {
		stack = fmt.Sprintf("%+v", e.StackTrace())
	}

	h.Logger.WithError(err).WithField("request", rid).WithField("stack", stack).Errorln("An error occurred while sending the response.")
	h.WriteCode(
		message,
		rid,
		code,
		&jsonError{
			Message:   err.Error(),
		},
	)
}
