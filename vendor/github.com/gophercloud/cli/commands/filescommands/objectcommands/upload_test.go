package objectcommands

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/gophercloud/cli/handler"
	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
	osObjects "github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	th "github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/testhelper/client"
	"github.com/gophercloud/cli/output"
)

func newUpCmd(fs *flag.FlagSet) *commandUpload {
	return &commandUpload{Ctx: &handler.Context{
		CLIContext: cli.NewContext(cli.NewApp(), fs, nil),
	}}
}

func TestUploadContext(t *testing.T) {
	cmd := newUpCmd(flag.NewFlagSet("flags", 1))
	th.AssertDeepEquals(t, cmd.Ctx, cmd.Context())
}

func TestUploadKeys(t *testing.T) {
	cmd := &commandUpload{}
	th.AssertDeepEquals(t, keysUpload, cmd.Keys())
}

func TestUploadServiceClientType(t *testing.T) {
	cmd := &commandUpload{}
	th.AssertEquals(t, serviceClientType, cmd.ServiceClientType())
}

func TestUploadErrWhenCtnrMissing(t *testing.T) {
	fs := flag.NewFlagSet("flags", 1)

	err := newUpCmd(fs).HandleFlags(&handler.Resource{})

	expected := output.ErrMissingFlag{Msg: "--container is required."}
	th.AssertDeepEquals(t, expected, err)
}

func TestUploadHandlePipe(t *testing.T) {
	cmd := &commandUpload{}

	actual := &handler.Resource{
		Params: &paramsUpload{},
	}

	f, err := ioutil.TempFile("", "bar")
	th.AssertNoErr(t, err)
	defer os.Remove(f.Name())

	err = cmd.HandlePipe(actual, f.Name())
	th.AssertNoErr(t, err)
}

func TestUploadExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/foo/bar", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		w.Header().Add("Content-Type", "text/plain")
		hash := md5.New()
		io.WriteString(hash, "hodor")
		localChecksum := hash.Sum(nil)
		w.Header().Set("ETag", fmt.Sprintf("%x", localChecksum))
		w.WriteHeader(201)
		fmt.Fprintf(w, `hodor`)
	})

	fs := flag.NewFlagSet("flags", 1)
	cmd := newUpCmd(fs)
	cmd.Ctx.ServiceClient = client.ServiceClient()

	res := &handler.Resource{
		Params: &paramsUpload{
			container: "foo",
			object:    "bar",
			stream:    strings.NewReader("hodor"),
			opts:      osObjects.CreateOpts{},
		},
	}

	cmd.Execute(res)

	th.AssertNoErr(t, res.Err)
	th.AssertEquals(t, "Successfully uploaded object [bar] to container [foo]\n", res.Result)
}
