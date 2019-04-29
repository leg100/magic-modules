package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceUsableSubnetworks_basic(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUsableSubnetworksConfig(project),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceUsableSubnetworksCheck("data.google_compute_usable_subnetworks.subnetworks", "google_compute_subnetwork.subnetwork"),
				),
			},
		},
	})
}

func testAccDataSourceUsableSubnetworksCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[data_source_name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", data_source_name)
		}

		rs, ok := s.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("can't find %s in state", resource_name)
		}

		ds_attr := ds.Primary.Attributes
		rs_attr := rs.Primary.Attributes

		ds_attr_subnetworks, ok := ds_attr["subnetworks"]
		if !ok {
			return fmt.Errorf("expected attribute `subnetworks` not set on resource %s", data_source_name)
		}

		if num_of_subnetworks := len(ds_attr_subnetworks); num_of_subnetworks != 1 {
			return fmt.Errorf("expected attribute `subnetworks` of resource %s to be a list of one element, but contains %d elements", data_source_name, num_of_subnetworks)
		}

		subnetwork_attrs_to_test := []string{
			"name",
			"region",
			"project",
			"self_link",
			"ip_cidr_range",
			"network",
			"secondary_ip_range",
		}

		for _, attr_to_check := range subnetwork_attrs_to_test {
			if ds_attr_subnetworks[0][attr_to_check] != rs_attr[attr_to_check] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr_to_check,
					ds_attr[attr_to_check],
					rs_attr[attr_to_check],
				)
			}
		}

		if v1RsNetwork := ConvertSelfLinkToV1(rs_attr["network"]); ds_attr["network"] != v1RsNetwork {
			return fmt.Errorf("network is %s; want %s", ds_attr["network"], v1RsNetwork)
		}

		return nil
	}
}

func testAccCheckUsableSubnetworksConfig(project string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "network" {
	name = "%s"
	description = "my-description"
}

resource "google_compute_subnetwork" "subnetwork" {
	name = "subnetwork-test"
	ip_cidr_range = "10.0.0.0/24"
	network  = "${google_compute_network.foobar.self_link}"
	secondary_ip_range {
		range_name = "tf-test-secondary-range"
		ip_cidr_range = "192.168.1.0/24"
	}
}

data "google_compute_usable_subnetworks" "subnetworks" {
	filter = "name:subnetwork-test network:${google_compute_network.network.self_link}"
}
`, acctest.RandomWithPrefix("network-test"))
}
