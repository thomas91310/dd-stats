// +build unit

package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetric(t *testing.T) {
	componentName := "s3_downloader"
	metricName := "count"
	m := NewMetric(DefaultAppNamespace, componentName, metricName)
	assert.Equal(t, DefaultAppNamespace, m.namespace)
	assert.Equal(t, componentName, m.component)
	assert.Equal(t, metricName, m.name)
}

func TestFormatMetric(t *testing.T) {
	componentName := "s3_downloader"
	metricName := "count"
	m := NewMetric(DefaultAppNamespace, componentName, metricName)
	metric := fmt.Sprintf("%v.%v.%v", DefaultAppNamespace, componentName, metricName)
	assert.Equal(t, metric, m.String())
}
