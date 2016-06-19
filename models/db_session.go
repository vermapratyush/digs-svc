package models

import (
	"gopkg.in/mgo.v2"
	"fmt"
)

var DefaultDatabase = "heroku_qnx0661v"
var Session, _ = mgo.Dial(fmt.Sprintf("mongodb://node-js:node-js@ds015194.mlab.com:15194/%s", DefaultDatabase))
