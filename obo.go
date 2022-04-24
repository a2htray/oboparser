package oboparser

import (
	"encoding/json"
)

type StanzaType uint

const (
	TermStanza StanzaType = iota
	TypedefStanza
	InstanceStanza
)

func (s StanzaType) MarshalJSON() ([]byte, error) {
	switch s {
	case TermStanza:
		return []byte("\"Term\""), nil
	case TypedefStanza:
		return []byte("\"Typedef\""), nil
	case InstanceStanza:
		return []byte("\"Instance\""), nil
	}
	return nil, nil
}

type Stanza struct {
	Type   StanzaType        `json:"type"`
	Values map[string]string `json:"values"`
}

var headerFields = []string{
	"format-version",
	"data-version",
	"date",
	"default-namespace",
	"saved-by",
	"auto-generated-by",
	"subsetdef",
	"import",
	"synonymtypedef",
	"idspace",
	"default-relationship-id",
	"idmapping",
	"remark",
}

var termFields = []string{
	"id",
	"is_anonymous",
	"name",
	"namespace",
	"alt_id",
	"def",
	"comment",
	"subset",
	"synonym",
	"xref",
	"is_a",
	"derived",
	"intersection_of",
	"union_of",
	"disjoint_from",
	"relationship",
	"is_obsolete",
	"replaced_by",
	"consider",
}

var typedefFields = []string{
	"id",
	"is_anonymous",
	"name",
	"namespace",
	"alt_id",
	"def",
	"comment",
	"subset",
	"synonym",
	"xref",
	"domain",
	"range",
	"is_anti_symmetric",
	"is_cyclic",
	"is_reflexive",
	"is_symmetric",
	"is_transitive",
	"is_a",
	"inverse",
	"transitive_over",
	"relationship",
	"is_metadata_tag",
	"is_obsolete",
	"replaced_by",
	"consider",
}

var instanceFields = []string{
	"id",
	"is_anonymous",
	"name",
	"namespace",
	"alt_id",
	"comment",
	"synonym",
	"xref",
	"instance_of",
	"property_value",
	"is_obsolete",
	"replaced_by",
	"consider",
}

func ExtendHeaderFields(fields []string) {
	headerFields = append(headerFields, fields...)
}

func ExtendTermFields(fields []string) {
	termFields = append(termFields, fields...)
}

func ExtendTypedefFields(fields []string) {
	typedefFields = append(typedefFields, fields...)
}

func ExtendInstanceFields(fields []string) {
	instanceFields = append(instanceFields, fields...)
}

func newOBO() *OBO {
	return &OBO{
		headerFields: make([]string, 0),
		headerMap:    make(map[string][]string),
		Stanzas:      make([]*Stanza, 0),
	}
}

type OBO struct {
	headerFields []string
	headerMap    map[string][]string
	Stanzas      []*Stanza
}

func (o *OBO) HeaderValue(field string) string {
	if stringIn(field, o.headerFields) {
		return o.headerMap[field][0]
	}
	return ""
}

func (o *OBO) HeaderValues(field string) []string {
	if stringIn(field, o.headerFields) {
		return o.headerMap[field]
	}
	return []string{}
}

func (o *OBO) SetHeader(field string, value string) {
	if !stringIn(field, headerFields) { return }

	if !stringIn(field, o.headerFields) {
		o.headerFields = append(o.headerFields, field)
		o.headerMap[field] = make([]string, 0)
	}
	o.headerMap[field] = append(o.headerMap[field], value)
}

type jsonOBO map[string]interface{}

func createJSONOBO(o OBO) jsonOBO {
	oboJSON := make(map[string]interface{})
	for field, value := range o.headerMap {
		if len(value) > 1 {
			oboJSON[field] = value
		} else {
			oboJSON[field] = value[0]
		}
	}
	oboJSON["stanzas"] = o.Stanzas
	return oboJSON
}

func (o OBO) JSON() string {
	oboJSON := createJSONOBO(o)
	bs, err := json.Marshal(oboJSON)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

func (o OBO) JSONIdent(prefix string, indent string) string {
	oboJSON := createJSONOBO(o)
	bs, err := json.MarshalIndent(oboJSON, prefix, indent)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

