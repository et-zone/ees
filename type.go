package ees


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

