package domain


type UserLoginResponse struct {
	StatusCode int32
	SessionId string `json:"sessionId" bson:"sessionId"`
	UserId string `json:"userId" bson:"userId"`
}

type MessageReceivedResponse struct {
	StatusCode int32 `json:"statusCode" bson:"statusCode"`
}

type ErrorResponse struct {
	StatusCode int32 `json:"statusCode" bson:"statusCode"`
	ErrorCode int32 `json:"errorCode" bson:"errorCode"`
	Message string `json:"message" bson:"message"`
}

type MessageSendResponse struct {
	Sent bool `json:"sent"`
}

type MessageGetResponse struct {
	UID string `json:"userId" bson:"userId"`
	From string `json:"name" bson:"name"`
	Message string `json:"body" bson:"body"`
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
	ProfilePicture string `json:"picture" bson:"picture"`
}