package domain

type BaseRequest struct {
	UserAgent string
	SessionID string
}

type UserLoginRequest struct {
	BaseRequest
	FBID string `json:"fbid" bson:"fbid"`
	Locale string `json:"locale" bson:"locale"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName string `json:"lastName" bson:"lastName"`
	Email string `json:"email" bson:"email"`
	About string `json:"about" bson:"about"`
	ProfilePicture string `bson:"profilePicture" json:"profilePicture"`
	FBVerified string `bson:"fbVerified" json:"fbVerified"`
	AccessToken string `json:"accessToken" bson:"accessToken"`
}

type MessageSendRequest struct {
	BaseRequest
	Body      string `json:"body" bson:"body"`
	Location  Coordinate `json:"location" bson:"location"`
	Reach int64 `json:"reach" bson:"reach"`
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
}

type Coordinate struct  {
	Latitude float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

type MessageGetRequest struct {
	BaseRequest
	Username  string `json:"username" bson:"username"`
	Location  []float64 `json:"location" bson:"location"`
	Distance  int64 `json:"distance" bson:"distance"`
}