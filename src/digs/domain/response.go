package domain

type UserLoginResponse struct {
	StatusCode int32
	SessionId  string `json:"sessionId" bson:"sessionId"`
	UserId     string `json:"userId" bson:"userId"`
	Settings   SettingResponse `json:"settings" bson:"settings"`
	Verified   bool `json:"verified" bson:"verified"`
}

type MessageReceivedResponse struct {
	RequestId  string `json:"messageId" bson:"messageId"`
	StatusCode int32 `json:"statusCode" bson:"statusCode"`
	GID        string `json:"groupId" bson:"groupId"`
}

type GenericResponse struct {
	RequestId   string `json:"requestId" bson:"requestId"`
	StatusCode  int32 `json:"statusCode" bson:"statusCode"`
	MessageCode int32 `json:"messageCode" bson:"messageCode"`
	Message     string `json:"message" bson:"message"`
}

type MessageSendResponse struct {
	RequestId string `json:"requestId" bson:"requestId"`
	Sent      bool `json:"sent"`
}

type MessageGetResponse struct {
	UID            string `json:"userId" bson:"userId"`
	MID            string `json:"messageId" bson:"messageId"`
	GID            string `json:"groupId" bson:"groupId"`
	From           string `json:"name" bson:"name"`
	About          string `json:"about" bson:"about"`
	IsGroup        bool `json:"isGroup" bson:"isGroup"`
	Message        string `json:"body" bson:"body"`
	Verified       bool `json:"verified" bson:"verified"`
	Timestamp      int64 `json:"timestamp" bson:"timestamp"`
	ProfilePicture string `json:"picture" bson:"picture"`
}

type PersonResponse struct {
	Name           string `json:"name" bson:"name"`
	UID            string `json:"userId" bson:"userId"`
	GID            string `json:"groupId" bson:"groupId"`
	About          string `json:"about" bson:"about"`
	Verified       bool `json:"verified" bson:"verified"`
	UnreadCount    int64 `json:"unreadCount" bson:"unreadCount"`
	ActiveState    string `json:"state" bson:"state"`
	Activity       string `json:"activity" bson:"activity"`
	MemberCount    int `json:"memberCount" bson:"memberCount"`
	ProfilePicture string `json:"picture" bson:"picture"`
	IsGroup        bool `json:"isGroup" bson:"isGroup"`
}

type MessagePushResponse struct {
	MessageGetResponse
}

type SettingResponse struct {
	Range            float64 `json:"messageRange" bson:"messageRange"`
	PublicProfile    bool `json:"publicProfile" bson:"publicProfile"`
	PushNotification bool `json:"enableNotification" bson:"enableNotification"`
}

type CreateGroupResponse struct {
	GID          string `json:"groupId" bson:"groupId"`
	GroupName    string `json:"groupName" bson:"groupName"`
	GroupAbout   string `json:"groupAbout" bson:"groupAbout"`
	GroupPicture string `json:"groupPicture" bson:"groupPicture"`
	Messages     []MessageGetResponse `json:"messages" bson:"messages"`
}

type MessagePutResponse struct {
	ResourceUrl string `json:"resourceUrl" bson:"resourceUrl"`
}

type GroupDetail struct {
	GID          string `json:"groupId" bson:"groupId"`
	Users        []PersonResponse `json:"users" bson:"users"`
	GroupName    string `json:"groupName" bson:"groupName"`
	GroupAbout   string `json:"groupAbout" bson:"groupAbout"`
	GroupPicture string `json:"groupPicture" bson:"groupPicture"`
}