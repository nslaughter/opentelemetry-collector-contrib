// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceBuilder(t *testing.T) {
	for _, test := range []string{"default", "all_set", "none_set"} {
		t.Run(test, func(t *testing.T) {
			cfg := loadResourceAttributesConfig(t, test)
			rb := NewResourceBuilder(cfg)
			rb.SetAzuremonitorSubscriptionID("azuremonitor.subscription_id-val")
			rb.SetAzuremonitorTenantID("azuremonitor.tenant_id-val")
			rb.SetCloudAccountID("cloud.account.id-val")
			rb.SetCloudAvailabilityZone("cloud.availability_zone-val")
			rb.SetCloudPlatform("cloud.platform-val")
			rb.SetCloudProvider("cloud.provider-val")
			rb.SetCloudRegion("cloud.region-val")
			rb.SetCloudResourceID("cloud.resource_id-val")

			res := rb.Emit()
			assert.Equal(t, 0, rb.Emit().Attributes().Len()) // Second call should return empty Resource

			switch test {
			case "default":
				assert.Equal(t, 7, res.Attributes().Len())
			case "all_set":
				assert.Equal(t, 8, res.Attributes().Len())
			case "none_set":
				assert.Equal(t, 0, res.Attributes().Len())
				return
			default:
				assert.Failf(t, "unexpected test case: %s", test)
			}

			val, ok := res.Attributes().Get("azuremonitor.subscription_id")
			assert.True(t, ok)
			if ok {
				assert.EqualValues(t, "azuremonitor.subscription_id-val", val.Str())
			}
			val, ok = res.Attributes().Get("azuremonitor.tenant_id")
			assert.True(t, ok)
			if ok {
				assert.EqualValues(t, "azuremonitor.tenant_id-val", val.Str())
			}
			val, ok = res.Attributes().Get("cloud.account.id")
			assert.True(t, ok)
			if ok {
				assert.EqualValues(t, "cloud.account.id-val", val.Str())
			}
			val, ok = res.Attributes().Get("cloud.availability_zone")
			assert.Equal(t, test == "all_set", ok)
			if ok {
				assert.EqualValues(t, "cloud.availability_zone-val", val.Str())
			}
			val, ok = res.Attributes().Get("cloud.platform")
			assert.True(t, ok)
			if ok {
				assert.EqualValues(t, "cloud.platform-val", val.Str())
			}
			val, ok = res.Attributes().Get("cloud.provider")
			assert.True(t, ok)
			if ok {
				assert.EqualValues(t, "cloud.provider-val", val.Str())
			}
			val, ok = res.Attributes().Get("cloud.region")
			assert.True(t, ok)
			if ok {
				assert.EqualValues(t, "cloud.region-val", val.Str())
			}
			val, ok = res.Attributes().Get("cloud.resource_id")
			assert.True(t, ok)
			if ok {
				assert.EqualValues(t, "cloud.resource_id-val", val.Str())
			}
		})
	}
}
