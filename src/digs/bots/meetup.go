package bots

import (
	"github.com/afex/hystrix-go/hystrix"
	"digs/common"
	"time"
	"net/http"
	"digs/logger"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

func GetMeetup(longitude, latitude, radius float64) MeetupResults {

	results := MeetupResults{}
	_ = hystrix.Do(common.MeetupAPI, func() error {
		today := time.Now().Unix() * 1000
		tomorrow := today + (1000 * 24 * 60 * 60)
		url := fmt.Sprintf("https://api.meetup.com/2/concierge?lon=%f&lat=%f&radius=%f&time=%d,%d&sign=true&key=%s", longitude, latitude, radius * 0.621371 / 1000, today, tomorrow, common.Meetup_API_KEY)
		logger.Debug("MeetupAPI|URL=", url)

		resp, err := http.Get(url)
		defer resp.Body.Close()
		if err != nil {
			logger.Error("MeetupError|URL=", url, "|err=", err)
			return err
		}


		err = json.NewDecoder(resp.Body).Decode(&results)
		if err != nil {
			logger.Error("MeetupDeserializeError|URL=", url, "|err=", err)
			return err
		}

		return nil
	}, nil)

	return results
}

func ShortenUrl(url string) string {
	result := ""
	err := hystrix.Do(common.BitlyAPI, func() error {
		url := "http://api.bitly.com/v3/shorten?login=" + common.Bitly_Login + "&apiKey=" + common.Bitly_API_KEY + "&longUrl=" + url + "&format=txt"
		logger.Debug("BitLyAPI|URL=", url)
		resp, err := http.Get(url)
		defer resp.Body.Close()

		if err != nil {
			logger.Error("BITLYError|URL=", url,"|err=", err)
			return err
		}
		data, err := ioutil.ReadAll(resp.Body)
		result = string(data)
		return err
	}, nil)

	if err != nil {
		result = url
	}
	return result
}