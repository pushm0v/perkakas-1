package elastic

// ESDynamicTemplate https://www.elastic.co/guide/en/elasticsearch/reference/current/dynamic-templates.html
type MatchConditions struct {
	Match            string       `json:"match,omitempty"`
	UnMatch          string       `json:"unmatch,omitempty"`
	PathMatch        string       `json:"path_match,omitempty"`
	MatchMappingType string       `json:"match_mapping_type,omitempty"`
	Mapping          MatchMapping `json:"mapping,omitempty"`
}

type MatchMapping struct {
	Type            string           `json:"type"`
	Norms           bool             `json:"norms"`
	Fields          map[string]Field `json:"fields,omitempty"`
	IncludeInParent bool             `json:"include_in_parent,omitempty"`
}

type Field struct {
	Type        string `json:"type"`
	IgnoreAbove int    `json:"ignore_above,omitempty"`
}
