// data_base.go ---
//
// Filename: data_base.go
// Description:
// Author: Mateo Rodriguez Ripolles
// Maintainer:
// Created: lun. févr. 18 16:25:03 2019 (+0100)
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

// postgre tag array function

// CREATE FUNCTION get_tag_arr(integer) RETURNS text[] AS $$
// SELECT ARRAY(SELECT t.tag_name FROM uncharblog.tags AS t
//	INNER JOIN uncharblog.post_tag_ref AS pt ON t.tag_id = pt.tag_id
//	INNER JOIN uncharblog.posts AS p ON p.post_id = pt.post_id
//	WHERE p.post_id in ($1))
// $$ LANGUAGE SQL;

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

const DB_USER = "uncharblog"
const DB_NAME = "pam_test"
const DB_PASS = "K5N3gwww5U8Yxfcv"
const DB_HOST = "192.168.27.122"
const CONN_STR = "dbname=" + DB_NAME +
	" user=" + DB_USER +
	" password=" + DB_PASS +
	" host=" + DB_HOST

type pgDB struct {
	Db                         *sql.DB
	SqlGetAllPosts             *sql.Stmt
	SqlGetPost                 *sql.Stmt
	SqlUpdateAddPost           *sql.Stmt
	SqlUpdatePost              *sql.Stmt
	SqlCreateTagNameByPostFunc *sql.Stmt
	SqlCreateTagIdByPostFunc   *sql.Stmt
	SqlPostsByTag              *sql.Stmt
	SqlDelPostTags             *sql.Stmt
	SqlAddPostTags             *sql.Stmt
	SqlAddTag                  *sql.Stmt
}

func ConnectDB(ConnStr string) (*pgDB, error) {
	db, err := sql.Open("postgres", ConnStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to %s\n", ConnStr)
		log.Print(err)
		return nil, err
	} else {
		p := &pgDB{Db: db}
		if err = p.Db.Ping(); err != nil {
			return nil, err
		}
		if err = p.PrepareSqlStatements(); err != nil {
			return nil, err
		}
		if _, err = p.SqlCreateTagNameByPostFunc.Exec(); err != nil {
			return nil, err
		}
		if _, err = p.SqlCreateTagIdByPostFunc.Exec(); err != nil {
			return nil, err
		}
		return p, nil
	}
}

func (p *pgDB) PrepareSqlStatements() error {
	var err error

	if p.SqlGetAllPosts, err = p.Db.Prepare("SELECT post_id, post_title, post_path, TO_CHAR(creation_date, 'dd-mon-YYYY'), TO_CHAR(update_date, 'dd-mon-YYYY'), post_link, get_tag_id_by_post(post_id), get_tag_name_by_post(post_id) FROM uncharblog.posts ORDER BY creation_date DESC NULLS LAST LIMIT $1"); err != nil {
		return err
	}
	if p.SqlGetPost, err = p.Db.Prepare("SELECT post_id, post_title, post_path, TO_CHAR(creation_date, 'dd-mon-YYYY'), TO_CHAR(update_date, 'dd-mon-YYYY'), post_link, get_tag_id_by_post(post_id), get_tag_name_by_post(post_id) FROM uncharblog.posts WHERE post_id=$1"); err != nil {
		return err
	}
	if p.SqlUpdateAddPost, err = p.Db.Prepare("with updated as (UPDATE uncharblog.posts SET post_title=$2, update_date=$3, post_link=$4 WHERE post_id=$1) INSERT INTO uncharblog.posts (post_title, creation_date, update_date, post_link) SELECT $2, $3, $3, $4 WHERE NOT EXISTS (SELECT 1 FROM uncharblog.posts WHERE post_id=$1) RETURNING post_id;"); err != nil {
		return err
	}
	if p.SqlUpdatePost, err = p.Db.Prepare("UPDATE uncharblog.posts SET post_path=$2 WHERE post_id=$1"); err != nil {
		return err
	}
	if p.SqlCreateTagNameByPostFunc, err = p.Db.Prepare("CREATE OR REPLACE FUNCTION get_tag_name_by_post(integer) RETURNS text[] AS $$ SELECT ARRAY(SELECT t.tag_name FROM uncharblog.tags AS t INNER JOIN uncharblog.post_tag_ref AS pt ON t.tag_id = pt.tag_id INNER JOIN uncharblog.posts AS p ON p.post_id = pt.post_id WHERE p.post_id in ($1)) $$ LANGUAGE SQL;"); err != nil {
		return err
	}
	if p.SqlCreateTagIdByPostFunc, err = p.Db.Prepare("CREATE OR REPLACE FUNCTION get_tag_id_by_post(integer) RETURNS int[] AS $$ SELECT ARRAY(SELECT t.tag_id FROM uncharblog.tags AS t INNER JOIN uncharblog.post_tag_ref AS pt ON t.tag_id = pt.tag_id INNER JOIN uncharblog.posts AS p ON p.post_id = pt.post_id	WHERE p.post_id in ($1)) $$ LANGUAGE SQL;"); err != nil {
		return err
	}
	if p.SqlPostsByTag, err = p.Db.Prepare("SELECT t.post_id, t.post_title, t.post_path, TO_CHAR(t.creation_date, 'dd-mon-YYYY'), TO_CHAR(t.update_date, 'dd-mon-YYYY'), post_link, get_tag_id_by_post(t.post_id), get_tag_name_by_post(t.post_id) FROM uncharblog.posts AS t INNER JOIN uncharblog.post_tag_ref AS pt ON t.post_id = pt.post_id INNER JOIN uncharblog.tags AS p ON p.tag_id = pt.tag_id WHERE p.tag_id in ($1) ORDER BY t.creation_date DESC NULLS LAST LIMIT $2"); err != nil {
		return err
	}
	if p.SqlDelPostTags, err = p.Db.Prepare("DELETE FROM uncharblog.post_tag_ref WHERE post_id = $1"); err != nil {
		return err
	}
	if p.SqlAddPostTags, err = p.Db.Prepare("INSERT INTO uncharblog.post_tag_ref VALUES ($1, $2)"); err != nil {
		return err
	}
	if p.SqlAddTag, err = p.Db.Prepare("WITH s AS (SELECT tag_id FROM uncharblog.tags WHERE tag_name = $1), i AS (INSERT INTO uncharblog.tags(tag_name) SELECT $1 WHERE NOT EXISTS (SELECT 1 FROM s) RETURNING tag_id) SELECT tag_id FROM i UNION ALL SELECT tag_id FROM s"); err != nil {
		return err
	}
	return nil
}

//
// data_base.go ends here
