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
	BoolTrue = true
	BoolFalse = false
)

type Mapping struct {
	Dynamic string            `json:"dynamic,omitempty"`//default true
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