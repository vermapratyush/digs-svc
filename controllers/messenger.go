package controllers

import (
	"digs/domain"
	"encoding/json"
	"strings"
	"strconv"
	"digs/models"
	"errors"
)

type MessengerController struct {
	HttpBaseController
}

func (this *MessengerController) Post()  {
	var request domain.MessageSendRequest
	this.Super(&request.BaseRequest)

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &request)
	if err != nil {
		this.Serve500(err)
		return
	}
	request.SessionID = this.Ctx.Input.Header("SID")
	request.UserAgent = this.Ctx.Input.UserAgent()

	//Write to database
	_, err = models.CreateMessage(request.Username, request.Location, request.Body)
	if err != nil {
		this.Serve500(err)
		return
	}
	this.Serve304()
}

func (this *MessengerController) Get()  {
	var request domain.MessageGetRequest
	this.Super(&request.BaseRequest)

	err := populateGetParams(this, &request)
	if err != nil {
		this.Serve500(err)
		return
	}
	request.SessionID = this.Ctx.Input.Header("SID")
	request.UserAgent = this.Ctx.Input.UserAgent()
	//Get from database

	messages, err := models.GetMessages(request.Distance, request.Location)
	if err != nil {
		this.Serve500(err)
		return
	}
	this.Serve200(messages)
}

func populateGetParams(this *MessengerController, req *domain.MessageGetRequest) error {
	var err error
	req.Username = this.GetString("username")
	req.Distance, err = strconv.ParseInt(this.GetString("distance"), 10, 64)
	if err != nil {
		return errors.New("Distance invalid.")
	}
	locationParam := this.GetString("location")
	locationArray := strings.Split(locationParam, ",")
	if len(locationArray) != 2 {
		return errors.New("Location format invalid. Please specify longitude,latitude")
	}
	req.Location.Coordinates = []float64{-1.0, -1.0}
	req.Location.Coordinates[0], err = strconv.ParseFloat(locationArray[0], 64)
	if err != nil {
		return errors.New("Location format invalid. Please specify longitude correctly")
	}
	req.Location.Coordinates[1], err = strconv.ParseFloat(locationArray[1], 64)
	if err != nil {
		return errors.New("Location format invalid. Please specify latitude correctly")
	}
	return err
}

