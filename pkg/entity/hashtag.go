package entity

type Hashtag struct {
	HashtagID string `bson:"_id" json:"_id"`
	Keyword   string `bson:"keyword" json:"keyword"`
}

type TimelineResponse struct {
	Data     []Timeline `json:"data"`
	NextPage string     `bson:"next_page" json:"next_page"`
}

type Timeline struct {
	ID           string `bson:"_id" json:"_id"`
	Text         string `bson:"text" json:"text"`
	UserID       string `bson:"user_id" json:"user_id"`
	Likes        int    `bson:"likes" json:"likes"`
	ParentThread string `bson:"parent_thread" json:"parent_thread"`
	RepostCount  int    `bson:"repost_count" json:"repost_count"`
}

type UserInfo struct {
	ID              string `bson:"_id" json:"_id"`
	DisplayName     string `bson:"display_name" json:"display_name"`
	Username        string `bson:"username" json:"username"`
	ProfileImageURL string `bson:"profile_image_url" json:"profile_image_url"`
	Description     string `bson:"description" json:"description"`
	Follower        int    `bson:"follower" json:"follower"`
	Following       int    `bson:"following" json:"following"`
}

type Job struct {
	Thread  Thread      `json:"timeline"`
	Account AccountInfo `json:"account"`
}
