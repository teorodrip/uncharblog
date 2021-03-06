### Makefile ---
##
## Filename: Makefile
## Description:
## Author: Mateo Rodriguez Ripolles
## Maintainer:
## Created: lun. févr. 18 09:25:47 2019 (+0100)
## Version:
## Package-Requires: ()
## Last-Updated:
##           By:
##     Update #: 0
## URL:
## Doc URL:
## Keywords:
## Compatibility:
##
######################################################################
##
### Commentary:
##
##
##
######################################################################
##
### Change Log:
##
##
######################################################################
##
## This program is free software: you can redistribute it and/or modify
## it under the terms of the GNU General Public License as published by
## the Free Software Foundation, either version 3 of the License, or (at
## your option) any later version.
##
## This program is distributed in the hope that it will be useful, but
## WITHOUT ANY WARRANTY; without even the implied warranty of
## MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
## General Public License for more details.
##
## You should have received a copy of the GNU General Public License
## along with GNU Emacs.  If not, see <https://www.gnu.org/licenses/>.
##
######################################################################
##
### Code:

.PHONY: all clean fclean re

SHELL = /bin/bash

NAME = uncharblog

CC = go

FUNCS = uncharblog.go \
	uncharserver.go \
	dir_utils.go \
	data_base.go

SRC_DIR = src/

all: $(NAME)

SRC = $(addprefix $(SRC_DIR), $(FUNCS))

$(NAME): $(SRC)
	@printf "Building uncharblog...\n"
	@$(CC) build -o $(NAME) $(SRC)
	@printf "Done.\n"

fclean:
	@rm -f $(NAME)

re: fclean
	@make

######################################################################
### Makefile ends here
