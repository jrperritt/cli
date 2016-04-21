package stackcommands

import (
	"github.com/gophercloud/cli/commandoptions"
	"github.com/gophercloud/cli/handler"
	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
	osStacks "github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/openstack/orchestration/v1/stacks"
	"github.com/gophercloud/cli/util"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", "[--name <stackName> | --id <stackID> | --stdin name]"),
	Description: "Retrieve a deployed stack",
	Action:      actionGet,
	Flags:       commandoptions.CommandFlags(flagsGet, nil),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGet, keysGet))
	},
}

func flagsGet() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` or `stdin`isn't provided] The stack name.",
		},
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` or `stdin` isn't provided] The stack id.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: name.",
		},
	}
}

type paramsGet struct {
	stackName string
	stackID   string
}

var keysGet = []string{"Capabilities", "CreationTime", "Description", "DisableRollback", "ID", "Links", "NotificationTopics", "Outputs", "Parameters", "Name", "Status", "StatusReason", "Tags", "TemplateDescription", "Timeout", "UpdatedTime"}

type commandGet handler.Command

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

func (command *commandGet) HandlePipe(resource *handler.Resource, item string) error {
	name, id, err := IDAndName(command.Ctx.ServiceClient, item, "")
	if err != nil {
		return err
	}
	resource.Params = &paramsGet{
		stackName: name,
		stackID:   id,
	}
	return nil
}

func (command *commandGet) HandleSingle(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	name := c.String("name")
	id := c.String("id")
	name, id, err := IDAndName(command.Ctx.ServiceClient, name, id)
	if err != nil {
		return err
	}
	resource.Params = &paramsGet{
		stackName: name,
		stackID:   id,
	}
	return nil
}

func (command *commandGet) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsGet)
	stackName := params.stackName
	stackID := params.stackID

	stack, err := osStacks.Get(command.Ctx.ServiceClient, stackName, stackID).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = stack
}

func (command *commandGet) PreCSV(resource *handler.Resource) error {
	resource.Result = stackSingle(resource.Result)
	resource.FlattenMap("Parameters")
	resource.FlattenMap("Outputs")
	resource.FlattenMap("Links")
	resource.FlattenMap("NotificationTopics")
	resource.FlattenMap("Capabilities")
	return nil
}

func (command *commandGet) PreTable(resource *handler.Resource) error {
	return command.PreCSV(resource)
}

func (command *commandGet) StdinField() string {
	return "name"
}
