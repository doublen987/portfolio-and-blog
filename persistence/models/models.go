package models

type Collection struct {
	ID string `bson:"_id", json:"id"`
}

type Post struct {
	ID                string   `bson:"_id",json:"id"`
	Title             string   `bson:"title", json:"title"`
	Content           string   `bson:"content", json:"content"`
	Description       string   `bson:"description", json: "description"`
	Thumbnail         string   `bson:"thumbnail", json:"thumbnail"`
	PublishTimestamp  string   `bson:"publishtimestamp", json:"publishtimestamp"`
	LastEditTimestamp string   `bson:"lastedittimestamp", json:"lastedittimestamp"`
	Hidden            bool     `bson:"hidden", json:"hidden"`
	Published         bool     `bson:"published", json:"published"`
	Tags              []string `bson:"tags", json:"tags"`
}

type Project struct {
	ID          string   `bson:"_id", json:"id"`
	Title       string   `bson:"title", json:"title"`
	Description string   `bson:"description", json:"description"`
	Link        string   `bson:"link", json:"link"`
	Thumbnail   string   `bson:"thumbnail", json:"thumbnail"`
	Timestamp   string   `bson:"timestamp", json:"timestamp"`
	Tags        []string `bson:"tags", json:"tags"`
}

type Link struct {
	ID          string `bson:"_id",json:"id"`
	Title       string `bson:"title",json:"title"`
	Description string `bson:"description",json:"description"`
}

type TimelineEvent struct {
	ID          string `bson:"_id", json"id"`
	Title       string `bson:"title", json:"title"`
	Description string `bson:"description", json:"description"`
	Image       string `bson:"image", json:"image"`
}
