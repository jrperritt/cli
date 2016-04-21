package objectcommands

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gophercloud/cli/commandoptions"
	"github.com/gophercloud/cli/handler"
	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
	osObjects "github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/gophercloud/cli/util"
)

var upload = cli.Command{
	Name:        "upload",
	Usage:       util.Usage(commandPrefix, "upload", "--container <containerName> --name <objectName>"),
	Description: "Uploads an object",
	Action:      actionUpload,
	Flags:       commandoptions.CommandFlags(flagsUpload, keysUpload),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsUpload, keysUpload))
	},
}

func flagsUpload() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "container",
			Usage: "[required] The name of the container to upload the object into.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided with value of 'file'] The name the object should have in the Cloud Files container.",
		},
		cli.StringFlag{
			Name:  "content",
			Usage: "[optional; required if `file` or `stdin` isn't provided] The string contents to upload.",
		},
		cli.StringFlag{
			Name:  "file",
			Usage: "[optional; required if `content` or `stdin` isn't provided] The file name containing the contents to upload.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `file` or `content` isn't provided] The field being piped to STDIN, if any. Valid values are: file, content.",
		},
		cli.StringFlag{
			Name:  "content-type",
			Usage: "[optional] The Content-Type header.",
		},
		cli.IntFlag{
			Name:  "content-length",
			Usage: "[optional] The Content-Length header.",
		},
		cli.IntFlag{
			Name:  "delete-after",
			Usage: "[optional] Delete the object after this many seconds.",
		},
		cli.IntFlag{
			Name:  "delete-at",
			Usage: "[optional] Delete the object at this (Unix Epoch) timestamp.",
		},
		cli.StringFlag{
			Name:  "metadata",
			Usage: "[optional] A comma-separated string of key=value pairs.",
		},
	}
}

var keysUpload = []string{}

type paramsUpload struct {
	container string
	object    string
	stream    io.ReadSeeker
	opts      osObjects.CreateOpts
}

type commandUpload handler.Command

func actionUpload(c *cli.Context) {
	command := &commandUpload{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandUpload) Context() *handler.Context {
	return command.Ctx
}

func (command *commandUpload) Keys() []string {
	return keysUpload
}

func (command *commandUpload) ServiceClientType() string {
	return serviceClientType
}

func (command *commandUpload) HandleFlags(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"container"})
	if err != nil {
		return err
	}

	c := command.Ctx.CLIContext
	containerName := c.String("container")

	if err := CheckContainerExists(command.Ctx.ServiceClient, containerName); err != nil {
		return err
	}

	opts := osObjects.CreateOpts{
		ContentLength: int64(c.Int("content-length")),
		ContentType:   c.String("content-type"),
		DeleteAt:      c.Int("delete-at"),
		DeleteAfter:   c.Int("delete-after"),
	}

	if c.IsSet("metadata") {
		metadata, err := command.Ctx.CheckKVFlag("metadata")
		if err != nil {
			return err
		}
		opts.Metadata = metadata
	}

	resource.Params = &paramsUpload{
		container: containerName,
		opts:      opts,
	}

	return nil
}

func (command *commandUpload) HandlePipe(resource *handler.Resource, item string) error {
	readSeeker, err := os.Open(item)
	if err != nil {
		return err
	}
	resource.Params.(*paramsUpload).object = readSeeker.Name()
	resource.Params.(*paramsUpload).stream = readSeeker
	return nil
}

func (command *commandUpload) HandleSingle(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsUpload).object = command.Ctx.CLIContext.String("name")

	if command.Ctx.CLIContext.IsSet("file") {
		readSeeker, err := os.Open(command.Ctx.CLIContext.String("file"))
		if err != nil {
			return err
		}
		resource.Params.(*paramsUpload).stream = readSeeker
	} else if command.Ctx.CLIContext.IsSet("content") {
		readSeeker := strings.NewReader(command.Ctx.CLIContext.String("content"))
		resource.Params.(*paramsUpload).stream = readSeeker
	} else {
		return fmt.Errorf("One of `--file` and `--content` must be provided if not piping to STDIN.")
	}
	return nil
}

func (command *commandUpload) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsUpload)

	defer func() {
		if closeable, ok := params.stream.(io.ReadCloser); ok {
			closeable.Close()
		}
	}()

	containerName := params.container
	objectName := params.object
	stream := params.stream
	opts := params.opts

	rawResponse := objects.Create(command.Ctx.ServiceClient, containerName, objectName, stream, opts)
	if rawResponse.Err != nil {
		resource.Err = rawResponse.Err
		return
	}
	resource.Result = fmt.Sprintf("Successfully uploaded object [%s] to container [%s]\n", objectName, containerName)
}

func (command *commandUpload) StdinField() string {
	return "file"
}

func (command *commandUpload) StreamField() string {
	return "content"
}

func (command *commandUpload) HandleStreamPipe(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}
	resource.Params.(*paramsUpload).object = command.Ctx.CLIContext.String("name")
	resource.Params.(*paramsUpload).stream = os.Stdin
	return nil
}
