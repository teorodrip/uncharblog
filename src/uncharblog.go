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
)

const STYLE_SHEETS_LOCAL_PATH = "./src/stylesheets/"
const STYLE_SHEETS_URL_PATH = "/stylesheets/"
const FONTS_LOCAL_PATH = "./resources/fonts/"
const FONTS_URL_PATH = "/fonts/"
const POST_LOCAL_PATH = "/home/unchartech_2/go/src/uncharblog/resources/posts/"

type Post struct {
	Id           string
	Title        string
	CreationDate string
	UpdateDate   string
	Fil          File
}

type Page struct {
	Title string
	Body  []byte
	List  []Post
}

type File struct {
	Path string
	Body []byte
}

func (f *File) SaveFile() error {
	return (ioutil.WriteFile(f.Path, f.Body, 0600))
}

func (f *File) LoadFile() error {
	var err error

	f.Body, err = ioutil.ReadFile(f.Path)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	s, err := NewUncharServer()
	if err != nil {
		log.Fatal(err)
	}
	s.Start()
}

//
// uncharblog.go ends here
