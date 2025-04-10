// Generated by PMS #582
package eip

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceGlobalEipQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGlobalEipQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"type": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the resource type.`,
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the resources list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the quota type.`,
						},
						"used": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the used num.`,
						},
						"quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the total quotas.`,
						},
						"min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the min quotas.`,
						},
					},
				},
			},
		},
	}
}

type GlobalEipQuotasDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newGlobalEipQuotasDSWrapper(d *schema.ResourceData, meta interface{}) *GlobalEipQuotasDSWrapper {
	return &GlobalEipQuotasDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceGlobalEipQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newGlobalEipQuotasDSWrapper(d, meta)
	lisGeiResQuoRst, err := wrapper.ListGeipResourceQuotas()
	if err != nil {
		return diag.FromErr(err)
	}

	err = wrapper.listGeipResourceQuotasToSchema(lisGeiResQuoRst)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	return nil
}

// @API EIP GET /v3/{domain_id}/geip/quotas
func (w *GlobalEipQuotasDSWrapper) ListGeipResourceQuotas() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "geip")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{domain_id}/geip/quotas"
	uri = strings.ReplaceAll(uri, "{domain_id}", w.Config.DomainID)
	params := map[string]any{
		"type": w.ListToArray("type"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		Request().
		Result()
}

func (w *GlobalEipQuotasDSWrapper) listGeipResourceQuotasToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("resources", schemas.SliceToList(body.Get("quotas.resources"),
			func(resources gjson.Result) any {
				return map[string]any{
					"type":  resources.Get("type").Value(),
					"used":  resources.Get("used").Value(),
					"quota": resources.Get("quota").Value(),
					"min":   resources.Get("min").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
