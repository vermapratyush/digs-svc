package domain


type UserLoginResponse struct {
	StatusCode int32
	SessionId string
	Name string
	Email string
	About string
}

type MessageReceivedResponse struct {
	StatusCode int32
}

type ErrorResponse struct {
	StatusCode int32
	Message string
}

type MessageSendResponse struct {
	Sent bool `json:"sent"`
}

type MessageGetResponse struct {
	UID string `json:"uid" bson:"uid"`
	From string `json:"from" bson:"from"`
	Message string `json:"message" bson:"message"`
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
}