package encode

import (
	"errors"
	"fmt"
	"reflect"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	ErrValueNotIntOrFloat = errors.New("value is not of type int or float64")
)

type ExportMetricsServiceRequest struct {
	// An array of ResourceMetrics.
	// For data coming from a single resource this array will typically contain one
	// element. Intermediary nodes (such as OpenTelemetry Collector) that receive
	// data from multiple origins typically batch the data before forwarding further and
	// in that case this array will contain multiple elements.
	ResourceMetrics []*ResourceMetrics `protobuf:"bytes,1,rep,name=resource_metrics,json=resourceMetrics,proto3" json:"resource_metrics,omitempty"`
}

// ResourceMetrics is a collection of ScopeMetrics and the associated Resource
// that created them.
type ResourceMetrics struct {
	// Resource represents the entity that collected the metrics.
	Resource Attributes `protobuf:"name=resource,json=resource,proto3" json:"resource"`

	// ScopeMetrics are the collection of metrics with unique Scopes.
	ScopeMetrics []ScopeMetrics `protobuf:"name=scope_metrics,json=scope_metrics,proto3" json:"scope_metrics"`
}

type Attributes struct {
	Attributes *resource.Resource `protobuf:"name=attributes,json=attributes,proto3" json:"attributes"`
}

// ScopeMetrics is a collection of Metrics Produces by a Meter.
type ScopeMetrics struct {
	// Scope is the Scope that the Meter was created with.
	Scope instrumentation.Scope
	// Metrics are a list of aggregations created by the Meter.
	Metrics []Metrics
}

// Metrics is a collection of one or more aggregated timeseries from an Instrument.
type Metrics struct {
	// Name is the name of the Instrument that created this data.
	Name string
	// Description is the description of the Instrument, which can be used in documentation.
	Description string
	// Unit is the unit in which the Instrument reports.
	Unit string
	// Data is the aggregated data from an Instrument.
	Summary *Summary `protobuf:"bytes,11,opt,name=summary,proto3,oneof" json:"summary,omitempty"`

	Sum *Sum `protobuf:"bytes,6,opt,name=sum,proto3,oneof" json:"sum,omitempty"`

	Histogram *Histogram `protobuf:"bytes,8,opt,name=histogram,proto3,oneof" json:"histogram,omitempty"`

	Gauge *Gauge `protobuf:"bytes,4,opt,name=gauge,proto3,oneof" json:"gauge,omitempty"`
}

type Histogram struct {
	DataPoints []HistogramDataPoint `protobuf:"bytes,1,rep,name=data_points,json=dataPoints,proto3" json:"data_points,omitempty"`

	AggregationTemporality AggregationTemporality `protobuf:"varint,2,opt,name=aggregation_temporality,json=aggregationTemporality,proto3,enum=opentelemetry.proto.metrics.v1.AggregationTemporality" json:"aggregation_temporality,omitempty"`
}

type HistogramDataPoint struct {
	Attributes []attribute.KeyValue

	StartTimeUnixNano uint64 `protobuf:"fixed64,2,opt,name=start_time_unix_nano,json=startTimeUnixNano,proto3" json:"start_time_unix_nano,omitempty"`

	TimeUnixNano uint64 `protobuf:"fixed64,3,opt,name=time_unix_nano,json=timeUnixNano,proto3" json:"time_unix_nano,omitempty"`

	Count uint64 `protobuf:"fixed64,4,opt,name=count,proto3" json:"count,omitempty"`

	Sum float64 `protobuf:"fixed64,5,opt,name=sum,proto3" json:"sum,omitempty"`

	BucketCounts []uint64 `protobuf:"fixed64,6,rep,packed,name=bucket_counts,json=bucketCounts,proto3" json:"bucket_counts,omitempty"`

	ExplicitBounds []float64 `protobuf:"fixed64,7,rep,packed,name=explicit_bounds,json=explicitBounds,proto3" json:"explicit_bounds,omitempty"`
}

type Sum struct {
	DataPoints []DataPoint `protobuf:"bytes,1,rep,name=data_points,json=dataPoints,proto3" json:"data_points,omitempty"`

	IsMonotonic bool `protobuf:"varint,3,opt,name=is_monotonic,json=isMonotonic,proto3" json:"is_monotonic,omitempty"`

	AggregationTemporality AggregationTemporality `protobuf:"varint,2,opt,name=aggregation_temporality,json=aggregationTemporality,proto3,enum=opentelemetry.proto.metrics.v1.AggregationTemporality" json:"aggregation_temporality,omitempty"`
}

type AggregationTemporality int32

type DataPoint struct {
	StartTimeUnixNano uint64 `protobuf:"fixed64,2,opt,name=start_time_unix_nano,json=startTimeUnixNano,proto3" json:"start_time_unix_nano,omitempty"`

	TimeUnixNano uint64 `protobuf:"fixed64,3,opt,name=time_unix_nano,json=timeUnixNano,proto3" json:"time_unix_nano,omitempty"`

	Value int64 `protobuf:"fixed64,4,opt,name=value,proto3" json:"value,omitempty"`

	Attributes []attribute.KeyValue
}

type Summary struct {
	// DataPoints are the individual aggregated measurements with unique
	DataPoints []SummaryDataPoint `protobuf:"name=data_points,json=data_points,proto3" json:"data_points"`
}

type SummaryDataPoint struct {
	Attributes        []attribute.KeyValue
	StartTimeUnixNano uint64                             `protobuf:"fixed64,2,opt,name=start_time_unix_nano,json=startTimeUnixNano,proto3" json:"start_time_unix_nano,omitempty"`
	TimeUnixNano      uint64                             `protobuf:"fixed64,3,opt,name=time_unix_nano,json=timeUnixNano,proto3" json:"time_unix_nano,omitempty"`
	Count             uint64                             `protobuf:"fixed64,4,opt,name=count,proto3" json:"count,omitempty"`
	Sum               float64                            `protobuf:"fixed64,5,opt,name=sum,proto3" json:"sum,omitempty"`
	QuantileValues    []SummaryDataPoint_ValueAtQuantile `protobuf:"bytes,6,rep,name=quantiles,json=quantiles,proto3" json:"quantiles,omitempty"`
	Flags             uint32                             `protobuf:"varint,8,opt,name=flags,proto3" json:"flags,omitempty"`
}

type SummaryDataPoint_ValueAtQuantile struct {
	Quantile float64 `protobuf:"fixed64,1,opt,name=quantile,proto3" json:"quantile,omitempty"`
	Value    float64 `protobuf:"fixed64,2,opt,name=value,proto3" json:"value,omitempty"`
}

type NumberDataPoint struct {
	Attributes        []attribute.KeyValue `protobuf:"bytes,7,rep,name=attributes,proto3" json:"attributes"`
	StartTimeUnixNano uint64               `protobuf:"fixed64,2,opt,name=start_time_unix_nano,json=startTimeUnixNano,proto3" json:"start_time_unix_nano,omitempty"`
	TimeUnixNano      uint64               `protobuf:"fixed64,3,opt,name=time_unix_nano,json=timeUnixNano,proto3" json:"time_unix_nano,omitempty"`
	//Value             isNumberDataPoint_Value `protobuf_oneof:"value"`
	Value    interface{} `json:"value,omitempty"`
	AsDouble *float64    `json:"as_double,omitempty"`
	AsInt    *int        `json:"as_int,omitempty"`
	// (Optional) List of exemplars collected from
	// measurements that were used to form the data point
	// Exemplars []Exemplar `protobuf:"bytes,5,rep,name=exemplars,proto3" json:"exemplars"`
	// Flags that apply to this specific data point.  See DataPointFlags
	// for the available flags and their meaning.
	// Flags uint32 `protobuf:"varint,8,opt,name=flags,proto3" json:"flags,omitempty"`
}

type Gauge struct {
	DataPoints []*NumberDataPoint `protobuf:"bytes,1,rep,name=data_points,json=dataPoints,proto3" json:"data_points,omitempty"`
}

func (n *NumberDataPoint) SetValue(i interface{}) error {
	switch reflect.TypeOf(i).Kind() {
	case reflect.Int:
		fmt.Println("int")
		v, ok := i.(int)
		if !ok {
			return ErrValueNotIntOrFloat
		}

		n.AsInt = &v
	case reflect.Float64:
		v, ok := i.(float64)
		if !ok {
			return ErrValueNotIntOrFloat
		}

		n.AsDouble = &v
	}
	return nil
}
