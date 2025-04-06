package mongodb

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/felipecooper/log-horizon/business/domain/mlog"
	"github.com/felipecooper/log-horizon/foundation/compress"
	"github.com/felipecooper/log-horizon/foundation/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	log         logger.Logger
	db          *mongo.Database
	collection  *mongo.Collection
	compressor  compress.Compressor
	compression bool
	exportPath  string
}

type Config struct {
	DatabaseName     string
	CollectionName   string
	URI              string
	CompressionLevel int
	ExportPath       string
}

func NewStore(ctx context.Context, log logger.Logger, cfg Config) (*Store, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, fmt.Errorf("connecting to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("pinging MongoDB: %w", err)
	}

	collection := client.Database(cfg.DatabaseName).Collection(cfg.CollectionName)

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "timestamp", Value: 1},
			{Key: "level", Value: 1},
		},
		Options: options.Index().SetBackground(true),
	}

	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Error(ctx, "failed to create index", "error", err)
	}

	return &Store{
		log:         log,
		db:          client.Database(cfg.DatabaseName),
		collection:  collection,
		compressor:  compress.NewGzipCompressor(),
		compression: true,
		exportPath:  cfg.ExportPath,
	}, nil
}

func (s *Store) Write(ctx context.Context, log *mlog.Log) error {
	if s.compression && len(log.Message) > 100 {
		compressed, err := s.compressor.Compress([]byte(log.Message))
		if err != nil {
			s.log.Error(ctx, "failed to compress log message", "error", err)
		} else {
			log.Message = string(compressed)
			log.Compressed = true
			log.CompressedAt = time.Now()
		}
	}

	_, err := s.collection.InsertOne(ctx, log)
	if err != nil {
		s.log.Error(ctx, "failed to insert log in MongoDB", "error", err)
		return err
	}
	return nil
}

func (s *Store) Search(ctx context.Context, criteria mlog.SearchCriteria) (mlog.SearchResult, error) {
	filter := s.buildFilter(criteria)

	pageSize := criteria.PageSize
	if pageSize <= 0 {
		pageSize = 50
	}

	findOptions := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(pageSize)).
		SetSkip(int64(criteria.Page * pageSize))

	cursor, err := s.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return mlog.SearchResult{}, err
	}
	defer cursor.Close(ctx)

	var logs []mlog.Log
	for cursor.Next(ctx) {
		var log mlog.Log
		if err := cursor.Decode(&log); err != nil {
			continue
		}

		if log.Compressed {
			decompressed, err := s.compressor.Decompress([]byte(log.Message))
			if err == nil {
				log.Message = string(decompressed)
			}
		}

		logs = append(logs, log)
	}

	totalCount, err := s.Count(ctx, criteria)
	if err != nil {
		s.log.Error(ctx, "failed to count logs", "error", err)
	}

	hasMore := (criteria.Page+1)*pageSize < totalCount
	nextPage := criteria.Page + 1
	if !hasMore {
		nextPage = criteria.Page
	}

	return mlog.SearchResult{
		Logs:     logs,
		Total:    totalCount,
		HasMore:  hasMore,
		NextPage: nextPage,
	}, nil
}

func (s *Store) ExportToFile(ctx context.Context, criteria mlog.SearchCriteria) (string, int64, error) {
	filter := s.buildFilter(criteria)
	findOptions := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})

	cursor, err := s.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return "", 0, err
	}
	defer cursor.Close(ctx)

	filename := fmt.Sprintf("logs_export_%d.txt", time.Now().Unix())
	filepath := filepath.Join(s.exportPath, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	var size int64
	for cursor.Next(ctx) {
		var log mlog.Log
		if err := cursor.Decode(&log); err != nil {
			continue
		}

		if log.Compressed {
			decompressed, err := s.compressor.Decompress([]byte(log.Message))
			if err == nil {
				log.Message = string(decompressed)
			}
		}

		line := fmt.Sprintf("[%s] [%s] %s\n", log.Timestamp.Format(time.RFC3339), log.Level, log.Message)
		n, err := file.WriteString(line)
		if err != nil {
			s.log.Error(ctx, "error writing to export file", "error", err)
			continue
		}
		size += int64(n)
	}

	return filename, size, nil
}

func (s *Store) Count(ctx context.Context, criteria mlog.SearchCriteria) (int, error) {
	filter := s.buildFilter(criteria)
	count, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (s *Store) buildFilter(criteria mlog.SearchCriteria) bson.M {
	filter := bson.M{}

	timeFilter := bson.M{}
	if !criteria.TimeRange.StartTime.IsZero() {
		timeFilter["$gte"] = criteria.TimeRange.StartTime
	}
	if !criteria.TimeRange.EndTime.IsZero() {
		timeFilter["$lte"] = criteria.TimeRange.EndTime
	}
	if len(timeFilter) > 0 {
		filter["timestamp"] = timeFilter
	}

	if criteria.Level != "" {
		filter["level"] = criteria.Level
	}

	return filter
}

var _ mlog.Store = (*Store)(nil)
