package models

import "errors"

var ErrBadConfigPort = errors.New("port must be upper than 0")
