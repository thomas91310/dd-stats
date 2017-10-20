package main

import "github.com/thomas91310/datadog-go/statsd"

func main() {
	dvc, err := NewDataDogSVC(
		ClientEnable,
		DefaultAddress,
		"YO_NAMESPACE",
		NoDefaultTags,
	)
	if err != nil {
		panic(err)
	}
	dvc.SendEvent(
		"S3_DOWNLOADER",
		"COUNTER",
		"Starting to download events from S3",
		NoDefaultTags,
		statsd.Info,
		statsd.Low,
	)

	dvc.SendEvent(
		"S3_DOWNLOADER",
		"EXIT",
		"Done downloading events from S3",
		NoDefaultTags,
		statsd.Warning,
		statsd.Normal,
	)
}
