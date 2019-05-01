// The metrics package is an abstraction of different metrics and analytics packages using Prometheus.
package metrics

import "context"

// A Bridge connects Prometheus metrics to a push-based analytics system
type Bridge interface {
	Push(context.Context) error
}
