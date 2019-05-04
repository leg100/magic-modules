package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
)

func dataSourceGoogleComputeUsableSubnetworks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceComputeUsableSubnetworksRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnetworks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"self_link": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_cidr_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secondary_ip_range": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"range_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_cidr_range": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceComputeUsableSubnetworksRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	subnetworks := make([]map[string]interface{}, 0)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	d.SetId(project)

	filter := d.Get("filter").(string)

	token := ""
	for paginate := true; paginate; {
		call := config.clientCompute.Subnetworks.ListUsable(project)
		call = call.PageToken(token)
		call = call.Filter(filter)

		res, err := call.Do()
		if err != nil {
			return fmt.Errorf("Error retrieving usable subnetworks: %s", err)
		}

		pageSubnetworks := flattenUsableSubnetworks(res.Items)
		subnetworks = append(subnetworks, pageSubnetworks...)

		token := res.NextPageToken
		paginate = token != ""
	}

	if err := d.Set("subnetworks", subnetworks); err != nil {
		return fmt.Errorf("Error retrieving subnetworks: %s", err)
	}

	return nil
}

func flattenUsableSubnetworks(usableSubnetworks []*compute.UsableSubnetwork) []map[string]interface{} {
	usableSubnetworksSchema := make([]map[string]interface{}, 0, len(usableSubnetworks))
	for _, usableSubnetwork := range usableSubnetworks {
		self_link_segments := strings.Split(usableSubnetwork.Subnetwork, "/")

		data := map[string]interface{}{
			"name":               self_link_segments[10],
			"region":             self_link_segments[8],
			"project":            self_link_segments[6],
			"self_link":          usableSubnetwork.Subnetwork,
			"ip_cidr_range":      usableSubnetwork.IpCidrRange,
			"network":            usableSubnetwork.Network,
			"secondary_ip_range": flattenUsableSecondaryRanges(usableSubnetwork.SecondaryIpRanges),
		}

		usableSubnetworksSchema = append(usableSubnetworksSchema, data)
	}

	return usableSubnetworksSchema
}

func flattenUsableSecondaryRanges(secondaryRanges []*compute.UsableSubnetworkSecondaryRange) []map[string]interface{} {
	secondaryRangesSchema := make([]map[string]interface{}, 0, len(secondaryRanges))
	for _, secondaryRange := range secondaryRanges {
		data := map[string]interface{}{
			"range_name":    secondaryRange.RangeName,
			"ip_cidr_range": secondaryRange.IpCidrRange,
		}

		secondaryRangesSchema = append(secondaryRangesSchema, data)
	}

	return secondaryRangesSchema
}
