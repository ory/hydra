package cli

import (
	"context"
	"fmt"
	"github.com/ory/hydra/driver"
	"github.com/ory/hydra/driver/configuration"
	"github.com/ory/viper"
	"github.com/ory/x/flagx"
	"github.com/ory/x/logrusx"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

type MigrateHandlerFizz struct{}

func (h *MigrateHandlerFizz) MigrateSQL(cmd *cobra.Command, args []string) {
	var d driver.Driver

	if flagx.MustGetBool(cmd, "read-from-env") {
		d = driver.NewDefaultDriver(logrusx.New(), false, nil, "", "", "", false)
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
		d = driver.NewDefaultDriver(logrusx.New(), false, nil, "", "", "", false)
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

	// convert migration tables
	if err := migrateOldMigrationTables(conn); err != nil {
		fmt.Printf("Could not convert the migration table:\n%v\n", err)
		os.Exit(1)
		return
	}

	// print migration status
	fmt.Println("The following migration is planned:")
	fmt.Println("")
	if err := p.MigrationStatus(context.Background(), os.Stdout); err != nil {
		fmt.Printf("Could not get the migration status:\n%v\n", errors.WithStack(err))
		os.Exit(1)
		return
	}

	if !flagx.MustGetBool(cmd, "yes") {
		fmt.Println("")
		fmt.Println("To skip the next question use flag --yes (at your own risk).")
		if !askForConfirmation("Do you wish to execute this migration plan?") {
			fmt.Println("Migration aborted.")
			return
		}
	}

	// apply migrations
	if err := p.MigrateUp(context.Background()); err != nil {
		fmt.Printf("Could not apply migrations:\n%v\n", errors.WithStack(err))
	}

	fmt.Println("Successfully applied migrations!")

	return
}
