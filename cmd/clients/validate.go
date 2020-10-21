package clients

import (
	"bytes"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/ory/jsonschema/v3"
	"github.com/ory/x/cmdx"
	"github.com/ory/x/viperx"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
)

var clientsValidateCmd = &cobra.Command{
	Use:   "validate <client.json> [<client2.json> ...]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Validate client JSON files",
	Long:  fmt.Sprintf("Validate client JSON files that can be used with other commands. %s", helperStdInFile),
	RunE: func(cmd *cobra.Command, args []string) (retErr error) {
		for _, fn := range args {
			var r io.Reader
			if fn == "-" {
				r = cmd.InOrStdin()
			} else {
				var err error
				r, err = os.Open(fn)
				if err != nil {
					_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not open %s: %s\n", fn, err)
					retErr = cmdx.FailSilently(cmd)
					continue
				}
			}

			data, err := ioutil.ReadAll(r)
			if err != nil {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Could not read %s: %s\n", fn, err)
				retErr = cmdx.FailSilently(cmd)
				continue
			}

			if err := validateClient(cmd, fn, string(data)); errors.Is(err, cmdx.ErrNoPrintButFail) {
				retErr = err
			} else if err != nil {
				// unexpected error should fail the command immediately
				return err
			}
		}

		return
	},
}

var (
	schemasBox = packr.New("schemas", "../../.schema")
	schemas    = make(map[string]*jsonschema.Schema)
)

const (
	apiSwaggerSchemaName   = "api.swagger.json"
	modelsOAuth2ClientPath = "#/definitons/oAuth2Client"
)

func validateClient(cmd *cobra.Command, src, client string) error {
	swaggerSchema, ok := schemas[modelsOAuth2ClientPath]
	if !ok {
		// get swagger schema
		sf, err := schemasBox.Open(apiSwaggerSchemaName)
		if err != nil {
			return errors.Wrap(err, "Could not open swagger schema. This is an error with the binary you use and should be reported. Thanks ;)")
		}

		// add swagger schema
		schemaCompiler := jsonschema.NewCompiler()
		err = schemaCompiler.AddResource(apiSwaggerSchemaName, sf)
		if err != nil {
			return errors.Wrap(err, "Could not add swagger schema to the schema compiler. This is an error with the binary you use and should be reported. Thanks ;)")
		}

		// compile swagger payload definition
		swaggerSchema, err = schemaCompiler.Compile(apiSwaggerSchemaName + modelsOAuth2ClientPath)
		if err != nil {
			return errors.Wrap(err, "Could not compile the identity schema. This is an error with the binary you use and should be reported. Thanks ;)")
		}
		// force additional properties to false because swagger does not render this
		swaggerSchema.AdditionalProperties = false
		schemas[modelsOAuth2ClientPath] = swaggerSchema
	}


	// validate against swagger definition
	var foundValidationErrors bool
	err := swaggerSchema.Validate(bytes.NewBufferString(client))
	if err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s: not valid\n", src)
		viperx.PrintHumanReadableValidationErrors(cmd.ErrOrStderr(), err)
		foundValidationErrors = true
	}

	if foundValidationErrors {
		return cmdx.FailSilently(cmd)
	}
	return nil
}
