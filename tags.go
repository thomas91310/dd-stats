package main

import "fmt"

var (
	//NoDefaultTags represents an empty slice of Datadog Tags
	//It is used when no tags need to be sent to the dogstatsd
	NoDefaultTags = DataDogTags{
		Tags: nil,
	}
)

// DataDogTag represents a datadog tag
type DataDogTag struct {
	name  string
	value string
}

//NewDataDogTag returns a new DataDogTag
func NewDataDogTag(name string, value string) DataDogTag {
	return DataDogTag{
		name:  name,
		value: value,
	}
}

// String returns how a tag is expected by DataDog
func (ddt DataDogTag) String() string {
	return fmt.Sprintf("%s:%s", ddt.name, ddt.value)
}

// DataDogTags represents the list entire list of tags
// for the metric that will be sent
type DataDogTags struct {
	Tags []DataDogTag
}

// Format returns the expected []string for DataDog sdk
func (ddt DataDogTags) Format() []string {
	tags := []string{}
	for _, tag := range ddt.Tags {
		tags = append(tags, tag.String())
	}
	return tags
}
