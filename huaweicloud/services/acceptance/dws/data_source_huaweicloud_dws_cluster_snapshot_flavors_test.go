package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccClusterSnapshotFlavors_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_dws_cluster_snapshot_flavors.test"
		dc     = acceptance.InitDataSourceCheck(dcName)

		rName = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterSnapshotFlavors_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "flavors.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.id"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.classify"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.version"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.default_node"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.max_node"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.code"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.status"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.min_node"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.flavor_code"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.volume_num"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.default_capacity"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.scenario"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.duplicate"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.attributes.#"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.product_versions.#"),
					resource.TestCheckResourceAttrSet(dcName, "flavors.0.volume_used.#"),
					resource.TestCheckOutput("is_existed_snapshot", "true"),
				),
			},
		},
	})
}

func testAccClusterSnapshotFlavors_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_snapshot" "test" {
  cluster_id  = "%[1]s"
  name        = "%[2]s"
  description = "Created by terraform script"
}

data "huaweicloud_dws_cluster_snapshot_flavors" "test" {
  depends_on = [
    huaweicloud_dws_snapshot.test
  ]

  snapshot_id = huaweicloud_dws_snapshot.test.id
}

locals {
  snapshot_ids = data.huaweicloud_dws_cluster_snapshot_flavors.test.flavors[*].id
}

output "is_existed_snapshot" {
  value = length(local.snapshot_ids) > 0
}
`, acceptance.HW_DWS_CLUSTER_ID, name)
}
