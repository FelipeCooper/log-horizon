package mlogapp

import (
	"context"
	"errors"
	"time"

	"github.com/felipecooper/log-horizon/app/sdk/proto/mlog"
	domain "github.com/felipecooper/log-horizon/business/domain/mlog"
	"github.com/felipecooper/log-horizon/foundation/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log  logger.Logger
	mlog *domain.Business
	mlog.UnimplementedLogWriterServer
	mlog.UnimplementedLogReaderServer
}

func NewApp(log logger.Logger, mlog *domain.Business) *App {
	return &App{
		log:  log,
		mlog: mlog,
	}
}

func (a *App) Register(ctx context.Context, req *mlog.NewLog) (*mlog.LogResponse, error) {
	log := NewLogFromProto(req)
	a.log.Info(ctx, "log received", "message", log.Message, "level", log.Level)

	domainLog, err := a.mlog.Register(ctx, log.Message, domain.Level(log.Level), log.Metadata)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidLevel) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "fail to register log")
	}

	return ToProtoResponse(domainLog), nil
}

func (a *App) Search(ctx context.Context, req *mlog.SearchQuery) (*mlog.Logs, error) {
	search := NewSearchFromProto(req)
	a.log.Info(ctx, "search request received",
		"startTime", search.StartTime,
		"endTime", search.EndTime,
		"level", search.Level,
	)

	result, err := a.mlog.Query(
		ctx,
		search.StartTime,
		search.EndTime,
		domain.Level(search.Level),
		search.Page,
		search.PageSize,
	)
	if err != nil {
		a.log.Error(ctx, "error searching logs", "error", err)
		if errors.Is(err, domain.ErrInvalidLevel) || errors.Is(err, domain.ErrInvalidTimeRange) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed on search logs")
	}

	return ToProtoLogs(result), nil
}

func (a *App) ExportToFile(ctx context.Context, req *mlog.SearchQuery) (*mlog.FileResponse, error) {
	search := NewSearchFromProto(req)
	a.log.Info(ctx, "export request received",
		"startTime", search.StartTime,
		"endTime", search.EndTime,
		"level", search.Level,
	)

	fileURL, fileSize, err := a.mlog.ExportToFile(
		ctx,
		search.StartTime,
		search.EndTime,
		domain.Level(search.Level),
	)
	if err != nil {
		a.log.Error(ctx, "error exporting logs to file", "error", err)
		if errors.Is(err, domain.ErrInvalidLevel) || errors.Is(err, domain.ErrInvalidTimeRange) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to export logs to file")
	}

	return ToProtoFileResponse(fileURL, fileSize), nil
}

func (a *App) StreamFile(req *mlog.SearchQuery, stream mlog.LogReader_StreamFileServer) error {
	ctx := stream.Context()
	search := NewSearchFromProto(req)
	a.log.Info(ctx, "stream request received",
		"startTime", search.StartTime,
		"endTime", search.EndTime,
		"level", search.Level,
	)

	pageSize := 100
	if search.PageSize > 0 && search.PageSize < pageSize {
		pageSize = search.PageSize
	}

	page := 0
	hasMore := true

	for hasMore {
		result, err := a.mlog.Query(
			ctx,
			search.StartTime,
			search.EndTime,
			domain.Level(search.Level),
			page,
			pageSize,
		)
		if err != nil {
			a.log.Error(ctx, "error streaming logs", "error", err, "page", page)
			if errors.Is(err, domain.ErrInvalidLevel) || errors.Is(err, domain.ErrInvalidTimeRange) {
				return status.Error(codes.InvalidArgument, err.Error())
			}
			return status.Error(codes.Internal, "failed on search logs")
		}

		if err := stream.Send(ToProtoLogs(result)); err != nil {
			a.log.Error(ctx, "error sending stream chunk", "error", err, "page", page)
			return status.Error(codes.Internal, "failed on send stream chunk")
		}

		hasMore = result.HasMore
		page = result.NextPage

		if hasMore {
			time.Sleep(10 * time.Millisecond)
		}
	}

	return nil
}
