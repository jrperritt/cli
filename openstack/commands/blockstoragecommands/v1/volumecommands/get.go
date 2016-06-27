package volumecommands

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumes"
)

var _ lib.PipeCommander = new(commandGet)

type commandGet struct {
	openstack.CommandUtil
	VolumeV1Command
}

var get = cli.Command{
	Name:         "get",
	Usage:        util.Usage(commandPrefix, "get", "[--id <volumeID> | --name <volumeName> | --stdin id]"),
	Description:  "Gets a volume",
	Action:       actionGet,
	Flags:        openstack.CommandFlags(flagsGet, []string{}),
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsGet) },
}

func actionGet(ctx *cli.Context) {
	c := new(commandGet)
	c.Context = ctx
	lib.Run(ctx, c)
}

var flagsGet = []cli.Flag{
	cli.StringFlag{
		Name:  "id",
		Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the volume.",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the volume.",
	},
	cli.StringFlag{
		Name:  "stdin",
		Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id",
	},
}

func (c *commandGet) HandleFlags() error {
	return nil
}

func (c *commandGet) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *commandGet) HandleSingle() (interface{}, error) {
	id, err := c.IDOrName(volumes.IDFromName)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (c *commandGet) Execute(item interface{}, out chan interface{}) {
	defer func() {
		close(out)
	}()
	var m map[string]map[string]interface{}
	err := volumes.Get(c.ServiceClient, item.(string)).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["volume"]
	default:
		out <- err
	}
}

func (c *commandGet) PipeFieldOptions() []string {
	return []string{"id"}
}

/*
func (c *commandGet) PreCSV() error {
	resource.FlattenMap("Metadata")
	resource.FlattenMap("Attachments")
	return nil
}

func (c *commandGet) PreTable() error {
	return command.PreCSV(resource)
}
*/
