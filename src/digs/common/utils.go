package common

import "math"

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


// http://en.wikipedia.org/wiki/Haversine_formula
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2 - la1) + math.Cos(la1) * math.Cos(la2) * hsin(lo2 - lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta / 2), 2)
}
