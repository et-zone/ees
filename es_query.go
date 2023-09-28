package ees

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/et-zone/ees/elasticsql"
	"github.com/olivere/elastic/v7"
	"reflect"
	"strings"
)

type Eelastic struct {
	client *elastic.Client
}

func (e *Eelastic) Count(ctx context.Context, indexName ...string) (int64, error) {
	count, err := e.client.Count(indexName...).Do(ctx)
	if err != nil {
		return 0, err
	}
	return count, err
}

// between 支持数字类型,支持指定字段查询
func (e *Eelastic) QuerySql(ctx context.Context, sql string, result interface{}) (total int64, err error) {
	if sql == "" {
		return 0, errors.New("sql is ''")
	}
	if !strings.Contains(sql, "limit") && !strings.Contains(sql, "LIMIT") {
		sql += " limit " + DefaultLimit
	}
	dsl, table, err := elasticsql.Convert(sql)
	if err != nil {
		return 0, err
	}
	rep, err := e.client.Search(table).Source(dsl).Do(ctx)
	if err != nil {
		return 0, err
	}

	v := reflect.ValueOf(result)

	if v.IsNil() {
		//return 0,errors.New("out []interface{} is nil ")
		return 0, errors.New("result can not nil")
	}

	if v.Kind() != reflect.Ptr {
		//return 0,errors.New("out []interface{} is nil ")
		return 0, errors.New("result must a addr of Slice")
	}

	len := len(rep.Hits.Hits)

	if len == 0 {
		return 0, nil
	}
	if v.Elem().Kind() == reflect.Slice {
		dataJson := "["
		for i, v := range rep.Hits.Hits {
			if i == len-1 {
				dataJson += string(v.Source)
			} else {
				dataJson += string(v.Source) + ","
			}
		}
		dataJson += "]"
		err := json.Unmarshal([]byte(dataJson), result)
		if err != nil {
			return 0, err
		}
		return rep.Hits.TotalHits.Value, err

	}

	if len > 0 {
		err := json.Unmarshal(rep.Hits.Hits[0].Source, result)
		if err != nil {
			return 0, err
		}
		return 1, err
	}

	return 0, err
}

func (e *Eelastic) Search(ctx context.Context, query *EQueryReq) (*elastic.SearchResult, error) {
	cli := e.client.Search(query.Index...).Query(query.Query.query)
	for k, v := range query.aggregationMap {
		cli = cli.Aggregation(k, v)
	}

	if query.fetchSourceContext != nil {
		cli = cli.FetchSourceContext(query.fetchSourceContext)
	}
	if query.collapseField != "" {
		cli = cli.Collapse(elastic.NewCollapseBuilder(query.collapseField))
	}

	return cli.SortBy(query.sort...).From(query.from).Size(query.size).Pretty(true).Do(ctx)
}

func (e *Eelastic) QueryAggregations(ctx context.Context, query *EQueryReq) (*elastic.SearchResult, error) {
	cli := e.client.Search(query.Index...).Query(query.Query.query).From(0).Size(0)
	for k, v := range query.aggregationMap {
		cli = cli.Aggregation(k, v)
	}
	return cli.Pretty(true).Do(ctx)
}

func (e *Eelastic) GetIndexDetail(ctx context.Context, index ...string) (map[string]interface{}, error) {
	//GetMaping
	if len(index) == 0 {
		return nil, nil
	}
	ret, err := e.client.GetMapping().Index(index...).Pretty(true).Do(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			return nil, nil
		}
		return nil, err
	}
	return ret, err
}

// if version =8.*  some version can not use it to search
func (e *Eelastic) MultiSearch(ctx context.Context, req MultiQueryReq) (*elastic.MultiSearchResult, error) {
	//index:=reqList.Index
	list := []*elastic.SearchRequest{}
	for _, v := range req.Req {
		request := e.getSearchRequest(v)
		list = append(list, request)
	}
	return e.client.MultiSearch().Index(req.Index).Add(list...).Pretty(true).Do(ctx)
}

func (e *Eelastic) getSearchRequest(query *EQueryReq) *elastic.SearchRequest {

	request := elastic.NewSearchRequest().
		Index(query.Index...).
		Query(query.Query.query).
		SortBy(query.sort...).
		From(query.from).
		Size(query.size)
	for name, aggregation := range query.aggregationMap {
		if name != "" {
			request = request.Aggregation(name, aggregation)
		}
	}

	if query.fetchSourceContext != nil {
		request = request.FetchSourceContext(query.fetchSourceContext)
	}
	if query.collapseField != "" {
		request = request.Collapse(elastic.NewCollapseBuilder(query.collapseField))
	}
	return request
}
