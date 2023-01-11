package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"

	"github.com/doublen987/Projects/MySite/server/persistence/models"
	"github.com/rs/xid"
)

var (
	ErrPostNotFound  = errors.New("Cannot find post")
	ErrInvalidPostId = errors.New("Invalid post id")
)

type MongoUser struct {
	ID          primitive.ObjectID `bson:"_id", json:"ID"`
	Username    string             `bson:"username", json:"username"`
	Password    string             `bson:"password", json:"password"`
	Description string             `bson:"description", json:"description"`
	Thumbnail   string             `bson:"thumbnail", json:"thumbnail"`
}

type MongoPost struct {
	ID                 primitive.ObjectID `bson:"_id", json:"ID"`
	Title              string             `bson: "title,omitempty", json:"title"`
	Content            string             `bson: "content,omitempty", json:"content"`
	Description        string             `bson: "description,omitempty", json:"description"`
	Thumbnail          string             `bson: "thumbnail,omitempty, json: "thumbnail"`
	ThumbnailStretched bool               `bson:"thumbnailstretched", json:"thumbnailstretched"`
	PublishTimestamp   time.Time          `bson:"publishtimestamp", json:"publishtimestamp"`
	LastEditTimestamp  time.Time          `bson:"lastedittimestamp", json:"lastedittimestamp"`
	Hidden             bool               `bson:"hidden", json:"hidden"`
	Published          bool               `bson:"published", json:"published"`
	PublishDate        time.Time          `bson:"publishdate", json:"publishdate"`
	Tags               []string           `bson:"tags,omitempty", json:"tags,omitempty"`
}

type MongoProject struct {
	ID                 primitive.ObjectID `bson:"_id" json:"ID"`
	Title              string             `bson:"title" json:"title"`
	Link               string             `bson:"link" json:"link"`
	Description        string             `bson:"description" json:"description"`
	Thumbnail          string             `bson:"thumbnail" json:"thumbnail"`
	ThumbnailStretched bool               `bson:"thumbnailstretched" json:"thumbnailstretched"`
	Timestamp          time.Time          `bson:"timestamp" json:"timestamp"`
	Tags               []string           `bson:"tags,omitempty", json:"tags,omitempty"`
}

type MongoLink struct {
	ID          primitive.ObjectID `bson:"_id" json:"ID"`
	Title       string             `bson: "title,omitempty" json:"title"`
	Description string             `bson: "description,omitempty" json:"description"`
}

type MongoKnowledgeTimelineEvent struct {
	ID          primitive.ObjectID `bson:"_id" json:"ID"`
	Title       string             `bson:"title,omitempty" json:"title"`
	Description string             `bson:"description,omitempty" json:"description"`
	Image       string             `bson:"image,omitempty" json:"image"`
}

type MongoPage struct {
	ID       string                   `bson:"ID" json:"ID"`
	Name     string                   `bson:"name" json:"name"`
	Homepage bool                     `bson:"homepage" json:"homepage"`
	Sections []map[string]interface{} `bson:"sections" json:"section"`
}

type MongodbHandler struct {
	Session      *mongo.Client
	SettingsID   string
	DatabaseName string
}

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

type nullawareStrDecoder struct{}

func (nullawareStrDecoder) DecodeValue(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Kind() != reflect.String {
		return errors.New("bad type or not settable")
	}
	var str string
	var err error
	switch vr.Type() {
	case bsontype.String:
		if str, err = vr.ReadString(); err != nil {
			return err
		}
	case bsontype.Null: // THIS IS THE MISSING PIECE TO HANDLE NULL!
		if err = vr.ReadNull(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("cannot decode %v into a string type", vr.Type())
	}

	val.SetString(str)
	return nil
}

func NewMongodbHandler(connection string, databaseName string) (*MongodbHandler, error) {
	fmt.Println("Connecting to mongodb: " + connection)
	//s, err := mgo.Dial(connection)
	ctx := context.Background()
	s, err := mongo.NewClient(options.Client().ApplyURI(connection).SetRegistry(
		bson.NewRegistryBuilder().RegisterTypeDecoder(reflect.TypeOf(""), nullawareStrDecoder{}).Build(),
	))
	if err != nil {
		fmt.Println(err)
		return &MongodbHandler{}, err
	}
	err = s.Connect(ctx)
	if err != nil {
		fmt.Println(err)
		return &MongodbHandler{}, err
	}

	currentTime := time.Now()
	connectionEstablished := false
	for !connectionEstablished {
		diff := time.Since(currentTime)
		if diff.Seconds() > 1 {
			fmt.Println("Trying to connect to database")
			myCtx, _ := context.WithTimeout(context.Background(), 3*time.Second)
			err = s.Ping(myCtx, readpref.Primary())
			if err == nil {
				connectionEstablished = true
				fmt.Println("Connected to database")
			}
			currentTime = time.Now()
		}
	}

	handler := MongodbHandler{
		Session:      s,
		SettingsID:   "settings",
		DatabaseName: databaseName,
	}

	if err == nil {
		if count, err := handler.CountElements(context.Background(), "settings"); err == nil && count > 0 {
			handler.AddSettings(context.Background(), models.Settings{
				WebsiteName:      "MySite",
				Logo:             "",
				BackgroundColor1: "#26021d",
				BackgroundColor2: "#33112b",
				BackgroundColor3: "#58234c",
				TextColor1:       "#d08a12",
				TextColor2:       "#fbc23f",
				TextColor3:       "#ffffff",
			})
		}

		if count, err := handler.CountElements(context.Background(), "tags"); err == nil && count == 0 {
			handler.RemoveAllTags(context.Background())
			handler.AddTag(context.Background(), models.Tag{
				Name:               "MongoDB",
				Thumbnail:          "MongoDB-logo.svg?storagetype=filesystem",
				ThumbnailStretched: false,
			})
		}

		if count, err := handler.CountElements(context.Background(), "pages"); err == nil && count == 0 {
			handler.RemoveAll(context.Background(), "pages")
			sections := make([]interface{}, 2)

			bla := make(map[string]string)
			bla["type"] = "image"
			bla["filename"] = "Eternal-logo.svg?storagetype=filesystem"
			sections[0] = bla

			bla1 := make(map[string]string)
			bla1["type"] = "text"
			bla1["header"] = "Welocme to Eternal!"
			bla1["content"] = "The best Content Management System for building portfolios."
			sections[1] = bla1
			handler.AddPage(context.Background(), models.Page{
				Name:     "Welcome to Eternal",
				Homepage: true,
				Sections: sections,
			})
		}
	}

	return &handler, err
}

func (handler *MongodbHandler) CountElements(ctx context.Context, collectionName string) (int64, error) {
	s := handler.Session
	count, err := s.Database(handler.DatabaseName).Collection(collectionName).CountDocuments(ctx, bson.D{})
	return count, err
}
func (handler *MongodbHandler) RemoveAll(ctx context.Context, collectionName string) error {
	s := handler.Session

	err := s.Database(handler.DatabaseName).Collection(collectionName).Drop(ctx)
	return err
}
func (handler *MongodbHandler) AddUser(ctx context.Context, user models.User) error {
	s := handler.Session

	newID := primitive.NewObjectID()

	newUser := MongoUser{
		ID:          newID,
		Username:    user.Username,
		Password:    getHash([]byte(user.Password)),
		Description: user.Description,
		Thumbnail:   user.Thumbnail,
	}

	_, err := s.Database(handler.DatabaseName).Collection("users").InsertOne(ctx, newUser)
	return err
}
func (handler *MongodbHandler) RemoveUser(ctx context.Context, userID string) error {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return err
	// }

	if oid, err := primitive.ObjectIDFromHex(userID); err == nil {
		_, err := s.Database(handler.DatabaseName).Collection("users").DeleteOne(ctx, bson.D{{"_id", oid}})
		return err
	} else {
		return err
	}
}
func (handler *MongodbHandler) UpdateUser(ctx context.Context, user models.User) error {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return models.Post{}, err
	// }
	var updateFields bson.D

	var fields map[string]interface{}
	fields = make(map[string]interface{})
	fields["username"] = user.Username
	if user.Password != "" {
		fields["password"] = getHash([]byte(user.Password))
	}
	fields["description"] = user.Description
	fields["thumbnail"] = user.Thumbnail

	for key, value := range fields {
		updateFields = append(updateFields, bson.E{key, value})
	}

	fmt.Println(primitive.ObjectIDFromHex(user.ID))
	if id, err := primitive.ObjectIDFromHex(user.ID); err == nil {
		_, err := s.Database(handler.DatabaseName).Collection("users").UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set", updateFields}})

		return err
	} else {
		return ErrInvalidPostId
	}
}
func (handler *MongodbHandler) GetUsers(ctx context.Context) ([]models.User, error) {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)

	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "Username", Value: 1}})

	if searchTerm := ctx.Value("search-term"); searchTerm != nil && searchTerm != "" {
		filter = bson.M{
			"$or": []bson.M{
				{
					"username": bson.M{
						"$regex": primitive.Regex{
							Pattern: searchTerm.(string),
							Options: "i",
						},
					},
				},
			},
		}
	}

	musers := []MongoUser{}
	users := []models.User{}
	cursor, err := s.Database(handler.DatabaseName).Collection("users").Find(ctx, filter, findOptions)
	if err != nil {
		return users, err
	}
	cursor.All(ctx, &musers)
	defer cursor.Close(ctx)
	for _, muser := range musers {
		users = append(users, models.User{
			ID:          muser.ID.Hex(),
			Username:    muser.Username,
			Password:    muser.Password,
			Description: muser.Description,
			Thumbnail:   muser.Thumbnail,
		})
	}
	return users, err
}
func (handler *MongodbHandler) Authenticate(ctx context.Context, username string, password string) (bool, error) {
	s := handler.Session

	muser := MongoUser{}
	cursor, err := s.Database(handler.DatabaseName).Collection("users").Find(ctx, bson.D{{"username", username}})
	fmt.Println(getHash([]byte(password)))
	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		err = cursor.Decode(&muser)
		if err != nil {
			return false, err
		}
		if bcrypt.CompareHashAndPassword([]byte(muser.Password), []byte(password)) == nil {
			return true, nil
		}
		return false, nil
	}
	return false, ErrPostNotFound
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

	mpostTags := []string{}
	for _, tag := range post.Tags {
		mpostTags = append(mpostTags, tag.ID)
	}

	newPost := MongoPost{
		ID:                 newID,
		Title:              post.Title,
		Content:            post.Content,
		Description:        post.Description,
		Thumbnail:          post.Thumbnail,
		ThumbnailStretched: post.ThumbnailStretched,
		LastEditTimestamp:  timeOfAddage,
		Hidden:             post.Hidden,
		Published:          false,
		Tags:               mpostTags,
	}

	_, err := s.Database(handler.DatabaseName).Collection("posts").InsertOne(ctx, newPost)
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
	fields["thumbnailstretched"] = post.ThumbnailStretched
	fields["lastedittimestamp"] = time.Now()
	fields["hidden"] = post.Hidden

	mpostTags := []string{}
	for _, tag := range post.Tags {
		mpostTags = append(mpostTags, tag.ID)
	}

	fmt.Println(mpostTags)
	fields["tags"] = mpostTags

	for key, value := range fields {
		updateFields = append(updateFields, bson.E{key, value})
	}

	fmt.Println(primitive.ObjectIDFromHex(post.ID))
	if id, err := primitive.ObjectIDFromHex(post.ID); err == nil {
		_, err := s.Database(handler.DatabaseName).Collection("posts").UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set", updateFields}})

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
		_, err := s.Database(handler.DatabaseName).Collection("posts").ReplaceOne(ctx, bson.D{{"_id", id}}, newPost)
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
		_, err := s.Database(handler.DatabaseName).Collection("posts").DeleteOne(ctx, bson.D{{"_id", oid}})
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
		_, err := s.Database(handler.DatabaseName).Collection("posts").UpdateOne(ctx, bson.D{{"_id", oid}}, bson.D{{"$set",
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

	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "lastedittimestamp", Value: -1}})
	var perPage int64 = 10

	if page := ctx.Value("page"); page != nil && page != "" {

		pageInt, err := strconv.Atoi(page.(string))
		if err != nil {
			fmt.Println(err)
		} else {
			findOptions.SetLimit(perPage)
			findOptions.SetSkip((int64(pageInt) - 1) * perPage)
		}
	}

	if searchTerm := ctx.Value("search-term"); searchTerm != nil && searchTerm != "" {
		filter = bson.M{
			"$or": []bson.M{
				{
					"title": bson.M{
						"$regex": primitive.Regex{
							Pattern: searchTerm.(string),
							Options: "i",
						},
					},
				},
				{
					"description": bson.M{
						"$regex": primitive.Regex{
							Pattern: searchTerm.(string),
							Options: "i",
						},
					},
				},
			},
		}
	}

	if ctx.Value("published") == true {
		filter["published"] = true
	}
	if ctx.Value("hidden") == false {
		filter["hidden"] = false
	}
	posts := []models.Post{}

	tags := []models.Tag{}
	tagmap := map[string]models.Tag{}
	cursorTags, err := s.Database(handler.DatabaseName).Collection("tags").Find(ctx, bson.D{})
	if err != nil {
		return posts, err
	}
	cursorTags.All(ctx, &tags)
	cursorTags.Close(ctx)
	for _, tag := range tags {
		tagmap[tag.ID] = tag
	}

	mposts := []MongoPost{}

	cursor, err := s.Database(handler.DatabaseName).Collection("posts").Find(ctx, filter, findOptions)
	if err != nil {
		fmt.Println("Error retrieving posts")
		return posts, err
	}

	//fmt.Println(tagmap)
	cursor.All(ctx, &mposts)
	defer cursor.Close(ctx)
	fmt.Println(mposts)
	for _, mpost := range mposts {

		newPostTags := []models.Tag{}
		for _, postTag := range mpost.Tags {
			postrealtag := tagmap[postTag]
			newPostTags = append(newPostTags, postrealtag)
		}

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
			Tags:              newPostTags,
		})
	}
	return posts, err
}
func (handler *MongodbHandler) GetPostsCount(ctx context.Context) (int64, error) {
	s := handler.Session

	filter := bson.M{}

	published := ctx.Value("published")
	if published == true {
		filter["published"] = true
	}
	hidden := ctx.Value("hidden")
	if hidden == false {
		filter["hidden"] = false
	}

	count, err := s.Database(handler.DatabaseName).Collection("posts").CountDocuments(ctx, filter)
	if err != nil {
		fmt.Println(err)
	}
	return count, err
}
func (handler *MongodbHandler) GetPost(ctx context.Context, id string) (models.Post, error) {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return models.Post{}, err
	// }

	filter := bson.M{}
	if ctx.Value("published") == true {
		filter["published"] = true
	}
	if ctx.Value("hidden") == false {
		filter["hidden"] = false
	}

	post := models.Post{}
	mpost := MongoPost{}
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		filter["_id"] = oid
		cursor, err := s.Database(handler.DatabaseName).Collection("posts").Find(ctx, filter)
		defer cursor.Close(ctx)
		if cursor.Next(ctx) {
			err = cursor.Decode(&mpost)
			if err != nil {
				return models.Post{}, err
			}
			post = models.Post{
				ID:                 mpost.ID.Hex(),
				Title:              mpost.Title,
				Content:            mpost.Content,
				Description:        mpost.Description,
				Thumbnail:          mpost.Thumbnail,
				ThumbnailStretched: mpost.ThumbnailStretched,
				PublishTimestamp:   mpost.PublishTimestamp.Format("January-02-2006"),
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

	mprojectTags := []string{}
	for _, tag := range project.Tags {
		mprojectTags = append(mprojectTags, tag.ID)
	}

	newProject := MongoProject{
		ID:                 newID,
		Title:              project.Title,
		Link:               project.Link,
		Description:        project.Description,
		Thumbnail:          project.Thumbnail,
		ThumbnailStretched: project.ThumbnailStretched,
		Timestamp:          time.Now(),
		Tags:               mprojectTags,
	}

	_, err := s.Database(handler.DatabaseName).Collection("projects").InsertOne(ctx, newProject)
	return err
}
func (handler *MongodbHandler) UpdateProject(ctx context.Context, project models.Project) (models.Project, error) {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return models.Post{}, err
	// }
	var updateFields bson.D

	var fields map[string]interface{}
	fields = make(map[string]interface{})
	fields["title"] = project.Title
	fields["description"] = project.Description
	fields["link"] = project.Link
	fields["thumbnail"] = project.Thumbnail
	fields["thumbnailstretched"] = project.ThumbnailStretched

	mpostTags := []string{}
	for _, tag := range project.Tags {
		mpostTags = append(mpostTags, tag.ID)
	}

	fmt.Println(mpostTags)
	fields["tags"] = mpostTags

	for key, value := range fields {
		updateFields = append(updateFields, bson.E{key, value})
	}

	fmt.Println(primitive.ObjectIDFromHex(project.ID))
	if id, err := primitive.ObjectIDFromHex(project.ID); err == nil {
		_, err := s.Database(handler.DatabaseName).Collection("projects").UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set", updateFields}})

		return models.Project{}, err
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
		_, err := s.Database(handler.DatabaseName).Collection("projects").DeleteOne(ctx, bson.D{{"_id", oid}})
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
	tags := []models.Tag{}

	tagmap := map[string]models.Tag{}
	cursorTags, err := s.Database(handler.DatabaseName).Collection("tags").Find(ctx, bson.D{})
	if err != nil {
		return projects, err
	}
	cursorTags.All(ctx, &tags)
	cursorTags.Close(ctx)
	for _, tag := range tags {
		tagmap[tag.ID] = tag
	}

	cursor, err := s.Database(handler.DatabaseName).Collection("projects").Find(ctx, bson.D{})
	if err != nil {
		return projects, err
	}
	cursor.All(ctx, &mprojects)
	cursor.Close(ctx)

	for _, mproject := range mprojects {

		newPostTags := []models.Tag{}
		for _, postTag := range mproject.Tags {
			postrealtag := tagmap[postTag]
			newPostTags = append(newPostTags, postrealtag)
		}

		projects = append(projects, models.Project{
			ID:                 mproject.ID.Hex(),
			Title:              mproject.Title,
			Link:               mproject.Link,
			Description:        mproject.Description,
			Thumbnail:          mproject.Thumbnail,
			ThumbnailStretched: mproject.ThumbnailStretched,
			Timestamp:          mproject.Timestamp.Format("Jan/06/15"),
			Tags:               newPostTags,
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
		cursor, err := s.Database(handler.DatabaseName).Collection("projects").Find(ctx, bson.D{{"_id", oid}})
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

	filter := bson.M{}

	published := ctx.Value("published")
	if published == true {
		filter["published"] = true
	}
	hidden := ctx.Value("hidden")
	if hidden == false {
		filter["hidden"] = false
	}

	mlinks := []MongoLink{}
	links := []models.Link{}
	cursor, err := s.Database(handler.DatabaseName).Collection("posts").Find(ctx, filter)
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
	cursor, err := s.Database(handler.DatabaseName).Collection("knowledgetimelineevents").Find(ctx, bson.D{})
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

	_, err := s.Database(handler.DatabaseName).Collection("knowledgetimelineevents").InsertOne(ctx, newEvent)
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

		_, err = s.Database(handler.DatabaseName).Collection("knowledgetimelineevents").UpdateOne(ctx, bson.D{{"_id", updatedEvent.ID}}, updatedEvent)
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
		_, err := s.Database(handler.DatabaseName).Collection("knowledgetimelineevents").DeleteOne(ctx, bson.D{{"_id", oid}})
		return err
	} else {
		return err
	}
}
func (handler *MongodbHandler) RemoveAllTags(ctx context.Context) error {
	s := handler.Session

	err := s.Database(handler.DatabaseName).Collection("tags").Drop(ctx)
	return err
}
func (handler *MongodbHandler) AddTag(ctx context.Context, tag models.Tag) error {
	s := handler.Session

	//newID := primitive.NewObjectID()
	guid := xid.New()
	tag.ID = guid.String()

	_, err := s.Database(handler.DatabaseName).Collection("tags").InsertOne(ctx, tag)
	return err
}
func (handler *MongodbHandler) RemoveTag(ctx context.Context, tagID string) error {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return err
	// }

	//if oid, err := primitive.ObjectIDFromHex(userID); err == nil {
	_, err := s.Database(handler.DatabaseName).Collection("tags").DeleteOne(ctx, bson.D{{"ID", tagID}})
	return err
}
func (handler *MongodbHandler) UpdateTag(ctx context.Context, tag models.Tag) error {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return models.Post{}, err
	// }
	var updateFields bson.D

	var fields map[string]interface{}
	fields = make(map[string]interface{})
	fields["name"] = tag.Name
	fields["content"] = tag.Content
	fields["description"] = tag.Description
	fields["thumbnail"] = tag.Thumbnail
	fields["thumbnailstretched"] = tag.ThumbnailStretched

	for key, value := range fields {
		updateFields = append(updateFields, bson.E{key, value})
	}

	// fmt.Println(primitive.ObjectIDFromHex(user.ID))
	// if id, err := primitive.ObjectIDFromHex(user.ID); err == nil {
	_, err := s.Database(handler.DatabaseName).Collection("tags").UpdateOne(ctx, bson.D{{"ID", tag.ID}}, bson.D{{"$set", updateFields}})

	return err
}
func (handler *MongodbHandler) GetTags(ctx context.Context) ([]models.Tag, error) {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)

	filter := bson.M{}
	findOptions := options.Find()

	tags := []models.Tag{}
	cursor, err := s.Database(handler.DatabaseName).Collection("tags").Find(ctx, filter, findOptions)
	if err != nil {
		return tags, err
	}
	cursor.All(ctx, &tags)
	defer cursor.Close(ctx)
	return tags, err
}
func (handler *MongodbHandler) AddPage(ctx context.Context, page models.Page) error {
	s := handler.Session

	// sections :=
	// for _, section := range page.Sections {
	// 	stackSection, ok := section.(models.StackSection)

	// }

	guid := xid.New()
	page.ID = guid.String()

	_, err := s.Database(handler.DatabaseName).Collection("pages").InsertOne(ctx, page)
	return err
}
func (handler *MongodbHandler) UpdatePage(ctx context.Context, page models.Page) error {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return models.Post{}, err
	// }
	var updateFields bson.D

	var fields map[string]interface{}
	fields = make(map[string]interface{})
	fields["name"] = page.Name
	fields["homepage"] = page.Homepage
	fields["sections"] = page.Sections

	for key, value := range fields {
		updateFields = append(updateFields, bson.E{key, value})
	}

	fmt.Println(page.ID)
	_, err := s.Database(handler.DatabaseName).Collection("pages").UpdateOne(ctx, bson.D{{"ID", page.ID}}, bson.D{{"$set", updateFields}})

	return err
}
func (handler *MongodbHandler) RemovePage(ctx context.Context, id string) error {
	s := handler.Session

	_, err := s.Database(handler.DatabaseName).Collection("pages").DeleteOne(ctx, bson.D{{"ID", id}})
	return err
}
func (handler *MongodbHandler) GetPages(ctx context.Context) ([]models.Page, error) {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)

	filter := bson.M{}
	pages := []models.Page{}

	cursor, err := s.Database(handler.DatabaseName).Collection("pages").Find(ctx, filter)
	if err != nil {
		fmt.Println("Error retrieving posts")
		return pages, err
	}

	//fmt.Println(tagmap)
	cursor.All(ctx, &pages)
	defer cursor.Close(ctx)
	return pages, err
}
func (handler *MongodbHandler) GetPage(ctx context.Context, id string) (models.Page, error) {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return models.Post{}, err
	// }

	filter := bson.M{}
	page := models.Page{}
	mpage := MongoPage{}
	if ctx.Value("homepage") == true {
		filter["homepage"] = true
	} else {
		filter["ID"] = id
	}
	cursor, err := s.Database(handler.DatabaseName).Collection("pages").Find(ctx, filter)
	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		err = cursor.Decode(&mpage)
		if err != nil {
			return models.Page{}, err
		}
		page.ID = mpage.ID
		page.Homepage = mpage.Homepage
		page.Name = mpage.Name
		for _, section := range mpage.Sections {
			page.Sections = append(page.Sections, section)
		}
		return page, nil
	}
	return models.Page{}, ErrPostNotFound
}
func (handler *MongodbHandler) AddVisit(ctx context.Context, visit models.Visit) error {
	fields := make(map[string]interface{})
	fields["pages."+visit.URL] = 1
	fields["countries."+visit.Country] = 1
	s := handler.Session
	opts := options.Update().SetUpsert(true)

	date := time.Now().Format("2006/01/02")

	_, err := s.Database(handler.DatabaseName).Collection("visits").UpdateOne(ctx, bson.D{{"date", date}}, bson.D{{"$inc", fields}}, opts)

	return err
}
func (handler *MongodbHandler) GetVisits(ctx context.Context) ([]models.VisitSummary, error) {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)

	filter := bson.D{}
	visits := []models.VisitSummary{}

	cursor, err := s.Database(handler.DatabaseName).Collection("visits").Find(ctx, filter)
	if err != nil {
		fmt.Println("Error retrieving posts")
		return visits, err
	}

	//fmt.Println(tagmap)
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var visit models.VisitSummary
		cursor.Decode(&visit)
		visits = append(visits, visit)
	}
	fmt.Println(visits)
	return visits, err
}
func (handler *MongodbHandler) AddSettings(ctx context.Context, settings models.Settings) error {
	s := handler.Session

	settings.ID = handler.SettingsID
	replaceModel := options.ReplaceOptions{}
	replaceModel.SetUpsert(true)
	//newID := primitive.NewObjectID()
	// cursor := s.Database(handler.databaseName).Collection("settings").FindOne(ctx, bson.D{{"ID", handler.settingsID}})
	// err := cursor.Decode(&settings)
	// if err != nil {
	_, err := s.Database(handler.DatabaseName).Collection("settings").ReplaceOne(ctx, bson.D{{"ID", handler.SettingsID}}, settings, &replaceModel)
	return err
	//}
	return nil
}
func (handler *MongodbHandler) UpdateSettings(ctx context.Context, settings models.Settings) error {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)
	// if err != nil {
	// 	return models.Post{}, err
	// }
	var updateFields bson.D

	var fields map[string]interface{}
	fields = make(map[string]interface{})
	fields["websiteName"] = settings.WebsiteName
	if settings.Logo != "" {
		fields["logo"] = settings.Logo
	}
	fields["backgroundColor1"] = settings.BackgroundColor1
	fields["backgroundColor2"] = settings.BackgroundColor2
	fields["backgroundColor3"] = settings.BackgroundColor3
	fields["textColor1"] = settings.TextColor1
	fields["textColor2"] = settings.TextColor2
	fields["textColor3"] = settings.TextColor3
	// if user.Password != "" {
	// 	fields["password"] = getHash([]byte(user.Password))
	// }

	for key, value := range fields {
		updateFields = append(updateFields, bson.E{key, value})
	}

	_, err := s.Database(handler.DatabaseName).Collection("settings").UpdateOne(ctx, bson.D{{"ID", handler.SettingsID}}, bson.D{{"$set", updateFields}})
	if err != nil {
		return ErrInvalidPostId
	}
	return nil
}
func (handler *MongodbHandler) GetSettings(ctx context.Context) (models.Settings, error) {
	s := handler.Session
	// err := s.Connect(ctx)
	// defer s.Disconnect(ctx)

	settings := models.Settings{}
	cursor := s.Database(handler.DatabaseName).Collection("settings").FindOne(ctx, bson.D{{"ID", handler.SettingsID}}, options.FindOne())

	err := cursor.Decode(&settings)
	if err != nil {
		return settings, err
	}

	return settings, nil
}
