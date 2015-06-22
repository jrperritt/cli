package keypaircommands

import (
	"github.com/codegangsta/cli"
	"github.com/fatih/structs"
	"github.com/jrperritt/rack/handler"
	"github.com/jrperritt/rack/util"
	osKeypairs "github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/rackspace/gophercloud/rackspace/compute/v2/keypairs"
)

var list = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Lists keypairs",
	Action:      actionList,
	Flags:       util.CommandFlags(flagsList, keysList),
	BashComplete: func(c *cli.Context) {
		util.CompleteFlags(util.CommandFlags(flagsList, keysList))
	},
}

func flagsList() []cli.Flag {
	return []cli.Flag{}
}

var keysList = []string{"Name", "Fingerprint"}

type paramsList struct{}

type commandList handler.Command

func actionList(c *cli.Context) {
	command := &commandList{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandList) Context() *handler.Context {
	return command.Ctx
}

func (command *commandList) Keys() []string {
	return keysList
}

func (command *commandList) ServiceClientType() string {
	return serviceClientType
}

func (command *commandList) HandleFlags(resource *handler.Resource) error {
	return nil
}

func (command *commandList) HandleSingle(resource *handler.Resource) error {
	return nil
}

func (command *commandList) Execute(resource *handler.Resource) {
	allPages, err := keypairs.List(command.Ctx.ServiceClient).AllPages()
	if err != nil {
		resource.Err = err
		return
	}
	keys, err := osKeypairs.ExtractKeyPairs(allPages)
	if err != nil {
		resource.Err = err
		return
	}
	result := make([]map[string]interface{}, len(keys))
	for j, key := range keys {
		result[j] = structs.Map(key)
	}
	resource.Result = result
}
