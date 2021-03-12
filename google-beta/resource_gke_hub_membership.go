// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGKEHubMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceGKEHubMembershipCreate,
		Read:   resourceGKEHubMembershipRead,
		Update: resourceGKEHubMembershipUpdate,
		Delete: resourceGKEHubMembershipDelete,

		Importer: &schema.ResourceImporter{
			State: resourceGKEHubMembershipImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of this entity type to be displayed on the console.`,
			},
			"external_id": {
				Type:     schema.TypeString,
				Required: true,
				Description: `An externally-generated and managed ID for this Membership.
This ID may still be modified after creation but it is not recommended to do so.
The ID must match the regex: '[a-zA-Z0-9][a-zA-Z0-9_\-\.]*' If this Membership
represents a Kubernetes cluster, this value should be set to the UUID of the
kube-system namespace object.`,
			},
			"membership_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The client-provided identifier of the membership.`,
			},
			"endpoint": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `If this Membership is a Kubernetes API server hosted on GKE, this is a self link to its GCP resource.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gke_cluster": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: `If this Membership is a Kubernetes API server hosted on GKE, this is a self link to its GCP resource.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_link": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										Description: `Self-link of the GCP resource for the GKE cluster.
For example: '//container.googleapis.com/projects/my-project/zones/us-west1-a/clusters/my-cluster'.
It can be at the most 1000 characters in length.  If the cluster is provisioned with Terraform,
this is '"//container.googleapis.com/${google_container_cluster.my-cluster.id}"'.`,
									},
								},
							},
						},
					},
				},
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: `Labels to apply to this membership.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unique identifier of the membership.`,
			},
			"unique_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Output only. Google-generated UUID for this resource.`,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceGKEHubMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	descriptionProp, err := expandGKEHubMembershipDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	labelsProp, err := expandGKEHubMembershipLabels(d.Get("labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}
	externalIdProp, err := expandGKEHubMembershipExternalId(d.Get("external_id"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("external_id"); !isEmptyValue(reflect.ValueOf(externalIdProp)) && (ok || !reflect.DeepEqual(v, externalIdProp)) {
		obj["externalId"] = externalIdProp
	}
	endpointProp, err := expandGKEHubMembershipEndpoint(d.Get("endpoint"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("endpoint"); !isEmptyValue(reflect.ValueOf(endpointProp)) && (ok || !reflect.DeepEqual(v, endpointProp)) {
		obj["endpoint"] = endpointProp
	}

	url, err := replaceVars(d, config, "https://gkehub.googleapis.com/v1beta1/projects/{{project}}/locations/global/memberships?membershipId={{membership_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Membership: %#v", obj)
	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Membership: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Membership: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Use the resource in the operation response to populate
	// identity fields and d.Id() before read
	var opRes map[string]interface{}
	err = gKEHubOperationWaitTimeWithResponse(
		config, res, &opRes, project, "Creating Membership", userAgent,
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create Membership: %s", err)
	}

	if err := d.Set("name", flattenGKEHubMembershipName(opRes["name"], d, config)); err != nil {
		return err
	}

	// This may have caused the ID to update - update it if so.
	id, err = replaceVars(d, config, "{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Membership %q: %#v", d.Id(), res)

	return resourceGKEHubMembershipRead(d, meta)
}

func resourceGKEHubMembershipRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "https://gkehub.googleapis.com/v1beta1/{{name}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Membership: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequest(config, "GET", billingProject, url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("GKEHubMembership %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Membership: %s", err)
	}

	if err := d.Set("name", flattenGKEHubMembershipName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Membership: %s", err)
	}
	if err := d.Set("description", flattenGKEHubMembershipDescription(res["description"], d, config)); err != nil {
		return fmt.Errorf("Error reading Membership: %s", err)
	}
	if err := d.Set("labels", flattenGKEHubMembershipLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Membership: %s", err)
	}
	if err := d.Set("external_id", flattenGKEHubMembershipExternalId(res["externalId"], d, config)); err != nil {
		return fmt.Errorf("Error reading Membership: %s", err)
	}
	if err := d.Set("unique_id", flattenGKEHubMembershipUniqueId(res["uniqueId"], d, config)); err != nil {
		return fmt.Errorf("Error reading Membership: %s", err)
	}
	if err := d.Set("endpoint", flattenGKEHubMembershipEndpoint(res["endpoint"], d, config)); err != nil {
		return fmt.Errorf("Error reading Membership: %s", err)
	}

	return nil
}

func resourceGKEHubMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Membership: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	labelsProp, err := expandGKEHubMembershipLabels(d.Get("labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}
	externalIdProp, err := expandGKEHubMembershipExternalId(d.Get("external_id"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("external_id"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, externalIdProp)) {
		obj["externalId"] = externalIdProp
	}

	url, err := replaceVars(d, config, "https://gkehub.googleapis.com/v1beta1/projects/{{project}}/locations/global/memberships/{{membership_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Membership %q: %#v", d.Id(), obj)
	updateMask := []string{}

	if d.HasChange("labels") {
		updateMask = append(updateMask, "labels")
	}

	if d.HasChange("external_id") {
		updateMask = append(updateMask, "externalId")
	}
	// updateMask is a URL parameter but not present in the schema, so replaceVars
	// won't set it
	url, err = addQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "PATCH", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf("Error updating Membership %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating Membership %q: %#v", d.Id(), res)
	}

	err = gKEHubOperationWaitTime(
		config, res, project, "Updating Membership", userAgent,
		d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return err
	}

	return resourceGKEHubMembershipRead(d, meta)
}

func resourceGKEHubMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Membership: %s", err)
	}
	billingProject = project

	url, err := replaceVars(d, config, "https://gkehub.googleapis.com/v1beta1/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting Membership %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "DELETE", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "Membership")
	}

	err = gKEHubOperationWaitTime(
		config, res, project, "Deleting Membership", userAgent,
		d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting Membership %q: %#v", d.Id(), res)
	return nil
}

func resourceGKEHubMembershipImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenGKEHubMembershipName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenGKEHubMembershipDescription(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenGKEHubMembershipLabels(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenGKEHubMembershipExternalId(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenGKEHubMembershipUniqueId(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenGKEHubMembershipEndpoint(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["gke_cluster"] =
		flattenGKEHubMembershipEndpointGkeCluster(original["gkeCluster"], d, config)
	return []interface{}{transformed}
}
func flattenGKEHubMembershipEndpointGkeCluster(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["resource_link"] =
		flattenGKEHubMembershipEndpointGkeClusterResourceLink(original["resourceLink"], d, config)
	return []interface{}{transformed}
}
func flattenGKEHubMembershipEndpointGkeClusterResourceLink(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func expandGKEHubMembershipDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandGKEHubMembershipLabels(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandGKEHubMembershipExternalId(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandGKEHubMembershipEndpoint(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedGkeCluster, err := expandGKEHubMembershipEndpointGkeCluster(original["gke_cluster"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGkeCluster); val.IsValid() && !isEmptyValue(val) {
		transformed["gkeCluster"] = transformedGkeCluster
	}

	return transformed, nil
}

func expandGKEHubMembershipEndpointGkeCluster(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedResourceLink, err := expandGKEHubMembershipEndpointGkeClusterResourceLink(original["resource_link"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedResourceLink); val.IsValid() && !isEmptyValue(val) {
		transformed["resourceLink"] = transformedResourceLink
	}

	return transformed, nil
}

func expandGKEHubMembershipEndpointGkeClusterResourceLink(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}