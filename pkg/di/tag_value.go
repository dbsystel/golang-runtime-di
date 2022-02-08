package di

import "strings"

const (
	TagPrefixQualifier = "qualifier="
	TagValueOptional   = "optional"
	AllQualifiers      = "*"
)

// TagValue represents the tag options for CDI
type TagValue struct {
	// Qualifier is the qualifier to be used
	Qualifier string
	// Required denotes if the injection is required and at least a single instance is necessary
	Required bool
}

func (v TagValue) IsAllQualifier() bool {
	return v.Qualifier == AllQualifiers
}

// TagValueFrom creates a new TagValue from the tag
func TagValueFrom(val string) TagValue {
	result := TagValue{Required: true}
	for _, part := range strings.Split(val, ",") {
		if part == TagValueOptional {
			result.Required = false
		}
		if strings.HasPrefix(part, TagPrefixQualifier) {
			result.Qualifier = part[len(TagPrefixQualifier):]
		}
	}
	return result
}
