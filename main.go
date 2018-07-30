package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

func main() {
	log.SetFlags(0)

	z := make(map[string]map[string]string)

	for _, source := range Sources {
		dd := Fetch(source)

		for image, d := range dd {
			if _, ok := z[image]; !ok {
				z[image] = make(map[string]string)
			}

			for k, v := range d {
				if v == "" {
					continue
				}

				if z[image][k] == "" {
					z[image][k] = v
				} else if z[image][k] == v {
					// do nothing.
				} else {
					// z[image][k] += " ;;; " + v
					z[image][k] += fmt.Sprintf(" ;;; [%s] %s", source.Name, v)
				}
			}
		}
	}

	var ids []string
	id2image := make(map[string]string)
	for image, d := range z {
		id := d["Identifier"]
		ids = append(ids, id)
		id2image[id] = image
	}
	sort.Strings(ids)

	for _, id := range ids {
		d := z[id2image[id]]

		fmt.Printf("\n\n%s\n\n", id)

		var keys []string
		for k, _ := range d {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fmt.Printf("\t%s :: \t%s\n", k, d[k])
		}
	}

	fd, err := os.Create("MERGED.txt")
	if err != nil {
		log.Fatalf("Cannot create MERGED.txt: %v", err)
	}
	w := csv.NewWriter(fd)
	for _, id := range ids {
		d := z[id2image[id]]
		id := strings.Split(d["Identifier"], " ")[0]
		imageURL := d["ImageURL"]

		err = w.Write([]string{id, imageURL})
		if err != nil {
			log.Fatalf("Cannot Write to MERGED.txt: %v", err)
		}
	}
	w.Flush()
	err = fd.Close()
	if err != nil {
		log.Fatalf("Cannot close MERGED.txt: %v", err)
	}
}

func Fetch(sd SourceDef) map[string]map[string]string {
	response, err := http.Get(sd.URL)
	if err != nil {
		log.Fatalf("Cannot http.get %s %q: %v", sd.Name, sd.URL, err)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("failure in ioutil.ReadAll of response body %s: %v", sd.Name, err)
	}

	err = ioutil.WriteFile(sd.Name+".txt", contents, 0644)

	buf := bytes.NewBuffer(contents)
	r := csv.NewReader(buf)
	r.Comment = '#'
	r.FieldsPerRecord = -1

	recs, err := r.ReadAll()
	if err != nil {
		log.Fatalf("failure in csv.ReadAll of response body %s: %v", sd.Name, err)
	}

	schema := strings.Split(sd.Schema, ",")
	n := len(schema)

	z := make(map[string]map[string]string)

	for i, rec := range recs {
		_ = i
		// println(i, rec)
		d := make(map[string]string)
		for j, word := range rec {
			if j < n {
				// println(i, j, schema[j], word)
				d[schema[j]] = word
			} else {
				// log.Printf("%q %d/%d EXTRA: %s", sd.Name, i, j, word)
				d["Extra"] = word
			}
		}

		image := d["ImageURL"]
		if image == "ImageURL" {
			continue  // This came from a schema line.
		}
		if image == "" {
			log.Fatalf("No ImageURL: %v", d)
		}
		z[image] = d
	}

	return z
}

type SourceDef struct {
	Name   string
	URL    string
	Schema string
}

var Sources = []SourceDef{
	{"Scott",
		"https://www.swharden.com/qrss/plus/grabbers.csv",
		"Identifier,Callsign,Title,Name,Location,WebSiteURL,ImageURL"},
	{"Henry",
		"https://docs.google.com/spreadsheets/d/1gjtBlXywKiKzC5nnVe7h7w9hTR18_o8TTP0XNVMMGyI/export?format=csv",
		"Timestamp,Comment,Identifier,ImageURL,Location,Callsign,Name,Title,WebSiteURL,RecentlyUp"},
	{"Andy",
		"http://qsl.net/g0ftd/activegrabberlist.txt",
		"Identifier,ImageURL,RecentlyUp"},
}
