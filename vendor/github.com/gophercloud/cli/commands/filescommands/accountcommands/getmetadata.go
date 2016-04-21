package accountcommands

import (
	"github.com/gophercloud/cli/commandoptions"
	"github.com/gophercloud/cli/handler"
	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
	"github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/openstack/objectstorage/v1/accounts"
	"github.com/gophercloud/cli/util"
)

var getMetadata = cli.Command{
	Name:        "get-metadata",
	Usage:       util.Usage(commandPrefix, "get-metadata", ""),
	Description: "Get metadata associated with the account.",
	Action:      actionGetMetadata,
	Flags:       commandoptions.CommandFlags(flagsGetMetadata, keysGetMetadata),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGetMetadata, keysGetMetadata))
	},
}

func flagsGetMetadata() []cli.Flag {
	return []cli.Flag{}
}

var keysGetMetadata = []string{"Metadata"}

type paramsGetMetadata struct {
	containerName string
}

type commandGetMetadata handler.Command

func actionGetMetadata(c *cli.Context) {
	command := &commandGetMetadata{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandGetMetadata) Context() *handler.Context {
	return command.Ctx
}

func (command *commandGetMetadata) Keys() []string {
	return keysGetMetadata
}

func (command *commandGetMetadata) ServiceClientType() string {
	return serviceClientType
}

func (command *commandGetMetadata) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsGetMetadata{}
	return nil
}

func (command *commandGetMetadata) Execute(resource *handler.Resource) {
	metadata, err := accounts.Get(command.Ctx.ServiceClient).ExtractMetadata()
	if err != nil {
		resource.Err = err
		return
	}

	resource.Result = metadata
}

func (command *commandGetMetadata) PreCSV(resource *handler.Resource) error {
	resource.Result = map[string]interface{}{
		"Metadata": resource.Result,
	}
	resource.Keys = []string{"Metadata"}
	resource.FlattenMap("Metadata")
	return nil
}

func (command *commandGetMetadata) PreTable(resource *handler.Resource) error {
	return command.PreCSV(resource)
}
