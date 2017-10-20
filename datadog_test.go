package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/thomas91310/datadog-go/statsd"

	"github.com/stretchr/testify/assert"
)

func TestDisableNewDataDogSVC(t *testing.T) {
	dds, err := NewDataDogSVC(ClientDisable, DefaultAddress, DefaultAppNamespace, NoDefaultTags)
	assert.Nil(t, err)
	assert.NotNil(t, dds)
	assert.Nil(t, dds.client)
	assert.Equal(t, ClientDisable, dds.disable)
	assert.Equal(t, DefaultRuntimeTick, dds.runtimeTick)
	assert.Equal(t, DefaultLimitConsecutiveErrors, dds.limitConsecutiveErrors)
}

func TestActiveNewDataDogSVC(t *testing.T) {
	statdsFakeClient := &statsd.Client{}
	dds, err := NewDataDogSVC(ClientEnable, DefaultAddress, DefaultAppNamespace, NoDefaultTags)
	assert.Nil(t, err)
	assert.NotNil(t, dds)
	assert.Equal(t, 0, len(dds.client.Tags))
	assert.Equal(t, DefaultAppNamespace, dds.namespace)
	assert.Equal(t, reflect.TypeOf(statdsFakeClient), reflect.TypeOf(dds.client))
	assert.Equal(t, ClientEnable, dds.disable)
	assert.Equal(t, DefaultRuntimeTick, dds.runtimeTick)
	assert.Equal(t, DefaultLimitConsecutiveErrors, dds.limitConsecutiveErrors)
}

func TestNewDataDogSVCBadStatsDServer(t *testing.T) {
	address := "cantfindme:"
	port := "8125"
	dds, err := NewDataDogSVC(ClientEnable, address+port, DefaultAppNamespace, NoDefaultTags)
	assert.NotNil(t, err)
	assert.Nil(t, dds)
	errSplit := strings.Split(fmt.Sprintf("%v", err), ".")
	assert.Equal(t, fmt.Sprintf("Error creating statsD client to %v", address+port), errSplit[0])
}

func TestDisableSendGauge(t *testing.T) {
	dds := DataDogSVC{
		client:                 nil,
		disable:                ClientDisable,
		namespace:              DefaultAppNamespace,
		runtimeTick:            DefaultRuntimeTick,
		limitConsecutiveErrors: DefaultLimitConsecutiveErrors,
	}
	assert.Nil(t, dds.SendGauge("s3_downloader", "count", 1, NoDefaultTags))
}

func TestDisableSendEvent(t *testing.T) {
	dds := DataDogSVC{
		client:                 nil,
		disable:                ClientDisable,
		namespace:              DefaultAppNamespace,
		runtimeTick:            DefaultRuntimeTick,
		limitConsecutiveErrors: DefaultLimitConsecutiveErrors,
	}
	assert.Nil(t, dds.SendEvent("pipeline", "exit", "Exiting the pipeline", NoDefaultTags, statsd.Warning, statsd.Low))
}
