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

func ConnectDB(ConnStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", ConnStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to %s\n", ConnStr)
		log.Print(err)
		return nil, err
	}
	return db, nil
}

func ExeIndQuery(Query string, args ...interface{}) (*sql.Rows, error) {
	ConnStr := "dbname=" + DB_NAME +
		" user=" + DB_USER +
		" password=" + DB_PASS +
		" host=" + DB_HOST
	db, err := ConnectDB(ConnStr)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(Query, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing the query:\n%s\nWith args: %s\n", Query, args)
		log.Print(err)
		return nil, err
	}
	return rows, nil
}

//
// data_base.go ends here
