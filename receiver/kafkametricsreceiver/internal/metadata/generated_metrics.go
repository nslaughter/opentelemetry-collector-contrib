// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
)

type metricKafkaBrokers struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills kafka.brokers metric with initial data.
func (m *metricKafkaBrokers) init() {
	m.data.SetName("kafka.brokers")
	m.data.SetDescription("Number of brokers in the cluster.")
	m.data.SetUnit("{brokers}")
	m.data.SetEmptySum()
	m.data.Sum().SetIsMonotonic(false)
	m.data.Sum().SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
}

func (m *metricKafkaBrokers) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntValue(val)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricKafkaBrokers) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricKafkaBrokers) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricKafkaBrokers(cfg MetricConfig) metricKafkaBrokers {
	m := metricKafkaBrokers{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricKafkaConsumerGroupLag struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills kafka.consumer_group.lag metric with initial data.
func (m *metricKafkaConsumerGroupLag) init() {
	m.data.SetName("kafka.consumer_group.lag")
	m.data.SetDescription("Current approximate lag of consumer group at partition of topic")
	m.data.SetUnit("1")
	m.data.SetEmptyGauge()
	m.data.Gauge().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricKafkaConsumerGroupLag) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, groupAttributeValue string, topicAttributeValue string, partitionAttributeValue int64) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Gauge().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntValue(val)
	dp.Attributes().PutStr("group", groupAttributeValue)
	dp.Attributes().PutStr("topic", topicAttributeValue)
	dp.Attributes().PutInt("partition", partitionAttributeValue)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricKafkaConsumerGroupLag) updateCapacity() {
	if m.data.Gauge().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Gauge().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricKafkaConsumerGroupLag) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Gauge().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricKafkaConsumerGroupLag(cfg MetricConfig) metricKafkaConsumerGroupLag {
	m := metricKafkaConsumerGroupLag{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricKafkaConsumerGroupLagSum struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills kafka.consumer_group.lag_sum metric with initial data.
func (m *metricKafkaConsumerGroupLagSum) init() {
	m.data.SetName("kafka.consumer_group.lag_sum")
	m.data.SetDescription("Current approximate sum of consumer group lag across all partitions of topic")
	m.data.SetUnit("1")
	m.data.SetEmptyGauge()
	m.data.Gauge().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricKafkaConsumerGroupLagSum) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, groupAttributeValue string, topicAttributeValue string) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Gauge().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntValue(val)
	dp.Attributes().PutStr("group", groupAttributeValue)
	dp.Attributes().PutStr("topic", topicAttributeValue)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricKafkaConsumerGroupLagSum) updateCapacity() {
	if m.data.Gauge().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Gauge().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricKafkaConsumerGroupLagSum) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Gauge().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricKafkaConsumerGroupLagSum(cfg MetricConfig) metricKafkaConsumerGroupLagSum {
	m := metricKafkaConsumerGroupLagSum{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricKafkaConsumerGroupMembers struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills kafka.consumer_group.members metric with initial data.
func (m *metricKafkaConsumerGroupMembers) init() {
	m.data.SetName("kafka.consumer_group.members")
	m.data.SetDescription("Count of members in the consumer group")
	m.data.SetUnit("{members}")
	m.data.SetEmptySum()
	m.data.Sum().SetIsMonotonic(false)
	m.data.Sum().SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricKafkaConsumerGroupMembers) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, groupAttributeValue string) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntValue(val)
	dp.Attributes().PutStr("group", groupAttributeValue)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricKafkaConsumerGroupMembers) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricKafkaConsumerGroupMembers) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricKafkaConsumerGroupMembers(cfg MetricConfig) metricKafkaConsumerGroupMembers {
	m := metricKafkaConsumerGroupMembers{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricKafkaConsumerGroupOffset struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills kafka.consumer_group.offset metric with initial data.
func (m *metricKafkaConsumerGroupOffset) init() {
	m.data.SetName("kafka.consumer_group.offset")
	m.data.SetDescription("Current offset of the consumer group at partition of topic")
	m.data.SetUnit("1")
	m.data.SetEmptyGauge()
	m.data.Gauge().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricKafkaConsumerGroupOffset) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, groupAttributeValue string, topicAttributeValue string, partitionAttributeValue int64) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Gauge().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntValue(val)
	dp.Attributes().PutStr("group", groupAttributeValue)
	dp.Attributes().PutStr("topic", topicAttributeValue)
	dp.Attributes().PutInt("partition", partitionAttributeValue)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricKafkaConsumerGroupOffset) updateCapacity() {
	if m.data.Gauge().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Gauge().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricKafkaConsumerGroupOffset) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Gauge().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricKafkaConsumerGroupOffset(cfg MetricConfig) metricKafkaConsumerGroupOffset {
	m := metricKafkaConsumerGroupOffset{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricKafkaConsumerGroupOffsetSum struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills kafka.consumer_group.offset_sum metric with initial data.
func (m *metricKafkaConsumerGroupOffsetSum) init() {
	m.data.SetName("kafka.consumer_group.offset_sum")
	m.data.SetDescription("Sum of consumer group offset across partitions of topic")
	m.data.SetUnit("1")
	m.data.SetEmptyGauge()
	m.data.Gauge().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricKafkaConsumerGroupOffsetSum) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, groupAttributeValue string, topicAttributeValue string) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Gauge().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntValue(val)
	dp.Attributes().PutStr("group", groupAttributeValue)
	dp.Attributes().PutStr("topic", topicAttributeValue)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricKafkaConsumerGroupOffsetSum) updateCapacity() {
	if m.data.Gauge().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Gauge().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricKafkaConsumerGroupOffsetSum) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Gauge().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricKafkaConsumerGroupOffsetSum(cfg MetricConfig) metricKafkaConsumerGroupOffsetSum {
	m := metricKafkaConsumerGroupOffsetSum{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricKafkaPartitionCurrentOffset struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills kafka.partition.current_offset metric with initial data.
func (m *metricKafkaPartitionCurrentOffset) init() {
	m.data.SetName("kafka.partition.current_offset")
	m.data.SetDescription("Current offset of partition of topic.")
	m.data.SetUnit("1")
	m.data.SetEmptyGauge()
	m.data.Gauge().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricKafkaPartitionCurrentOffset) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, topicAttributeValue string, partitionAttributeValue int64) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Gauge().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntValue(val)
	dp.Attributes().PutStr("topic", topicAttributeValue)
	dp.Attributes().PutInt("partition", partitionAttributeValue)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricKafkaPartitionCurrentOffset) updateCapacity() {
	if m.data.Gauge().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Gauge().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricKafkaPartitionCurrentOffset) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Gauge().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricKafkaPartitionCurrentOffset(cfg MetricConfig) metricKafkaPartitionCurrentOffset {
	m := metricKafkaPartitionCurrentOffset{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricKafkaPartitionOldestOffset struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills kafka.partition.oldest_offset metric with initial data.
func (m *metricKafkaPartitionOldestOffset) init() {
	m.data.SetName("kafka.partition.oldest_offset")
	m.data.SetDescription("Oldest offset of partition of topic")
	m.data.SetUnit("1")
	m.data.SetEmptyGauge()
	m.data.Gauge().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricKafkaPartitionOldestOffset) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, topicAttributeValue string, partitionAttributeValue int64) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Gauge().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntValue(val)
	dp.Attributes().PutStr("topic", topicAttributeValue)
	dp.Attributes().PutInt("partition", partitionAttributeValue)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricKafkaPartitionOldestOffset) updateCapacity() {
	if m.data.Gauge().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Gauge().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricKafkaPartitionOldestOffset) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Gauge().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricKafkaPartitionOldestOffset(cfg MetricConfig) metricKafkaPartitionOldestOffset {
	m := metricKafkaPartitionOldestOffset{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricKafkaPartitionReplicas struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills kafka.partition.replicas metric with initial data.
func (m *metricKafkaPartitionReplicas) init() {
	m.data.SetName("kafka.partition.replicas")
	m.data.SetDescription("Number of replicas for partition of topic")
	m.data.SetUnit("{replicas}")
	m.data.SetEmptySum()
	m.data.Sum().SetIsMonotonic(false)
	m.data.Sum().SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricKafkaPartitionReplicas) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, topicAttributeValue string, partitionAttributeValue int64) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntValue(val)
	dp.Attributes().PutStr("topic", topicAttributeValue)
	dp.Attributes().PutInt("partition", partitionAttributeValue)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricKafkaPartitionReplicas) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricKafkaPartitionReplicas) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricKafkaPartitionReplicas(cfg MetricConfig) metricKafkaPartitionReplicas {
	m := metricKafkaPartitionReplicas{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricKafkaPartitionReplicasInSync struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills kafka.partition.replicas_in_sync metric with initial data.
func (m *metricKafkaPartitionReplicasInSync) init() {
	m.data.SetName("kafka.partition.replicas_in_sync")
	m.data.SetDescription("Number of synchronized replicas of partition")
	m.data.SetUnit("{replicas}")
	m.data.SetEmptySum()
	m.data.Sum().SetIsMonotonic(false)
	m.data.Sum().SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricKafkaPartitionReplicasInSync) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, topicAttributeValue string, partitionAttributeValue int64) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntValue(val)
	dp.Attributes().PutStr("topic", topicAttributeValue)
	dp.Attributes().PutInt("partition", partitionAttributeValue)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricKafkaPartitionReplicasInSync) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricKafkaPartitionReplicasInSync) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricKafkaPartitionReplicasInSync(cfg MetricConfig) metricKafkaPartitionReplicasInSync {
	m := metricKafkaPartitionReplicasInSync{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

type metricKafkaTopicPartitions struct {
	data     pmetric.Metric // data buffer for generated metric.
	config   MetricConfig   // metric config provided by user.
	capacity int            // max observed number of data points added to the metric.
}

// init fills kafka.topic.partitions metric with initial data.
func (m *metricKafkaTopicPartitions) init() {
	m.data.SetName("kafka.topic.partitions")
	m.data.SetDescription("Number of partitions in topic.")
	m.data.SetUnit("{partitions}")
	m.data.SetEmptySum()
	m.data.Sum().SetIsMonotonic(false)
	m.data.Sum().SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
	m.data.Sum().DataPoints().EnsureCapacity(m.capacity)
}

func (m *metricKafkaTopicPartitions) recordDataPoint(start pcommon.Timestamp, ts pcommon.Timestamp, val int64, topicAttributeValue string) {
	if !m.config.Enabled {
		return
	}
	dp := m.data.Sum().DataPoints().AppendEmpty()
	dp.SetStartTimestamp(start)
	dp.SetTimestamp(ts)
	dp.SetIntValue(val)
	dp.Attributes().PutStr("topic", topicAttributeValue)
}

// updateCapacity saves max length of data point slices that will be used for the slice capacity.
func (m *metricKafkaTopicPartitions) updateCapacity() {
	if m.data.Sum().DataPoints().Len() > m.capacity {
		m.capacity = m.data.Sum().DataPoints().Len()
	}
}

// emit appends recorded metric data to a metrics slice and prepares it for recording another set of data points.
func (m *metricKafkaTopicPartitions) emit(metrics pmetric.MetricSlice) {
	if m.config.Enabled && m.data.Sum().DataPoints().Len() > 0 {
		m.updateCapacity()
		m.data.MoveTo(metrics.AppendEmpty())
		m.init()
	}
}

func newMetricKafkaTopicPartitions(cfg MetricConfig) metricKafkaTopicPartitions {
	m := metricKafkaTopicPartitions{config: cfg}
	if cfg.Enabled {
		m.data = pmetric.NewMetric()
		m.init()
	}
	return m
}

// MetricsBuilder provides an interface for scrapers to report metrics while taking care of all the transformations
// required to produce metric representation defined in metadata and user config.
type MetricsBuilder struct {
	config                             MetricsBuilderConfig // config of the metrics builder.
	startTime                          pcommon.Timestamp    // start time that will be applied to all recorded data points.
	metricsCapacity                    int                  // maximum observed number of metrics per resource.
	metricsBuffer                      pmetric.Metrics      // accumulates metrics data before emitting.
	buildInfo                          component.BuildInfo  // contains version information.
	metricKafkaBrokers                 metricKafkaBrokers
	metricKafkaConsumerGroupLag        metricKafkaConsumerGroupLag
	metricKafkaConsumerGroupLagSum     metricKafkaConsumerGroupLagSum
	metricKafkaConsumerGroupMembers    metricKafkaConsumerGroupMembers
	metricKafkaConsumerGroupOffset     metricKafkaConsumerGroupOffset
	metricKafkaConsumerGroupOffsetSum  metricKafkaConsumerGroupOffsetSum
	metricKafkaPartitionCurrentOffset  metricKafkaPartitionCurrentOffset
	metricKafkaPartitionOldestOffset   metricKafkaPartitionOldestOffset
	metricKafkaPartitionReplicas       metricKafkaPartitionReplicas
	metricKafkaPartitionReplicasInSync metricKafkaPartitionReplicasInSync
	metricKafkaTopicPartitions         metricKafkaTopicPartitions
}

// metricBuilderOption applies changes to default metrics builder.
type metricBuilderOption func(*MetricsBuilder)

// WithStartTime sets startTime on the metrics builder.
func WithStartTime(startTime pcommon.Timestamp) metricBuilderOption {
	return func(mb *MetricsBuilder) {
		mb.startTime = startTime
	}
}

func NewMetricsBuilder(mbc MetricsBuilderConfig, settings receiver.CreateSettings, options ...metricBuilderOption) *MetricsBuilder {
	mb := &MetricsBuilder{
		config:                             mbc,
		startTime:                          pcommon.NewTimestampFromTime(time.Now()),
		metricsBuffer:                      pmetric.NewMetrics(),
		buildInfo:                          settings.BuildInfo,
		metricKafkaBrokers:                 newMetricKafkaBrokers(mbc.Metrics.KafkaBrokers),
		metricKafkaConsumerGroupLag:        newMetricKafkaConsumerGroupLag(mbc.Metrics.KafkaConsumerGroupLag),
		metricKafkaConsumerGroupLagSum:     newMetricKafkaConsumerGroupLagSum(mbc.Metrics.KafkaConsumerGroupLagSum),
		metricKafkaConsumerGroupMembers:    newMetricKafkaConsumerGroupMembers(mbc.Metrics.KafkaConsumerGroupMembers),
		metricKafkaConsumerGroupOffset:     newMetricKafkaConsumerGroupOffset(mbc.Metrics.KafkaConsumerGroupOffset),
		metricKafkaConsumerGroupOffsetSum:  newMetricKafkaConsumerGroupOffsetSum(mbc.Metrics.KafkaConsumerGroupOffsetSum),
		metricKafkaPartitionCurrentOffset:  newMetricKafkaPartitionCurrentOffset(mbc.Metrics.KafkaPartitionCurrentOffset),
		metricKafkaPartitionOldestOffset:   newMetricKafkaPartitionOldestOffset(mbc.Metrics.KafkaPartitionOldestOffset),
		metricKafkaPartitionReplicas:       newMetricKafkaPartitionReplicas(mbc.Metrics.KafkaPartitionReplicas),
		metricKafkaPartitionReplicasInSync: newMetricKafkaPartitionReplicasInSync(mbc.Metrics.KafkaPartitionReplicasInSync),
		metricKafkaTopicPartitions:         newMetricKafkaTopicPartitions(mbc.Metrics.KafkaTopicPartitions),
	}

	for _, op := range options {
		op(mb)
	}
	return mb
}

// updateCapacity updates max length of metrics and resource attributes that will be used for the slice capacity.
func (mb *MetricsBuilder) updateCapacity(rm pmetric.ResourceMetrics) {
	if mb.metricsCapacity < rm.ScopeMetrics().At(0).Metrics().Len() {
		mb.metricsCapacity = rm.ScopeMetrics().At(0).Metrics().Len()
	}
}

// ResourceMetricsOption applies changes to provided resource metrics.
type ResourceMetricsOption func(pmetric.ResourceMetrics)

// WithResource sets the provided resource on the emitted ResourceMetrics.
// It's recommended to use ResourceBuilder to create the resource.
func WithResource(res pcommon.Resource) ResourceMetricsOption {
	return func(rm pmetric.ResourceMetrics) {
		res.CopyTo(rm.Resource())
	}
}

// WithStartTimeOverride overrides start time for all the resource metrics data points.
// This option should be only used if different start time has to be set on metrics coming from different resources.
func WithStartTimeOverride(start pcommon.Timestamp) ResourceMetricsOption {
	return func(rm pmetric.ResourceMetrics) {
		var dps pmetric.NumberDataPointSlice
		metrics := rm.ScopeMetrics().At(0).Metrics()
		for i := 0; i < metrics.Len(); i++ {
			switch metrics.At(i).Type() {
			case pmetric.MetricTypeGauge:
				dps = metrics.At(i).Gauge().DataPoints()
			case pmetric.MetricTypeSum:
				dps = metrics.At(i).Sum().DataPoints()
			}
			for j := 0; j < dps.Len(); j++ {
				dps.At(j).SetStartTimestamp(start)
			}
		}
	}
}

// EmitForResource saves all the generated metrics under a new resource and updates the internal state to be ready for
// recording another set of data points as part of another resource. This function can be helpful when one scraper
// needs to emit metrics from several resources. Otherwise calling this function is not required,
// just `Emit` function can be called instead.
// Resource attributes should be provided as ResourceMetricsOption arguments.
func (mb *MetricsBuilder) EmitForResource(rmo ...ResourceMetricsOption) {
	rm := pmetric.NewResourceMetrics()
	ils := rm.ScopeMetrics().AppendEmpty()
	ils.Scope().SetName("otelcol/kafkametricsreceiver")
	ils.Scope().SetVersion(mb.buildInfo.Version)
	ils.Metrics().EnsureCapacity(mb.metricsCapacity)
	mb.metricKafkaBrokers.emit(ils.Metrics())
	mb.metricKafkaConsumerGroupLag.emit(ils.Metrics())
	mb.metricKafkaConsumerGroupLagSum.emit(ils.Metrics())
	mb.metricKafkaConsumerGroupMembers.emit(ils.Metrics())
	mb.metricKafkaConsumerGroupOffset.emit(ils.Metrics())
	mb.metricKafkaConsumerGroupOffsetSum.emit(ils.Metrics())
	mb.metricKafkaPartitionCurrentOffset.emit(ils.Metrics())
	mb.metricKafkaPartitionOldestOffset.emit(ils.Metrics())
	mb.metricKafkaPartitionReplicas.emit(ils.Metrics())
	mb.metricKafkaPartitionReplicasInSync.emit(ils.Metrics())
	mb.metricKafkaTopicPartitions.emit(ils.Metrics())

	for _, op := range rmo {
		op(rm)
	}

	if ils.Metrics().Len() > 0 {
		mb.updateCapacity(rm)
		rm.MoveTo(mb.metricsBuffer.ResourceMetrics().AppendEmpty())
	}
}

// Emit returns all the metrics accumulated by the metrics builder and updates the internal state to be ready for
// recording another set of metrics. This function will be responsible for applying all the transformations required to
// produce metric representation defined in metadata and user config, e.g. delta or cumulative.
func (mb *MetricsBuilder) Emit(rmo ...ResourceMetricsOption) pmetric.Metrics {
	mb.EmitForResource(rmo...)
	metrics := mb.metricsBuffer
	mb.metricsBuffer = pmetric.NewMetrics()
	return metrics
}

// RecordKafkaBrokersDataPoint adds a data point to kafka.brokers metric.
func (mb *MetricsBuilder) RecordKafkaBrokersDataPoint(ts pcommon.Timestamp, val int64) {
	mb.metricKafkaBrokers.recordDataPoint(mb.startTime, ts, val)
}

// RecordKafkaConsumerGroupLagDataPoint adds a data point to kafka.consumer_group.lag metric.
func (mb *MetricsBuilder) RecordKafkaConsumerGroupLagDataPoint(ts pcommon.Timestamp, val int64, groupAttributeValue string, topicAttributeValue string, partitionAttributeValue int64) {
	mb.metricKafkaConsumerGroupLag.recordDataPoint(mb.startTime, ts, val, groupAttributeValue, topicAttributeValue, partitionAttributeValue)
}

// RecordKafkaConsumerGroupLagSumDataPoint adds a data point to kafka.consumer_group.lag_sum metric.
func (mb *MetricsBuilder) RecordKafkaConsumerGroupLagSumDataPoint(ts pcommon.Timestamp, val int64, groupAttributeValue string, topicAttributeValue string) {
	mb.metricKafkaConsumerGroupLagSum.recordDataPoint(mb.startTime, ts, val, groupAttributeValue, topicAttributeValue)
}

// RecordKafkaConsumerGroupMembersDataPoint adds a data point to kafka.consumer_group.members metric.
func (mb *MetricsBuilder) RecordKafkaConsumerGroupMembersDataPoint(ts pcommon.Timestamp, val int64, groupAttributeValue string) {
	mb.metricKafkaConsumerGroupMembers.recordDataPoint(mb.startTime, ts, val, groupAttributeValue)
}

// RecordKafkaConsumerGroupOffsetDataPoint adds a data point to kafka.consumer_group.offset metric.
func (mb *MetricsBuilder) RecordKafkaConsumerGroupOffsetDataPoint(ts pcommon.Timestamp, val int64, groupAttributeValue string, topicAttributeValue string, partitionAttributeValue int64) {
	mb.metricKafkaConsumerGroupOffset.recordDataPoint(mb.startTime, ts, val, groupAttributeValue, topicAttributeValue, partitionAttributeValue)
}

// RecordKafkaConsumerGroupOffsetSumDataPoint adds a data point to kafka.consumer_group.offset_sum metric.
func (mb *MetricsBuilder) RecordKafkaConsumerGroupOffsetSumDataPoint(ts pcommon.Timestamp, val int64, groupAttributeValue string, topicAttributeValue string) {
	mb.metricKafkaConsumerGroupOffsetSum.recordDataPoint(mb.startTime, ts, val, groupAttributeValue, topicAttributeValue)
}

// RecordKafkaPartitionCurrentOffsetDataPoint adds a data point to kafka.partition.current_offset metric.
func (mb *MetricsBuilder) RecordKafkaPartitionCurrentOffsetDataPoint(ts pcommon.Timestamp, val int64, topicAttributeValue string, partitionAttributeValue int64) {
	mb.metricKafkaPartitionCurrentOffset.recordDataPoint(mb.startTime, ts, val, topicAttributeValue, partitionAttributeValue)
}

// RecordKafkaPartitionOldestOffsetDataPoint adds a data point to kafka.partition.oldest_offset metric.
func (mb *MetricsBuilder) RecordKafkaPartitionOldestOffsetDataPoint(ts pcommon.Timestamp, val int64, topicAttributeValue string, partitionAttributeValue int64) {
	mb.metricKafkaPartitionOldestOffset.recordDataPoint(mb.startTime, ts, val, topicAttributeValue, partitionAttributeValue)
}

// RecordKafkaPartitionReplicasDataPoint adds a data point to kafka.partition.replicas metric.
func (mb *MetricsBuilder) RecordKafkaPartitionReplicasDataPoint(ts pcommon.Timestamp, val int64, topicAttributeValue string, partitionAttributeValue int64) {
	mb.metricKafkaPartitionReplicas.recordDataPoint(mb.startTime, ts, val, topicAttributeValue, partitionAttributeValue)
}

// RecordKafkaPartitionReplicasInSyncDataPoint adds a data point to kafka.partition.replicas_in_sync metric.
func (mb *MetricsBuilder) RecordKafkaPartitionReplicasInSyncDataPoint(ts pcommon.Timestamp, val int64, topicAttributeValue string, partitionAttributeValue int64) {
	mb.metricKafkaPartitionReplicasInSync.recordDataPoint(mb.startTime, ts, val, topicAttributeValue, partitionAttributeValue)
}

// RecordKafkaTopicPartitionsDataPoint adds a data point to kafka.topic.partitions metric.
func (mb *MetricsBuilder) RecordKafkaTopicPartitionsDataPoint(ts pcommon.Timestamp, val int64, topicAttributeValue string) {
	mb.metricKafkaTopicPartitions.recordDataPoint(mb.startTime, ts, val, topicAttributeValue)
}

// Reset resets metrics builder to its initial state. It should be used when external metrics source is restarted,
// and metrics builder should update its startTime and reset it's internal state accordingly.
func (mb *MetricsBuilder) Reset(options ...metricBuilderOption) {
	mb.startTime = pcommon.NewTimestampFromTime(time.Now())
	for _, op := range options {
		op(mb)
	}
}
