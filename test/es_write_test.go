package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/et-zone/ees/es"
	"github.com/olivere/elastic/v7"
	"strconv"
	"testing"
	"time"
)

var host = "https://dev-elastic.carsome.dev"
var table = "lo_test"
var uname = "elastic"
var pwd = "s6G74t5hww202p1QETAH20ab"

type Result struct {
	ID    int64  `json:"id,omitempty"`
	NameK string `json:"name_keyword,omitempty"`      //Keyword
	NameT string `json:"name_text_default,omitempty"` //text
	//NameA string   `json:"name_text_analyzer,omitempty"`//text Analyzer
	//NameS string   `json:"name_text_s_analyzer,omitempty"`//text sarrch Analyzer
	High float64  `json:"high,omitempty"`
	TNow string   `json:"ntime,omitempty"`
	List []string `json:"lis,omitempty"` //
}

func TestInsertAll(t *testing.T) {
	names := []string{"zhang san", "li si", "wang wu", "li ming", "liu xiang", "liu qiang dong", "zhou xun", "yang mi", "sun qi", "guo"}
	lists := [][]string{[]string{"aa"}, []string{"aa", "bb"}, []string{"cc"}, []string{"bb"}, []string{"aa", "cc"}, []string{"cc"}, []string{"dd"}, []string{"ee"}, []string{"aa"}, []string{}}
	ctx := context.TODO()

	data := map[string]interface{}{}
	for i := 0; i < 10; i++ {
		nearbyObj := &Result{ //
			ID: 211158 + int64(i),
			//NameA: names[i],
			NameK: names[i],
			//NameS: names[i],
			NameT: names[i],
			High:  12.30,
			TNow:  time.Now().Format("2006-01-02 15:04:05"),
			List:  lists[i],
		}
		data[fmt.Sprintf("%v", nearbyObj.ID)] = nearbyObj
	}
	_, err := es.UpsertAll(ctx, table, data)
	fmt.Println(err)
}

func TestInsertOne(t *testing.T) {
	ctx := context.TODO()
	nearbyObj := &Result{ //
		ID:    111158,
		High:  12.30,
		NameT: "aa",
		TNow:  "2021-08-05 00:12:20",
		List:  []string{"bb", "cc"},
	}

	isok, err := es.UpsertOne(ctx, table, fmt.Sprintf("%v", nearbyObj.ID), nearbyObj)

	if err != nil {
		fmt.Printf("failed | err : %s\n", err)
		return
	}
	fmt.Println(isok)
}

func TestDelOneByID(t *testing.T) {
	//删表
	ctx := context.TODO()
	delRep, err := es.DelItemByID(ctx, table, strconv.Itoa(211158))
	if err != nil {
		fmt.Printf("delete es data failed | err : %s\n", err)
		return
	}
	fmt.Println(err)
	fmt.Printf("%+v\n", delRep)
}

func TestDelOneByIDs(t *testing.T) {
	//删表
	ctx := context.TODO()
	delRep, err := es.DelItemByIDs(ctx, table, strconv.Itoa(211159), strconv.Itoa(211160))
	if err != nil {
		fmt.Printf("delete es data failed | err : %s\n", err)
		return
	}
	fmt.Println(err)
	fmt.Printf("%+v\n", delRep)
}

func TestDelTable(t *testing.T) {
	//删表
	ctx := context.TODO()
	delRep, err := es.DelTable(ctx, table)
	if err != nil {
		fmt.Printf("delete es data failed | err : %s\n", err)
		return
	}
	fmt.Println(err)
	fmt.Printf("%+v\n", delRep)
}

func TestUpdate(t *testing.T) {
	ctx := context.TODO()
	isok, err := es.UpsertOne(ctx, table, fmt.Sprintf("%v", 211158), map[string]interface{}{
		"high": 22,
	})
	fmt.Println(err)
	fmt.Println(isok)

}

// 结构化创建table 支持 7.*版本
func TestInitTable(t *testing.T) {
	//支持7.*
	ctx := context.TODO()
	mapping := es.NewMapping()
	mapping.SetDynamic(es.Dynamic.False())
	mapping.SetField("id", es.NewField().SetType(es.Type.Long()))
	mapping.SetField("name", es.NewField().SetType(es.Type.Text()).SetSearchAnalyzer(es.IkMaxWord()))
	mapping.SetField("high", es.NewField().SetType(es.Type.Float()))
	mapping.SetField("ntime", es.NewField().SetType(es.Type.Date()).SetFormat(es.DateTimeFormat))
	mapping.SetField("obj", es.NewField().SetType(es.Type.Object()))

	b, _ := json.Marshal(mapping.Mappings())
	fmt.Println("mapping===", string(b))

	isok, err := es.InitTable(ctx, table, mapping.Mappings())
	fmt.Println(err)
	fmt.Println(isok)
	if err != nil {
		return
	}

	//	查看
	mp, err := es.GetTableDetail(ctx, table)

	fmt.Println(err)
	c, _ := json.Marshal(mp)
	fmt.Println(string(c))

}

func TestMain(m *testing.M) {
	//初始化es Client  -- SetSniff 集群使用
	err := es.InitESClient(elastic.SetURL(host), elastic.SetSniff(false), elastic.SetBasicAuth(uname, pwd))
	if err != nil {
		fmt.Printf("create es client failed | err : %s\n", err)
		panic(err)
		return
	}
	m.Run()
}
