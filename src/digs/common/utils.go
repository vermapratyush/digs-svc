package common

const (
	PushNotification_API_KEY = "AIzaSyCMYdgQUqL8X7D5OaY7hvADMOQzA6WaqPI"
	MessageBatchSize = 50

	//Hystrix Commands
	MessageWrite = "MessageWrite"
	MessageGetAll = "MessageGetAll"
	Notification = "Notification"
	UserAccount = "UserAccount"
	UserAccountGetAll = "UserAccountGetAll"
	SessionWrite = "SessionWrite"
	SessionGet = "SessionGet"
	SessionDel = "SessionDel"
	FeedAdd = "FeedInsert"
	FeedGet = "FeedGet"
	LocationUpdate = "LocationUpdate"
	LocationGet = "LocationGet"
	LocationUserFind = "LocationUserFind"
	AndroidPush = "AndroidPush"
	IOSPush = "IOSPush"
)

func GetName(firstName, lastName string) string {
	return firstName + " " + lastName
}

func IndexOf(haystack []string, needle string) int {
	for idx, h := range (haystack) {
		if h == needle {
			return idx
		}
	}
	return -1
}


