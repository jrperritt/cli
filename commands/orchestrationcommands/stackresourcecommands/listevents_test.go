package stackresourcecommands

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	osStackEvents "github.com/rackspace/rack/internal/github.com/gophercloud/gophercloud/openstack/orchestration/v1/stackevents"
	th "github.com/rackspace/rack/internal/github.com/gophercloud/gophercloud/testhelper"
	"github.com/rackspace/rack/internal/github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListEventsContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandListEvents{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestListEventsKeys(t *testing.T) {
	cmd := &commandListEvents{}
	expected := keysListEvents
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestListEventsServiceClientType(t *testing.T) {
	cmd := &commandListEvents{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}

func TestListEventsHandleFlags(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	flagset.String("stack-name", "", "")
	flagset.String("stack-id", "", "")
	flagset.String("name", "", "")
	flagset.Set("stack-name", "stack1")
	flagset.Set("stack-id", "id1")
	flagset.Set("name", "resource1")
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandListEvents{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}

	expected := &handler.Resource{
		Params: &paramsListEvents{
			stackName:    "stack1",
			stackID:      "id1",
			resourceName: "resource1",
		},
	}
	actual := &handler.Resource{
		Params: &paramsListEvents{},
	}
	err := cmd.HandleFlags(actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expected.Params.(*paramsListEvents).stackName, actual.Params.(*paramsListEvents).stackName)
	th.AssertEquals(t, expected.Params.(*paramsListEvents).stackID, actual.Params.(*paramsListEvents).stackID)
	th.AssertEquals(t, expected.Params.(*paramsListEvents).resourceName, actual.Params.(*paramsListEvents).resourceName)
}

func TestListEventsExecute(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/stacks/stack1/id1/resources/resource1/events", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"events": [{"event_time": "2014-06-03T20:59:46", "resource_name":"resource1"}]}`)
	})
	cmd := &commandListEvents{
		Ctx: &handler.Context{
			ServiceClient: client.ServiceClient(),
		},
	}
	actual := &handler.Resource{
		Params: &paramsListEvents{
			stackName:    "stack1",
			stackID:      "id1",
			resourceName: "resource1",
			opts:         &osStackEvents.ListResourceEventsOpts{},
		},
	}
	cmd.Execute(actual)
	th.AssertNoErr(t, actual.Err)
}
