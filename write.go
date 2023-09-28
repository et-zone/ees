package ees

//版本7.x
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"strings"
)

const (
	CreateSucc   = "created"
	UpdateSucc   = "updated"
	DeleteSucc   = "deleted"
	Noop         = "noop"
	DefaultLimit = "1000"
)

// init client
func NewClient(opts ...elastic.ClientOptionFunc) (client *Eelastic, err error) {
	//client, err = elastic.NewClient(
	//	elastic.SetURL(host),elastic.SetSniff(false),
	//)
	opts = append(opts, elastic.SetTraceLog(new(tracelog)))
	cli, err := elastic.NewClient(opts...)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't connect to elasticsearch | err : %s \n", err))
	}

	log.Println("connect to elasticsearch success")

	return &Eelastic{client: cli}, nil
}

type tracelog struct{}

//实现输出
func (tracelog) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

//Here we think table = index, all _type = default("_doc")
func (e *Eelastic) InitTable(ctx context.Context, index string, mappings interface{}) (bool, error) {
	//_type default = "_doc"

	isExit, _ := e.client.IndexExists(index).Do(ctx)
	if isExit {
		return true, nil
	}

	rep, err := e.client.CreateIndex(index).BodyJson(mappings).Do(ctx)
	if err != nil {
		return false, err
	}
	return rep.ShardsAcknowledged, err
}

func (e *Eelastic) UpsertOne(ctx context.Context, index string, id string, value interface{}) (isok bool, err error) {
	doc := elastic.NewBulkUpdateRequest().Index(index).Id(id).Doc(value).DocAsUpsert(true)
	bulk := e.client.Bulk()
	rep, err := bulk.Add(doc).Do(ctx)

	res := rep.Items[0]
	if rep != nil && rep.Errors {
		e, _ := json.Marshal(res["update"].Error)
		return false, errors.New(string(e))
	}
	return true, err
}

// values key is id
func (e *Eelastic) UpsertAll(ctx context.Context, index string, values map[string]interface{}) (int64, error) {
	docs := []elastic.BulkableRequest{}
	for k, v := range values {
		doc := elastic.NewBulkUpdateRequest().Index(index).Id(k).Doc(v).DocAsUpsert(true)
		docs = append(docs, doc)
	}
	bulk := e.client.Bulk()
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
func (e *Eelastic) GetItemByID(index string, id interface{}) (rep *elastic.GetResult, err error) {
	rep, err = e.client.Get().
		Index(index).
		Id(fmt.Sprintf("%v", id)).
		Do(context.Background())
	return
}

func (e *Eelastic) Client() *elastic.Client {
	return e.client
}

//删除数据
func (e *Eelastic) DelItemByID(ctx context.Context, index string, id string) (isDel bool, err error) {
	if index == "" || id == "" {
		return false, errors.New("table or id not null")
	}
	rep, err := e.client.Delete().
		Index(index).
		Id(fmt.Sprintf("%v", id)).
		//Refresh("wait_for").
		Do(ctx)

	if rep != nil && rep.Result == DeleteSucc {
		return true, err
	}
	return false, nil
}

//ids can int or string
func (e *Eelastic) DelItemByIDs(ctx context.Context, index string, ids ...interface{}) (count int64, err error) {
	if index == "" {
		return 0, errors.New("table not null")
	}
	if len(ids) == 0 {
		return 0, errors.New("ids not null")
	}
	//t := time.Now()
	q := elastic.NewTermsQuery("_id", ids...)
	rep, err := e.client.DeleteByQuery(index).Query(q).Do(ctx)
	//fmt.Println("======", time.Since(t))
	if err != nil {
		return 0, err
	}
	return rep.Deleted, err

}

func (e *Eelastic) DelTable(ctx context.Context, index string) (isok bool, err error) {
	if index == "" {
		return false, errors.New("table not null")
	}
	rep, err := e.client.DeleteIndex(index).Do(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			return true, nil
		}
		return false, err
	}
	return rep.Acknowledged, err

}

func (e *Eelastic) SetIndexAlias(ctx context.Context, index, AliasName string) (bool, error) {
	r, err := e.client.Alias().Add(index, AliasName).Do(ctx)
	if err != nil {
		return false, err
	}
	return r.ShardsAcknowledged, err
}

func (e *Eelastic) DelIndexAlias(ctx context.Context, index, AliasName string) (bool, error) {
	r, err := e.client.Alias().Remove(index, AliasName).Do(ctx)
	if err != nil {
		return false, err
	}
	return r.ShardsAcknowledged, err
}
