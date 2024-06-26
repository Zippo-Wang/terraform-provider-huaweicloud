// Generated by PMS #229
package lts

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
)

func DataSourceLtsSearchCriteria() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLtsSearchCriteriaRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"log_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the log group.`,
			},
			"search_criteria": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All search criteria that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_stream_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the log stream.`,
						},
						"log_stream_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the log stream.`,
						},
						"criteria": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of the search criteria under specified log stream.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the search criterion.`,
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the search criterion.`,
									},
									"criteria": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The content of search criterion.`,
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the search criterion.`,
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

type SearchCriteriaDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newSearchCriteriaDSWrapper(d *schema.ResourceData, meta interface{}) *SearchCriteriaDSWrapper {
	return &SearchCriteriaDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceLtsSearchCriteriaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newSearchCriteriaDSWrapper(d, meta)
	lisQueAllSeaCriRst, err := wrapper.ListQueryAllSearchCriterias()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listQueryAllSearchCriteriasToSchema(lisQueAllSeaCriRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API LTS GET /v1.0/{project_id}/lts/groups/{group_id}/search-criterias
func (w *SearchCriteriaDSWrapper) ListQueryAllSearchCriterias() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "lts")
	if err != nil {
		return nil, err
	}

	uri := "/v1.0/{project_id}/lts/groups/{group_id}/search-criterias"
	uri = strings.ReplaceAll(uri, "{group_id}", w.Get("log_group_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *SearchCriteriaDSWrapper) listQueryAllSearchCriteriasToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("search_criteria", schemas.SliceToList(body.Get("search_criterias"),
			func(searchCriteria gjson.Result) any {
				return map[string]any{
					"log_stream_id":   searchCriteria.Get("log_stream_id").Value(),
					"log_stream_name": searchCriteria.Get("log_stream_name").Value(),
					"criteria": schemas.SliceToList(searchCriteria.Get("criterias"),
						func(criteria gjson.Result) any {
							return map[string]any{
								"id":       criteria.Get("id").Value(),
								"name":     criteria.Get("name").Value(),
								"criteria": criteria.Get("criteria").Value(),
								"type":     criteria.Get("search_type").Value(),
							}
						},
					),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
