package storage

import "errors"

var (
	ErrNotFound   = errors.New("no url for this alias")
	ErrAliasExist = errors.New("this alias is already exists")
)
