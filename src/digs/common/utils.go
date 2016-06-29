package common

const (
	PushNotification_API_KEY = "AIzaSyAv23-LpWCS97b1CR0nV-JioSLk6MrM0_U"
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
