package models

import (
	"github.com/afex/hystrix-go/hystrix"
	"digs/common"
	"fmt"
	"digs/logger"
	"encoding/json"
	"net/http"
)

type FourSquareResponse struct {
	Response FourSquareVenues `json:"response" bson:"response"`
}

type FourSquareVenues struct {
	Venues []FourSquareVenue `json:"venues" bson:"venues"`
	Venue  FourSquareVenue `json:"venue" bson:"venue"`
}

type FourSquareVenue struct {
	Id string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	VenueStats FourSquareVenueStats `json:"stats" bson:"VenueStats"`
}

type FourSquareVenueStats struct {
	CheckinsCount int64 `json:"checkinsCount" bson:"checkinsCount"`
	UsersCount int64 `json:"usersCount" bson:"usersCount"`
	TipCount int64 `json:"tipCount" bson:"tipCount"`
}


func SearchFourSquareVenue(longitude, latitude, radius float64) FourSquareResponse {
	results := FourSquareResponse{}

	_ = hystrix.Do(common.FourSquareAPI, func() error {

		url := fmt.Sprintf("https://api.foursquare.com/v2/venues/search?ll=%f,%f&radius=%f&client_id=%s&client_secret=%s&v=%s",
			latitude, longitude, radius,
			common.FourSquare_API_CLIENT_ID, common.FourSquare_API_CLIENT_SECRET, common.FourSquare_API_CLIENT_VERSION)
		logger.Debug("FourSquareAPI|URL=", url)

		resp, err := http.Get(url)
		defer resp.Body.Close()
		if err != nil {
			logger.Error("FourSquareError|URL=", url, "|err=", err)
			return err
		}

		err = json.NewDecoder(resp.Body).Decode(&results)
		if err != nil {
			logger.Error("FourSquareDeserializeError|URL=", url, "|err=", err)
			return err
		}

		return nil
	}, nil)

	return results
}

func GetFourSquareVenue(id string) FourSquareResponse {
	results := FourSquareResponse{}

	_ = hystrix.Do(common.MeetupAPI, func() error {

		url := fmt.Sprintf("https://api.foursquare.com/v2/venues/%s?client_id=%s&client_secret=%s&v=%s",
			id, common.FourSquare_API_CLIENT_ID, common.FourSquare_API_CLIENT_SECRET, common.FourSquare_API_CLIENT_VERSION)
		logger.Debug("FourSquareAPI|URL=", url)

		resp, err := http.Get(url)
		defer resp.Body.Close()
		if err != nil {
			logger.Error("FourSquareError|URL=", url, "|err=", err)
			return err
		}


		err = json.NewDecoder(resp.Body).Decode(&results)
		if err != nil {
			logger.Error("FourSquareDeserializeError|URL=", url, "|err=", err)
			return err
		}

		return nil
	}, nil)

	return results
}