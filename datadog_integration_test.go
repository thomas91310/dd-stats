// +build integration

package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDisableAppMetricsAfterNErrors(t *testing.T) {
	dds, err := NewDataDogSVC(ClientEnable, DefaultAddress, DefaultAppNamespace, NoDefaultTags)
	assert.Nil(t, err)
	assert.Equal(t, DefaultRuntimeTick, dds.runtimeTick)
	assert.Equal(t, DefaultLimitConsecutiveErrors, dds.limitConsecutiveErrors)

	//tick = 1second (every 1 second we'll send runtime metrics)
	//limitConsecutiveErrors = 3 (limit of consecutive errors before we stop sending metrics to dogstatsd)
	dds.limitConsecutiveErrors = 3
	dds.runtimeTick = 1
	go dds.AppMetrics("pipeline", "s3://ttheissier/2017-09-10-05")
	time.Sleep(6 * time.Second)
	assert.Equal(t, ClientDisable, dds.disable)
}
