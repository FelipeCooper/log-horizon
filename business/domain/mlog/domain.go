package mlog

import (
	"context"
)

type Writer interface {
	Write(ctx context.Context, log *Log) error
}

type Reader interface {
	Search(ctx context.Context, criteria SearchCriteria) (SearchResult, error)
	Count(ctx context.Context, criteria SearchCriteria) (int, error)
}

type Exporter interface {
	ExportToFile(ctx context.Context, criteria SearchCriteria) (string, int64, error)
}

type Store interface {
	Writer
	Reader
	Exporter
}
