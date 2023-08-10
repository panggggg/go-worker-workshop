package entity

type ThreadResponse struct {
	Data     []Thread `json:"data"`
	NextPage string   `bson:"next_page" json:"next_page"`
}

type Thread struct {
	ID           string `bson:"_id" json:"_id"`
	Text         string `bson:"text" json:"text"`
	UserID       string `bson:"user_id" json:"user_id"`
	Likes        int    `bson:"likes" json:"likes"`
	ParentThread string `bson:"parent_thread" json:"parent_thread"`
	RepostCount  int    `bson:"repost_count" json:"repost_count"`
}
