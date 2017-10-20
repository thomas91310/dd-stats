package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTag(t *testing.T) {
	tagName := "host-id"
	tagValue := "27"
	ddt := NewDataDogTag(tagName, tagValue)
	assert.Equal(t, tagName, ddt.name)
	assert.Equal(t, tagValue, ddt.value)
}

func TestFormatTag(t *testing.T) {
	tagName := "host-id"
	tagValue := "27"
	ddt := NewDataDogTag(tagName, tagValue)
	tag := fmt.Sprintf("%s:%s", tagName, tagValue)
	assert.Equal(t, tag, ddt.String())
}

func TestFormatMultipleTags(t *testing.T) {
	ddts := DataDogTags{
		Tags: []DataDogTag{
			NewDataDogTag("host-id", "27"),
			NewDataDogTag("bucket-name", "s3://yo/2017-10-10-07-45"),
			NewDataDogTag("collection", "yo"),
		},
	}

	expected := []string{
		"host-id:27",
		"bucket-name:s3://yo/2017-10-10-07-45",
		"collection:yo",
	}

	assert.Equal(t, true, reflect.DeepEqual(expected, ddts.Format()))
}
