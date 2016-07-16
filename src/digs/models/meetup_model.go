package models


type MeetupEvents []MeetupEvent
type ByYesCount struct{ MeetupEvents }
func (s MeetupEvents) Len() int      { return len(s) }
func (s MeetupEvents) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByYesCount) Less(i, j int) bool { return s.MeetupEvents[i].YesCount > s.MeetupEvents[j].YesCount }

type MeetupResults struct {
	Results []MeetupEvent `json:"results" bson:"results"`
	MetaData Meta `json:"meta" bson:"meta"`
}

type Meta struct {
	NextPage string `json:"next" bson:"next"`
	TotalCount int64 `json:"total_count" bson:"total_count"`
	Count int64 `json:"count" bson:"count"`
}

type MeetupEvent struct {
	EventLatLong Venue `json:"venue" bson:"venue"`
	YesCount int64 `json:"yes_rsvp_count" bson:"yes_rsvp_count"`
	Title string `json:"name" bson:"name"`
	EventUrl string `json:"event_url" bson:"event_url"`
}

type Venue struct  {
	Longitude float64 `json:"lon" bson:"lon"`
	Latitude float64 `json:"lat" bson:"lat"`
}