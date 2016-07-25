package helper

import (
	"github.com/moul/http2curl"
	"net/http"
	"github.com/go-errors/errors"
)

func DoDryRequest(dry bool, req *http.Request) error {
	if dry {
		command, err := http2curl.GetCurlCommand(req)
		if err != nil {
			return errors.New(err)
		}

		return errors.Errorf("Because you are using the dry option, the request will not be executed. The curl equivalent of this command is: \n\n%s\n", command)
	}
	return nil
}