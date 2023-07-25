package util

import "io"

type NamedReadSeeker interface {
	io.ReadSeeker
	Name() string
}
