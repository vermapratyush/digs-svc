package common

import (
	"math"
	"digs/domain"
)

const (
	PushNotification_API_KEY = "AIzaSyCMYdgQUqL8X7D5OaY7hvADMOQzA6WaqPI"
	Meetup_API_KEY = "06f4d95e5f7e7e6775682b25106f17"
	FourSquare_API_CLIENT_ID = "SHZCEZO4D141UZ1LOT52SDKHPC3VFBUMBCJPVCUJT4INBIXF"
	FourSquare_API_CLIENT_SECRET = "TKBLBXNKG5SQF0VZ3JTJEQLXAZWEYKXRJXWEFQOSBIK3IOY1"
	FourSquare_API_CLIENT_VERSION = "20160728"
	Bitly_Login = "vermapratyush"
	Bitly_API_KEY = "R_4d0b283da71e99ec6c0eb703f289982f"
	MessageBatchSize = 50
	AmazonS3BucketName = "powow-file-sharing"

	//Constant-Variables
	DefaultReach = 10000000.0

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
	FeedDel = "FeedDel"
	LocationUpdate = "LocationUpdate"
	LocationGet = "LocationGet"
	LocationUserFind = "LocationUserFind"
	UserGroup = "UserGroup"
	UserGroupBatch = "UserGroupBatch"
	AndroidPush = "AndroidPush"
	AmazonS3 = "AmazonS3"
	IOSPush = "IOSPush"
	MeetupAPI = "MeetupAPI"
	FourSquareAPI = "FourSquareAPI"
	BitlyAPI = "BitlyAPI"
)

func GetStringArrayAsMap(array []string) (map[string]struct{}) {
	arrayAsMap := make(map[string]struct{})
	for _, user := range(array) {
		arrayAsMap[user] = struct{}{}
	}
	return arrayAsMap
}

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

func IsUserBlocked(blockedUsers []string, fromUID string) bool {
	for _, blockedUser := range(blockedUsers) {
		if fromUID == blockedUser {
			return true
		}
	}
	return false
}

func DistanceLong(long1, lat1, long2, lat2 float64) float64 {
	pointA := domain.Coordinate{
		Latitude:lat1,
		Longitude:long1,
	}
	pointB := domain.Coordinate{
		Latitude:lat2,
		Longitude:long2,
	}
	return Distance(&pointA, &pointB)
}
// http://en.wikipedia.org/wiki/Haversine_formula
func Distance(pointA, pointB *domain.Coordinate) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = pointA.Latitude * math.Pi / 180
	lo1 = pointA.Longitude * math.Pi / 180
	la2 = pointB.Latitude * math.Pi / 180
	lo2 = pointB.Longitude * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2 - la1) + math.Cos(la1) * math.Cos(la2) * hsin(lo2 - lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta / 2), 2)
}
