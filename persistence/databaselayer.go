package persistence

import (
	"context"
	"errors"

	"github.com/doublen987/Projects/MySite/server/persistence/models"
	"github.com/doublen987/Projects/MySite/server/persistence/mongodb"
)

const (
	MYSQL uint8 = iota
	SQLITE
	POSTGRESQL
	MONGODB
)

// const (
// 	FILESYSTEM uint8 = iota
// )

type DBTYPE string
type STORAGETYPE string

type DBHandler interface {
	GetLinks(ctx context.Context) ([]models.Link, error)
	GetPosts(ctx context.Context) ([]models.Post, error)
	GetPost(ctx context.Context, id string) (models.Post, error)
	AddPost(ctx context.Context, post models.Post) error
	UpdatePost(ctx context.Context, post models.Post) (models.Post, error)
	ReplacePost(ctx context.Context, post models.Post) error
	RemovePost(ctx context.Context, id string) error
	PublishPost(ctx context.Context, id string) error
	GetProjects(ctx context.Context) ([]models.Project, error)
	GetProject(ctx context.Context, id string) (models.Project, error)
	AddProject(ctx context.Context, post models.Project) error
	UpdateProject(ctx context.Context, post models.Project) (models.Project, error)
	RemoveProject(ctx context.Context, id string) error
	GetKnowledgeTimelineEvents(ctx context.Context) ([]models.TimelineEvent, error)
	AddKnowledgeTimelineEvent(ctx context.Context, event models.TimelineEvent) error
	UpdateKnowledgeTimelineEvent(ctx context.Context, event models.TimelineEvent) (models.TimelineEvent, error)
	RemoveKnowledgeTimelineEvent(ctx context.Context, id string) error
}

var DBTypeNotSupported = errors.New("The Database type provided is not supported...")

func GetDataBaseHandler(dbtype uint8, connection string) (DBHandler, error) {
	switch dbtype {
	case MONGODB:
		return mongodb.NewMongodbHandler(connection)
	}
	return nil, DBTypeNotSupported
}
