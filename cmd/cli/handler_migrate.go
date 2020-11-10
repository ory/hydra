package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/ory/x/errorsx"

	"github.com/ory/x/cmdx"

	"github.com/spf13/cobra"

	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/viper"
	"github.com/ory/x/flagx"
	"github.com/ory/x/logrusx"
)

type MigrateHandler struct{}

func newMigrateHandler() *MigrateHandler {
	return &MigrateHandler{}
}

func (h *MigrateHandler) MigrateSQL(cmd *cobra.Command, args []string) {
	var d driver.Driver

	if flagx.MustGetBool(cmd, "read-from-env") {
		d = driver.NewDefaultDriver(logrusx.New("", ""), false, nil, "", "", "", false)
		if len(d.Configuration().DSN()) == 0 {
			fmt.Println(cmd.UsageString())
			fmt.Println("")
			fmt.Println("When using flag -e, environment variable DSN must be set")
			os.Exit(1)
			return
		}
	} else {
		if len(args) != 1 {
			fmt.Println(cmd.UsageString())
			os.Exit(1)
			return
		}
		viper.Set(configuration.ViperKeyDSN, args[0])
		d = driver.NewDefaultDriver(logrusx.New("", ""), false, nil, "", "", "", false)
	}

	p := d.Registry().Persister()
	conn := p.Connection(context.Background())
	if conn == nil {
		fmt.Println(cmd.UsageString())
		fmt.Println("")
		fmt.Printf("Migrations can only be executed against a SQL-compatible driver but DSN is not a SQL source.\n")
		os.Exit(1)
		return
	}

	if err := conn.Open(); err != nil {
		fmt.Printf("Could not open the database connection:\n%+v\n", err)
		os.Exit(1)
		return
	}

	// convert migration tables
	if err := p.PrepareMigration(context.Background()); err != nil {
		fmt.Printf("Could not convert the migration table:\n%+v\n", err)
		os.Exit(1)
		return
	}

	// print migration status
	fmt.Println("The following migration is planned:")
	fmt.Println("")
	if err := p.MigrationStatus(context.Background(), os.Stdout); err != nil {
		fmt.Printf("Could not get the migration status:\n%+v\n", errorsx.WithStack(err))
		os.Exit(1)
		return
	}

	if !flagx.MustGetBool(cmd, "yes") {
		fmt.Println("")
		fmt.Println("To skip the next question use flag --yes (at your own risk).")
		if !cmdx.AskForConfirmation("Do you wish to execute this migration plan?", nil, nil) {
			fmt.Println("Migration aborted.")
			return
		}
	}

	// apply migrations
	if err := p.MigrateUp(context.Background()); err != nil {
		fmt.Printf("Could not apply migrations:\n%+v\n", errorsx.WithStack(err))
	}

	fmt.Println("Successfully applied migrations!")
}
