/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package metricsx

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"time"

	"io"
	"math/rand"

	"github.com/pborman/uuid"
	"github.com/segmentio/analytics-go"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

type MetricsManager struct {
	sync.RWMutex    `json:"-"`
	start           time.Time          `json:"-"`
	Segment         analytics.Client   `json:"-"`
	Logger          logrus.FieldLogger `json:"-"`
	shouldCommit    bool               `json:"-"`
	salt            string
	whitelistedURLs []string  `json:"-"`
	sampling        float64   `json:"-"`
	rng             io.Reader `json:"-"`

	ID               string            `json:"id"`
	UpTime           int64             `json:"uptime"`
	MemoryStatistics *MemoryStatistics `json:"memory"`
	BuildVersion     string            `json:"buildVersion"`
	BuildHash        string            `json:"buildHash"`
	BuildTime        string            `json:"buildTime"`
	InstanceID       string            `json:"instanceId"`
	ServiceName      string            `json:"serviceName"`
}

func Hash(value string) string {
	hash := sha256.New()
	hash.Write([]byte(value))
	return hex.EncodeToString(hash.Sum(nil))
}

func NewMetricsManager(
	id string,
	enable bool,
	writeKey string,
	whitelistedURLs []string,
	logger logrus.FieldLogger,
	serviceName string,
	sampling float64,
	endpoint string,
) *MetricsManager {
	segment, err := analytics.NewWithConfig(writeKey, analytics.Config{
		Interval:  time.Hour * 24,
		BatchSize: 100,
	})

	if err != nil {
		logger.WithError(err).Fatalf("Unable to initialise segment.")
		return nil
	}

	mm := &MetricsManager{
		InstanceID:       uuid.New(),
		Segment:          segment,
		Logger:           logger,
		MemoryStatistics: &MemoryStatistics{},
		ID:               id,
		start:            time.Now().UTC(),
		salt:             uuid.New(),
		shouldCommit:     enable,
		whitelistedURLs:  whitelistedURLs,
		ServiceName:      serviceName,
		sampling:         sampling,
	}
	return mm
}

func NewMetricsManagerWithConfig(
	id string,
	enable bool,
	writeKey string,
	whitelistedURLs []string,
	logger logrus.FieldLogger,
	serviceName string,
	sampling float64,
	config analytics.Config,
) *MetricsManager {
	segment, err := analytics.NewWithConfig(writeKey, config)

	if err != nil {
		logger.WithError(err).Fatalf("Unable to initialise segment.")
		return nil
	}

	mm := &MetricsManager{
		InstanceID:       uuid.New(),
		Segment:          segment,
		Logger:           logger,
		MemoryStatistics: &MemoryStatistics{},
		ID:               id,
		start:            time.Now().UTC(),
		salt:             uuid.New(),
		shouldCommit:     enable,
		whitelistedURLs:  whitelistedURLs,
		ServiceName:      serviceName,
		sampling:         sampling,
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

	if err := Retry(sw.Logger, time.Minute*5, time.Hour, func() error {
		return sw.Segment.Enqueue(analytics.Identify{
			UserId: sw.ID,
			Traits: analytics.NewTraits().
				Set("goarch", runtime.GOARCH).
				Set("goos", runtime.GOOS).
				Set("numCpu", runtime.NumCPU()).
				Set("runtimeVersion", runtime.Version()).
				Set("version", version).
				Set("Hash", hash).
				Set("buildTime", buildTime).
				Set("service", sw.ServiceName).
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
			Event:      "memstats",
			Properties: analytics.Properties(sw.MemoryStatistics.ToMap()),
			Context:    &analytics.Context{IP: net.IPv4(0, 0, 0, 0)},
		}); err != nil {
			sw.Logger.WithError(err).Debug("Could not commit anonymized telemetry data")
		} else {
			sw.Logger.Debug("Telemetry data transmitted")
		}
		time.Sleep(time.Hour * 24)
	}
}

func (sw *MetricsManager) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, r)

	if !sw.shouldCommit || rand.Float64() > sw.sampling {
		return
	}

	latency := time.Now().UTC().Sub(start.UTC()) / time.Millisecond

	scheme := "https:"
	if r.TLS == nil {
		scheme = "http:"
	}

	path := sw.anonymizePath(r.URL.Path, sw.salt)
	query := sw.anonymizeQuery(r.URL.Query(), sw.salt)

	// Collecting request info
	res := rw.(negroni.ResponseWriter)
	status := res.Status()
	size := res.Size()

	sw.Segment.Enqueue(analytics.Page{
		UserId: sw.ID,
		Name:   path,
		Properties: analytics.
			NewProperties().
			SetURL(scheme+"//"+sw.ID+path+"?"+query).
			SetPath(path).
			SetName(path).
			Set("status", status).
			Set("size", size).
			Set("latency", latency).
			Set("instance", sw.InstanceID).
			Set("service", sw.ServiceName).
			Set("method", r.Method),
		Context: &analytics.Context{IP: net.IPv4(0, 0, 0, 0)},
	})
}

func (sw *MetricsManager) anonymizePath(path string, salt string) string {
	paths := sw.whitelistedURLs
	path = strings.ToLower(path)

	for _, p := range paths {
		p = strings.ToLower(p)
		if len(path) == len(p) && path[:len(p)] == strings.ToLower(p) {
			return p
		} else if len(path) > len(p) && path[:len(p)+1] == strings.ToLower(p)+"/" {
			return path[:len(p)] + "/" + Hash(path[len(p):]+"|"+salt)
		}
	}

	return ""
}

func (sw *MetricsManager) anonymizeQuery(query url.Values, salt string) string {
	for _, q := range query {
		for i, s := range q {
			if s != "" {
				s = Hash(s + "|" + salt)
				q[i] = s
			}
		}
	}
	return query.Encode()
}
