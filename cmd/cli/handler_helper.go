package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ory/hydra/pkg"
	hydra "github.com/ory/hydra/sdk/go/hydra/swagger"
)

func checkResponse(response *hydra.APIResponse, err error, expectedStatusCode int) {
	pkg.Must(err, "Could not validate token: %s", err)

	if response.StatusCode != expectedStatusCode {
		fmt.Printf("Command failed because status code %d was expeceted but code %d was received.", expectedStatusCode, response.StatusCode)
		os.Exit(1)
		return
	}
}

func formatResponse(response interface{}) string {
	out, err := json.MarshalIndent(response, "", "\t")
	pkg.Must(err, `Command failed because an error ("%s") occurred while prettifying output.`, err)
	return string(out)
}
