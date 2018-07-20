package storage

import (
	"time"
)

type FileInfo interface {
	Name() string
	Size() int64
	ModTime() time.Time
	MimeType() string
	Hash() string
}
