package alicloud

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudPvtzZonesDataSource_basic(t *testing.T) {
	var pvtzZone pvtz.DescribeZoneInfoResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAlicloudPvtzZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudPvtzZoneDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.basic", &pvtzZone),
					testAccCheckAlicloudDataSourceID("data.alicloud_pvtz_zones.keyword"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zones.keyword", "zones.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_pvtz_zones.keyword", "zones.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zones.keyword", "zones.0.name", "tf-testacc.test.com"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zones.keyword", "zones.0.remark", ""),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zones.keyword", "zones.0.record_count", "0"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zones.keyword", "zones.0.is_ptr", "false"),
					resource.TestCheckResourceAttrSet("data.alicloud_pvtz_zones.keyword", "zones.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_pvtz_zones.keyword", "zones.0.update_time"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zones.keyword", "zones.0.bind_vpcs.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudPvtzZonesDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAlicloudPvtzZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudPvtzZoneDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_pvtz_zones.keyword"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zones.keyword", "zones.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_pvtz_zones.keyword", "zones.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_pvtz_zones.keyword", "zones.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_pvtz_zones.keyword", "zones.0.remark"),
					resource.TestCheckNoResourceAttr("data.alicloud_pvtz_zones.keyword", "zones.0.record_count"),
					resource.TestCheckNoResourceAttr("data.alicloud_pvtz_zones.keyword", "zones.0.is_ptr"),
					resource.TestCheckNoResourceAttr("data.alicloud_pvtz_zones.keyword", "zones.0.creation_time"),
					resource.TestCheckNoResourceAttr("data.alicloud_pvtz_zones.keyword", "zones.0.update_time"),
					resource.TestCheckNoResourceAttr("data.alicloud_pvtz_zones.keyword", "zones.0.bind_vpcs"),
				),
			},
		},
	})
}

const testAccCheckAlicloudPvtzZoneDataSourceBasic = `
resource "alicloud_pvtz_zone" "basic" {
	name = "tf-testacc.test.com"
}
data "alicloud_pvtz_zones" "keyword" {
	keyword = "${alicloud_pvtz_zone.basic.name}"
}
`

const testAccCheckAlicloudPvtzZoneDataSourceEmpty = `
data "alicloud_pvtz_zones" "keyword" {
	keyword = "tf-testacc-fake-name"
}
`
