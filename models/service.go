package models

type Service struct {
	Name    string   `json:"name" bson:"name"`
	URL     string   `json:"url" bson:"url"`
	Headers []string `json:"headers" bson:"headers"`
}
