package prometheus

// Outputs Prometheus metrics
//
// swagger:route GET /metrics/prometheus admin prometheus
//
// Get Snapshot Metrics from the Hydra Service.
//
// If you're using k8s, you can then add annotations to your deployment like so:
//
// ```
// metadata:
//  annotations:
//    prometheus.io/port: "4445"
//      prometheus.io/path: "/metrics/prometheus"
// ```
//
// If the service supports TLS Edge Termination, this endpoint does not require the
// `X-Forwarded-Proto` header to be set.
//
//     Produces:
//     - plain/text
//
//     Responses:
//       200: emptyResponse
func swaggerPublicPrometheus() {}
