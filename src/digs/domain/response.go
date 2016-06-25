package domain

type UserLoginResponse struct {
	StatusCode int32
	SessionId string `json:"sessionId" bson:"sessionId"`
	UserId string `json:"userId" bson:"userId"`
}

type MessageReceivedResponse struct {
	RequestId string `json:"requestId" bson:"requestId"`
	StatusCode int32 `json:"statusCode" bson:"statusCode"`
}

type ErrorResponse struct {
	RequestId string `json:"requestId" bson:"requestId"`
	StatusCode int32 `json:"statusCode" bson:"statusCode"`
	ErrorCode int32 `json:"errorCode" bson:"errorCode"`
	Message string `json:"message" bson:"message"`
}

type MessageSendResponse struct {
	RequestId string `json:"requestId" bson:"requestId"`
	Sent bool `json:"sent"`
}

type MessageGetResponse struct {
	RequestId string `json:"requestId" bson:"requestId"`
	UID string `json:"userId" bson:"userId"`
	From string `json:"name" bson:"name"`
	Message string `json:"body" bson:"body"`
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
	ProfilePicture string `json:"picture" bson:"picture"`
}

type MessagePushResponse struct {
	MessageGetResponse
}