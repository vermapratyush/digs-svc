package domain

type BaseRequest struct {
	UserAgent string
	SessionID string
}

type UserLoginRequest struct {
	BaseRequest
	FirstName string `json:"firstName" bson:"firstName"`
	LastName string `json:"lastName" bson:"lastName"`
	Email string `json:"email" bson:"email"`
	About string `json:"about" bson:"about"`
	AccessToken string `json:"accessToken" bson:"accessToken"`
}

type MessageSendRequest struct {
	BaseRequest
	Body      string `json:"body" bson:"body"`
	Username  string `json:"username" bson:"username"`
	Location  []float64 `json:"location" bson:"location"`
}

type MessageGetRequest struct {
	BaseRequest
	Username  string `json:"username" bson:"username"`
	Location  []float64 `json:"location" bson:"location"`
	Distance  int64 `json:"distance" bson:"distance"`
}