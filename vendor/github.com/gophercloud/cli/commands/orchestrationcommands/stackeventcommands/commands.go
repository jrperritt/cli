package stackeventcommands

import "github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"

var commandPrefix = "orchestration event"
var serviceClientType = "orchestration"

// Get returns all the commands allowed for an `orchestration event` request.
func Get() []cli.Command {
	return []cli.Command{
		get,
	}
}
