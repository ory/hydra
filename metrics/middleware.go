package metrics

import (
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/pborman/uuid"
	"github.com/segmentio/analytics-go"
	"github.com/urfave/negroni"
	//"github.com/ory/hydra/cmd"
	"github.com/Sirupsen/logrus"
	"github.com/ory/hydra/pkg"
)

type MetricsManager struct {
	sync.RWMutex `json:"-"`
	ID           string `json:"id"`
	UpTime       int64  `json:"uptime"`
	*Snapshot
	Segment      *analytics.Client  `json:"-"`
	Logger       logrus.FieldLogger `json:"-"`
	start        time.Time          `json:"-"`
	buildVersion string
	buildHash    string
	buildTime    string
}

func NewMetricsManager(l logrus.FieldLogger) *MetricsManager {
	l.Info("Setting up telemetry - for more information please visit https://ory.gitbooks.io/hydra/content/telemetry.html")
	mm := &MetricsManager{
		Snapshot: &Snapshot{
			Metrics:     newMetrics(),
			HTTPMetrics: newHttpMetrics(),
			Paths:       map[string]*PathMetrics{},
		},
		ID:      uuid.New(),
		Segment: analytics.New("JYilhx5zP8wrzfykUinXrSUbo5cRA3aA"),
		Logger:  l,
		start:   time.Now(),
	}
	return mm
}

func (sw *MetricsManager) UpdateUpTime() {
	sw.Lock()
	defer sw.Unlock()
	sw.UpTime = int64(time.Now().Sub(sw.start) / time.Second)

}

const (
	defaultWait   = time.Minute * 15
	keepAliveWait = time.Minute * 5
)

func (sw *MetricsManager) RegisterSegment(version, hash, buildTime string) {
	time.Sleep(defaultWait)
	pkg.Retry(sw.Logger, time.Minute*2, defaultWait, func() error {
		return sw.Segment.Identify(&analytics.Identify{
			AnonymousId: sw.ID,
			Traits: map[string]interface{}{
				"goarch":         runtime.GOARCH,
				"goos":           runtime.GOOS,
				"numCpu":         runtime.NumCPU(),
				"runtimeVersion": runtime.Version(),
				"version":        version,
				"hash":           hash,
				"buildTime":      buildTime,
			},
			Context: map[string]interface{}{
				"ip": "0.0.0.0",
			},
		})
	})
}

func (sw *MetricsManager) TickKeepAlive() {
	time.Sleep(defaultWait)
	for {
		if err := sw.Segment.Track(&analytics.Track{
			Event:       "keep-alive",
			AnonymousId: sw.ID,
			Properties:  map[string]interface{}{"nonInteraction": 1},
			Context:     map[string]interface{}{"ip": "0.0.0.0"},
		}); err != nil {
			logrus.WithError(err).Debugf("Could not commit anonymized telemetry data")
		}
		time.Sleep(keepAliveWait)
	}
}

func (sw *MetricsManager) CommitTelemetry() {
	for {
		time.Sleep(defaultWait)
		if err := sw.Segment.Track(&analytics.Track{
			Event:       "telemetry",
			AnonymousId: sw.ID,
			Properties: map[string]interface{}{
				"upTime":    sw.UpTime,
				"requests":  sw.Requests,
				"responses": sw.Responses,
				"paths":     sw.Paths,
				"methods":   sw.Methods,
				"sizes":     sw.Sizes,
				"status":    sw.Status,
				"latencies": sw.Latencies,
				"raw":       sw,
			},
			Context: map[string]interface{}{
				"ip": "0.0.0.0",
			},
		}); err != nil {
			logrus.WithError(err).Debugf("Could not commit anonymized telemetry data")
		}
	}
}

func (sw *MetricsManager) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	method := r.Method
	path := r.RequestURI

	go func() {
		sw.Lock()
		defer sw.Unlock()
		sw.Snapshot.AddRequest()
		sw.Snapshot.AddMethodRequest(method)
		sw.Snapshot.Path(r.RequestURI).AddRequest()
		sw.Snapshot.Path(r.RequestURI).AddMethodRequest(method)
	}()

	// Latency
	start := time.Now()
	next(rw, r)
	latency := time.Now().Sub(start) / time.Millisecond

	// Collecting request info
	res := rw.(negroni.ResponseWriter)
	status := res.Status()
	size := res.Size()

	sw.Lock()
	defer sw.Unlock()
	sw.Snapshot.AddResponse()
	sw.Snapshot.AddMethodResponse(method)
	sw.Snapshot.AddSize(size)
	sw.Snapshot.AddStatus(status)
	sw.Snapshot.AddLatency(latency)

	sw.Snapshot.Path(path).AddResponse()
	sw.Snapshot.Path(path).AddMethodResponse(method)
	sw.Snapshot.Path(path).AddSize(size)
	sw.Snapshot.Path(path).AddStatus(status)
	sw.Snapshot.Path(path).AddLatency(latency)

	sw.Snapshot.Path(path).StatusMetrics(status).AddLatency(latency)
	sw.Snapshot.Path(path).MethodMetrics(method).AddLatency(latency)
	sw.Snapshot.Path(path).SizeMetrics(size).AddLatency(latency)
	sw.Snapshot.StatusMetrics(status).AddLatency(latency)
	sw.Snapshot.MethodMetrics(method).AddLatency(latency)
	sw.Snapshot.SizeMetrics(size).AddLatency(latency)
}
