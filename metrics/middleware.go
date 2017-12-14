// Copyright © 2017 Aeneas Rekkas <aeneas+oss@aeneas.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/segmentio/analytics-go"
	"github.com/urfave/negroni"
	//"github.com/ory/hydra/cmd"
	"crypto/sha512"
	"encoding/base64"
	"strings"

	"github.com/ory/hydra/pkg"
	"github.com/pborman/uuid"
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
	issuerURL    string
	databaseURL  string
	internalID   string
}

func shouldCommit(issuerURL string, databaseURL string) bool {
	return !(databaseURL == "" || databaseURL == "memory" || issuerURL == "" || strings.Contains(issuerURL, "localhost"))
}

func identify(issuerURL string) string {
	hash := sha512.New()
	hash.Write([]byte(issuerURL))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

func NewMetricsManager(issuerURL string, databaseURL string, l logrus.FieldLogger) *MetricsManager {
	l.Info("Setting up telemetry - for more information please visit https://ory.gitbooks.io/hydra/content/telemetry.html")

	mm := &MetricsManager{
		Snapshot: &Snapshot{
			MemorySnapshot: &MemorySnapshot{},
			ID:             identify(issuerURL),
			Metrics:        newMetrics(),
			HTTPMetrics:    newHttpMetrics(),
			Paths:          map[string]*PathMetrics{},
			start:          time.Now().UTC(),
		},
		internalID:  uuid.New(),
		Segment:     analytics.New("h8dRH3kVCWKkIFWydBmWsyYHR4M0u0vr"),
		Logger:      l,
		issuerURL:   issuerURL,
		databaseURL: databaseURL,
	}
	return mm
}

const (
	defaultWait   = time.Minute * 15
	keepAliveWait = time.Minute * 5
)

func (sw *MetricsManager) RegisterSegment(version, hash, buildTime string) {
	if !shouldCommit(sw.issuerURL, sw.databaseURL) {
		sw.Logger.Info("Detected local environment, skipping telemetry commit")
		return
	}

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
				"instanceId":     sw.internalID,
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
	if !shouldCommit(sw.issuerURL, sw.databaseURL) {
		sw.Logger.Info("Detected local environment, skipping telemetry commit")
		return
	}

	time.Sleep(defaultWait)
	for {
		if err := sw.Segment.Track(&analytics.Track{
			Event:       "keep-alive",
			AnonymousId: sw.ID,
			Properties:  map[string]interface{}{"instanceId": sw.internalID},
			Context:     map[string]interface{}{"ip": "0.0.0.0"},
		}); err != nil {
			sw.Logger.WithError(err).Debug("Could not send telemetry keep alive")
		}
		sw.Logger.Debug("Transmitted telemetry heartbeat (keep-alive)")
		time.Sleep(keepAliveWait)
	}
}

func (sw *MetricsManager) CommitTelemetry() {
	if !shouldCommit(sw.issuerURL, sw.databaseURL) {
		sw.Logger.Info("Detected local environment, skipping telemetry commit")
		return
	}

	for {
		time.Sleep(defaultWait)
		sw.Update()
		if err := sw.Segment.Track(&analytics.Track{
			Event:       "telemetry",
			AnonymousId: sw.ID,
			Properties: map[string]interface{}{
				"upTime":     sw.UpTime,
				"requests":   sw.Requests,
				"responses":  sw.Responses,
				"paths":      sw.Paths,
				"methods":    sw.Methods,
				"sizes":      sw.Sizes,
				"status":     sw.Status,
				"latencies":  sw.Latencies,
				"raw":        sw,
				"memory":     sw.MemorySnapshot,
				"instanceId": sw.internalID,
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

	sw.Lock()
	sw.Snapshot.AddRequest()
	sw.Snapshot.AddMethodRequest(method)
	sw.Snapshot.Path(r.RequestURI).AddRequest()
	sw.Snapshot.Path(r.RequestURI).AddMethodRequest(method)
	sw.Unlock()

	// Latency
	start := time.Now().UTC()
	next(rw, r)
	latency := time.Now().UTC().Sub(start) / time.Millisecond

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
