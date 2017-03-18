// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package panel

import (
	"html/template"
	"net/http"
	"log"
	"sophie"
	"time"
)


var templates = template.Must(template.ParseFiles("src/panel/templates/index.html"))
var mr = new(sophie.Master)

func index(w http.ResponseWriter, r *http.Request) {
	mr.ElapsedTime = time.Since(mr.StartTime).Seconds()
	templates.ExecuteTemplate(w, "index.html", mr)
}

func StartServer(mrPassed *sophie.Master) {
	mr = mrPassed

	http.HandleFunc("/", index)

	fs := http.FileServer(http.Dir("src/panel/public/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Printf("Sophie Web UI available at http://%s:%s\n", "localhost", "8000")

	http.ListenAndServe(":8000", nil)
}