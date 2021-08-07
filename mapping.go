package ees

//版本7.x
const (
	DEFAULT_NUM = 0
	//type
	TYPE = Type(DEFAULT_NUM)
	//dynamic type
	DYNAMIC = Dynamic(DEFAULT_NUM)
	//index type
	INDEX = Index(DEFAULT_NUM)

	//yyyy-MM-dd HH:mm:ss
	DATE_TIME_FORMAT = "yyyy-MM-dd HH:mm:ss"
	//yyyy-MM-dd
	DATE_FORMAT = "yyyy-MM-dd"
	//HH:mm:ss
	TIME_FORMAT = "HH:mm:ss"
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

//return mappings interface{}
func (m *Mapping) Mappings() interface{} {

	return map[string]Mapping{
		"mappings": *m,
	}
}

type Index int

func (index Index) Enabled() bool {
	return true
}
func (index Index) Disabled() bool {
	return false
}

type Type int

func (f Type) Long() string {
	return "long"
}

func (f Type) Float() string {
	return "float"
}

//text类型自动会分词，支持模糊搜索，使用=匹配即可
func (f Type) Text() string {
	return "text"
}

func (f Type) Keyword() string {
	return "keyword"
}

func (f Type) Date() string {
	return "date"
}

func (f Type) Boolean() string {
	return "boolean"
}

func (f Type) Geo() string {
	return "geo_point"
}

func (f Type) IP() string {
	return "ip"
}

type Dynamic string

func (d Dynamic) True() string {
	return "true"
}

func (d Dynamic) False() string {
	return "false"
}

func (d Dynamic) Strict() string {
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

type Analyzer int

//索引时分词
func (a Analyzer) IkMaxWord() string {
	return "ik_max_word"
}

//搜索时分词
func (a Analyzer) IkSmart() string {
	return "ik_smart"
}

//搜索阶段分词，会覆盖Analyzer的属性
type SearchAnalyzer int

//索引时分词
func (a SearchAnalyzer) IkSmart() string {
	return "ik_smart"
}

type Field struct {
	Type           string `json:"type"`
	Index          bool   `json:"index,omitempty"`           //默认选true，支持搜索,旧版本使用:"not_analyzed"
	Format         string `json:"format,omitempty"`          //时间类型格式化
	Analyzer       string `json:"analyzer,omitempty"`        //索引存储阶段和搜索阶段都分词，索引时用ik_max_word，搜索时分词器用ik_smart
	SearchAnalyzer string `json:"search_analyzer,omitempty"` //搜索阶段分词，会覆盖上面的属性 ，"search_analyzer": "ik_smart"
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

func (f *Field) SetIndex(index bool) *Field {
	f.Index = index
	return f
}

func (f *Field) SetFormat(fmat string) *Field {

	if f.Type != TYPE.Date() {
		panic("set format date error")
	}
	f.Format = fmat
	return f
}

func (f *Field) SetAnalyzer(a string) *Field {
	f.Analyzer = a
	return f
}

func (f *Field) SetSearchAnalyzer(s string) *Field {
	f.SearchAnalyzer = s
	return f
}

type Geo struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
