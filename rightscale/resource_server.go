package rightscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rightscale/terraform-provider-rightscale/rightscale/rsc"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Read:   resourceRead,
		Exists: resourceExists,
		Delete: resourceDelete,
		Create: resourceCreateServer(serverWriteFields),
		Update: resourceUpdateFunc(serverWriteFields),

		Schema: map[string]*schema.Schema{
			"deployment_href": &schema.Schema{
				Description: "ID of deployment in which to create server",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"description": &schema.Schema{
				Description: "description of server",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"instance": &schema.Schema{
				Description: "server instance details",
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Elem:        resourceInstance(),
			},
			"name": &schema.Schema{
				Description: "name of server",
				Type:        schema.TypeString,
				Required:    true,
			},
			"optimized": &schema.Schema{
				Description: "A flag indicating whether Instances of this Server should be optimized for high-performance volumes (e.g. Volumes supporting a specified number of IOPS). Not supported in all Clouds.",
				Type:        schema.TypeBool,
				Optional:    true,
			},

			// Read-only fields
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_instance": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeMap},
				Computed: true,
			},
			"next_instance": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func serverWriteFields(d *schema.ResourceData) rsc.Fields {
	fields := rsc.Fields{}
	// construct 'instance' hash so we end up with a server WITH a running instance
	if i, ok := d.GetOk("instance"); ok {
		fields["instance"] = instanceWriteFieldsFromMap(i.([]interface{})[0].(map[string]interface{}))
	}
	// TODO - Is this needed still?
	//fields["instance"] = i.([]interface{})[0].(map[string]interface{})
	//if fields["instance"].(map[string]interface{})["associate_public_ip_address"].(bool) == false {
	//  delete(fields["instance"].(map[string]interface{}), "associate_public_ip_address")
	//}

	if o, ok := d.GetOk("optimized"); ok {
		if o.(bool) {
			fields["optimized"] = "true"
		} else {
			fields["optimized"] = "false"
		}
	}
	for _, f := range []string{
		"deployment_href", "description", "name", "resource_group_href",
		"server_tag_scope",
	} {
		if v, ok := d.GetOk(f); ok {
			fields[f] = v
		}
	}
	return rsc.Fields{"server": fields}
}

func resourceCreateServer(fieldsFunc func(*schema.ResourceData) rsc.Fields) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		client := m.(rsc.Client)
		res, err := client.Create("rs_cm", "servers", fieldsFunc(d))
		if err != nil {
			return err
		}
		for k, v := range res.Fields {
			d.Set(k, v)
		}
		d.SetId(res.Locator.Namespace + ":" + res.Locator.Href)
		return nil
	}
}
