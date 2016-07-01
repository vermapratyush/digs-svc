package domain

type UserLoginResponse struct {
	StatusCode int32
	SessionId string `json:"sessionId" bson:"sessionId"`
	UserId string `json:"userId" bson:"userId"`
	Settings SettingResponse `json:"settings" bson:"settings"`
}

type MessageReceivedResponse struct {
	RequestId string `json:"messageId" bson:"messageId"`
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
	UID string `json:"userId" bson:"userId"`
	MID string `json:"messageId" bson:"messageId"`
	From string `json:"name" bson:"name"`
	Message string `json:"body" bson:"body"`
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
	ProfilePicture string `json:"picture" bson:"picture"`
}

type PersonResponse struct {
	Name string `json:"name" bson:"name"`
	UID string `json:"userId" bson:"userId"`
	About string `json:"about" bson:"about"`
	Activity string `json:"activity" bson:"activity"`
	ProfilePicture string `json:"picture" bson:"picture"`
}

type MessagePushResponse struct {
	MessageGetResponse
}

type SettingResponse struct {
	Range float64 `json:"messageRange" bson:"messageRange"`
	PublicProfile bool `json:"publicProfile" bson:"publicProfile"`
	PushNotification bool `json:"enableNotification" bson:"enableNotification"`

}