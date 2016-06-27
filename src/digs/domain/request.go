package domain

//Except for Setting Controller this is used.
type BaseRequest struct {
	HeaderUserAgent string
	HeaderSessionID string
	RequestId string `json:"requestId" bson:"requestId"`
	SessionID string `json:"sessionId" bson:"sessionId"`
}

type NotificationRequest struct {
	BaseRequest
	NotificationID string `json:"notificationId" bson:"notificationId"`
	OSType string `json:"os" bson:"os"`
}

type UserLogoutRequest struct {
	BaseRequest
	NotificationId string `json:"notificationId" bson:"notificationId"`
}

type UserLoginRequest struct {
	BaseRequest
	FBID string `json:"fbid" bson:"fbid"`
	Locale string `json:"locale" bson:"locale"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName string `json:"lastName" bson:"lastName"`
	Email string `json:"email" bson:"email"`
	About string `json:"about" bson:"about"`
	ProfilePicture string `bson:"picture" json:"picture"`
	AccessToken string `json:"accessToken" bson:"accessToken"`
	FBVerified bool `bson:"fbVerified" json:"fbVerified"`
}

type MessageSendRequest struct {
	BaseRequest
	Body      string `json:"body" bson:"body"`
	Location  Coordinate `json:"location" bson:"location"`
	Reach int64 `json:"reach" bson:"reach"`
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
}

type Coordinate struct  {
	BaseRequest
	Latitude float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

type MessageGetRequest struct {
	BaseRequest
	Username  string `json:"username" bson:"username"`
	Location  []float64 `json:"location" bson:"location"`
	Distance  int64 `json:"distance" bson:"distance"`
}

type SettingRequest struct {
	Range int64 `json:"messageRange" bson:"messageRange"`
	PublicProfile bool `json:"publicProfile" bson:"publicProfile"`
	PushNotification bool `json:"notificationEnabled" bson:"notificationEnabled"`

}