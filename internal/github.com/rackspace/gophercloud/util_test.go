package gophercloud

import (
	"testing"

	th "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/testhelper"
)

func TestWaitFor(t *testing.T) {
	err := WaitFor(5, func() (bool, error) {
		return true, nil
	})
	th.CheckNoErr(t, err)
}
