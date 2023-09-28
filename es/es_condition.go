package es

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"reflect"
	"strings"
)

const AggregationSize = 1000
const CollapseSize = 1000
const (
	SortAsc  = "ASC"
	SortDesc = "DESC"
)

type EQueryReq struct {
	Index              []string
	Query              EQuery
	From               int
	Size               int
	Sort               []elastic.Sorter
	Source             string
	FetchSourceContext *elastic.FetchSourceContext
	CollapseField      string
	AggregationMap     map[string]elastic.Aggregation //aliseName:field
}

func NewEQueryReq(index ...string) *EQueryReq {
	return &EQueryReq{
		Index:          index,
		Query:          EQuery{elastic.NewBoolQuery()},
		AggregationMap: map[string]elastic.Aggregation{},
	}
}

func (e *EQueryReq) Aggregation(aliseName string, fieldName string) *EQueryReq {
	if aliseName == "" || fieldName == "" {
		return e
	}
	e.AggregationMap[aliseName] = elastic.NewTermsAggregation().Field(fieldName).Size(AggregationSize)
	return e
}

func (e *EQueryReq) SetPage(pageNo, pageSize int) *EQueryReq {
	e.From = (pageNo - 1) * pageSize
	e.Size = pageSize
	if e.From <= 0 {
		e.From = 0
	}
	return e
}

//fetchSource = true ==>need show fields
func (e *EQueryReq) SetFetchSource(fetchSource bool) *EQueryReq {
	if e.FetchSourceContext == nil {
		e.FetchSourceContext = elastic.NewFetchSourceContext(fetchSource)
	} else {
		e.FetchSourceContext.SetFetchSource(fetchSource)
	}
	return e
}

func (e *EQueryReq) SetIncludesSource(includes ...string) *EQueryReq {
	if e.FetchSourceContext == nil {
		e.FetchSourceContext = elastic.NewFetchSourceContext(true)
	}
	e.FetchSourceContext.Include(includes...)
	return e
}

func (e *EQueryReq) SetExcludesSource(excludes ...string) *EQueryReq {
	if e.FetchSourceContext == nil {
		e.FetchSourceContext = elastic.NewFetchSourceContext(true)
	}
	e.FetchSourceContext.Exclude(excludes...)
	return e
}

func (e *EQueryReq) SetCollapseField(collapseField string) *EQueryReq {
	e.CollapseField = collapseField
	return e
}

/*
	eg: elastic.NewFieldSort("field").Desc()
*/
func (e *EQueryReq) SetSort(sortField elastic.Sorter) *EQueryReq {
	e.Sort = append(e.Sort, sortField)
	return e
}

func BuildQuery(params interface{}) (query *elastic.BoolQuery) {
	query = elastic.NewBoolQuery()
	var brandModelFilter *elastic.BoolQuery
	object := reflect.ValueOf(params)
	elems := object.Elem()
	typeOfObject := elems.Type()
	for i := 0; i < elems.NumField(); i++ {
		fieldName := typeOfObject.Field(i).Name
		field := elems.Field(i)
		if field.IsZero() {
			continue
		}
		fieldName = strings.ToLower(fieldName[0:1]) + fieldName[1:]
		switch field.Kind() {
		case reflect.String:
			query.Must(elastic.NewMatchPhraseQuery(fieldName, field.Interface()))
		case reflect.Uint64:
			q := elastic.NewTermQuery(fieldName, field.Interface())
			query.Must(q)
		case reflect.Array, reflect.Slice:
			if field.Len() > 0 {
				values := make([]interface{}, 0)
				for j := 0; j < field.Len(); j++ {
					values = append(values, field.Index(j).Interface())
				}
				q := elastic.NewTermsQuery(fieldName, values...)
				query.Must(q)
			}
		case reflect.Struct:
			structElems := field.Convert(field.Type())
			typeOfStruct := structElems.Type()
			q := elastic.NewRangeQuery(fieldName)
			for j := 0; j < structElems.NumField(); j++ {
				f := structElems.Field(j)
				fName := typeOfStruct.Field(j).Name
				if fName == "Min" && !f.IsZero() {
					q.Gte(f.Interface())
				}
				if fName == "Max" && !f.IsZero() {
					q.Lte(f.Interface())
				}
			}
			query.Must(q)
		case reflect.Ptr:
			if !field.IsNil() {
				structElems := field.Convert(field.Type())
				if !structElems.IsZero() {
					q := elastic.NewTermQuery(fieldName, structElems.Elem().Interface())
					query.Must(q)
				}

			}
		default:
			fmt.Println("ok")
		}
	}
	if brandModelFilter != nil {
		query.Must(brandModelFilter)
	}
	return query
}

type EQuery struct {
	*elastic.BoolQuery
}

func (e *EQuery) And() {

}
