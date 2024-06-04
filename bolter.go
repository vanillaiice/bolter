package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	kval "github.com/kval-access-language/kval-boltdb"
)

// instructionLine is used to display instructions in the terminal.
const instructionLine = "> Enter bucket to explore (CTRL-X to quit, CTRL-B to go back, ENTER to go back to ROOT Bucket):"

// goingBack is used to display 'going back' in the terminal.
const goingBack = "> Going back..."

// Impl holds the implementation details.
type Impl struct {
	kb     kval.Kvalboltdb // kb is the connection to the bolt database.
	fmt    Formatter       // fmt is the formatter interface.
	bucket string          // bucket is the bucket name.
	loc    string          // loc is the current location in the database.
	cache  string          // cache is the previous location in the database.
	root   bool            // root is used to determine if we are in the ROOT bucket.
}

// initDB initializes the boltdb connection.
func (i *Impl) initDB(file string) {
	var err error
	i.kb, err = kval.Connect(file)
	if err != nil {
		log.Fatal(err)
	}
}

// updateLoc updates the current location in the database.
func (i *Impl) updateLoc(bucket string, goBack bool) string {
	if bucket == i.cache {
		i.loc = bucket
		return i.loc
	}

	if goBack {
		s := strings.Split(i.loc, ">>")
		i.loc = strings.Join(s[:len(s)-1], ">>")
		i.bucket = strings.Trim(s[len(s)-2], " ")
		return i.loc
	}

	if i.loc == "" {
		i.loc = bucket
		i.bucket = bucket
	} else {
		i.loc = i.loc + " >> " + bucket
		i.bucket = bucket
	}

	return i.loc
}

// listBucketItems lists the items in a bucket.
func (i *Impl) listBucketItems(bucket string, goBack bool) {
	items := []Item{}
	getQuery := i.updateLoc(bucket, goBack)

	if getQuery != "" {
		fmt.Fprintf(os.Stdout, "Query: "+getQuery+"\n\n")
		res, err := kval.Query(i.kb, "GET "+getQuery)
		if err != nil {
			if err.Error() == "No Keys: There are no key::value pairs in this bucket" {
				fmt.Fprintf(os.Stdout, "> There are no key::value pairs in this bucket\n")
				if i.root {
					i.listBuckets()
					return
				}
				i.listBucketItems(i.loc, true)
			} else if err.Error() != "Cannot GOTO bucket, bucket not found" {
				log.Fatal(err)
			} else {
				fmt.Fprintf(os.Stdout, "> Bucket not found\n")
				if i.root {
					i.listBuckets()
					return
				}
				i.listBucketItems(i.loc, true)
			}
		}

		if len(res.Result) == 0 {
			fmt.Fprintf(os.Stdout, "Invalid request.\n\n")
			i.listBucketItems(i.cache, false)
			return
		}

		for k, v := range res.Result {
			if v == kval.Nestedbucket {
				k = k + "*"
				v = ""
			}
			items = append(items, Item{Key: string(k), Value: string(v)})
		}

		fmt.Fprintf(os.Stdout, "Bucket: %s\n", bucket)

		i.fmt.DumpBucketItems(os.Stdout, i.bucket, items)

		i.root = false

		i.cache = getQuery

		outputInstructionline()
	}
}

// listBuckets lists the buckets in the database.
func (i *Impl) listBuckets() {
	i.root = true
	i.loc = ""

	buckets := []Bucket{}

	res, err := kval.Query(i.kb, "GET _")
	if err != nil {
		log.Fatal(err)
	}

	for k := range res.Result {
		buckets = append(buckets, Bucket{Name: string(k) + "*"})
	}

	fmt.Fprint(os.Stdout, "DB Layout:\n\n")

	i.fmt.DumpBuckets(os.Stdout, buckets)

	outputInstructionline()
}

// readInput reads user input from stdin.
func (i *Impl) readInput() {
	i.listBuckets()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Fprintln(os.Stdout, "")

		bucket := scanner.Text()

		switch bucket {
		case "\x18":
			return
		case "\x02":
			if !strings.Contains(i.loc, "") || !strings.Contains(i.loc, ">>") {
				fmt.Fprintf(os.Stdout, "%s\n", goingBack)
				i.loc = ""
				i.listBuckets()
			} else {
				i.listBucketItems(bucket, true)
			}
		case "":
			i.listBuckets()
		default:
			i.listBucketItems(bucket, false)
		}

		bucket = ""
	}
}

// Formatter interface.
type Formatter interface {
	DumpBuckets(io.Writer, []Bucket)
	DumpBucketItems(io.Writer, string, []Item)
}

// Item is a key/value pair in a bucket.
type Item struct {
	Key   string
	Value string
}

// Bucket is a bucket in the database.
type Bucket struct {
	Name string
}

// outputInstructionline outputs the instruction line.
func outputInstructionline() {
	fmt.Fprintf(os.Stdout, "\n%s\n\n", instructionLine)
}
