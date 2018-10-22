package cmdx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func Must(err error, message string, args ...interface{}) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, message+"\n", args...)
	os.Exit(1)
}

func CheckResponse(err error, expectedStatusCode int, response *http.Response) {
	Must(err, "Command failed because error \"%s\" occurred.\n", err)

	if response.StatusCode != expectedStatusCode {
		out, _ := ioutil.ReadAll(response.Body)
		Fatalf(
			`Command failed because status code %d was expected but code %d was received.

Response payload:

%s`,
			expectedStatusCode,
			response.StatusCode,
			FormatResponse(json.RawMessage(out)),
		)
	}
}

func FormatResponse(response interface{}) string {
	out, err := json.MarshalIndent(response, "", "\t")
	Must(err, `Command failed because an error ("%s") occurred while prettifying output.`, err)
	return string(out)
}

func Fatalf(message string, args ...interface{}) {
	if len(args) > 0 {
		fmt.Fprintf(os.Stderr, message+"\n", args...)
	} else {
		fmt.Fprintln(os.Stderr, message)
	}
	os.Exit(1)
}
