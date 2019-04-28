package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
)

func dataSourceGoogleComputeUsableSubnetworks() *schema.Resource {
	return &schema.Resource{
		Read: datasourceGoogleComputeUsableSubnetsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"self_link": {
							Type:		schema.TypeString,
							Computed:	true,
						},
					},
				},
			},
		},
	}
}

func datasourceGoogleComputeUsableSubnetsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	subnets := make([]interface{}, 0)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	d.SetId(project)

	// TODO: filter

	token := ""
	for paginate := true; paginate; {
		res, err := config.clientCompute.Subnetworks.ListUsable(project).PageToken(token).Do()
		if err != nil {
			return fmt.Errorf("Error retrieving usable subnets: %s", err)
		}

		pageSubnets := flattenUsableSubnets(res.Items)
		subnets = append(subnets, pageSubnets...)

		token := res.NextPageToken
		paginate = token != ""
	}

	if err := d.Set("subnets", subnets); err != nil {
		return fmt.Errorf("Error retrieving subnets: %s", err)
	}

	return nil
}

func flattenUsableSubnets(usableSubnets []*compute.UsableSubnetwork) interface{} {
	attrs := make([]interface{}, 0, len(usableSubnets))
	for _, usableSubnet := range usableSubnets {
		attrs = append(attrs, map[string]interface{}{
			"self_link": usableSubnet.Subnetwork,
		})
	}

	return attrs
}
