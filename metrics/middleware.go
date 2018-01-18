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
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/ory/hydra/client"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/oauth2"
	"github.com/ory/hydra/pkg"
	"github.com/ory/hydra/policy"
	"github.com/ory/hydra/warden/group"
	"github.com/pborman/uuid"
	"github.com/segmentio/analytics-go"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

type MetricsManager struct {
	sync.RWMutex                    `json:"-"`
	start        time.Time          `json:"-"`
	Segment      analytics.Client   `json:"-"`
	Logger       logrus.FieldLogger `json:"-"`
	issuerURL    string             `json:"-"`
	databaseURL  string             `json:"-"`
	shouldCommit bool               `json:"-"`

	ID               string            `json:"id"`
	UpTime           int64             `json:"uptime"`
	MemoryStatistics *MemoryStatistics `json:"memory"`
	BuildVersion     string            `json:"buildVersion"`
	BuildHash        string            `json:"buildHash"`
	BuildTime        string            `json:"buildTime"`
	InstanceID       string            `json:"instanceId"`
}

func shouldCommit(issuerURL string, databaseURL string) bool {
	return !(databaseURL == "" || databaseURL == "memory" || issuerURL == "" || strings.Contains(issuerURL, "localhost"))
}

func identify(issuerURL string) string {
	hash := sha256.New()
	hash.Write([]byte(issuerURL))
	return hex.EncodeToString(hash.Sum(nil))
}

func NewMetricsManager(issuerURL string, databaseURL string, l logrus.FieldLogger) *MetricsManager {
	l.Info("Setting up telemetry - for more information please visit https://ory.gitbooks.io/hydra/content/telemetry.html")

	mm := &MetricsManager{
		InstanceID:       uuid.New(),
		Segment:          analytics.New("h8dRH3kVCWKkIFWydBmWsyYHR4M0u0vr"),
		Logger:           l,
		issuerURL:        issuerURL,
		databaseURL:      databaseURL,
		MemoryStatistics: &MemoryStatistics{},
		ID:               identify(issuerURL),
		start:            time.Now().UTC(),
		shouldCommit:     shouldCommit(issuerURL, databaseURL),
	}
	return mm
}

func (sw *MetricsManager) RegisterSegment(version, hash, buildTime string) {
	sw.Lock()
	defer sw.Unlock()

	if !sw.shouldCommit {
		sw.Logger.Info("Detected local environment, skipping telemetry commit")
		return
	}

	if err := pkg.Retry(sw.Logger, time.Minute*5, time.Hour, func() error {
		return sw.Segment.Enqueue(analytics.Identify{
			UserId: sw.ID,
			Traits: analytics.NewTraits().
				Set("goarch", runtime.GOARCH).
				Set("goos", runtime.GOOS).
				Set("numCpu", runtime.NumCPU()).
				Set("runtimeVersion", runtime.Version()).
				Set("version", version).
				Set("hash", hash).
				Set("buildTime", buildTime).
				Set("instanceId", sw.InstanceID),
			Context: &analytics.Context{
				IP: net.IPv4(0, 0, 0, 0),
			},
		})
	}); err != nil {
		sw.Logger.WithError(err).Debug("Could not commit anonymized environment information")
	}
	sw.Logger.Debug("Transmitted anonymized environment information")
}

func (sw *MetricsManager) CommitMemoryStatistics() {
	if !sw.shouldCommit {
		sw.Logger.Info("Detected local environment, skipping telemetry commit")
		return
	}

	for {
		sw.MemoryStatistics.Update()
		if err := sw.Segment.Enqueue(analytics.Track{
			UserId:     sw.ID,
			Event:      "stats.memory",
			Properties: analytics.Properties(sw.MemoryStatistics.ToMap()),
			Context:    &analytics.Context{IP: net.IPv4(0, 0, 0, 0)},
		}); err != nil {
			sw.Logger.WithError(err).Debug("Could not commit anonymized telemetry data")
		} else {
			sw.Logger.Debug("Telemetry data transmitted")
		}
		time.Sleep(time.Hour)
	}
}

func (sw *MetricsManager) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	scheme := "https:"
	if r.TLS == nil {
		scheme = "http:"
	}

	start := time.Now().UTC()
	path := anonymizePath(r.URL.Path)

	next(rw, r)

	if !sw.shouldCommit {
		return
	}

	latency := time.Now().UTC().Sub(start) / time.Millisecond

	// Collecting request info
	res := rw.(negroni.ResponseWriter)
	status := res.Status()
	size := res.Size()

	sw.Segment.Enqueue(analytics.Page{
		UserId: sw.ID,
		Name:   path,
		Properties: analytics.
			NewProperties().
			SetURL(scheme + "//" + sw.ID + path).
			SetPath(path).
			SetName(path).
			Set("status", status).
			Set("size", size).
			Set("latency", latency).
			Set("method", r.Method),
		Context: &analytics.Context{IP: net.IPv4(0, 0, 0, 0)},
	})
}

func anonymizePath(path string) string {
	paths := []string{
		client.ClientsHandlerPath,
		jwk.KeyHandlerPath,
		jwk.WellKnownKeysPath,
		oauth2.DefaultConsentPath,
		oauth2.TokenPath,
		oauth2.AuthPath,
		oauth2.UserinfoPath,
		oauth2.WellKnownPath,
		oauth2.IntrospectPath,
		oauth2.RevocationPath,
		oauth2.ConsentRequestPath,
		policy.PolicyHandlerPath,
		"/warden/token/allowed",
		"/warden/allowed",
		group.GroupsHandlerPath,
		"/health/status",
		"/",
	}

	for _, p := range paths {
		if len(path) >= len(p) && strings.ToLower(path[:len(p)]) == strings.ToLower(p) {
			return p
		}
	}

	return ""
}
