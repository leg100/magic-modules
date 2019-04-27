package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceUsableSubnets_basic(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUsableSubnetsConfig(project),
				Check: resource.ComposeTestCheckFunc(
					// TODO: check whether we can improve upon this test (do we know what the subnets are called in advance?)
					// We can't guarantee no project won't have our project ID as a prefix, so we'll check set-ness rather than correctness
					resource.TestCheckResourceAttrSet("data.google_compute_usable_subnetworks.subnets", "subnets.0.self_link"),
				),
			},
		},
	})
}

func testAccCheckUsableSubnetsConfig(project string) string {
	return fmt.Sprintf(`
data "google_compute_usable_subnetworks" "subnets" {
  project = "%s"
}
`, project)
}
