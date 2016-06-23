Slack API in Go [![GoDoc](https://godoc.org/github.com/nlopes/slack?status.svg)](https://godoc.org/github.com/nlopes/slack) [![Build Status](https://travis-ci.org/nlopes/slack.svg)](https://travis-ci.org/nlopes/slack)
===============

This library supports most if not all of the `api.slack.com` REST
calls, as well as the Real-Time Messaging protocol over websocket, in
a fully managed way.


Note: If you just updated from master and it broke your implementation, please check [0.0.1](https://github.com/nlopes/slack/releases/tag/v0.0.1)

## Installing

### *go get*

    $ go get github.com/nlopes/slack

## Example

### Getting all groups

    import (
		"fmt"

		"github.com/nlopes/slack"
	)

    func main() {
		api := slack.New("YOUR_TOKEN_HERE")
		// If you set debugging, it will log all requests to the console
		// Useful when encountering issues
		// api.SetDebug(true)
		groups, err := api.GetGroups(false)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		for _, group := range groups {
			fmt.Printf("ID: %s, Name: %s\n", group.ID, group.Name)
		}
	}

### Getting User Information

    import (
	    "fmt"

	    "github.com/nlopes/slack"
    )

    func main() {
	    api := slack.New("YOUR_TOKEN_HERE")
	    user, err := api.GetUserInfo("U023BECGF")
	    if err != nil {
		    fmt.Printf("%s\n", err)
		    return
	    }
	    fmt.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
    }

## Minimal RTM usage:

See https://github.com/nlopes/slack/blob/master/examples/websocket/websocket.go


## Contributing

You are more than welcome to contribute to this project.  Fork and
make a Pull Request, or create an Issue if you see any problem.

## License

BSD 2 Clause license
