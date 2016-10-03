package flavor

import (
	"github.com/gophercloud/cli/lib/traits"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "compute flavor"

type FlavorV2Command struct {
	traits.Commandable
	traits.Computeable
}

// Get returns all the commands allowed for a `compute flavor` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
	}
}
