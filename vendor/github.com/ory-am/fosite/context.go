package fosite

import "golang.org/x/net/context"

func NewContext() context.Context {
	return context.Background()
}
