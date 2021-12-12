package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/doublen987/Projects/MySite/server/persistence/models"
)

var (
	ErrPostNotFound  = errors.New("Cannot find post")
	ErrInvalidPostId = errors.New("Invalid post id")
)

type MongoPost struct {
	ID                primitive.ObjectID `bson:"_id", json:"ID"`
	Title             string             `bson: "title,omitempty", json:"title"`
	Content           string             `bson: "content,omitempty", json:"content"`
	Description       string             `bson: "description,omitempty", json:"description"`
	Thumbnail         string             `bson: "thumbnail,omitempty, json: "thumbnail"`
	PublishTimestamp  time.Time          `bson:"publishtimestamp", json:"publishtimestamp"`
	LastEditTimestamp time.Time          `bson:"lastedittimestamp", json:"lastedittimestamp"`
	Hidden            bool               `bson:"hidden", json:"hidden"`
	Published         bool               `bson:"published", json:"published"`
	PublishDate       time.Time          `bson:"publishdate", json:"publishdate"`
	Tags              []string           `bson:"tags", json:"tags"`
}

type MongoProject struct {
	ID          primitive.ObjectID `bson:"_id", json:"ID"`
	Title       string             `bson: "title,omitempty", json:"title"`
	Link        string             `bson: "link,omitempty", json:"link"`
	Description string             `bson:"description,omitempty", json:"description"`
	Thumbnail   string             `bson:"thumbnail,omitempty", json:"description"`
	Timestamp   time.Time          `bson:"timestamp", json:"timestamp"`
	Tags        []string           `bson:"tags", json:"tags"`
}

type MongoLink struct {
	ID          primitive.ObjectID `bson:"_id", json:"ID"`
	Title       string             `bson: "title,omitempty", json:"title"`
	Description string             `bson: "description,omitempty", json:"description"`
}

type MongoKnowledgeTimelineEvent struct {
	ID          primitive.ObjectID `bson:"_id", json:"ID"`
	Title       string             `bson:"title,omitempty", json:"title"`
	Description string             `bson:"description,omitempty", json:"description"`
	Image       string             `bson:"image,omitempty", json:"image"`
}

type MongodbHandler struct {
	Session *mongo.Client
}

func NewMongodbHandler(connection string) (*MongodbHandler, error) {
	fmt.Println("Connecting to mongodb: " + connection)
	//s, err := mgo.Dial(connection)
	ctx := context.Background()
	s, err := mongo.NewClient(options.Client().ApplyURI(connection))
	if err != nil {
		fmt.Println(err)
		return &MongodbHandler{}, err
	}
	err = s.Connect(ctx)
	return &MongodbHandler{
		Session: s,
	}, err
}
func (handler *MongodbHandler) AddPost(ctx context.Context, post models.Post) error {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return err
	// }
	newID := primitive.NewObjectID()

	timeOfAddage := time.Now()

	newPost := MongoPost{
		ID:                newID,
		Title:             post.Title,
		Content:           post.Content,
		Description:       post.Description,
		Thumbnail:         post.Thumbnail,
		LastEditTimestamp: timeOfAddage,
		Hidden:            post.Hidden,
		Published:         false,
		Tags:              post.Tags,
	}

	_, err := s.Database("MojSajt").Collection("posts").InsertOne(ctx, newPost)
	return err
}
func (handler *MongodbHandler) UpdatePost(ctx context.Context, post models.Post) (models.Post, error) {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return models.Post{}, err
	// }
	var updateFields bson.D

	var fields map[string]interface{}
	fields = make(map[string]interface{})
	fields["title"] = post.Title
	fields["content"] = post.Content
	fields["description"] = post.Description
	fields["thumbnail"] = post.Thumbnail
	fields["lastedittimestamp"] = time.Now()
	fields["hidden"] = post.Hidden

	for key, value := range fields {
		updateFields = append(updateFields, bson.E{key, value})
	}

	fmt.Println(primitive.ObjectIDFromHex(post.ID))
	if id, err := primitive.ObjectIDFromHex(post.ID); err == nil {
		_, err := s.Database("MojSajt").Collection("posts").UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set", updateFields}})

		return models.Post{}, err
	} else {
		return models.Post{}, ErrInvalidPostId
	}
}
func (handler *MongodbHandler) ReplacePost(ctx context.Context, post models.Post) error {
	s := handler.Session

	fmt.Println(primitive.ObjectIDFromHex(post.ID))
	if id, err := primitive.ObjectIDFromHex(post.ID); err == nil {

		newPost := MongoPost{
			ID:                id,
			Title:             post.Title,
			Content:           post.Content,
			Description:       post.Description,
			LastEditTimestamp: time.Now(),
			Hidden:            post.Hidden,
			Published:         post.Published,
			Thumbnail:         post.Thumbnail,
		}
		_, err := s.Database("MojSajt").Collection("posts").ReplaceOne(ctx, bson.D{{"_id", id}}, newPost)
		return err
	} else {
		return ErrInvalidPostId
	}
}
func (handler *MongodbHandler) RemovePost(ctx context.Context, id string) error {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return err
	// }

	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		_, err := s.Database("MojSajt").Collection("posts").DeleteOne(ctx, bson.D{{"_id", oid}})
		return err
	} else {
		return err
	}
}
func (handler *MongodbHandler) PublishPost(ctx context.Context, id string) error {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return err
	// }
	timestamp := time.Now()

	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		_, err := s.Database("MojSajt").Collection("posts").UpdateOne(ctx, bson.D{{"_id", oid}}, bson.D{{"$set",
			bson.M{
				"published":         true,
				"publishtimestamp":  timestamp,
				"lastedittimestamp": timestamp,
			}}})
		return err
	} else {
		return err
	}
}
func (handler *MongodbHandler) GetPosts(ctx context.Context) ([]models.Post, error) {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)

	mposts := []MongoPost{}
	posts := []models.Post{}
	cursor, err := s.Database("MojSajt").Collection("posts").Find(ctx, bson.D{})
	if err != nil {
		return posts, err
	}
	fmt.Println(time.Now().Format("January-02-2006"))
	cursor.All(ctx, &mposts)
	defer cursor.Close(ctx)
	for _, mpost := range mposts {
		posts = append(posts, models.Post{
			ID:                mpost.ID.Hex(),
			Title:             mpost.Title,
			Content:           mpost.Content,
			Description:       mpost.Description,
			Thumbnail:         mpost.Thumbnail,
			PublishTimestamp:  mpost.PublishTimestamp.Format("January-02-2006"),
			LastEditTimestamp: mpost.LastEditTimestamp.Format("January-02-2006"),
			Published:         mpost.Published,
			Hidden:            mpost.Hidden,
		})
	}
	return posts, err
}
func (handler *MongodbHandler) GetPost(ctx context.Context, id string) (models.Post, error) {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return models.Post{}, err
	// }

	post := models.Post{}
	mpost := MongoPost{}
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		cursor, err := s.Database("MojSajt").Collection("posts").Find(ctx, bson.D{{"_id", oid}})
		defer cursor.Close(ctx)
		if cursor.Next(ctx) {
			err = cursor.Decode(&mpost)
			if err != nil {
				return models.Post{}, err
			}
			post = models.Post{
				ID:               mpost.ID.Hex(),
				Title:            mpost.Title,
				Content:          mpost.Content,
				Description:      mpost.Description,
				Thumbnail:        mpost.Thumbnail,
				PublishTimestamp: mpost.PublishTimestamp.Format("January-02-2006"),
			}
			return post, nil
		}
		return models.Post{}, ErrPostNotFound
	} else {
		return models.Post{}, ErrInvalidPostId
	}
}
func (handler *MongodbHandler) AddProject(ctx context.Context, project models.Project) error {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return err
	// }

	newID := primitive.NewObjectID()

	newProject := MongoProject{
		ID:          newID,
		Title:       project.Title,
		Link:        project.Link,
		Description: project.Description,
		Thumbnail:   project.Thumbnail,
		Timestamp:   time.Now(),
		Tags:        project.Tags,
	}

	_, err := s.Database("MojSajt").Collection("projects").InsertOne(ctx, newProject)
	return err
}
func (handler *MongodbHandler) UpdateProject(ctx context.Context, project models.Project) (models.Project, error) {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return models.Project{}, err
	// }

	fmt.Println(primitive.ObjectIDFromHex(project.ID))
	if oid, err := primitive.ObjectIDFromHex(project.ID); err == nil {
		updatedProject := MongoProject{
			ID:          oid,
			Title:       project.Title,
			Link:        project.Link,
			Description: project.Description,
			Thumbnail:   project.Thumbnail,
		}

		_, err = s.Database("MojSajt").Collection("projects").ReplaceOne(ctx, bson.D{{"_id", updatedProject.ID}}, updatedProject)
		return project, err
	} else {
		return models.Project{}, ErrInvalidPostId
	}
}
func (handler *MongodbHandler) RemoveProject(ctx context.Context, id string) error {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return err
	// }

	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		_, err := s.Database("MojSajt").Collection("projects").DeleteOne(ctx, bson.D{{"_id", oid}})
		return err
	} else {
		return err
	}
}
func (handler *MongodbHandler) GetProjects(ctx context.Context) ([]models.Project, error) {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return []models.Project{}, err
	// }

	mprojects := []MongoProject{}
	projects := []models.Project{}
	cursor, err := s.Database("MojSajt").Collection("projects").Find(ctx, bson.D{})
	if err != nil {
		return projects, err
	}
	cursor.All(ctx, &mprojects)
	cursor.Close(ctx)

	for _, mproject := range mprojects {
		projects = append(projects, models.Project{
			ID:          mproject.ID.Hex(),
			Title:       mproject.Title,
			Link:        mproject.Link,
			Description: mproject.Description,
			Thumbnail:   mproject.Thumbnail,
			Timestamp:   mproject.Timestamp.Format("Jan/06/15"),
		})
	}
	return projects, err
}
func (handler *MongodbHandler) GetProject(ctx context.Context, id string) (models.Project, error) {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return models.Project{}, err
	// }

	p := models.Project{}
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		cursor, err := s.Database("MojSajt").Collection("projects").Find(ctx, bson.D{{"_id", oid}})
		defer cursor.Close(ctx)
		if err != nil {
			return models.Project{}, err
		}
		if cursor.Next(ctx) {
			err = cursor.Decode(&p)
			if err != nil {
				return models.Project{}, err
			}
			return p, nil
		}
		return models.Project{}, ErrPostNotFound
	} else {
		return models.Project{}, ErrInvalidPostId
	}
}
func (handler *MongodbHandler) GetLinks(ctx context.Context) ([]models.Link, error) {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return []models.Link{}, err
	// }

	mlinks := []MongoLink{}
	links := []models.Link{}
	cursor, err := s.Database("MojSajt").Collection("posts").Find(ctx, bson.D{})
	cursor.All(ctx, &mlinks)
	cursor.Close(ctx)
	for _, mlink := range mlinks {
		links = append(links, models.Link{ID: mlink.ID.Hex(), Title: mlink.Title, Description: mlink.Description})
	}
	return links, err
}
func (handler *MongodbHandler) GetKnowledgeTimelineEvents(ctx context.Context) ([]models.TimelineEvent, error) {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return []models.Project{}, err
	// }

	mevents := []MongoKnowledgeTimelineEvent{}
	events := []models.TimelineEvent{}
	cursor, err := s.Database("MojSajt").Collection("knowledgetimelineevents").Find(ctx, bson.D{})
	cursor.All(ctx, &mevents)
	cursor.Close(ctx)
	for _, mevent := range mevents {
		events = append(events, models.TimelineEvent{ID: mevent.ID.Hex(), Title: mevent.Title, Description: mevent.Description, Image: mevent.Image})
	}
	return events, err
}
func (handler *MongodbHandler) AddKnowledgeTimelineEvent(ctx context.Context, event models.TimelineEvent) error {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return err
	// }

	newID := primitive.NewObjectID()

	newEvent := MongoKnowledgeTimelineEvent{
		ID:          newID,
		Title:       event.Title,
		Description: event.Description,
		Image:       event.Image,
	}

	_, err := s.Database("MojSajt").Collection("knowledgetimelineevents").InsertOne(ctx, newEvent)
	return err
}
func (handler *MongodbHandler) UpdateKnowledgeTimelineEvent(ctx context.Context, event models.TimelineEvent) (models.TimelineEvent, error) {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return models.Project{}, err
	// }

	fmt.Println(primitive.ObjectIDFromHex(event.ID))
	if oid, err := primitive.ObjectIDFromHex(event.ID); err == nil {
		updatedEvent := MongoKnowledgeTimelineEvent{
			ID:          oid,
			Title:       event.Title,
			Description: event.Description,
			Image:       event.Image,
		}

		_, err = s.Database("MojSajt").Collection("knowledgetimelineevents").UpdateOne(ctx, bson.D{{"_id", updatedEvent.ID}}, updatedEvent)
		return event, err
	} else {
		return models.TimelineEvent{}, ErrInvalidPostId
	}
}
func (handler *MongodbHandler) RemoveKnowledgeTimelineEvent(ctx context.Context, id string) error {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return err
	// }

	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		_, err := s.Database("MojSajt").Collection("knowledgetimelineevents").DeleteOne(ctx, bson.D{{"_id", oid}})
		return err
	} else {
		return err
	}
}
