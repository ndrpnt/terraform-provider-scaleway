package scaleway

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	account "github.com/scaleway/scaleway-sdk-go/api/account/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func resourceScalewayAccountSSKKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalewayAccountSSHKeyCreate,
		ReadContext:   resourceScalewayAccountSSHKeyRead,
		UpdateContext: resourceScalewayAccountSSHKeyUpdate,
		DeleteContext: resourceScalewayAccountSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(defaultAccountSSHKeyTimeout),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the SSH key",
			},
			"public_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The public SSH key",
				// We don't consider trailing \n as diff
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return strings.Trim(oldValue, "\n") == strings.Trim(newValue, "\n")
				},
			},
			"organization_id": organizationIDSchema(),
			"project_id":      projectIDSchema(),
		},
	}
}

func resourceScalewayAccountSSHKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	accountAPI := accountAPI(meta)

	res, err := accountAPI.CreateSSHKey(&account.CreateSSHKeyRequest{
		Name:      d.Get("name").(string),
		PublicKey: strings.Trim(d.Get("public_key").(string), "\n"),
		ProjectID: expandStringPtr(d.Get("project_id")),
	}, scw.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(res.ID)

	return resourceScalewayAccountSSHKeyRead(ctx, d, meta)
}

func resourceScalewayAccountSSHKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	accountAPI := accountAPI(meta)

	res, err := accountAPI.GetSSHKey(&account.GetSSHKeyRequest{
		SSHKeyID: d.Id(),
	}, scw.WithContext(ctx))
	if err != nil {
		if is404Error(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	_ = d.Set("name", res.Name)
	_ = d.Set("public_key", res.PublicKey)
	_ = d.Set("organization_id", res.OrganizationID)
	_ = d.Set("project_id", res.ProjectID)

	return nil
}

func resourceScalewayAccountSSHKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	accountAPI := accountAPI(meta)

	if d.HasChange("name") {
		_, err := accountAPI.UpdateSSHKey(&account.UpdateSSHKeyRequest{
			SSHKeyID: d.Id(),
			Name:     expandUpdatedStringPtr(d.Get("name")),
		}, scw.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceScalewayAccountSSHKeyRead(ctx, d, meta)
}

func resourceScalewayAccountSSHKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	accountAPI := accountAPI(meta)

	err := accountAPI.DeleteSSHKey(&account.DeleteSSHKeyRequest{
		SSHKeyID: d.Id(),
	}, scw.WithContext(ctx))
	if err != nil && !is404Error(err) {
		return diag.FromErr(err)
	}

	return nil
}
