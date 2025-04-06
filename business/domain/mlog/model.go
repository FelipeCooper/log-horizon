package mlog

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Level string

const (
	Error Level = "error"
	Warn  Level = "warn"
	Debug Level = "debug"
	Info  Level = "info"
)

type Log struct {
	ID           ulid.ULID
	Message      string
	Timestamp    time.Time
	Level        Level
	Metadata     map[string]string
	Compressed   bool
	CompressedAt time.Time
}

type TimeRange struct {
	StartTime time.Time
	EndTime   time.Time
}

type SearchCriteria struct {
	TimeRange TimeRange
	Level     Level
	PageSize  int
	Page      int
}

type SearchResult struct {
	Logs     []Log
	Total    int
	HasMore  bool
	NextPage int
}

func (l Level) IsValid() bool {
	switch l {
	case Error, Warn, Debug, Info:
		return true
	}
	return false
}
