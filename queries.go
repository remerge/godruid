package godruid

import (
	"encoding/json"
)

// Check http://druid.io/docs/0.6.154/Querying.html#query-operators for detail description.

// The Query interface stands for any kinds of druid query.
type Query interface {
	Setup()
	OnResponse(content []byte) error
}

// ---------------------------------
// GroupBy Query
// ---------------------------------

type QueryGroupBy struct {
	QueryType        string                 `json:"queryType"`
	DataSource       string                 `json:"dataSource"`
	Dimensions       []DimSpec              `json:"dimensions"`
	Granularity      Granlarity             `json:"granularity"`
	LimitSpec        *Limit                 `json:"limitSpec,omitempty"`
	Having           *Having                `json:"having,omitempty"`
	Filter           *Filter                `json:"filter,omitempty"`
	Aggregations     []Aggregation          `json:"aggregations"`
	PostAggregations []PostAggregation      `json:"postAggregations,omitempty"`
	Intervals        []string               `json:"intervals"`
	Context          map[string]interface{} `json:"context,omitempty"`
	QueryResult      []GroupbyItem          `json:"-"`
}

type GroupbyItem struct {
	Version   string                 `json:"version"`
	Timestamp string                 `json:"timestamp"`
	Event     map[string]interface{} `json:"event"`
}

func (q *QueryGroupBy) Setup() { q.QueryType = "groupBy" }
func (q *QueryGroupBy) OnResponse(content []byte) error {
	res := new([]GroupbyItem)
	err := json.Unmarshal(content, res)
	if err != nil {
		return err
	}
	q.QueryResult = *res
	return nil
}

// ---------------------------------
// Search Query
// ---------------------------------

type QuerySearch struct {
	QueryType        string                 `json:"queryType"`
	DataSource       string                 `json:"dataSource"`
	Granularity      Granlarity             `json:"granularity"`
	Filter           *Filter                `json:"filter,omitempty"`
	Intervals        []string               `json:"intervals"`
	SearchDimensions []string               `json:"searchDimensions,omitempty"`
	Query            *SearchQuery           `json:"query"`
	Sort             *SearchSort            `json:"sort"`
	Context          map[string]interface{} `json:"context,omitempty"`

	QueryResult []SearchItem `json:"-"`
}

type SearchItem struct {
	Timestamp string     `json:"timestamp"`
	Result    []DimValue `json:"result"`
}

type DimValue struct {
	Dimension string `json:"dimension"`
	Value     string `json:"value"`
}

func (q *QuerySearch) Setup() { q.QueryType = "search" }
func (q *QuerySearch) OnResponse(content []byte) error {
	res := new([]SearchItem)
	err := json.Unmarshal(content, res)
	if err != nil {
		return err
	}
	q.QueryResult = *res
	return nil
}

// ---------------------------------
// SegmentMetadata Query
// ---------------------------------

type QuerySegmentMetadata struct {
	QueryType  string                 `json:"queryType"`
	DataSource string                 `json:"dataSource"`
	Intervals  []string               `json:"intervals"`
	ToInclude  *ToInclude             `json:"toInclude,omitempty"`
	Merge      interface{}            `json:"merge,omitempty"`
	Context    map[string]interface{} `json:"context,omitempty"`

	QueryResult []SegmentMetaData `json:"-"`
}

type SegmentMetaData struct {
	Id        string                `json:"id"`
	Intervals []string              `json:"intervals"`
	Columns   map[string]ColumnItem `json:"columns"`
}

type ColumnItem struct {
	Type        string      `json:"type"`
	Size        int         `json:"size"`
	Cardinality interface{} `json:"cardinality"`
}

func (q *QuerySegmentMetadata) setup() { q.QueryType = "segmentMetadata" }
func (q *QuerySegmentMetadata) onResponse(content []byte) error {
	res := new([]SegmentMetaData)
	err := json.Unmarshal(content, res)
	if err != nil {
		return err
	}
	q.QueryResult = *res
	return nil
}

// ---------------------------------
// TimeBoundary Query
// ---------------------------------

type QueryTimeBoundary struct {
	QueryType  string                 `json:"queryType"`
	DataSource string                 `json:"dataSource"`
	Bound      string                 `json:"bound,omitempty"`
	Context    map[string]interface{} `json:"context,omitempty"`

	QueryResult []TimeBoundaryItem `json:"-"`
}

type TimeBoundaryItem struct {
	Timestamp string       `json:"timestamp"`
	Result    TimeBoundary `json:"result"`
}

type TimeBoundary struct {
	MinTime string `json:"minTime"`
	MaxTime string `json:"minTime"`
}

func (q *QueryTimeBoundary) setup() { q.QueryType = "timeBoundary" }
func (q *QueryTimeBoundary) onResponse(content []byte) error {
	res := new([]TimeBoundaryItem)
	err := json.Unmarshal(content, res)
	if err != nil {
		return err
	}
	q.QueryResult = *res
	return nil
}

// ---------------------------------
// Timeseries Query
// ---------------------------------

type QueryTimeseries struct {
	QueryType        string                 `json:"queryType"`
	DataSource       string                 `json:"dataSource"`
	Granularity      Granlarity             `json:"granularity"`
	Filter           *Filter                `json:"filter,omitempty"`
	Aggregations     []Aggregation          `json:"aggregations"`
	PostAggregations []PostAggregation      `json:"postAggregations,omitempty"`
	Intervals        []string               `json:"intervals"`
	Context          map[string]interface{} `json:"context,omitempty"`

	QueryResult []Timeseries `json:"-"`
}

type Timeseries struct {
	Timestamp string                 `json:"timestamp"`
	Result    map[string]interface{} `json:"result"`
}

func (q *QueryTimeseries) Setup() { q.QueryType = "timeseries" }
func (q *QueryTimeseries) OnResponse(content []byte) error {
	res := new([]Timeseries)
	err := json.Unmarshal(content, res)
	if err != nil {
		return err
	}
	q.QueryResult = *res
	return nil
}

// ---------------------------------
// TopN Query
// ---------------------------------

type QueryTopN struct {
	QueryType        string                 `json:"queryType"`
	DataSource       string                 `json:"dataSource"`
	Granularity      Granlarity             `json:"granularity"`
	Dimension        DimSpec                `json:"dimension"`
	Threshold        int                    `json:"threshold"`
	Metric           *TopNMetric            `json:"metric"`
	Filter           *Filter                `json:"filter,omitempty"`
	Aggregations     []Aggregation          `json:"aggregations"`
	PostAggregations []PostAggregation      `json:"postAggregations,omitempty"`
	Intervals        []string               `json:"intervals"`
	Context          map[string]interface{} `json:"context,omitempty"`

	QueryResult []TopNItem `json:"-"`
}

type TopNItem struct {
	Timestamp string                   `json:"timestamp"`
	Result    []map[string]interface{} `json:"result"`
}

func (q *QueryTopN) Setup() { q.QueryType = "topN" }
func (q *QueryTopN) OnResponse(content []byte) error {
	res := new([]TopNItem)
	err := json.Unmarshal(content, res)
	if err != nil {
		return err
	}
	q.QueryResult = *res
	return nil
}
