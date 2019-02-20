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
	"log"
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
	Db        *pgDB
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
	if server.Db, err = ConnectDB(CONN_STR); err != nil {
		return nil, err
	}
	return server, nil
}

func (s *UncharServer) Start() {
	http.Handle(STYLE_SHEETS_URL_PATH, http.StripPrefix(STYLE_SHEETS_URL_PATH, http.FileServer(http.Dir(STYLE_SHEETS_LOCAL_PATH))))
	http.Handle(FONTS_URL_PATH, http.StripPrefix(FONTS_URL_PATH, http.FileServer(http.Dir(FONTS_LOCAL_PATH))))
	http.HandleFunc(VIEW_TAG, s.MakeHandler(s.ViewHandler))
	http.HandleFunc(EDIT_TAG, s.MakeHandler(s.EditHandler))
	http.HandleFunc(SAVE_TAG, s.MakeHandler(s.SaveHandler))
	http.HandleFunc(INDEX_TAG, s.MakeHandler(s.IndexHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
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

func (s *UncharServer) LoadRow(rows *sql.Rows, values ...interface{}) (error, bool) {
	var err error
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

	err := s.Db.SqlGetPost.QueryRow(title).Scan(&p.Id, &p.Title, &p.Fil.Path, &p.CreationDate, &p.UpdateDate)
	if err != nil {
		http.Redirect(w, r, EDIT_TAG+title, http.StatusFound)
		return
	}
	p.Fil.LoadFile()
	s.RenderTemplate(w, "view", p)
}

func (s *UncharServer) EditHandler(w http.ResponseWriter, r *http.Request, title string) {
	p := Post{Id: "0"}
	err := s.Db.SqlGetPost.QueryRow(title).Scan(&p.Id, &p.Title, &p.Fil.Path, &p.CreationDate, &p.UpdateDate)
	if err != nil && err != sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err == nil {
		p.Fil.LoadFile()
	}
	s.RenderTemplate(w, "edit", p)
}

func (s *UncharServer) SaveHandler(w http.ResponseWriter, r *http.Request, title string) {
	var p Post

	p.Id = title
	p.Fil.Body = []byte(r.FormValue(POST_BODY_TAG))
	p.Title = r.FormValue(POST_TITLE_TAG)
	err := s.Db.SqlUpdateAddPost.QueryRow(p.Id, p.Title).Scan(&p.Id)
	if err != nil && err != sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	p.Fil.Path = POST_LOCAL_PATH + p.Id + ".txt"
	if err == nil {
		_, err = s.Db.SqlUpdatePost.Exec(p.Id, p.Fil.Path)
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

	rows, err := s.Db.SqlGetAllPosts.Query(POST_LIMIT)
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
		err := rows.Scan(&(index.List[i].Id), &(index.List[i].Title), &(index.List[i].Fil.Path), &(index.List[i].CreationDate), &(index.List[i].UpdateDate))
		fmt.Printf("Date: %s\n", index.List[i].CreationDate)
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
