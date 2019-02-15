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
	"html/template"
	"net/http"
	"regexp"
)

const VIEW_TAG = "/view/"
const EDIT_TAG = "/edit/"
const SAVE_TAG = "/save/"
const POST_BODY_TAG = "body"
const HTML_DIR = "./src/html/"

type UncharServer struct {
	Templates *template.Template
	ValidPath *regexp.Regexp
}

func NewUncharServer() (*UncharServer, error) {
	file_names, err := GetFilesFromDir(HTML_DIR)
	if err != nil {
		return nil, err
	}
	server := &UncharServer{
		Templates: template.Must(template.ParseFiles(file_names...)),
		ValidPath: regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")}
	return server, nil
}

func (s *UncharServer) RenderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	var err error

	err = s.Templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *UncharServer) MakeHandler(fn func(w http.ResponseWriter, r *http.Request, title string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m []string

		m = s.ValidPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
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

//
// uncharserver.go ends here
