package securitygrouprulecommands

import (
	"fmt"

	"github.com/gophercloud/cli/commandoptions"
	"github.com/gophercloud/cli/handler"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
	securityGroupRules "github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
)

var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", "[--id <securityGroupRuleID> | --stdin id]"),
	Description: "Deletes an existing security group rule",
	Action:      actionDelete,
	Flags:       commandoptions.CommandFlags(flagsDelete, keysDelete),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsDelete, keysDelete))
	},
}

func flagsDelete() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` isn't provided] The ID of the security group rule.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
}

var keysDelete = []string{}

type paramsDelete struct {
	securityGroupRuleID string
}

type commandDelete handler.Command

func actionDelete(c *cli.Context) {
	command := &commandDelete{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandDelete) Context() *handler.Context {
	return command.Ctx
}

func (command *commandDelete) Keys() []string {
	return keysDelete
}

func (command *commandDelete) ServiceClientType() string {
	return serviceClientType
}

func (command *commandDelete) HandleFlags(resource *handler.Resource) error {
	resource.Params = &paramsDelete{}
	return nil
}

func (command *commandDelete) HandlePipe(resource *handler.Resource, item string) error {
	resource.Params.(*paramsDelete).securityGroupRuleID = item
	return nil
}

func (command *commandDelete) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"id"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsDelete).securityGroupRuleID = command.Ctx.CLIContext.String("id")
	return nil
}

func (command *commandDelete) Execute(resource *handler.Resource) {
	securityGroupRuleID := resource.Params.(*paramsDelete).securityGroupRuleID
	err := securityGroupRules.Delete(command.Ctx.ServiceClient, securityGroupRuleID).ExtractErr()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = fmt.Sprintf("Successfully deleted security group rule [%s]\n", securityGroupRuleID)
}

func (command *commandDelete) StdinField() string {
	return "id"
}
