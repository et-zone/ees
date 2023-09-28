package ees

import (
	"github.com/olivere/elastic/v7"
)

//size
const (
	AggregationSize = 1000
	CollapseSize    = 1000
)

//sort
const (
	SortAsc  = "ASC"
	SortDesc = "DESC"
)

func NewQuery() Query {
	return Query{elastic.NewBoolQuery()}
}

type EQueryReq struct {
	Index              []string
	Query              Query
	from               int
	size               int
	sort               []elastic.Sorter
	fetchSourceContext *elastic.FetchSourceContext
	collapseField      string
	aggregationMap     map[string]elastic.Aggregation //aliseName:field
}

func NewEQueryReq(index ...string) *EQueryReq {
	return &EQueryReq{
		Index:          index,
		Query:          Query{elastic.NewBoolQuery()},
		aggregationMap: map[string]elastic.Aggregation{},
	}
}

func (e *EQueryReq) SetAggregation(aliseName string, fieldName string) *EQueryReq {
	if aliseName == "" || fieldName == "" {
		return e
	}
	e.aggregationMap[aliseName] = elastic.NewTermsAggregation().Field(fieldName).Size(AggregationSize)
	return e
}

func (e *EQueryReq) SetPage(pageNo, pageSize int) *EQueryReq {
	e.from = (pageNo - 1) * pageSize
	e.size = pageSize
	if e.from <= 0 {
		e.from = 0
	}
	return e
}

//fetchSource = true ==>need show fields
func (e *EQueryReq) SetFetchSource(fetchSource bool) *EQueryReq {
	if e.fetchSourceContext == nil {
		e.fetchSourceContext = elastic.NewFetchSourceContext(fetchSource)
	} else {
		e.fetchSourceContext.SetFetchSource(fetchSource)
	}
	return e
}

func (e *EQueryReq) SetIncludesSource(includes ...string) *EQueryReq {
	if e.fetchSourceContext == nil {
		e.fetchSourceContext = elastic.NewFetchSourceContext(true)
	}
	e.fetchSourceContext.Include(includes...)
	return e
}

func (e *EQueryReq) SetExcludesSource(excludes ...string) *EQueryReq {
	if e.fetchSourceContext == nil {
		e.fetchSourceContext = elastic.NewFetchSourceContext(true)
	}
	e.fetchSourceContext.Exclude(excludes...)
	return e
}

func (e *EQueryReq) SetCollapseField(collapseField string) *EQueryReq {
	e.collapseField = collapseField
	return e
}

/*
	eg: elastic.NewFieldSort("field").Desc()
	desc = 4,3,2,1
*/
func (e *EQueryReq) SetSort(sortField ...elastic.Sorter) *EQueryReq {
	e.sort = append(e.sort, sortField...)
	return e
}

type Query struct {
	query *elastic.BoolQuery
}

// A=true and B=true
func (e *Query) And(fields ...elastic.Query) *Query {
	e.query.Must(fields...)
	return e
}

// A!=True and B!=true
func (e *Query) Not(fields ...elastic.Query) *Query {
	e.query.MustNot(fields...)
	return e
}

// A=true or B=true
func (e *Query) OR(fields ...elastic.Query) *Query {
	e.query.Should(fields...)
	return e
}

func FieldIN(fieldName string, val ...interface{}) *elastic.TermsQuery {
	return elastic.NewTermsQuery(fieldName, val...)
}

func FieldEq(fieldName string, val interface{}) *elastic.TermQuery {
	return elastic.NewTermQuery(fieldName, val)
}

func FieldRange(fieldName string) *elastic.RangeQuery {
	return elastic.NewRangeQuery(fieldName)
}

func FieldNotNull(fieldName string) *elastic.ExistsQuery {
	return elastic.NewExistsQuery(fieldName)
}
func FieldLike(fieldName, val string) *elastic.FuzzyQuery {
	return elastic.NewFuzzyQuery(fieldName, val)
}

// text type to Analyzer
func FieldAnalyzer(fieldName, val string) *elastic.MatchPhraseQuery {
	return elastic.NewMatchPhraseQuery(fieldName, val)
}

//top by ID
func FieldTopID(query Query, ids ...string) *elastic.PinnedQuery {
	return elastic.NewPinnedQuery().Ids(ids...).Organic(query.query)
}

func FieldStartWith(fieldName, subString string) *elastic.PrefixQuery {
	return elastic.NewPrefixQuery(fieldName, subString)
}

//list query
func FieldArray(fieldName string) *elastic.SliceQuery {
	//The id of the slice
	//The maximum number of slices
	return elastic.NewSliceQuery().Field(fieldName)
}
