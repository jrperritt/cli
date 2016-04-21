package securitygroupcommands

import (
	"github.com/gophercloud/cli/commandoptions"
	"github.com/gophercloud/cli/handler"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
	securityGroups "github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
)

var create = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", "--name <securityGroupName>"),
	Description: "Creates a security group",
	Action:      actionCreate,
	Flags:       commandoptions.CommandFlags(flagsCreate, keysCreate),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsCreate, keysCreate))
	},
}

func flagsCreate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The name for the security group.",
		},
		cli.StringFlag{
			Name:  "description",
			Usage: "[optional] A description for the security group.",
		},
	}
}

var keysCreate = []string{"ID", "Name"}

type paramsCreate struct {
	opts *osSecurityGroups.CreateOpts
}

type commandCreate handler.Command

func actionCreate(c *cli.Context) {
	command := &commandCreate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandCreate) Context() *handler.Context {
	return command.Ctx
}

func (command *commandCreate) Keys() []string {
	return keysCreate
}

func (command *commandCreate) ServiceClientType() string {
	return serviceClientType
}

func (command *commandCreate) HandleFlags(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}

	c := command.Ctx.CLIContext
	opts := &osSecurityGroups.CreateOpts{
		Name:        c.String("name"),
		Description: c.String("description"),
	}

	resource.Params = &paramsCreate{
		opts: opts,
	}

	return nil
}

func (command *commandCreate) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsCreate).opts
	securityGroup, err := securityGroups.Create(command.Ctx.ServiceClient, *opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = securityGroupSingle(securityGroup)
}
