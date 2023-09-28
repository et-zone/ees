package es

//版本7.x
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

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
func InitTable(ctx context.Context, index string, mappings interface{}) (bool, error) {
	//_type default = "_doc"

	isExit, _ := client.IndexExists(index).Do(ctx)
	if isExit {
		return true, nil
	}

	rep, err := client.CreateIndex(index).BodyJson(mappings).Do(ctx)
	if err != nil {
		return false, err
	}
	return rep.ShardsAcknowledged, err
}

func UpsertOne(ctx context.Context, index string, id string, value interface{}) (isok bool, err error) {
	doc := elastic.NewBulkUpdateRequest().Index(index).Id(id).Doc(value).DocAsUpsert(true)
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
func UpsertAll(ctx context.Context, index string, values map[string]interface{}) (int64, error) {
	docs := []elastic.BulkableRequest{}
	for k, v := range values {
		doc := elastic.NewBulkUpdateRequest().Index(index).Id(k).Doc(v).DocAsUpsert(true)
		docs = append(docs, doc)
	}
	bulk := client.Bulk()
	rep, err := bulk.Add(docs...).Do(ctx)

	if err != nil {
		return 0, err
	}

	if rep != nil && rep.Errors {
		e, _ := json.Marshal(rep)
		return 0, errors.New(string(e))
	}
	return int64(len(rep.Items)), err
}

// get data
func GetItemByID(index string, id interface{}) (rep *elastic.GetResult, err error) {
	rep, err = client.Get().
		Index(index).
		Id(fmt.Sprintf("%v", id)).
		Do(context.Background())
	return
}

func Client() *elastic.Client {
	return client
}

//删除数据
func DelItemByID(ctx context.Context, index string, id string) (isDel bool, err error) {
	if index == "" || id == "" {
		return false, errors.New("table or id not null")
	}
	rep, err := client.Delete().Refresh("wait_for").
		Index(index).
		Id(fmt.Sprintf("%v", id)).
		Do(ctx)

	if rep != nil && rep.Result == DeleteSucc {
		return true, err
	}
	return false, nil
}

//ids can int or string
func DelItemByIDs(ctx context.Context, index string, ids ...interface{}) (count int64, err error) {
	if index == "" {
		return 0, errors.New("table not null")
	}
	if len(ids) == 0 {
		return 0, errors.New("ids not null")
	}
	t := time.Now()
	q := elastic.NewTermsQuery("_id", ids...)
	rep, err := client.DeleteByQuery(index).Refresh("wait_for").Query(q).Do(ctx)
	fmt.Println("======", time.Since(t))
	if err != nil {
		return 0, err
	}
	return rep.Deleted, err

}

func DelTable(ctx context.Context, index string) (isok bool, err error) {
	if index == "" {
		return false, errors.New("table not null")
	}
	rep, err := client.DeleteIndex(index).Do(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			return true, nil
		}
		return false, err
	}
	return rep.Acknowledged, err

}

func SetIndexAlias(ctx context.Context, index, AliasName string) (bool, error) {
	r, err := client.Alias().Add(index, AliasName).Do(ctx)
	if err != nil {
		return false, err
	}
	return r.ShardsAcknowledged, err
}

func DelIndexAlias(ctx context.Context, index, AliasName string) (bool, error) {
	r, err := client.Alias().Remove(index, AliasName).Do(ctx)
	if err != nil {
		return false, err
	}
	return r.ShardsAcknowledged, err
}
