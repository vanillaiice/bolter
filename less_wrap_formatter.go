package main

import (
	"io"
	"os/exec"
)

// MoreWrapFormatter is a formatter that prints in more format.
type MoreWrapFormatter struct {
	formatter Formatter
}

// wrapDump pipes the output of the dump function into more.
func (mf MoreWrapFormatter) wrapDump(w io.Writer, dump func(io.Writer)) {
	lessCmd := exec.Command("more")
	pipeR, pipeW := io.Pipe()

	go func() {
		dump(pipeW)
		pipeW.Close()
	}()

	lessCmd.Stdin = pipeR
	lessCmd.Stdout = w

	lessCmd.Run()
}

// DumpBuckets prints in more format.
func (mf MoreWrapFormatter) DumpBuckets(w io.Writer, buckets []Bucket) {
	mf.wrapDump(w, func(w io.Writer) {
		mf.formatter.DumpBuckets(w, buckets)
	})
}

// DumpBucketItems prints the bucket's items in more format.
func (mf MoreWrapFormatter) DumpBucketItems(w io.Writer, bucket string, items []Item) {
	mf.wrapDump(w, func(w io.Writer) {
		mf.formatter.DumpBucketItems(w, bucket, items)
	})
}
