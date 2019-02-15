// uncharblog.go ---
//
// Filename: uncharblog.go
// Description:
// Author: Mateo Rodriguez Ripolles
// Maintainer:
// Created: ven. févr. 15 10:29:04 2019 (+0100)
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
	"io/ioutil"
	"log"
	"net/http"
)

const STYLE_SHEETS_LOCAL_PATH = "./src/stylesheets/"
const STYLE_SHEETS_URL_PATH = "/stylesheets/"
const TEXT_LOCAL_PATH = "./resources/text/"

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) Save() error {
	var filename string

	filename = p.Title + ".txt"
	return (ioutil.WriteFile(filename, p.Body, 0600))
}

func LoadPage(title string) (*Page, error) {
	var filename string
	var body []byte
	var err error

	filename = title + ".txt"
	body, err = ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func main() {
	s, err := NewUncharServer()
	if err != nil {
		log.Fatal(err)
	}
	http.Handle(STYLE_SHEETS_URL_PATH, http.StripPrefix(STYLE_SHEETS_URL_PATH, http.FileServer(http.Dir(STYLE_SHEETS_LOCAL_PATH))))
	http.HandleFunc(VIEW_TAG, s.MakeHandler(s.ViewHandler))
	http.HandleFunc(EDIT_TAG, s.MakeHandler(s.EditHandler))
	http.HandleFunc(SAVE_TAG, s.MakeHandler(s.SaveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

//
// uncharblog.go ends here