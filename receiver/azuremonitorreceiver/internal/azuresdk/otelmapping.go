package azuresdk

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
)

func init() {
	mustHandleAllPossibleMetricUnitValues()
}

// ensure the Azure MetricUnit -> OtelUnits mapping is complete
func mustHandleAllPossibleMetricUnitValues() {
	var missing []string
	for _, v := range azquery.PossibleMetricUnitValues() {
		if _, ok := azureMonitorUnits[v]; !ok {
			missing = append(missing, string(v))
		}
	}
	if len(missing) > 0 {
		panic("Missing MetricUnit -> OtelUnits mappings for: " + strings.Join(missing, ", "))
	}
}

// We may also need to provide return of []struct{name: string, value: *float64}
func ExtractAggregations(m *azquery.MetricValue) map[azquery.AggregationType]*float64 {
	// TODO: instead of requiring the SDK packages in the scraper code we could handle all metric
	// values here over here and in receiver code use strings or a type alias
	aggs := make(map[azquery.AggregationType]*float64)
	if m.Average != nil {
		aggs[azquery.AggregationTypeAverage] = m.Average
	}
	if m.Count != nil {
		aggs[azquery.AggregationTypeCount] = m.Count
	}
	if m.Maximum != nil {
		aggs[azquery.AggregationTypeMaximum] = m.Maximum
	}
	if m.Minimum != nil {
		aggs[azquery.AggregationTypeMinimum] = m.Minimum
	}
	if m.Total != nil {
		aggs[azquery.AggregationTypeTotal] = m.Total
	}
	return aggs
}

// see: https://github.com/Azure/azure-sdk-for-go/blob/498a2eff41ecc251950868d9810c61229f5a9050/sdk/monitor/azquery/constants.go#L121-L134
var azureMonitorUnits = map[azquery.MetricUnit]string{
	// mapping from AzureMonitor units to OpenTelemetry units, which are strings
	// in "The Unified Code for Units of Measure" (UCUM) format
	azquery.MetricUnitBitsPerSecond: "b/s",
	azquery.MetricUnitByteSeconds: "By.s",
	azquery.MetricUnitBytes: "By",
	azquery.MetricUnitBytesPerSecond: "By/s",
	azquery.MetricUnitCores: "{cores}",
	azquery.MetricUnitCount: "1",
	azquery.MetricUnitCountPerSecond: "1/s",
	azquery.MetricUnitMilliCores: "{mcores}",
	azquery.MetricUnitMilliSeconds: "ms",
	azquery.MetricUnitNanoCores: "{ncores}",
	azquery.MetricUnitPercent: "%",
	azquery.MetricUnitSeconds: "s",
	azquery.MetricUnitUnspecified: "",
}

// ToOTelUnits maps Azure Monitor units to OpenTelemetry units
func ToOTelUnits(a azquery.MetricUnit) string {
	// not checking ok because use of azquery.MetricUnit SDK defined type enum ensures we handle all possible values
	// since we init with mustHandleAllPossibleMetricUnitValues()
	v := azureMonitorUnits[a]
	return v
}
