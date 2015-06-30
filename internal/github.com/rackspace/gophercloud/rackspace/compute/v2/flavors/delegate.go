package flavors

import (
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud"
	os "github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
	"github.com/jrperritt/rack/internal/github.com/rackspace/gophercloud/pagination"
)

// ListOpts helps control the results returned by the List() function. For example, a flavor with a
// minDisk field of 10 will not be returned if you specify MinDisk set to 20.
type ListOpts struct {

	// MinDisk and MinRAM, if provided, elide flavors that do not meet your criteria.
	MinDisk int `q:"minDisk"`
	MinRAM  int `q:"minRam"`

	// Marker specifies the ID of the last flavor in the previous page.
	Marker string `q:"marker"`

	// Limit instructs List to refrain from sending excessively large lists of flavors.
	Limit int `q:"limit"`
}

// ToFlavorListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToFlavorListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// ListDetail enumerates the server images available to your account.
func ListDetail(client *gophercloud.ServiceClient, opts os.ListOptsBuilder) pagination.Pager {
	return os.ListDetail(client, opts)
}

// Get returns details about a single flavor, identity by ID.
func Get(client *gophercloud.ServiceClient, id string) os.GetResult {
	return os.Get(client, id)
}

// ExtractFlavors interprets a page of List results as Flavors.
func ExtractFlavors(page pagination.Page) ([]os.Flavor, error) {
	return os.ExtractFlavors(page)
}
