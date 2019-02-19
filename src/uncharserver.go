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
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
)

const VIEW_TAG = "/view/"
const EDIT_TAG = "/edit/"
const SAVE_TAG = "/save/"
const INDEX_TAG = "/"
const POST_BODY_TAG = "post_body"
const POST_TITLE_TAG = "post_title"
const HTML_DIR = "./src/html/"

const POST_LIMIT = 10

type UncharServer struct {
	Templates *template.Template
	ValidPath *regexp.Regexp
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
		ValidPath: regexp.MustCompile("^/((edit|save|view)/([a-zA-Z0-9]+))?$")}
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
		if m == nil || len(m) != 4 {
			fmt.Fprintf(os.Stderr, "Invalid URL catched: %s\n", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[3])
	}
}

func (s *UncharServer) LoadData(rows *sql.Rows, values ...interface{}) (error, bool) {
	col_names, err := rows.Columns()
	if err != nil || len(col_names) != len(values) {
		return err, false
	}
	if rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			return err, false
		}
		return nil, true
	}
	return nil, false
}

func (s *UncharServer) ViewHandler(w http.ResponseWriter, r *http.Request, title string) {
	var p Post
	var ok bool

	q := "SELECT post_id, post_title, post_path FROM uncharblog.posts WHERE post_id=$1"
	rows, err := ExeIndQuery(q, title)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer rows.Close()
	err, ok = s.LoadData(rows, &p.Id, &p.Title, &p.Fil.Path)
	if err != nil || !ok {
		http.Redirect(w, r, EDIT_TAG+title, http.StatusFound)
		return
	}
	p.Fil.LoadFile()
	s.RenderTemplate(w, "view", p)
}

func (s *UncharServer) EditHandler(w http.ResponseWriter, r *http.Request, title string) {
	var p Post
	var ok bool

	q := "SELECT post_title, post_path FROM uncharblog.posts WHERE post_id=$1"
	rows, err := ExeIndQuery(q, title)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer rows.Close()
	if len(title) > 0 {
		p.Id = title
	} else {
		p.Id = "0"
	}
	err, ok = s.LoadData(rows, &p.Title, &p.Fil.Path)
	if err == nil && ok {
		p.Fil.LoadFile()
	}
	s.RenderTemplate(w, "edit", p)
}

func (s *UncharServer) SaveHandler(w http.ResponseWriter, r *http.Request, title string) {
	var p Post
	var new_id string
	var ok bool

	p.Id = title
	p.Fil.Body = []byte(r.FormValue(POST_BODY_TAG))
	p.Title = r.FormValue(POST_TITLE_TAG)
	q := `UPDATE uncharblog.posts SET post_title='` + p.Title + `' WHERE post_id=` + p.Id + `;
				INSERT INTO uncharblog.posts (post_title)
							 SELECT '` + p.Title + `'
							 WHERE NOT EXISTS (SELECT 1 FROM uncharblog.posts WHERE post_id=` + p.Id + `)
							 RETURNING post_id;`
	rows, err := ExeIndQuery(q) //, p.Id, p.Title)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer rows.Close()
	if err, ok = s.LoadData(rows, &new_id); ok {
		p.Id = new_id
	}
	p.Fil.Path = POST_LOCAL_PATH + p.Id + ".txt"
	if ok {
		_, err = ExeIndQuery("UPDATE uncharblog.posts SET post_path=$2 WHERE post_id=$1", p.Id, p.Fil.Path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
	}
	err = p.Fil.SaveFile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, VIEW_TAG+p.Id, http.StatusFound)
}

func (s *UncharServer) IndexHandler(w http.ResponseWriter, r *http.Request, title string) {
	var i int
	var index Page

	rows, err := ExeIndQuery("SELECT * FROM uncharblog.posts")
	if err != nil {
		return
	}
	defer rows.Close()
	index.Title = "UncharBlog"
	index.Body = []byte("Awesome things.")
	index.List = make([]Post, 0, POST_LIMIT)
	i = 0
	for rows.Next() && i < POST_LIMIT {
		index.List = index.List[:(i + 1)]
		err := rows.Scan(&(index.List[i].Id), &(index.List[i].Title), &(index.List[i].Fil.Path))
		if err != nil {
			break
		}
		index.List[i].Fil.LoadFile()
		i++
	}
	s.RenderTemplate(w, "index", index)
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
