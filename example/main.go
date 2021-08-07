package main

//版本7.x
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/et-zone/ees"
	"github.com/xwb1989/sqlparser"

	// "github.com/et-zone/etest/mock"
	"github.com/olivere/elastic/v7"
)

var host = "http://127.0.0.1:9200"
var table = "stu"

type Result struct {
	ID   int64    `json:"id"`
	Name string   `json:"name"`
	High float64  `json:"high"`
	TNow string   `json:"ntime"`
	List []string `json:"lis"`
	Flag string   `json:"flag"` //kwd
	IsOK bool     `json:"isok"`
	Geo  ees.Geo  `json:"mgeo"`
	IP   string   `json:"ip"`
}

func main() {

	//初始化es Client
	err := ees.InitESClient(elastic.SetURL(host), elastic.SetSniff(false))
	if err != nil {
		fmt.Printf("create es client failed | err : %s\n", err)
		return
	}
	//ctx := context.TODO()
	//***********创建index（table)************
	//6.X 没有string类型了
	// mapping := `{
	//    "mappings": {
	//        "dynamic": "false",
	//        "properties": {
	//            "id": {
	//                "type": "long",
	//                "index":true
	//            },
	//            "name": {

	//                "type": "text"
	//            },
	//            "high": {
	//                "type": "float"
	//            },
	//            "ntime": {
	//                "type": "date",
	//                "format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
	//            },
	//            "lis": {
	//                "type": "text"
	//            },
	//            "isok": {
	//                "type": "boolean"
	//            },
	//            "flag": {
	//                "type": "keyword"
	//            },
	//            "mgeo": {
	//                "type": "geo_point"
	//            },
	//            "ip": {
	//                "type": "ip"
	//            }
	//        }
	//    }
	// }`
	// isok, err := ees.InitTable(ctx, table, mapping)
	// if err != nil {
	// 	fmt.Println("===1===", err.Error())
	// }
	// fmt.Println(isok)
	//初始化实体对象并添加到es

	//***********Upsert One************
	////
	//ctx := context.TODO()
	//nearbyObj := &Result{
	//	ID:   211158,
	//	Name: "bbdd",
	//	High: 12.30,
	//	TNow: "2021-08-05 00:12:20",
	//	List: []string{"aa", "bb", "cc"},
	//	Flag: "ff gg",
	//	IsOK: true,
	//	Geo:  ees.Geo{10.0001,22.003},
	//	IP:   "127.0.0.1",
	//}
	//isok, err := ees.UpsertOneESData(ctx, table, fmt.Sprintf("%v", nearbyObj.ID), nearbyObj)
	//
	//if err != nil {
	//	fmt.Printf("failed | err : %s\n", err)
	//	return
	//}
	//fmt.Println(isok)


	//***********ADD ALL************
	//ctx := context.TODO()
	//datas := map[string]interface{}{}
	//for i := 35; i <= 36; i++ {
	//	nearbyObj := &Result{
	//		ID:   111111 + int64(i),
	//		//Name: mock.Mock.DefaultString(),
	//		Name: "aabba",
	//		High: 10.30 + float64(i),
	//		TNow: mock.Mock.DefaultDateTime(),
	//		List: []string{"aa", "bb", "cc"},
	//		Flag: "ls",
	//		IsOK: false,
	//		Geo:  ees.Geo{10.0001, 22.0003},
	//		IP:   "127.0.0.7",
	//	}
	//	datas[fmt.Sprintf("%v", nearbyObj.ID)] = nearbyObj
	//}
	//
	//c, err := ees.UpsertAllESData(ctx, table, datas)
	//
	//if err != nil {
	//	fmt.Printf("failed | err : %s\n", err)
	//	return
	//}
	//fmt.Println(c)

	//***********Select ALL************
	//sql查询
	//sql:="select * from stu where id in(\"112\",\"114\",\"115\")"
	//sql := "select * from stu limit 100"
	// sql:="select * from stu where icon between 0 and 1000 "
	// sql:="select * from stu where id > 0 and id< 1000 limit 10"

	sql := "select * from stu where name = 'aa bb'  limit 100"
	dat := &[]Result{}
	ctx := context.TODO()
	cout, err := ees.SelectSql(ctx, sql, dat)
	if err != nil {
		fmt.Printf("search es data failed | err : %s\n", err)
		return
	}
	fmt.Println(cout)

	for _, d := range *dat {
		b, _ := json.Marshal(d)
		fmt.Println(string(b))
	}

	////删除es数据
	// delRep, err := ees.DeleteESItemByID(ctx, table, 111111)
	// if err != nil {
	// 	fmt.Printf("delete es data failed | err : %s\n", err)
	// 	return
	// }
	// fmt.Printf("%+v\n", delRep)

	//批量删除
	// delRep, err := ees.DeleteESItemByIDs(ctx, table, []interface{}{111118, 111119})
	// if err != nil {
	// 	fmt.Printf("delete es data failed | err : %s\n", err)
	// 	return
	// }
	// fmt.Printf("%+v\n", delRep)

	//删表
	//delRep, err := ees.DeleteTable(ctx, table)
	//if err != nil {
	//	fmt.Printf("delete es data failed | err : %s\n", err)
	//	return
	//}
	//fmt.Printf("%+v\n", delRep)

	//mp, err := ees.GetTableInfo(ctx, table)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//b, _ := json.Marshal(mp)
	//fmt.Println(string(b))

	//Test_INIT_Table()

	//SqlTest()
}

//结构化创建table
func Test_INIT_Table() {
	ctx := context.TODO()
	mapping := ees.NewMapping()
	mapping.SetDynamic(ees.DYNAMIC.False())
	mapping.SetField("id", ees.NewField().SetType(ees.TYPE.Long()))
	mapping.SetField("name", ees.NewField().SetType(ees.TYPE.Text()))
	mapping.SetField("high", ees.NewField().SetType(ees.TYPE.Float()))
	mapping.SetField("ntime", ees.NewField().SetType(ees.TYPE.Date()).SetFormat(ees.DATE_TIME_FORMAT))
	mapping.SetField("lis", ees.NewField().SetType(ees.TYPE.Text()))
	mapping.SetField("isok", ees.NewField().SetType(ees.TYPE.Boolean()))
	mapping.SetField("flag", ees.NewField().SetType(ees.TYPE.Keyword()))
	mapping.SetField("mgeo", ees.NewField().SetType(ees.TYPE.Geo()))
	mapping.SetField("ip", ees.NewField().SetType(ees.TYPE.IP()))

	b, _ := json.Marshal(mapping.Mappings())
	fmt.Println("mapping===", string(b))
	isok, err := ees.InitTable(ctx, table, mapping.Mappings())
	if err != nil {
		fmt.Println("===1===", err.Error())
	}
	fmt.Println(isok)

	//删表
	// delRep, err := ees.DeleteTable(ctx, table)
	// if err != nil {
	// 	fmt.Printf("delete es data failed | err : %s\n", err)
	// 	return
	// }
	// fmt.Printf("%+v\n", delRep)

	//	查看
	mp, err := ees.GetTableInfo(ctx, table)
	if err != nil {
		fmt.Println(err.Error())
	}
	c, _ := json.Marshal(mp)
	fmt.Println(string(c))

}

/*
== 说明
{
    "mappings": {
        "dynamic": false,
        "properties": {
            "id": {
                "type": "long",
                "index": true //默认是true（表示字段支持搜索匹配）
            },
            "name": {
                "type": "text"//支持分词后的全文检索
                "analyzer": "ik_max_word",    //索引存储阶段和搜索阶段都分词，索引时用ik_max_word，搜索时分词器用ik_smart
                "search_analyzer": "ik_smart",   //搜索阶段分词，会覆盖上面的属性
                "fields": { //通过不同的方法索引相同的字段通常非常有用，支持使用其他类型搜索如text转keyword，
		            "raw": {
		              "type":  "keyword"
		            },
				"index_options": "positions",//text类型支持，默认是positions，记录内容越多，占据空间越大,6.0版本开始弃用
            },
            "high": {
                "type": "float"
            },
            "ntime": {
                "type": "date",
                "format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
            },
            "lis": {
                "type": "text"//支持分词后的全文检索
            },
            "isok": {
                "type": "boolean"
            },
            "flag": {
                "type": "keyword"//支持全文文本搜索（相当于字符串是一个整体，不拆分的单元）
            },
            "mgeo": {
                "type": "geo_point"//地理位置类型
            },
            "ip": {
                "type": "ip"//IP类型
            }
        }
    }
}

https://www.cnblogs.com/haixiang/p/12040272.html

//

dynamic属性默认为true，新增字段时会自动创建mapping
dynamic属性被设置为false时，新增字段不会创建mapping，但是数据会存储，无法根据字段条件查询，但是该字段会会被match_all查询处理
dynamic属性被设置为strict时，数据写入直接出错

*/

func SqlTest(){
	sql:="select * from stu"
	ret,err:=sqlparser.Parse(sql)
	if err!=nil{
		fmt.Println(err.Error())
	}
	re:=sqlparser.GetBindvars(ret)

	fmt.Println(re)

}
