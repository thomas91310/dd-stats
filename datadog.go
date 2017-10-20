package main

import (
	"fmt"
	"runtime"
	"time"

	logging "github.com/op/go-logging"
	"github.com/thomas91310/datadog-go/statsd"
)

var log = logging.MustGetLogger("stats")

const (
	//ClientDisable is used when we don't want to send anything to dogstatsd
	ClientDisable = true
	//ClientEnable is used when we want to send metrics to dogstatsd
	ClientEnable = false
	//DefaultAddress is the default statsd address
	DefaultAddress = "localhost:8125"
	//DefaultAppNamespace is the default dogstatsd namespace used
	DefaultAppNamespace = "article_metrics"
	//DefaultRuntimeTick is how frequent we send runtime stats via AppMetrics goroutine
	DefaultRuntimeTick = 1
	//DefaultLimitConsecutiveErrors defines the limit of consecutive errors that need to happen
	//before we stop sending stats to dogstatsd
	DefaultLimitConsecutiveErrors = 10
)

// DataDogSVC is a service that holds a DataDog client. The client connects to the agent running
// on the host on port :8125
type DataDogSVC struct {
	client                 *statsd.Client
	disable                bool
	namespace              string
	runtimeTick            int
	limitConsecutiveErrors int
}

// NewDataDogSVC creates a new DataDogSVC
func NewDataDogSVC(disable bool, address string, namespace string, tags DataDogTags) (*DataDogSVC, error) {
	if disable {
		return &DataDogSVC{
			client:                 nil,
			disable:                disable,
			namespace:              DefaultAppNamespace,
			limitConsecutiveErrors: DefaultLimitConsecutiveErrors,
			runtimeTick:            DefaultRuntimeTick,
		}, nil
	}
	statsDClient, err := statsd.New(address)
	if err != nil {
		return nil, fmt.Errorf("Error creating statsD client to %v. Got %v", address, err)
	}
	statsDClient.Tags = tags.Format()

	return &DataDogSVC{
		client:                 statsDClient,
		disable:                disable,
		namespace:              namespace,
		limitConsecutiveErrors: DefaultLimitConsecutiveErrors,
		runtimeTick:            DefaultRuntimeTick,
	}, nil
}

// SendGauge sends the gauge to dogstatsd
func (dds DataDogSVC) SendGauge(componentName string, metricName string, metricValue float64, extraTags DataDogTags) error {
	if dds.disable {
		return nil
	}
	m := NewMetric(dds.namespace, componentName, metricName)
	tags := extraTags.Format()

	err := dds.client.Gauge(m.String(), metricValue, tags, 1)
	if err != nil {
		return fmt.Errorf("Error sending metric %v with value %v. Got %v", m.String(), metricValue, err)
	}
	return nil
}

// SendHistogram sends an histogram to dogstatsd
func (dds DataDogSVC) SendHistogram(componentName string, metricName string, metricValue float64, extraTags DataDogTags) error {
	if dds.disable {
		return nil
	}
	m := NewMetric(dds.namespace, componentName, metricName)
	tags := extraTags.Format()

	err := dds.client.Histogram(m.String(), metricValue, tags, 1)
	if err != nil {
		return fmt.Errorf("Error sending histogram %v with value %v. Got %v", m.String(), metricValue, err)
	}
	return nil
}

//SendEvent sends an event to dogstatsd
func (dds DataDogSVC) SendEvent(componentName string, eventName string, eventDescription string, extraTags DataDogTags, alertType statsd.EventAlertType, eventPriority statsd.EventPriority) error {
	if dds.disable {
		return nil
	}

	m := NewMetric(dds.namespace, componentName, eventName)
	tags := extraTags.Format()

	ev := &statsd.Event{
		Title:     m.String(),
		Text:      eventDescription,
		Timestamp: time.Now(),
		AlertType: alertType,
		Priority:  eventPriority,
		Tags:      tags,
	}

	err := dds.client.Event(ev)
	if err != nil {
		return fmt.Errorf("Error sending event %v. Got %v", ev, err)
	}
	return nil
}

//SendMemStats sends low level memory metrics of the application
func (dds DataDogSVC) SendMemStats(componentName string, extraTags DataDogTags) error {
	if dds.disable {
		return nil
	}
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	//these are the bytes that were allocated and still in use
	err := dds.SendGauge(componentName, "mem_alloc", float64(mem.Alloc), extraTags)
	if err != nil {
		return err
	}
	//total bytes allocated by the program
	err = dds.SendGauge(componentName, "mem_total_alloc", float64(mem.TotalAlloc), extraTags)
	if err != nil {
		return err
	}
	//being used on the heap right now
	err = dds.SendGauge(componentName, "mem_heap_alloc", float64(mem.HeapAlloc), extraTags)
	if err != nil {
		return err
	}
	//this includes what is being used by the heap and what has been reclaimed by the OS but not given back out
	err = dds.SendGauge(componentName, "mem_heap_sys", float64(mem.HeapSys), extraTags)
	if err != nil {
		return err
	}
	//this is the number of bytes that were returned to the OS
	err = dds.SendGauge(componentName, "mem_heap_released", float64(mem.HeapReleased), extraTags)
	if err != nil {
		return err
	}
	//this is the number of bytes that are not used by the application but retained by the GC
	err = dds.SendGauge(componentName, "mem_heap_idle", float64(mem.HeapIdle), extraTags)
	if err != nil {
		return err
	}
	//this sends the number of goroutines that currently exist.
	dds.SendGauge(componentName, "num_goroutines", float64(runtime.NumGoroutine()), extraTags)
	if err != nil {
		return err
	}
	return nil
}

//AppMetrics sends low level metrics every 5 seconds
func (dds *DataDogSVC) AppMetrics(componentName string, bucketName string) {
	if dds.disable {
		return
	}

	updateInterval := time.Tick(time.Duration(dds.runtimeTick) * time.Second)
	consecutiveErrors := 0
loop:
	for {
		select {
		case <-updateInterval:
			err := dds.SendMemStats(
				componentName,
				NoDefaultTags,
			)
			if err == nil {
				consecutiveErrors = 0
			} else {
				log.Infof("Error sending memory stats: %v", err)
				if consecutiveErrors == dds.limitConsecutiveErrors {
					log.Infof("Quit sending stats to Dogstatsd. Last error: %v", err)
					dds.disable = ClientDisable
					break loop
				}
				consecutiveErrors++
			}
		}
	}
}
