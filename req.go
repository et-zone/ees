package ees

import "github.com/olivere/elastic/v7"

/*
	Automatic execution of analyses
*/

const (
	ActionAnd = 1
	ActionOr  = 2
	ActionNot = 3
)

type QueryCondition struct {
	Index    []string `json:"index"`
	Fields   []*Field `json:"fields"`
	Sort     []Sort   `json:"sort"`
	Page     int      `json:"page"`
	PageSize int      `json:"pageSize"`
	Source   []string `json:"source"` //only return this fields,default all
}

type Field struct {
	Action int32  `json:"action"`
	Name   string `json:"name"`
	Alise  string `json:"alise"`
	Value  Value  `json:"value"`
}

func NewQueryCondition(index ...string) *QueryCondition {
	return &QueryCondition{
		Index: index,
	}
}
func (q *QueryCondition) AddField(field ...*Field) {
	q.Fields = append(q.Fields, field...)
}
func (q *QueryCondition) SetPage(pageNo, pageSize int) {
	q.Page = pageNo
	q.PageSize = pageSize
}

func NewQFieldAnd(name string) *Field {
	return &Field{
		Name:   name,
		Action: ActionAnd,
	}
}

func NewQFieldNot(name string) *Field {
	return &Field{
		Name:   name,
		Action: ActionNot,
	}
}

func NewQFieldOr(name string) *Field {
	return &Field{
		Name:   name,
		Action: ActionOr,
	}
}

type Value struct {
	Vals []interface{}
	Eq   *Eq  `json:"eq,omitempty"` //Vals first one
	In   *In  `json:"in,omitempty"`
	Gt   *Gt  `json:"gt,omitempty"`  //Vals first one
	Gte  *Gte `json:"gte,omitempty"` //Vals first one
	Lt   *Lt  `json:"lt,omitempty"`  //Vals first one
	Lte  *Lte `json:"lte,omitempty"` //Vals first one
}

func (f *Field) SetEq(val interface{}) *Field {
	f.Value.Eq = &Eq{val}

	return f
}

func (f *Field) SetIn(vals []interface{}) *Field {
	f.Value.In = &In{vals}
	return f
}

func (f *Field) SetGt(val interface{}) *Field {
	f.Value.Gt = &Gt{val}
	return f
}
func (f *Field) SetGte(val interface{}) *Field {
	f.Value.Gte = &Gte{val}
	return f
}
func (f *Field) SetLt(val interface{}) *Field {
	f.Value.Lt = &Lt{val}
	return f
}
func (f *Field) SetLte(val interface{}) *Field {
	f.Value.Lte = &Lte{val}
	return f
}

func (f *Field) NoCondition() bool {
	return f.Value.Lte == nil && f.Value.Lt == nil && f.Value.Gt == nil && f.Value.Gte == nil && f.Value.In == nil && f.Value.Eq == nil
}

type Eq struct {
	Value interface{}
}

type In struct {
	Value []interface{}
}

// >
type Gt struct {
	Value interface{}
}

// >=
type Gte struct {
	Value interface{}
}

// <
type Lt struct {
	Value interface{}
}

// <=
type Lte struct {
	Value interface{}
}

type Sort struct {
	Name string `json:"name"`
	Desc bool   `json:"desc"` //default false
}

func BuildEQueryReq(q *QueryCondition) *EQueryReq {
	if q == nil || len(q.Index) == 0 {
		return nil
	}
	req := NewEQueryReq(q.Index...)
	for _, field := range q.Fields {
		if field.NoCondition() || field.Name == "" || field.Action == 0 {
			continue
		}
		switch field.Action {
		case ActionAnd:
			switch {
			case field.Value.Eq != nil:
				req.Query.And(FieldEq(field.Name, field.Value.Eq.Value))
			case field.Value.In != nil:
				req.Query.And(FieldIN(field.Name, field.Value.In.Value...))
			case field.Value.Gt != nil || field.Value.Gte != nil || field.Value.Lt != nil || field.Value.Lte != nil:
				ranges := FieldRange(field.Name)
				if field.Value.Gt != nil && field.Value.Gt.Value != nil {
					ranges.Gt(field.Value.Gt.Value)
				}
				if field.Value.Gte != nil && field.Value.Gte.Value != nil {
					ranges.Gte(field.Value.Gte.Value)
				}
				if field.Value.Lte != nil && field.Value.Lte.Value != nil {
					ranges.Lte(field.Value.Lte.Value)
				}
				if field.Value.Lt != nil && field.Value.Lt.Value != nil {
					ranges.Lt(field.Value.Lt.Value)
				}
				req.Query.And(ranges)
			default:
			}
		case ActionOr:
			switch {
			case field.Value.Eq != nil:
				req.Query.OR(FieldEq(field.Name, field.Value.Eq.Value))
			case field.Value.In != nil:
				req.Query.OR(FieldIN(field.Name, field.Value.In.Value...))
			case field.Value.Gt != nil || field.Value.Gte != nil || field.Value.Lt != nil || field.Value.Lte != nil:
				ranges := FieldRange(field.Name)
				if field.Value.Gt != nil && field.Value.Gt.Value != nil {
					ranges.Gt(field.Value.Gt.Value)
				}
				if field.Value.Gte != nil && field.Value.Gte.Value != nil {
					ranges.Gte(field.Value.Gte.Value)
				}
				if field.Value.Lte != nil && field.Value.Lte.Value != nil {
					ranges.Lte(field.Value.Lte.Value)
				}
				if field.Value.Lt != nil && field.Value.Lt.Value != nil {
					ranges.Lt(field.Value.Lt.Value)
				}
				req.Query.OR(ranges)
			default:
			}
		case ActionNot:
			switch {
			case field.Value.Eq != nil:
				req.Query.Not(FieldEq(field.Name, field.Value.Eq.Value))
			case field.Value.In != nil:
				req.Query.Not(FieldIN(field.Name, field.Value.In.Value...))
			case field.Value.Gt != nil || field.Value.Gte != nil || field.Value.Lt != nil || field.Value.Lte != nil:
				ranges := FieldRange(field.Name)
				if field.Value.Gt != nil && field.Value.Gt.Value != nil {
					ranges.Gt(field.Value.Gt.Value)
				}
				if field.Value.Gte != nil && field.Value.Gte.Value != nil {
					ranges.Gte(field.Value.Gte.Value)
				}
				if field.Value.Lte != nil && field.Value.Lte.Value != nil {
					ranges.Lte(field.Value.Lte.Value)
				}
				if field.Value.Lt != nil && field.Value.Lt.Value != nil {
					ranges.Lt(field.Value.Lt.Value)
				}
				req.Query.Not(ranges)
			default:
			}
		}

	}
	for _, sort := range q.Sort {
		if sort.Name != "" {
			if sort.Desc {
				req.SetSort(elastic.NewFieldSort(sort.Name).Desc())
			} else {
				req.SetSort(elastic.NewFieldSort(sort.Name))
			}
		}
	}
	req.SetPage(q.Page, q.PageSize)
	req.SetIncludesSource(q.Source...)
	return req
}
