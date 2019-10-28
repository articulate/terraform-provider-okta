package okta

import (
	"fmt"

	"github.com/articulate/terraform-provider-okta/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var (
	sourceSchema = &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
				"type": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
				"name": &schema.Schema{
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}

	mappingResource = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The mapping property key.",
			},
			"expression": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"push_status": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      dontPush,
				ValidateFunc: validation.StringInSlice([]string{push, dontPush}, false),
			},
		},
	}
)

const (
	push     = "PUSH"
	dontPush = "DONT_PUSH"
)

func resourceOktaProfileMapping() *schema.Resource {
	return &schema.Resource{
		Create: resourceProfileMappingCreate,
		Read:   resourceProfileMappingRead,
		Update: resourceProfileMappingUpdate,
		Delete: resourceProfileMappingDelete,
		Exists: resourceProfileMappingExists,

		Schema: map[string]*schema.Schema{
			"source_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The source id of the mapping to manage.",
			},
			"delete_when_absent": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When turned on this flag will trigger the provider to delete mapping properties that are not defined in config. By default, we do not delete missing properties.",
			},
			"source": sourceSchema,
			"target": sourceSchema,
			"mappings": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     mappingResource,
			},
		},
	}
}

func buildMappingProperties(rawMap interface{}) map[string]*sdk.MappingProperty {
	res := map[string]*sdk.MappingProperty{}

	if m, ok := rawMap.(map[string]interface{}); ok {
		k := m["id"].(string)

		res[k] = &sdk.MappingProperty{
			Expression: m["expression"].(string),
			PushStatus: m["push_status"].(string),
		}
	}

	return res
}

func buildMapping(d *schema.ResourceData) sdk.Mapping {
	return sdk.Mapping{
		ID:         d.Id(),
		Properties: buildMappingProperties(d.Get("mappings")),
	}
}

func getProfileMapping(d *schema.ResourceData, m interface{}) (*sdk.Mapping, error) {
	client := getSupplementFromMetadata(m)
	mapping, resp, err := client.GetProfileMapping(d.Id())

	if is404(resp.StatusCode) {
		return nil, nil
	}

	return mapping, err
}

func resourceProfileMappingCreate(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	mapping, _, err := client.GetProfileMappingBySourceId(d.Get("source_id").(string))

	if err != nil || mapping == nil {
		return fmt.Errorf("failed to retrieve source, which is required to track mappings in state, error: %v", err)
	}

	d.SetId(mapping.ID)
	newMapping := buildMapping(d)

	_, _, err = client.UpdateMapping(mapping.ID, newMapping, nil)

	if err != nil {
		return err
	}

	return resourceProfileMappingRead(d, m)
}

func resourceProfileMappingDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceProfileMappingExists(d *schema.ResourceData, m interface{}) (bool, error) {
	m, err := getProfileMapping(d, m)

	return err == nil && m != nil, err
}

func resourceProfileMappingRead(d *schema.ResourceData, m interface{}) error {
	mapping, err := getProfileMapping(d, m)

	if mapping == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	d.Set("source", flattenMappingSource(mapping.Source))
	d.Set("target", flattenMappingSource(mapping.Target))
	d.Set("mappings", flattenMappingProperties(mapping.Properties))

	return nil
}

func flattenMappingProperties(src map[string]*sdk.MappingProperty) *schema.Set {
	arr := []interface{}{}

	for k, v := range src {
		arr = append(arr, map[string]interface{}{
			"id":          k,
			"push_status": v.PushStatus,
			"expression":  v.Expression,
		})
	}

	return schema.NewSet(schema.HashResource(mappingResource), arr)
}

func flattenMappingSource(src *sdk.MappingSource) map[string]interface{} {
	return map[string]interface{}{
		"id":   src.ID,
		"name": src.Name,
		"type": src.Type,
	}
}

func resourceProfileMappingUpdate(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	newMapping := buildMapping(d)

	_, _, err := client.UpdateMapping(d.Id(), newMapping, nil)

	if err != nil {
		return err
	}

	return resourceProfileMappingRead(d, m)
}
