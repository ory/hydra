package prometheus

// Outputs Prometheus metrics
//
// swagger:route GET /metrics/prometheus admin prometheus
//
// Get snapshot metrics from the Hydra service. If you're using k8s, you can then add annotations to
// your deployment like so:
//
// ```
// metadata:
//  annotations:
//    prometheus.io/port: "4445"
//      prometheus.io/path: "/metrics/prometheus"
// ```
//
//     Produces:
//     - plain/text
//
//     Responses:
//       200: emptyResponse
func swaggerPublicPrometheus() {}
