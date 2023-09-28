package main

//版本7.x
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/et-zone/ees/elasticsql"

	"testing"
)

func TestShowTableDetail(t *testing.T) {
	//	查看
	ctx := context.TODO()
	mp, err := es.GetIndexDetail(ctx, table)
	if err != nil {
		fmt.Println(err.Error())
	}
	c, _ := json.Marshal(mp)
	fmt.Println(string(c))
}

func TestShowCount(t *testing.T) {
	//	查看
	ctx := context.TODO()
	mp, err := es.Count(ctx, table)
	fmt.Println(err)
	fmt.Println(mp)

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

func TestSql(t *testing.T) {
	sql := "select * from stu where id=1 and name='fff' and age=22 or sex=1 order by id desc limit 20"

	dsl, _, _ := elasticsql.Convert(sql)
	fmt.Println(dsl)
}
func TestQuerySql(t *testing.T) {
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
	total, err := es.QuerySql(ctx, sql, dat)
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

}
