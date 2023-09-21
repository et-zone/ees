package main

//版本7.x
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/et-zone/ees"
	"github.com/et-zone/ees/elasticsql"
	"github.com/olivere/elastic/v7"
	"time"
)

var host = ""
var table = "gzy_test"

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

func Init() {
	//初始化es Client  -- SetSniff 集群使用
	err := ees.InitESClient(elastic.SetURL(host), elastic.SetSniff(false), elastic.SetBasicAuth("", ""))
	if err != nil {
		fmt.Printf("create es client failed | err : %s\n", err)
		panic(err)
		return
	}
}
func main() {
	//Test_INIT_Table()
	//insert()
	//show()
	//delTable()
	//insertAll()
	//query()
	//queryKeyword()
	//update()
	//SqlTest()
}

func delTable() {
	//删表
	ctx := context.TODO()
	delRep, err := ees.DelTable(ctx, table)
	if err != nil {
		fmt.Printf("delete es data failed | err : %s\n", err)
		return
	}
	fmt.Printf("%+v\n", delRep)
}

func insertAll() {
	names := []string{
		"zhang san",
		"li si",
		"wang wu",
		"li ming",
		"liu xiang",
		"liu qiang dong",
		"zhou xun",
		"yang mi",
		"sun qi",
		"guo"}
	lists := [][]string{
		[]string{"aa"},
		[]string{"aa", "bb"},
		[]string{"cc"},
		[]string{"bb"},
		[]string{"aa", "cc"},
		[]string{"cc"},
		[]string{"dd"},
		[]string{"ee"},
		[]string{"aa"},
		[]string{}}
	ctx := context.TODO()

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
		_, err := ees.UpsertOneESData(ctx, table, fmt.Sprintf("%v", nearbyObj.ID), nearbyObj)

		if err != nil {
			fmt.Printf("failed | err : %s\n", err)
		}
		if i == 3 || i == 7 {
			time.Sleep(time.Second)
		}
	}

}
func insert() {
	ctx := context.TODO()
	nearbyObj := &Result{ //
		ID:    211158,
		High:  12.30,
		NameT: "aa",
		TNow:  "2021-08-05 00:12:20",
		List:  []string{"bb", "cc"},
	}

	isok, err := ees.UpsertOneESData(ctx, table, fmt.Sprintf("%v", nearbyObj.ID), nearbyObj)

	if err != nil {
		fmt.Printf("failed | err : %s\n", err)
		return
	}
	fmt.Println(isok)
}

func update() {
	//update one field

	//func 1
	ctx := context.TODO()
	//nearbyObj := &Result{ //
	//	High: 19.30,
	//}
	//isok, err := ees.UpsertOneESData(ctx, table, fmt.Sprintf("%v", 211158), nearbyObj)

	//func 2
	isok, err := ees.UpsertOneESData(ctx, table, fmt.Sprintf("%v", 211158), map[string]interface{}{
		"high": 22,
	})

	if err != nil {
		fmt.Printf("failed | err : %s\n", err)
		return
	}
	fmt.Println(isok)

}

func show() {
	//	查看
	ctx := context.TODO()
	mp, err := ees.GetTableDetail(ctx, table)
	if err != nil {
		fmt.Println(err.Error())
	}
	c, _ := json.Marshal(mp)
	fmt.Println(string(c))
}

func query() {
	//***********ADD ALL************
	ctx := context.TODO()

	//***********Select ALL************
	//sql查询
	// keyword
	//sql := "select * from shop where name_keyword='guo' limit 10"
	//sql := "select * from shop where name_text_default='liu' limit 10"//自动分词
	//sql := "select * from shop where name_text_analyzer='dong' limit 10"//分词
	//sql := "select * from shop where name_text_s_analyzer='li' limit 10"//分词

	// 原始sql不支持array类型的dsl，需要自己写
	//sql := "select * from shop where name_text_default = 'li' limit 10"//分词 数组分词了,全文搜索配置fields
	sql := "select id from test_gzy " //分词 数组分词了,全文搜索配置fields

	//时间可以判断大小 ok
	//sql := "select * from shop where ntime >= '2022-08-09 17:25:15' limit 10"//分词

	// 数字大小判断 ok
	//sql := "select * from shop where high > 13 limit 10"//分词

	dat := &[]Result{}
	total, err := ees.QuerySql(ctx, sql, dat)
	if err != nil {
		fmt.Printf("search es data failed | err : %s\n", err)
		return
	}
	return
	fmt.Println("total =", total)

	//for _, d := range *dat {
	//	b, _ := json.Marshal(d)
	//	fmt.Println(string(b))
	//}
	b, _ := json.Marshal(dat)
	fmt.Println(string(b))

	////删除es数据
	// delRep, err := ees.DelESItemByID(ctx, table, 111111)
	// if err != nil {
	// 	fmt.Printf("delete es data failed | err : %s\n", err)
	// 	return
	// }
	// fmt.Printf("%+v\n", delRep)

	//批量删除
	// delRep, err := ees.DelESItemByIDs(ctx, table, []interface{}{111118, 111119})
	// if err != nil {
	// 	fmt.Printf("delete es data failed | err : %s\n", err)
	// 	return
	// }
	// fmt.Printf("%+v\n", delRep)

	//mp, err := ees.GetTableDetail(ctx, table)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//b, _ := json.Marshal(mp)
	//fmt.Println(string(b))

	//Test_INIT_Table()

	//SqlTest()
}

func queryKeyword() { //query more field like 全文检索
	ctx := context.TODO()
	r, err := ees.Client().Search("shop").Query(elastic.NewQueryStringQuery("aa").Field("name_text_default").Field("lis")).Do(ctx)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		b, _ := json.Marshal(r.Hits.Hits)
		fmt.Println(string(b))
	}
}

// 结构化创建table
func Test_INIT_Table() {
	ctx := context.TODO()
	mapping := ees.NewMapping()
	mapping.SetDynamic(ees.Dynamic.False())
	mapping.SetField("id", ees.NewField().SetType(ees.Type.Long()))
	mapping.SetField("name", ees.NewField().SetType(ees.Type.Text()).SetSearchAnalyzer(ees.IkMaxWord()))
	mapping.SetField("high", ees.NewField().SetType(ees.Type.Float()))
	mapping.SetField("ntime", ees.NewField().SetType(ees.Type.Date()).SetFormat(ees.DateTimeFormat))
	mapping.SetField("obj", ees.NewField().SetType(ees.Type.Object()))

	b, _ := json.Marshal(mapping.Mappings())
	fmt.Println("mapping===", string(b))
	return
	isok, err := ees.InitTable(ctx, table, mapping.Mappings())
	if err != nil {
		fmt.Println("===1===", err.Error())
		return
	}
	fmt.Println(isok)

	//删表
	// delRep, err := ees.DelTable(ctx, table)
	// if err != nil {
	// 	fmt.Printf("delete es data failed | err : %s\n", err)
	// 	return
	// }
	// fmt.Printf("%+v\n", delRep)

	//	查看
	mp, err := ees.GetTableDetail(ctx, table)
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

func SqlTest() {
	sql := "select * from stu where id=1 and name='fff' and age=22 or sex=1 order by id desc limit 20"

	dsl, _, _ := elasticsql.Convert(sql)
	fmt.Println(dsl)
}
