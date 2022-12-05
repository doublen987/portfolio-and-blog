package models

import (
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Collection struct {
	ID string `bson:"_id" json:"id"`
}

type User struct {
	ID          string `bson:"_id" json:"ID"`
	Username    string `bson:"username" json:"username"`
	Password    string `bson:"password" json:"password"`
	Description string `bson:"description" json:"description"`
	Thumbnail   string `bson:"thumbnail" json:"thumbnail"`
}

type Post struct {
	ID                 string `bson:"_id" json:"id"`
	Title              string `bson:"title" json:"title"`
	Content            string `bson:"content" json:"content"`
	Description        string `bson:"description" json:"description"`
	Thumbnail          string `bson:"thumbnail" json:"thumbnail"`
	ThumbnailStretched bool   `bson:"thumbnailstretched" json:"thumbnailstretched"`
	PublishTimestamp   string `bson:"publishtimestamp" json:"publishtimestamp"`
	LastEditTimestamp  string `bson:"lastedittimestamp" json:"lastedittimestamp"`
	Hidden             bool   `bson:"hidden" json:"hidden"`
	Published          bool   `bson:"published" json:"published"`
	Tags               []Tag  `bson:"tags" json:"tags"`
}

type Project struct {
	ID                 string `bson:"_id" json:"id"`
	Title              string `bson:"title" json:"title"`
	Description        string `bson:"description" json:"description"`
	Link               string `bson:"link" json:"link"`
	Thumbnail          string `bson:"thumbnail" json:"thumbnail"`
	ThumbnailStretched bool   `bson:"thumbnailstretched" json:"thumbnailstretched"`
	Timestamp          string `bson:"timestamp" json:"timestamp"`
	Tags               []Tag  `bson:"tags" json:"tags"`
}

func (p Project) String() string {
	return fmt.Sprintf("%s - %s", p.Title, p.Description)
}

type Link struct {
	ID          string `bson:"_id" json:"id"`
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
}

type TimelineEvent struct {
	ID          string `bson:"_id" json:"id"`
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Image       string `bson:"image" json:"image"`
}

type Tag struct {
	ID                 string `bson:"ID" json:"ID"`
	Name               string `bson:"header" json:"header"`
	Content            string `bson:"content" json:"content"`
	Description        string `bson:"description" json:"descrtiption"`
	Thumbnail          string `bson:"thumbnail" json:"thumbnail"`
	ThumbnailStretched bool   `bson:"thumbnailstretched" json:"thumbnailstretched"`
}

func (t Tag) Value() int {
	bla, _ := strconv.Atoi(t.ID)
	return bla
}

type PageSection interface {
	GetType() string
	GetID() string
}

type TextSection struct {
	ID      string `bson:"ID omitempty" json:"ID omitempty"`
	Header  string `bson:"header omitempty" json:"header omitempty"`
	Content string `bson:"content omitempty" json:"content omitempty"`
}

func (ts TextSection) GetType() string {
	return "text"
}

func (ss TextSection) GetID() string {
	return ss.ID
}

type ImageSection struct {
	ID    string `bson:"ID" json:"ID"`
	Image string `bson:"image" json:"image"`
}

func (is ImageSection) GetType() string {
	return "image"
}

func (ss ImageSection) GetID() string {
	return ss.ID
}

type TagSection struct {
	Name string `bson:"name" json:"name"`
	Tags []Tag  `bson:"tags" json:"tags"`
}

type StackSection struct {
	ID          string       `bson:"ID" json:"ID"`
	Name        string       `bson:"name" json:"name"`
	TagSections []TagSection `bson:"tagssections" json:"tagssections"`
}

func (ss StackSection) GetType() string {
	return "stack"
}

func (ss StackSection) GetID() string {
	return ss.ID
}

type Page struct {
	ID       string        `json:"ID" bson:"ID"`
	Name     string        `json:"name" bson:"name"`
	Homepage bool          `json:"homepage" bson:"homepage"`
	Sections []interface{} `json:"sections" bson:"sections"`
}

type Page2 struct {
	ID       string        `json:"ID" bson:"ID"`
	Name     string        `json:"name" bson:"name"`
	Homepage bool          `json:"homepage" bson:"homepage"`
	Sections []PageSection `json:"sections" bson:"sections"`
}

type Visit struct {
	URL     string `json:"url" bson:"url"`
	Country string `json:"country" bson:"country"`
}

type VisitSummary struct {
	ID        primitive.ObjectID `bson:"_id" json:"ID"`
	Date      string             `json:"date" bson:"date"`
	Pages     map[string]int     `json:"pages" bson:"pages"`
	Countries map[string]int     `json:"countries" bson:"countries"`
}

type SocialLink struct {
	Name string `bson:"name" json:"name"`
	Link string `bson:"link" json:"link"`
	Logo string `bson:"logo" json:"logo"`
}

type Settings struct {
	ID               string       `bson:"ID" json:"ID"`
	WebsiteName      string       `bson:"websiteName" json:"websiteName"`
	Logo             string       `bson:"logo" json:"logo"`
	BackgroundColor1 string       `bson:"backgroundColor1" json:"backgroundColor1"`
	BackgroundColor2 string       `bson:"backgroundColor2" json:"backgroundColor2"`
	BackgroundColor3 string       `bson:"backgroundColor3" json:"backgroundColor3"`
	TextColor1       string       `bson:"textColor1" json:"textColor1"`
	TextColor2       string       `bson:"textColor2" json:"textColor2"`
	TextColor3       string       `bson:"textColor3" json:"textColor3"`
	SocialLinks      []SocialLink `bson:"soclinks" json:"soclinks"`
}
