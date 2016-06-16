package domain

type Coordinate struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type BaseRequest struct {
}

type UserLoginRequest struct {
	UserAgent string
	SessionID string
	AccessToken string `json:"accessToken" bson:"accessToken"`
}

type MessageSendRequest struct {
	UserAgent string
	SessionID string
	Body      string `json:"body" bson:"body"`
	Username  string `json:"username" bson:"username"`
	Location  Coordinate `json:"location" bson:"location"`
}

type MessageGetRequest struct {
	UserAgent string
	SessionID string
	Username  string `json:"username" bson:"username"`
	Location  Coordinate `json:"location" bson:"location"`
	Distance  int64 `json:"distance" bson:"distance"`
}