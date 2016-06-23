package domain


type UserLoginResponse struct {
	StatusCode int32
	SessionId string `json:"sessionId" bson:"sessionId"`
}

type MessageReceivedResponse struct {
	StatusCode int32
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
	UID string `json:"uid" bson:"uid"`
	From string `json:"from" bson:"from"`
	Message string `json:"message" bson:"message"`
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
}