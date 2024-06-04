package main

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

// TableFormatter is a formatter that prints in table format.
type TableFormatter struct {
	noValues bool
}

// DumpBuckets prints in table format.
func (tf TableFormatter) DumpBuckets(w io.Writer, buckets []Bucket) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Buckets"})

	for _, b := range buckets {
		row := []string{b.Name}
		table.Append(row)
	}

	table.Render()
}

// DumpBucketItems prints the bucket's items in table format.
func (tf TableFormatter) DumpBucketItems(w io.Writer, bucket string, items []Item) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Key", "Value"})

	for _, item := range items {
		var row []string
		if tf.noValues {
			row = []string{item.Key, ""}
		} else {
			row = []string{item.Key, item.Value}
		}
		table.Append(row)
	}

	table.Render()
}
