package helper

import (
	"net/http"

	"github.com/moul/http2curl"
	"github.com/pkg/errors"
)

func DoDryRequest(dry bool, req *http.Request) error {
	if dry {
		command, err := http2curl.GetCurlCommand(req)
		if err != nil {
			return errors.Wrap(err, "")
		}

		return errors.Errorf("Because you are using the dry option, the request will not be executed. The curl equivalent of this command is: \n\n%s\n", command)
	}
	return nil
}
