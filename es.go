package ees

//版本7.x
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/et-zone/ees/elasticsql"

	"github.com/olivere/elastic/v7"
)

const (
	CreateSucc   = "created"
	UpdateSucc   = "updated"
	DeleteSucc   = "deleted"
	Noop         = "noop"
	DefaultLimit = "1000"
)

var client *elastic.Client

func stop() {
	client.Stop()
}

// init client
func InitESClient(opts ...elastic.ClientOptionFunc) (err error) {
	//client, err = elastic.NewClient(
	//	elastic.SetURL(host),elastic.SetSniff(false),
	//)
	client, err = elastic.NewClient(opts...)
	if err != nil {
		return errors.New(fmt.Sprintf("can't connect to elasticsearch | err : %s \n", err))
	}

	log.Println("connect to elasticsearch success")
	return
}

//Here we think table = index, all _type = default("_doc")
func InitTable(ctx context.Context, table string, mappings interface{}) (bool, error) {
	//_type default = "_doc"

	isExit, _ := client.IndexExists(table).Do(ctx)
	if isExit {
		return true, nil
	}

	rep, err := client.CreateIndex(table).BodyJson(mappings).Do(ctx)
	if err != nil {
		return false, err
	}
	return rep.ShardsAcknowledged, err
}

func UpsertOneESData(ctx context.Context, table string, id string, value interface{}) (isok bool, err error) {
	doc := elastic.NewBulkUpdateRequest().Index(table).Id(id).Doc(value).DocAsUpsert(true)
	bulk := client.Bulk()
	rep, err := bulk.Add(doc).Do(ctx)

	res := rep.Items[0]
	if rep != nil && rep.Errors {
		e, _ := json.Marshal(res["update"].Error)
		return false, errors.New(string(e))
	}
	return true, err
}


// values key is id
func UpsertAllESData(ctx context.Context, table string, values map[string]interface{}) (int64, error) {
	docs := []elastic.BulkableRequest{}
	for k, v := range values {
		doc := elastic.NewBulkUpdateRequest().Index(table).Id(k).Doc(v).DocAsUpsert(true)
		docs = append(docs, doc)
	}
	bulk := client.Bulk()
	rep, err := bulk.Add(docs...).Do(ctx)

	if err != nil {
		return 0, err
	}

	if rep != nil && rep.Errors {
		e, _ := json.Marshal(rep.Items[0]["update"].Error)
		return 0, errors.New(string(e))
	}
	return int64(len(rep.Items)), err
}

// get data
func getESItem(index ,id string) (rep *elastic.GetResult, err error) {
	rep, err = client.Get().
		Index(index).
		Id(id).
		Do(context.Background())
	return
}

// query
func queryESItem(index ...string) (rep *elastic.SearchResult, err error) {
	q := elastic.NewQueryStringQuery("icon:ccc")
	rep, err = client.Search(index...).
		Query(q).
		Do(context.Background())
	return
}

func Client()*elastic.Client{
	return client
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
		return rep.Hits.TotalHits.Value, nil
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

//删除数据
func DelESItemByID(ctx context.Context, tableName string, id int64) (isDel bool, err error) {
	if tableName == "" || id == 0 {
		return false, errors.New("table or id not null")
	}
	rep, err := client.Delete().
		Index(tableName).
		Id(fmt.Sprintf("%v", id)).
		Do(ctx)

	if rep != nil && rep.Result == DeleteSucc {
		return true, err
	}
	return false, nil
}

//ids can int or string
func DelESItemByIDs(ctx context.Context, tableName string, ids []interface{}) (count int64, err error) {
	if tableName == "" {
		return 0, errors.New("table not null")
	}
	if len(ids) == 0 {
		return 0, errors.New("ids not null")
	}
	q := elastic.NewTermsQuery("_id", ids...)
	rep, err := client.DeleteByQuery(tableName).Query(q).Do(ctx)
	if err != nil {
		return 0, err
	}
	return rep.Deleted, err

}

func DelTable(ctx context.Context, tableName string) (isok bool, err error) {
	if tableName == "" {
		return false, errors.New("table not null")
	}
	rep, err := client.DeleteIndex(tableName).Do(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			return true, nil
		}
		return false, err
	}
	return rep.Acknowledged, err

}

func GetTableDetail(ctx context.Context, tableName ...string) (map[string]interface{}, error) {
	//GetMaping
	if len(tableName) == 0 {
		return nil, nil
	}
	ret, err := client.GetMapping().Index(tableName...).Do(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			return nil, nil
		}
		return nil, err
	}
	return ret, err
}


func IndexSetAlias(ctx context.Context,indexName,AliasName string)(bool,error){
	r,err:=client.Alias().Add(indexName,AliasName).Do(ctx)
	if err!=nil{
		return false,err
	}
	return r.ShardsAcknowledged,err
}

func IndexDelAlias(ctx context.Context,indexName,AliasName string)(bool,error){
	r,err:=client.Alias().Remove(indexName,AliasName).Do(ctx)
	if err!=nil{
		return false,err
	}
	return r.ShardsAcknowledged,err
}