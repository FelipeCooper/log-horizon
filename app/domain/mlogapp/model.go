package mlogapp

import (
	"time"

	"github.com/felipecooper/log-horizon/app/sdk/proto/mlog"
	domain "github.com/felipecooper/log-horizon/business/domain/mlog"
)

type LogInput struct {
	Message   string
	Level     string
	Timestamp time.Time
	Metadata  map[string]string
}

func NewLogFromProto(proto *mlog.NewLog) LogInput {
	t := time.Unix(proto.Timestamp, 0)
	return LogInput{
		Message:   proto.Message,
		Level:     proto.Level,
		Timestamp: t,
		Metadata:  proto.Metadata,
	}
}

type SearchInput struct {
	StartTime time.Time
	EndTime   time.Time
	Level     string
	PageSize  int
	Page      int
	AsFile    bool
}

func NewSearchFromProto(proto *mlog.SearchQuery) SearchInput {
	return SearchInput{
		StartTime: time.Unix(proto.StartTime, 0),
		EndTime:   time.Unix(proto.EndTime, 0),
		Level:     proto.Level,
		PageSize:  int(proto.PageSize),
		Page:      int(proto.Page),
		AsFile:    proto.AsFile,
	}
}

func ToProtoLog(log domain.Log) *mlog.Log {
	return &mlog.Log{
		Id:        log.ID.String(),
		Message:   log.Message,
		Level:     string(log.Level),
		Timestamp: log.Timestamp.Unix(),
		Metadata:  log.Metadata,
	}
}

func ToProtoResponse(log domain.Log) *mlog.LogResponse {
	return &mlog.LogResponse{
		Id:     log.ID.String(),
		Status: "success",
	}
}

func ToProtoLogs(result domain.SearchResult) *mlog.Logs {
	protoLogs := make([]*mlog.Log, len(result.Logs))

	for i, log := range result.Logs {
		protoLogs[i] = ToProtoLog(log)
	}

	return &mlog.Logs{
		Logs:    protoLogs,
		Total:   int32(result.Total),
		HasMore: result.HasMore,
	}
}

func ToProtoFileResponse(fileURL string, fileSize int64) *mlog.FileResponse {
	return &mlog.FileResponse{
		FileUrl:     fileURL,
		FileSize:    fileSize,
		Compression: "gzip",
	}
}
