package ees

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
)

type Mapping struct {
	Dynamic string            `json:"dynamic,omitempty"`
	Fields  map[string]*Field `json:"properties"`
}

func NewMapping() *Mapping {
	return &Mapping{
		Fields: map[string]*Field{},
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
	m.Fields[fielName] = field
	return m
}

// Mappings return mappings interface{}
func (m *Mapping) Mappings() interface{} {
	return map[string]Mapping{
		"mappings": *m,
	}
}

type _index int

func (index _index) Enabled() bool {
	return true
}
func (index _index) Disabled() bool {
	return false
}

type _type int

func (f _type) Long() string {
	return "long"
}

func (f _type) Float() string {
	return "float"
}

//text类型自动会分词，支持模糊搜索，使用=匹配即可
func (f _type) Text() string {
	return "text"
}

//搜索完全匹配
func (f _type) Keyword() string {
	return "keyword"
}

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

/*
index_options对于6.0.0 中的数字字段，该参数已被弃用。

该index_options参数控制将哪些信息添加到倒排索引中，以用于搜索和突出显示目的。
它接受以下设置：
	docs  只有文档编号被索引。可以回答这个问题这个词是否存在于这个领域？

	freqs   文档编号和术语频率已编入索引。术语频率用于对重复术语的评分高于单个术语。

	positions   文档编号、术语频率和术语位置（或顺序）被编入索引。位置可用于 邻近或短语查询。

	offsets    文档编号、术语频率、位置以及开始和结束字符偏移（将术语映射回原始字符串）都被编入索引。统一荧光笔使用偏移来加速突出显示。
*/

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

type Fields struct {
	Type  string `json:"type"`
	Index bool   `json:"index,omitempty"` //默认选true，支持搜索,旧版本使用:"not_analyzed"
}

type Field struct {
	Type           string             `json:"type"`
	Index          bool               `json:"index,omitempty"`           //默认选true，支持搜索,旧版本使用:"not_analyzed"
	Format         string             `json:"format,omitempty"`          //时间类型格式化
	Analyzer       string             `json:"analyzer,omitempty"`        //写入时就进行分词，最大拆分=ik_max_word,粗略拆分=ik_smart  //text使用
	SearchAnalyzer string             `json:"search_analyzer,omitempty"` //搜索阶段进行分词，会覆盖上面的属性 ，"search_analyzer": "ik_smart"  //text使用
	Fields         *map[string]Fields `json:"fields,omitempty"`          //复合 字段，用于text 的全文搜索,如 type = keyword
	// IgnoreAbove int64  `json:"ignore_above,omitempty"` //对超过 ignore_above 的字符串，analyzer 不会进行处理
	// NullValue   string `json:"null_value,omitempty"` //支持字段为null，只有keyword类型支持，自定义Mapping常用参数，实际操作时不存储null数据
	// Store bool `json:"store,omitempty"` //(用于单独存储该field的原始值)默认情况下已存储
	// FieldData bool `json:"fielddata,omitempty"` //预加载，Fielddata会占用大量堆空间，尤其是在加载大量的文本字段时，默认禁用
}

func NewField() *Field {
	return &Field{}
}

func (f *Field) SetType(t string) *Field {
	f.Type = t
	return f
}

func (f *Field) SetFieldsType(t string) *Field {
	f.Fields = &map[string]Fields{
		row: {Type: t},
	}
	return f
}

func (f *Field) SetIndex(index bool) *Field {
	f.Index = index
	return f
}

func (f *Field) SetFormat(fmat string) *Field {

	if f.Type != Type.Date() {
		panic("set format date error")
	}
	f.Format = fmat
	return f
}

// text to ik when save
func (f *Field) SetAnalyzer(a string) *Field {
	f.Analyzer = a
	return f
}

// text to ik when search
func (f *Field) SetSearchAnalyzer(s string) *Field {
	f.SearchAnalyzer = s
	return f
}

type Geo struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
