package main

import "fmt"

//Metric constitutes a datadog metric name
type Metric struct {
	namespace string
	component string
	name      string
}

//NewMetric creates a new metric given a namespace and a component and a metric name
func NewMetric(namespace string, component string, name string) Metric {
	return Metric{
		namespace: namespace,
		component: component,
		name:      name,
	}
}

//String returns the name of the metric
func (m Metric) String() string {
	return fmt.Sprintf("%v.%v.%v", m.namespace, m.component, m.name)
}
