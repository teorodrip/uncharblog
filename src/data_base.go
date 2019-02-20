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
	Db               *sql.DB
	SqlGetAllPosts   *sql.Stmt
	SqlGetPost       *sql.Stmt
	SqlUpdateAddPost *sql.Stmt
	SqlUpdatePost    *sql.Stmt
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
		return p, nil
	}
}

func (p *pgDB) PrepareSqlStatements() error {
	var err error

	if p.SqlGetAllPosts, err = p.Db.Prepare("SELECT post_id, post_title, post_path, TO_CHAR(creation_date, 'dd-mon-YYYY'), TO_CHAR(update_date, 'dd-mon-YYYY') FROM uncharblog.posts ORDER BY creation_date DESC NULLS LAST LIMIT $1"); err != nil {
		return err
	}
	if p.SqlGetPost, err = p.Db.Prepare("SELECT post_id, post_title, post_path, TO_CHAR(creation_date, 'dd-mon-YYYY'), TO_CHAR(update_date, 'dd-mon-YYYY') FROM uncharblog.posts WHERE post_id=$1"); err != nil {
		return err
	}
	if p.SqlUpdateAddPost, err = p.Db.Prepare("with updated as (UPDATE uncharblog.posts SET post_title=$2, update_date=$3 WHERE post_id=$1) INSERT INTO uncharblog.posts (post_title, creation_date, update_date) SELECT $2, $3, $3 WHERE NOT EXISTS (SELECT 1 FROM uncharblog.posts WHERE post_id=$1) RETURNING post_id;"); err != nil {
		return err
	}
	if p.SqlUpdatePost, err = p.Db.Prepare("UPDATE uncharblog.posts SET post_path=$2 WHERE post_id=$1"); err != nil {
		return err
	}
	return nil
}

//
// data_base.go ends here
