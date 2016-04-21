package buildinfocommands

import (
	"github.com/gophercloud/cli/commandoptions"
	"github.com/gophercloud/cli/handler"
	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
	"github.com/gophercloud/cli/vendor/github.com/fatih/structs"
	"github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/openstack/orchestration/v1/buildinfo"
	"github.com/gophercloud/cli/util"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", ""),
	Description: "Retrieve build information for a Heat deployment",
	Action:      actionGet,
	Flags:       commandoptions.CommandFlags(flagsGet, nil),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGet, keysGet))
	},
}

var keysGet = []string{"API", "Engine", "FusionAPI"}

type commandGet handler.Command

func flagsGet() []cli.Flag {
	return []cli.Flag{}
}

func actionGet(c *cli.Context) {
	command := &commandGet{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandGet) Context() *handler.Context {
	return command.Ctx
}

func (command *commandGet) Keys() []string {
	return keysGet
}

func (command *commandGet) ServiceClientType() string {
	return serviceClientType
}

func (command *commandGet) HandleFlags(resource *handler.Resource) error {
	return nil
}

func (command *commandGet) Execute(resource *handler.Resource) {
	info, err := buildinfo.Get(command.Ctx.ServiceClient).Extract()
	if err != nil {
		resource.Err = err
	}
	result := structs.Map(info)
	for k, v := range result {
		result[k] = v.(map[string]interface{})["Revision"]
	}
	resource.Result = result
}
