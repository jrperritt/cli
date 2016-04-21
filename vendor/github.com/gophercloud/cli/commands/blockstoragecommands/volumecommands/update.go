package volumecommands

import (
	"github.com/gophercloud/cli/commandoptions"
	"github.com/gophercloud/cli/handler"
	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
	osVolumes "github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumes"
	"github.com/gophercloud/cli/util"
)

var update = cli.Command{
	Name:        "update",
	Usage:       util.Usage(commandPrefix, "update", "[--id <volumeID> | --name <volumeName>]"),
	Description: "Updates a volume",
	Action:      actionUpdate,
	Flags:       commandoptions.CommandFlags(flagsUpdate, keysUpdate),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsUpdate, keysUpdate))
	},
}

func flagsUpdate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` isn't provided] The ID of the volume.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` isn't provided] The name of the volume.",
		},
		cli.StringFlag{
			Name:  "rename",
			Usage: "[optional] A new name for this volume.",
		},
		cli.StringFlag{
			Name:  "description",
			Usage: "[optional] A new description for this volume.",
		},
	}
}

var keysUpdate = []string{"ID", "Name", "Description", "Size", "VolumeType", "SnapshotID", "Attachments"}

type paramsUpdate struct {
	volumeID string
	opts     *osVolumes.UpdateOpts
}

type commandUpdate handler.Command

func actionUpdate(c *cli.Context) {
	command := &commandUpdate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandUpdate) Context() *handler.Context {
	return command.Ctx
}

func (command *commandUpdate) Keys() []string {
	return keysUpdate
}

func (command *commandUpdate) ServiceClientType() string {
	return serviceClientType
}

func (command *commandUpdate) HandleFlags(resource *handler.Resource) error {
	volumeID, err := command.Ctx.IDOrName(osVolumes.IDFromName)
	if err != nil {
		return err
	}

	c := command.Ctx.CLIContext

	opts := &osVolumes.UpdateOpts{
		Name:        c.String("rename"),
		Description: c.String("description"),
	}

	resource.Params = &paramsUpdate{
		volumeID: volumeID,
		opts:     opts,
	}

	return nil
}

func (command *commandUpdate) Execute(resource *handler.Resource) {
	opts := resource.Params.(*paramsUpdate).opts
	volumeID := resource.Params.(*paramsUpdate).volumeID
	volume, err := osVolumes.Update(command.Ctx.ServiceClient, volumeID, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = volumeSingle(volume)
}

func (command *commandUpdate) PreCSV(resource *handler.Resource) error {
	resource.FlattenMap("Attachments")
	return nil
}

func (command *commandUpdate) PreTable(resource *handler.Resource) error {
	return command.PreCSV(resource)
}
