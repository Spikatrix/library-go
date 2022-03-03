package db

import "errors"

var (
	ErrBookNotFound     = errors.New("book not found")
	ErrBookDecodeFailed = errors.New("book decode failed")
	ErrBookQueryFailed  = errors.New("book query failed")
	ErrBookInsertFailed = errors.New("book insertion failed")
)
