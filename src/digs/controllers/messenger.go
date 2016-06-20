package controllers

import (
	"digs/domain"
	"strings"
	"strconv"
	"digs/models"
	"errors"
	"encoding/json"
)

type MessengerController struct {
	HttpBaseController
}

func (this *MessengerController) Post() {
	var request domain.MessageSendRequest
	this.Super(&request.BaseRequest)
	json.Unmarshal(this.Ctx.Input.RequestBody, &request)

	//Write to database
	_, err := models.CreateMessage(request.Username, request.Location[0], request.Location[1], request.Body)
	if err != nil {
		this.Serve500(err)
		return
	}
	this.Serve304()
}

func (this *MessengerController) Get() {
	var request domain.MessageGetRequest
	this.Super(&request.BaseRequest)

	err := populateGetParams(this, &request)
	if err != nil {
		this.Serve500(err)
		return
	}

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
	req.Location = make([]float64, 2)
	req.Location[0], err = strconv.ParseFloat(locationArray[0], 64)
	if err != nil {
		return errors.New("Location format invalid. Please specify longitude correctly")
	}
	req.Location[1], err = strconv.ParseFloat(locationArray[1], 64)
	if err != nil {
		return errors.New("Location format invalid. Please specify latitude correctly")
	}
	return err
}

