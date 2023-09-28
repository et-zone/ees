package ees

type Fields struct {
	Type  string `json:"type"`
	Index bool   `json:"index,omitempty"` //默认选true，支持搜索,旧版本使用:"not_analyzed"
}

type Field struct {
	Type           string             `json:"type"`
	Index          *bool              `json:"index,omitempty"`           //default true，支持搜索,旧版本使用:"not_analyzed"
	Store          *bool              `json:"store,omitempty"`           //default false ,need not set
	Format         string             `json:"format,omitempty"`          //时间类型格式化
	Analyzer       string             `json:"analyzer,omitempty"`        //写入时就进行分词，最大拆分=ik_max_word,粗略拆分=ik_smart  //text使用
	SearchAnalyzer string             `json:"search_analyzer,omitempty"` //搜索阶段进行分词，会覆盖上面的属性 ，"search_analyzer": "ik_smart"  //text使用
	Fields         *map[string]Fields `json:"fields,omitempty"`          //复合 字段，用于text 的全文搜索,如 type = keyword
	DocValues      *bool              `json:"doc_values,omitempty"`      //default true ,need not set
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

// ---- ---- ---- ---- ---- ---- ----

type _index int

func (index _index) Enabled() bool {
	return true
}
func (index _index) Disabled() bool {
	return false
}

type _type int

//number  64 byte int
func (f _type) Long() string {
	return "long"
}

//number 32 byte float
func (f _type) Float() string {
	return "float"
}

//number 32 byte int
func (f _type) Int() string {
	return "integer"
}

//number 16 byte int
func (f _type) Short() string {
	return "short"
}

//number 1byte
func (f _type) Byte() string {
	return "byte"
}

//number  64 byte ,float
func (f _type) Double() string {
	return "double"
}

//text type
func (f _type) Text() string {
	//类型自动会分词，支持模糊搜索，使用=匹配即可
	return "text"
}

//text type
func (f _type) TextMatchOnlyText() string {
	return "match_only_text"
}

// Keyword 搜索完全匹配
func (f _type) Keyword() string {
	return "keyword"
}

//Keyword
func (f _type) KeywordConstantKeyword() string {
	return "constant_keyword"
}

//Keyword
func (f _type) KeywordWildcard() string {
	return "wildcard"
}

//date type
func (f _type) Date() string {
	return "date"
}

func (f _type) Boolean() string {
	return "boolean"
}

func (f _type) Geo() string {
	return "geo_point"
}

func (f _type) IP() string {
	return "ip"
}

type _dynamic string

func (d _dynamic) True() string {
	return "true"
}

func (d _dynamic) False() string {
	return "false"
}

func (d _dynamic) Strict() string {
	return "strict"
}

//object type json string
func (f _type) Object() string {
	//类型自动会分词，支持模糊搜索，使用=匹配即可
	return "object"
}

// -------------------- map ----------------------

//版本7.x
const (
	DefaultNum = 0
	// type
	Type = _type(DefaultNum)
	// dynamic type
	Dynamic = _dynamic(DefaultNum)
	//index type
	Index = _index(DefaultNum)

	// yyyy-MM-dd HH:mm:ss
	DateTimeFormat = "yyyy-MM-dd HH:mm:ss"
	// yyyy-MM-dd
	DateFormat = "yyyy-MM-dd"
	//HH:mm:ss
	TimeFormat = "HH:mm:ss"
	row        = "row"
	BoolTrue   = true
	BoolFalse  = false
)

type Mapping struct {
	Dynamic string           `json:"dynamic,omitempty"` //default true
	Fields  map[string]Field `json:"properties"`
}

func NewMapping() *Mapping {
	return &Mapping{
		Fields: map[string]Field{},
	}
}

func (m *Mapping) SetDynamic(dynamic string) *Mapping {
	m.Dynamic = dynamic
	return m
}

func (m *Mapping) SetField(fielName string, field *Field) *Mapping {
	if fielName == "" {
		return m
	}
	m.Fields[fielName] = *field
	return m
}

// Mappings return mappings interface{}
func (m *Mapping) Mappings() interface{} {
	return map[string]Mapping{
		"mappings": *m,
	}
}

type IndexOptions int

func (o IndexOptions) Docs() string {
	return "docs"
}

func (o IndexOptions) Freqs() string {
	return "freqs"
}

func (o IndexOptions) Positions() string {
	return "positions"
}

func (o IndexOptions) Offsets() string {
	return "offsets"
}

// IkMaxWord  analyzer ,search_analyzer can use
func IkMaxWord() string {
	return "ik_max_word"
}

// IkSmart  analyzer ,search_analyzer can use
func IkSmart() string {
	return "ik_smart"
}

type Geo struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

/*
index_options对于6.0.0 中的数字字段，该参数已被弃用。

该index_options参数控制将哪些信息添加到倒排索引中，以用于搜索和突出显示目的。
它接受以下设置：
	docs  只有文档编号被索引。可以回答这个问题这个词是否存在于这个领域？

	freqs   文档编号和术语频率已编入索引。术语频率用于对重复术语的评分高于单个术语。

	positions   文档编号、术语频率和术语位置（或顺序）被编入索引。位置可用于 邻近或短语查询。

	offsets    文档编号、术语频率、位置以及开始和结束字符偏移（将术语映射回原始字符串）都被编入索引。统一荧光笔使用偏移来加速突出显示。
*/
