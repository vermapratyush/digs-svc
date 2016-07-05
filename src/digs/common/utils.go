package common

const (
	PushNotification_API_KEY = "AIzaSyDFEvXnsgfd8qhUzbQfD98GUbwpclGuJDc"
	MessageBatchSize = 50
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


