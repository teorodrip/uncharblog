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
)

const VIEW_TAG = "/view/"
const EDIT_TAG = "/edit/"
const SAVE_TAG = "/save/"

const POST_BODY_TAG = "body"

func RenderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	var title string
	var p *Page
	var err error

	title = r.URL.Path[len(VIEW_TAG):]
	p, err = LoadPage(title)
	if err != nil {
		http.Redirect(w, r, EDIT_TAG+title, http.StatusFound)
		return
	}
	RenderTemplate(w, "view", p)
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	var Title string
	var p *Page
	var err error

	Title = r.URL.Path[len(EDIT_TAG):]
	p, err = LoadPage(Title)
	if err != nil {
		p = &Page{Title: Title}
	}
	RenderTemplate(w, "edit", p)
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	var title, body string
	var p *Page

	title = r.URL.Path[len(SAVE_TAG):]
	body = r.FormValue(POST_BODY_TAG)
	p = &Page{Title: title, Body: []byte(body)}
	p.Save()
	http.Redirect(w, r, VIEW_TAG+title, http.StatusFound)
}

//
// uncharserver.go ends here
