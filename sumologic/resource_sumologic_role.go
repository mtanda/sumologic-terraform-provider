package sumologic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSumologicRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicRoleCreate,
		Read:   resourceSumologicRoleRead,
		Delete: resourceSumologicRoleDelete,
		Update: resourceSumologicRoleUpdate,
		Exists: resourceSumologicRoleExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "",
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "",
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "Etc/UTC",
			},
			"lookup_by_name": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  false,
			},
			"destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},
		},
	}
}

func resourceSumologicRoleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return err
	}

	collector, err := c.GetRole(id)

	if err != nil {
		return err
	}

	if collector == nil {
		log.Printf("[WARN] Role not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	d.Set("name", collector.Name)
	d.Set("description", collector.Description)
	d.Set("category", collector.Category)
	d.Set("timezone", collector.TimeZone)

	return nil
}

func resourceSumologicRoleDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Get("destroy").(bool) {
		id, _ := strconv.Atoi(d.Id())
		return c.DeleteRole(id)
	}

	return nil
}

func resourceSumologicRoleCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Get("lookup_by_name").(bool) {
		collector, err := c.GetRoleName(d.Get("name").(string))

		if err != nil {
			return err
		}

		if collector != nil {
			d.SetId(strconv.Itoa(collector.ID))
		}
	}

	if d.Id() == "" {
		id, err := c.CreateRole(Role{
			RoleType: "Hosted",
			Name:     d.Get("name").(string),
		})

		if err != nil {
			return err
		}

		d.SetId(strconv.Itoa(id))
	}

	return resourceSumologicRoleUpdate(d, meta)
}

func resourceSumologicRoleUpdate(d *schema.ResourceData, meta interface{}) error {

	collector := resourceToRole(d)

	c := meta.(*Client)
	err := c.UpdateRole(collector)

	if err != nil {
		return err
	}

	return resourceSumologicRoleRead(d, meta)
}

func resourceSumologicRoleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	c := meta.(*Client)

	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return false, fmt.Errorf("collector id should be an integer; got %s (%s)", d.Id(), err)
	}

	_, err = c.GetRole(id)

	return err == nil, nil
}

func resourceToRole(d *schema.ResourceData) Role {
	id, _ := strconv.Atoi(d.Id())

	return Role{
		ID:          id,
		RoleType:    "Hosted",
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Category:    d.Get("category").(string),
		TimeZone:    d.Get("timezone").(string),
	}
}
