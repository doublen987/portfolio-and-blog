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
	AddUser(ctx context.Context, user models.User) error
	GetUsers(ctx context.Context) ([]models.User, error)
	UpdateUser(ctx context.Context, user models.User) error
	RemoveUser(ctx context.Context, userID string) error
	Authenticate(ctx context.Context, username string, password string) (bool, error)
	GetPosts(ctx context.Context) ([]models.Post, error)
	GetPostsCount(ctx context.Context) (int64, error)
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
	AddTag(ctx context.Context, tag models.Tag) error
	RemoveTag(ctx context.Context, tagID string) error
	UpdateTag(ctx context.Context, tag models.Tag) error
	GetTags(ctx context.Context) ([]models.Tag, error)
	AddPage(ctx context.Context, page models.Page) error
	UpdatePage(ctx context.Context, page models.Page) error
	RemovePage(ctx context.Context, id string) error
	GetPage(ctx context.Context, id string) (models.Page, error)
	GetPages(ctx context.Context) ([]models.Page, error)
	AddVisit(ctx context.Context, visit models.Visit) error
	GetVisits(ctx context.Context) ([]models.VisitSummary, error)
	AddSettings(ctx context.Context, settings models.Settings) error
	UpdateSettings(ctx context.Context, settings models.Settings) error
	GetSettings(ctx context.Context) (models.Settings, error)
}

var DBTypeNotSupported = errors.New("The Database type provided is not supported...")

func GetDataBaseHandler(dbtype uint8, connection string) (DBHandler, error) {
	switch dbtype {
	case MONGODB:
		return mongodb.NewMongodbHandler(connection)
	}
	return nil, DBTypeNotSupported
}
