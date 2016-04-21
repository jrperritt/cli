package stackcommands

import (
	"io/ioutil"
	"path/filepath"

	"github.com/gophercloud/cli/commandoptions"
	"github.com/gophercloud/cli/handler"
	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
	osStacks "github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/openstack/orchestration/v1/stacks"
	"github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/openstack/orchestration/v1/stacks"
	"github.com/gophercloud/cli/util"
)

var adopt = cli.Command{
	Name:        "adopt",
	Usage:       util.Usage(commandPrefix, "adopt", "--name <stackName> --adopt-file <adoptFile>"),
	Description: "Creates a stack from existing resources",
	Action:      actionAdopt,
	Flags:       commandoptions.CommandFlags(flagsAdopt, nil),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsAdopt, keysAdopt))
	},
}

func flagsAdopt() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The stack name.",
		},
		cli.StringFlag{
			Name:  "adopt-file",
			Usage: "[required] Path to file specifying data required to adopt existing resources",
		},
		cli.StringFlag{
			Name:  "environment-file",
			Usage: "[optional] Path to the file or URL containing environment for the stack",
		},
		cli.IntFlag{
			Name:  "timeout",
			Usage: "[optional] Stack creation timeout in minutes.",
		},
		cli.BoolTFlag{
			Name:  "disable-rollback",
			Usage: "[optional] Disable rollback on create/update failure.",
		},
		cli.StringFlag{
			Name:  "parameters",
			Usage: "[optional] A comma-separated string of key=value pairs.",
		},
	}
}

type paramsAdopt struct {
	opts *osStacks.AdoptOpts
}

var keysAdopt = keysCreate

type commandAdopt handler.Command

func actionAdopt(c *cli.Context) {
	command := &commandAdopt{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandAdopt) Context() *handler.Context {
	return command.Ctx
}

func (command *commandAdopt) Keys() []string {
	return keysAdopt
}

func (command *commandAdopt) ServiceClientType() string {
	return serviceClientType
}

func (command *commandAdopt) HandleFlags(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name", "adopt-file"})
	if err != nil {
		return err
	}
	c := command.Ctx.CLIContext
	opts := &osStacks.AdoptOpts{
		Name: c.String("name"),
	}

	abs, err := filepath.Abs(c.String("adopt-file"))
	if err != nil {
		return err
	}
	adoptData, err := ioutil.ReadFile(abs)
	if err != nil {
		return err
	}
	opts.AdoptStackData = string(adoptData)

	if c.IsSet("disable-rollback") {
		disableRollback := c.Bool("disable-rollback")
		opts.DisableRollback = &disableRollback
	}

	if c.IsSet("environment-file") {
		opts.EnvironmentOpts = new(osStacks.Environment)
		opts.EnvironmentOpts.URL = c.String("environment-file")
	}

	if c.IsSet("timeout") {
		opts.Timeout = c.Int("timeout")
	}

	if c.IsSet("parameters") {
		parameters, err := command.Ctx.CheckKVFlag("parameters")
		if err != nil {
			return err
		}
		opts.Parameters = parameters
	}

	resource.Params = &paramsAdopt{
		opts: opts,
	}
	return nil
}

func (command *commandAdopt) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsAdopt).opts
	result, err := stacks.Adopt(command.Ctx.ServiceClient, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = result
}

func (command *commandAdopt) PreCSV(resource *handler.Resource) error {
	resource.Result = stackSingle(resource.Result)
	resource.FlattenMap("Links")
	return nil
}

func (command *commandAdopt) PreTable(resource *handler.Resource) error {
	return command.PreCSV(resource)
}
