package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceUsableSubnetworks_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUsableSubnetworksConfig(),
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

		rs_attr_keys := []string{
			"name",
			"region",
			"project",
			"self_link",
			"ip_cidr_range",
			"network",
			"secondary_ip_range",
		}

		var ds_attr_key string
		for _, rs_attr_key := range rs_attr_keys {
			ds_attr_key = fmt.Sprintf("subnetworks.0.%s", rs_attr_key)

			if ds_attr[ds_attr_key] != rs_attr[rs_attr_key] {
				return fmt.Errorf(
					"%s is %s; want %s",
					ds_attr_key,
					ds_attr[ds_attr_key],
					rs_attr[rs_attr_key],
				)
			}
		}

		if v1RsNetwork := ConvertSelfLinkToV1(rs_attr["network"]); ds_attr["subnetworks.0.network"] != v1RsNetwork {
			return fmt.Errorf("network is %s; want %s", ds_attr["subnetworks.0.network"], v1RsNetwork)
		}

		return nil
	}
}

func testAccCheckUsableSubnetworksConfig() string {
	return fmt.Sprintf(`
resource "google_compute_network" "network" {
	name = "%s"
	description = "my-description"
}

resource "google_compute_subnetwork" "subnetwork" {
	name = "subnetwork-test"
	ip_cidr_range = "10.0.0.0/24"
	network  = "${google_compute_network.network.self_link}"
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
