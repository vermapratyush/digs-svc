package domain

//Except for Setting Controller this is used.
type BaseRequest struct {
	HeaderUserAgent string
	HeaderSessionID string
	RequestId       string `json:"requestId" bson:"requestId"`
	SessionID       string `json:"sessionId" bson:"sessionId"`
}

type NotificationRequest struct {
	BaseRequest
	NotificationID    string `json:"notificationId" bson:"notificationId"`
	OldNotificationID string `json:"oldNotificationId" bson:"oldNotificationId"`
	OSType            string `json:"osType" bson:"osType"`
	AppVersion        string `json:"appVersion" bson:"appVersion"`
}

type UserLogoutRequest struct {
	BaseRequest
	NotificationId string `json:"notificationId" bson:"notificationId"`
}

type UserLoginRequest struct {
	BaseRequest
	FBID           string `json:"fbid" bson:"fbid"`
	Locale         string `json:"locale" bson:"locale"`
	FirstName      string `json:"firstName" bson:"firstName"`
	LastName       string `json:"lastName" bson:"lastName"`
	Email          string `json:"email" bson:"email"`
	About          string `json:"about" bson:"about"`
	ProfilePicture string `bson:"picture" json:"picture"`
	AccessToken    string `json:"accessToken" bson:"accessToken"`
}

type MessageSendRequest struct {
	BaseRequest
	Body      string `json:"body" bson:"body"`
	GID       string `json:"groupId" bson:"groupId"`
	Location  Coordinate `json:"location" bson:"location"`
	Reach     int64 `json:"reach" bson:"reach"`
	MID       string `json:"messageId" json:"messageId"`
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
}

type MessageTypingRequest struct {
	BaseRequest
	IsTyping bool `json:"isTyping" bson:"isTyping"`
	GID      string `json:"groupId" bson:"groupId"`
}

type Coordinate struct {
	BaseRequest
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

type MessageGetRequest struct {
	BaseRequest
	Username string `json:"username" bson:"username"`
	Location []float64 `json:"location" bson:"location"`
	Distance int64 `json:"distance" bson:"distance"`
}

type SettingRequest struct {
	BaseRequest
	Range            float64 `json:"messageRange" bson:"messageRange"`
	PublicProfile    bool `json:"publicProfile" bson:"publicProfile"`
	PushNotification bool `json:"enableNotification" bson:"enableNotification"`
}

type AbuseRequest struct {
	BaseRequest
	UID string `json:"userId" bson:"userId"`
	GID string `json:"groupId" bson:"groupId"`
	MID string `json:"messageId" bson:"messageId"`
}

type GroupCreateRequest struct {
	BaseRequest
	UIDS       []string `json:"userIds" bson:"userIds"`
	GID        string `json:"groupId" bson:"groupId"`
	GroupName  string `json:"groupName" bson:"groupName"`
	GroupAbout string `json:"groupAbout" bson:"groupAbout"`
	GroupPicture string `json:"groupPicture" bson:"groupPicture"`
	IsGroup    bool `json:"isGroup" bson:"isGroup"`
}

type UnreadRequest struct {
	BaseRequest
	GID string `json:"groupId" bson:"groupId"`
	MID string `json:"messageId" bson:"messageId"`
}

type MessagePinAddDeleteRequest struct {
	BaseRequest
	MID string `json:"messageId" bson:"messageId"`
}