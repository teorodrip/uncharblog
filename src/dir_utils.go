// dir_utils.go ---
//
// Filename: dir_utils.go
// Description:
// Author: Mateo Rodriguez Ripolles
// Maintainer:
// Created: ven. févr. 15 17:36:06 2019 (+0100)
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
	"strings"
)

func GetFilesFromDir(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	length := len(files)
	file_names := make([]string, length)
	for i := range files {
		file_names[i] = path + files[i].Name()
	}
	return file_names, nil
}

func SplitString(str, sep string) []string {
	n_sep := strings.Count(str, sep) + 1
	if n_sep == 1 && str == "" {
		return nil
	}
	ret := make([]string, n_sep)
	n_sep--
	i := 0
	for i < n_sep {
		j := 0
		m := strings.Index(str, sep)
		if m < 0 {
			break
		}
		for j < m && str[j:j+1] == " " {
			j++
		}
		ret[i] = str[j:m]
		str = str[m+len(sep):]
		i++
	}
	j := 0
	for j < len(str) && str[j:j+1] == " " {
		j++
	}
	ret[i] = str[j:]
	return ret[:i+1]
}

//
// dir_utils.go ends here
