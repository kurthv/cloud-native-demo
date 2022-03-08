// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/cypherfox/cloud-native-demo/pkg/k8s"
)

type Page struct {
	Title string
	Body  []byte
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
var k8sClient, err := k8s.NewKubeClient()

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}


func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func rootPage(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, fmt.Sprintf("Hello, Bug from %s \n", r.RemoteAddr))
	// io.WriteString(w, fmt.Sprintf("Current Number of Pods: %i \n", k8s.GetPods().length()))

	pods, err := k8sClient.GetPods()

	for _, pod := pods.Items {

		fmt.Printf("%s",pods)
	}
}

func main() {
	http.HandleFunc("/", rootPage)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
