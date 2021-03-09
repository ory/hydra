package cli

import (
	"context"
	"fmt"
	"github.com/ory/hydra/driver"
	"github.com/ory/x/configx"
	"github.com/ory/x/flagx"
	"github.com/spf13/cobra"
	"os"
	"time"
)

type JanitorHandler struct{}

func newJanitorHandler() *JanitorHandler {
	return &JanitorHandler{}
}

func (j *JanitorHandler) Purge(cmd *cobra.Command, args []string) {
	var d driver.Registry

	if flagx.MustGetBool(cmd, "read-from-env") {
		d = driver.New(
			driver.WithOptions(
				configx.SkipValidation(),
				configx.WithFlags(cmd.Flags()),
			),
			driver.DisableValidation(),
			driver.DisablePreloading())

		if len(d.Config().DSN()) == 0 {
			fmt.Println(fmt.Sprintf("%s\n%s", cmd.UsageString(),
				"When using flag -e, environment variable DSN must be set"))
			os.Exit(1)
			return
		}
	}

	p := d.Persister()

	var notAfter time.Time

	if keepYounger := flagx.MustGetString(cmd, ""); keepYounger != "" {
		// TODO: get configx/provider DurationF here to verify string to time.Time parse.
		d, _ := time.ParseDuration(keepYounger)
		notAfter = time.Now().Add(-d)
	} else {
		notAfter = time.Now()
	}

	conn := p.Connection(context.Background())

	if conn == nil {
		fmt.Println(fmt.Sprintf("%s\n%s\n", cmd.UsageString(),
			"Janitor can only be executed against a SQL-compatible driver but DSN is not a SQL source."))
		os.Exit(1)
		return
	}

	if err := conn.Open(); err != nil {
		fmt.Printf("Could not open the database connection:\n%+v\n", err)
		os.Exit(1)
		return
	}

	ctx := context.Background()

	if err := p.FlushInactiveAccessTokens(ctx, notAfter); err != nil {

	}

	if err := p.FlushInactiveLoginConsentRequests(ctx, notAfter); err != nil {

	}
}
