package ees

type Fields struct {
	Type  string `json:"type"`
	Index bool   `json:"index,omitempty"` //默认选true，支持搜索,旧版本使用:"not_analyzed"
}

type Field struct {
	Type  string `json:"type"`
	Index *bool  `json:"index,omitempty"` //default true，支持搜索,旧版本使用:"not_analyzed"
	Store *bool  `json:"store,omitempty"` //default false ,need not set
	Format         string             `json:"format,omitempty"`          //时间类型格式化
	Analyzer       string             `json:"analyzer,omitempty"`        //写入时就进行分词，最大拆分=ik_max_word,粗略拆分=ik_smart  //text使用
	SearchAnalyzer string             `json:"search_analyzer,omitempty"` //搜索阶段进行分词，会覆盖上面的属性 ，"search_analyzer": "ik_smart"  //text使用
	Fields         *map[string]Fields `json:"fields,omitempty"`          //复合 字段，用于text 的全文搜索,如 type = keyword
	DocValues *bool              `json:"doc_values,omitempty"` //default true ,need not set
	// IgnoreAbove int64  `json:"ignore_above,omitempty"` //对超过 ignore_above 的字符串，analyzer 不会进行处理
	// NullValue   string `json:"null_value,omitempty"` //支持字段为null，只有keyword类型支持，自定义Mapping常用参数，实际操作时不存储null数据
	// Store bool `json:"store,omitempty"` //(用于单独存储该field的原始值)默认情况下已存储
	// TextData bool `json:"fielddata,omitempty"` //预加载，Textdata会占用大量堆空间，尤其是在加载大量的文本字段时，默认禁用

}

//default type=text
func NewField() *Field {
	return &Field{
		Type: Type.Text(),
	}
}

//must do it.
func (f *Field) SetType(t string) *Field {
	f.Type = t
	return f
}

//type = keyword,text ，选填
func (f *Field) SetFields(t string) *Field {

	f.Fields = &map[string]Fields{
		row: {Type: t},
	}
	return f
}

//default true. type = text,keyward,boolean,date,geo_point,number
//有默认值，可以不处理
func (f *Field) CanIndex(can bool) *Field {
	f.Index = &can
	return f
}

//default false. type = text,keyword,binary,boolean,date,number
func (f *Field) CanStore(can bool) *Field {
	f.Store = &can
	return f
}

// type = keyword,boolean,binary,date,number
func (f *Field) CanDocValues(can bool) *Field {
	f.DocValues = &can
	return f
}

//  type = data ，一般 会主动选择设置此参数
func (f *Field) SetFormat(fmat string) *Field {
	if f.Type != Type.Date() {
		panic("set format date error")
	}
	f.Format = fmat
	return f
}

//type = text ,一般 会主动选择设置此参数
func (f *Field) SetAnalyzer(a string) *Field {
	f.Analyzer = a
	return f
}

//type = text ，一般 会主动选择设置此参数
func (f *Field) SetSearchAnalyzer(s string) *Field {
	f.SearchAnalyzer = s
	return f
}

