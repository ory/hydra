package contextx

import "context"

var RootContext = context.WithValue(context.Background(), "root", true) //nolint:staticcheck
