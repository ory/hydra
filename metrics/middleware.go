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
	"github.com/ory/hydra/pkg"
	"github.com/sirupsen/logrus"
)

type MetricsManager struct {
	sync.RWMutex `json:"-"`
	*Snapshot
	Segment      *analytics.Client  `json:"-"`
	Logger       logrus.FieldLogger `json:"-"`
	buildVersion string
	buildHash    string
	buildTime    string
}

func NewMetricsManager(l logrus.FieldLogger) *MetricsManager {
	l.Info("Setting up telemetry - for more information please visit https://ory.gitbooks.io/hydra/content/telemetry.html")
	mm := &MetricsManager{
		Snapshot: &Snapshot{
			MemorySnapshot: &MemorySnapshot{},
			ID:             uuid.New(),
			Metrics:        newMetrics(),
			HTTPMetrics:    newHttpMetrics(),
			Paths:          map[string]*PathMetrics{},
			start:          time.Now(),
		},
		Segment: analytics.New("JYilhx5zP8wrzfykUinXrSUbo5cRA3aA"),
		Logger:  l,
	}
	return mm
}

const (
	defaultWait   = time.Minute * 15
	keepAliveWait = time.Minute * 5
)

func (sw *MetricsManager) RegisterSegment(version, hash, buildTime string) {
	time.Sleep(defaultWait)
	if err := pkg.Retry(sw.Logger, time.Minute*2, defaultWait, func() error {
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
	}); err != nil {
		sw.Logger.WithError(err).Debug("Could not commit anonymized environment information")
	}
	sw.Logger.Debug("Transmitted anonymized environment information")
}

func (sw *MetricsManager) TickKeepAlive() {
	time.Sleep(defaultWait)
	for {
		if err := sw.Segment.Track(&analytics.Track{
			Event:       "keep-alive",
			AnonymousId: sw.ID,
			Properties:  map[string]interface{}{},
			Context:     map[string]interface{}{"ip": "0.0.0.0"},
		}); err != nil {
			sw.Logger.WithError(err).Debug("Could not send telemetry keep alive")
		}
		sw.Logger.Debug("Transmitted telemetry heartbeat (keep-alive)")
		time.Sleep(keepAliveWait)
	}
}

func (sw *MetricsManager) CommitTelemetry() {
	for {
		time.Sleep(defaultWait)
		sw.Update()
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
				"memory":    sw.MemorySnapshot,
			},
			Context: map[string]interface{}{
				"ip": "0.0.0.0",
			},
		}); err != nil {
			sw.Logger.WithError(err).Debug("Could not commit anonymized telemetry data")
		}
		sw.Logger.Debug("Telemetry data transmitted")
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
