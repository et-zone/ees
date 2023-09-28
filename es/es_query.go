package es

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/et-zone/ees/elasticsql"
	"github.com/olivere/elastic/v7"
	"reflect"
	"strings"
)

func Count(ctx context.Context, indexName ...string) (int64, error) {
	count, err := client.Count(indexName...).Do(ctx)
	if err != nil {
		return 0, err
	}
	return count, err
}

// between 支持数字类型,支持指定字段查询
func QuerySql(ctx context.Context, sql string, result interface{}) (total int64, err error) {
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
	rep, err := client.Search(table).Source(dsl).Do(ctx)
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

func Query(ctx context.Context, query EQueryReq) (*elastic.SearchResult, error) {
	cli := client.Search(query.Index...).Query(query.Query)
	for k, v := range query.AggregationMap {
		cli = cli.Aggregation(k, v)
	}

	if query.FetchSourceContext != nil {
		cli = cli.FetchSourceContext(query.FetchSourceContext)
	}
	if query.CollapseField != "" {
		cli = cli.Collapse(elastic.NewCollapseBuilder(query.CollapseField))
	}

	return cli.SortBy(query.Sort...).From(query.From).Size(query.Size).Pretty(true).Do(ctx)
}

func QueryAggregations(ctx context.Context, query EQueryReq) (*elastic.SearchResult, error) {
	cli := client.Search(query.Index...).Query(query.Query).From(0).Size(0)
	for k, v := range query.AggregationMap {
		cli = cli.Aggregation(k, v)
	}
	return cli.Pretty(true).Do(ctx)
}

func GetIndexDetail(ctx context.Context, index ...string) (map[string]interface{}, error) {
	//GetMaping
	if len(index) == 0 {
		return nil, nil
	}
	ret, err := client.GetMapping().Index(index...).Do(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			return nil, nil
		}
		return nil, err
	}
	return ret, err
}
