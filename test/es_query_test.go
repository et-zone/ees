package main

import (
	"context"
	"fmt"
	"github.com/et-zone/ees"
	"github.com/olivere/elastic/v7"
	"testing"
)

func TestSort(t *testing.T) {
	//	查看
	ctx := context.TODO()
	req := ees.NewEQueryReq(table)
	sort := elastic.NewFieldSort("CreateTime")
	req.SetSort(sort).SetPage(1, 10)
	result, err := es.Search(ctx, req)
	fmt.Println(err)
	fmt.Println(result)
}

func TestInclude(t *testing.T) {
	//	查看
	ctx := context.TODO()
	req := ees.NewEQueryReq(table)
	req.SetPage(1, 10).SetIncludesSource("carId", "carNo")
	result, err := es.Search(ctx, req)
	fmt.Println(err)
	fmt.Println(result)
}

func TestCollapse(t *testing.T) {
	//	查看
	ctx := context.TODO()
	req := ees.NewEQueryReq(table)
	//req.SetPage(1, 100).SetCollapseField("brandID").SetIncludesSource("brandName", "brandID")
	//req.SetPage(1, 100).SetCollapseField("brandName.keyword").SetIncludesSource("brandName", "brandID")
	req.SetPage(1, 100).SetCollapseField("brandName.keyword")
	result, err := es.Search(ctx, req)
	fmt.Println(err)
	fmt.Println(result)
}

func TestQueryWithAggregation(t *testing.T) {
	//	查看
	ctx := context.TODO()
	req := ees.NewEQueryReq(table)
	//req.SetPage(1, 10).SetAggregation("brandName", "brandName.keyword")
	req.SetPage(1, 10).SetAggregation("brandID", "brandID")
	result, err := es.Search(ctx, req)
	fmt.Println(err)
	fmt.Println(result)
}

func TestAggregation(t *testing.T) {
	//	查看
	ctx := context.TODO()
	req := ees.NewEQueryReq(table)
	//req.SetPage(1, 10).SetAggregation("brandName", "brandName.keyword")
	req.SetAggregation("brandID", "brandID")
	req.SetAggregationOpt("brandIDMax", elastic.NewMaxAggregation().Field("brandID"))

	result, err := es.QueryAggregations(ctx, req)
	fmt.Println(err)
	fmt.Println(result)
}

func TestQueryFieldAnd(t *testing.T) {
	//	查看
	ctx := context.TODO()
	req := ees.NewEQueryReq(table).SetPage(1, 10)
	//req.Query.And(ees.FieldIN("brandName.keyword", "Alfa Romeo"))
	//req.Query.And(ees.FieldEq("brandName.keyword", "Alfa Romeo"))
	//req.Query.And(ees.FieldLike("brandName", "Alfa"))
	//req.Query.And(ees.FieldAnalyzer("brandName", "Alfa"))
	//req.Query.And(ees.FieldNotNull("brandName"))
	//req.Query.And(ees.FieldStartWith("brandName.keyword", "Alf"))
	//req.Query.And(ees.FieldTopID(req.Query, "23293"))
	req.Query.And(ees.FieldRange("brandID").Gte(22655))
	result, err := es.Search(ctx, req)
	fmt.Println(err)
	fmt.Println(result)
}

func TestQueryField(t *testing.T) {
	//	查看
	ctx := context.TODO()
	req := ees.NewEQueryReq(table).SetPage(1, 10)
	req.Query.OR(ees.FieldIN("brandName.keyword", "Alfa Romeo"))
	//req.Query.OR(ees.FieldEq("brandName.keyword", "Alfa Romeo"))
	//req.Query.OR(ees.FieldLike("brandName", "Alfa"))
	//req.Query.OR(ees.FieldAnalyzer("brandName", "Alfa"))
	//req.Query.OR(ees.FieldNotNull("brandName"))
	//req.Query.OR(ees.FieldStartWith("brandName.keyword", "Alf"))
	//req.Query.OR(ees.FieldTopID(req.Query, "23293"))
	//req.Query.OR(ees.FieldRange("brandID").Gte(22655))
	result, err := es.Search(ctx, req)
	fmt.Println(err)
	fmt.Println(result)
}

func TestMultiSearch(t *testing.T) {
	//	查看
	ctx := context.TODO()
	req1 := ees.NewEQueryReq(table).SetPage(1, 10)
	req1.Query.And(ees.FieldIN("brandName.keyword", "Alfa Romeo"))

	req2 := ees.NewEQueryReq(table).SetPage(1, 5)
	req2.Query.And(ees.FieldIN("modelID", "8840"))
	req := ees.MultiQueryReq{
		Index: table,
		Req:   []*ees.EQueryReq{req1, req2},
	}
	result, err := es.MultiSearch(ctx, req)
	fmt.Println(err)
	fmt.Println(result)
}
