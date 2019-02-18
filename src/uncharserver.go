// uncharserver.go ---
//
// Filename: uncharserver.go
// Description:
// Author: Mateo Rodriguez Ripolles
// Maintainer:
// Created: ven. févr. 15 11:08:11 2019 (+0100)
// Version:
// Package-Requires: ()
// Last-Updated:
//           By:
//     Update #: 0
// URL:
// Doc URL:
// Keywords:
// Compatibility:
//
//

// Commentary:
//
//
//
//

// Change Log:
//
//
//
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or (at
// your option) any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with GNU Emacs.  If not, see <https://www.gnu.org/licenses/>.
//
//

// Code:

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

const VIEW_TAG = "/view/"
const EDIT_TAG = "/edit/"
const SAVE_TAG = "/save/"
const INDEX_TAG = "/"
const POST_BODY_TAG = "body"
const HTML_DIR = "./src/html/"

type UncharServer struct {
	Templates *template.Template
	ValidPath *regexp.Regexp
}

type Post struct {
	Id    string
	Title string
	Path  string
}

func NewUncharServer() (*UncharServer, error) {
	file_names, err := GetFilesFromDir(HTML_DIR)
	if err != nil {
		return nil, err
	}
	tmpl := template.New("uncharblog")
	tmpl.ParseFiles(file_names...)
	server := &UncharServer{
		Templates: tmpl,
		ValidPath: regexp.MustCompile("^/((edit|save|view)/[a-zA-Z0-9]+)?$")}
	return server, nil
}

func (s *UncharServer) RenderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	var err error

	err = s.Templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing %s template\n", tmpl+".html")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *UncharServer) MakeHandler(fn func(w http.ResponseWriter, r *http.Request, title string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m []string

		m = s.ValidPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			fmt.Fprintf(os.Stderr, "Invalid URL catched: %s\n", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func (s *UncharServer) ViewHandler(w http.ResponseWriter, r *http.Request, title string) {
	var p *Page
	var err error

	if err != nil {
		return
	}
	p, err = LoadPage(title)
	if err != nil {
		http.Redirect(w, r, EDIT_TAG+title, http.StatusFound)
		return
	}
	s.RenderTemplate(w, "view", p)
}

func (s *UncharServer) EditHandler(w http.ResponseWriter, r *http.Request, title string) {
	var p *Page
	var err error

	if err != nil {
		return
	}
	p, err = LoadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	s.RenderTemplate(w, "edit", p)
}

func (s *UncharServer) SaveHandler(w http.ResponseWriter, r *http.Request, title string) {
	var body string
	var p *Page
	var err error

	if err != nil {
		return
	}
	body = r.FormValue(POST_BODY_TAG)
	p = &Page{Title: title, Body: []byte(body)}
	err = p.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, VIEW_TAG+title, http.StatusFound)
}

func (s *UncharServer) IndexHandler(w http.ResponseWriter, r *http.Request, title string) {
	rows, err := ExeIndQuery("SELECT * FROM uncharblog.posts")
	if err != nil {
		return
	}
	defer rows.Close()
	var p Post

	for rows.Next() {
		rows.Scan(&p.Id, &p.Title, &p.Path)
		log.Println(p.Path)
	}
	// rawResult := make([][]byte, 3)
	// c := make([]interface{}, 3)
	// for i, _ := range rawResult {
	//	c[i] = &rawResult[i]
	// }
	// for rows.Next() {
	//	rows.Scan(c...)
	//	log.Println(int(rawResult[0]))
	// }
}

//
// uncharserver.go ends here
