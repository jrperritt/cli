package keypaircommands

import "github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"

var commandPrefix = "servers keypair"
var serviceClientType = "compute"

// Get returns all the commands allowed for a `compute keypairs` request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		remove,
		upload,
		generate,
	}
}
