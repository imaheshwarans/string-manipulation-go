package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"stringinator-go/config"
	"stringinator-go/constants"
	"stringinator-go/models"
	"stringinator-go/utils"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type StringInate struct{}

type Property struct {
	Count      int  `json:"count,omitempty"`
	Length     int  `json:"length,omitempty"`
	Palindrome bool `json:"palindrome,omitempty"`
	Shortest   bool `json:"shortest,omitempty"`
	Longest    bool `json:"longest,omitempty"`
}

type StringRequest struct {
	Input string `json:"input"`
}

type StringResponse struct {
	Value    string    `json:"string,omitempty"`
	Property *Property `json:"property,omitempty"`
}

var strList []string
var collection map[string]Property
var mu *sync.Mutex

var defaultLog *logrus.Logger
var sleepTime = 10 * time.Second
var shortest, longest, mostUsed string

func init() {
	config, _ := config.LoadConfiguration()
	defaultLog = utils.ConfigureLogs(config.LogLevel)

	defaultLog.Trace("controller/stringinate.go init() entering")
	defer defaultLog.Trace("controller/stringinate.go init() Leaving")
	mu = &sync.Mutex{}

	collection = make(map[string]Property)

	monitorStringCollection()
}

func allLongestStrings() {

	mu.Lock()
	defer mu.Unlock()

	defaultLog.Trace("controller/stringinate.go allLongestStrings() entering")
	defer defaultLog.Trace("controller/stringinate.go allLongestStrings() Leaving")

	defaultLog.Info("controller/stringinate.go allShortestStrings() proceeds calculating longest string")

	if longest != "" {
		property := collection[longest]
		property.Longest = false
		collection[longest] = property
	}

	if len(strList) > 0 {
		longest = strList[0]
	}

	for _, value := range strList {
		if len(value) > len(longest) {
			longest = value
		}
	}
	property := collection[longest]
	property.Longest = true
	collection[longest] = property

	defaultLog.Info("controller/stringinate.go allShortestStrings() going to sleep")
}

func allShortestStrings() {

	defaultLog.Trace("controller/stringinate.go allShortestStrings() entering")
	defer defaultLog.Trace("controller/stringinate.go allShortestStrings() Leaving")

	mu.Lock()
	defer mu.Unlock()

	defaultLog.Info("controller/stringinate.go allShortestStrings() proceeds calculating shortest string")

	if shortest != "" {
		property := collection[shortest]
		property.Shortest = false
		collection[shortest] = property
	}

	if len(strList) > 0 {
		shortest = strList[0]
	}

	for _, value := range strList {
		if len(value) < len(shortest) {
			shortest = value
		}
	}
	property := collection[shortest]
	property.Shortest = true
	collection[shortest] = property
}

func mostUsedString() {
	defaultLog.Trace("controller/stringinate.go mostUsedString() entering")
	defer defaultLog.Trace("controller/stringinate.go mostUsedString() Leaving")

	mu.Lock()
	defer mu.Unlock()

	defaultLog.Infof("controller/stringinate.go mostUsedString() proceeds calculating mostUsedString")

	if mostUsed != "" {
		mostUsed = ""
	}
	var count int

	if len(strList) > 0 {
		mostUsed = strList[0]
		count = collection[strList[0]].Count
	}

	for _, value := range strList {
		if count < collection[value].Count {
			count = collection[value].Count
			mostUsed = value
		}
	}

	defaultLog.Infof("controller/stringinate.go mostUsedString() most used string is %s", mostUsed)
}

func monitorStringCollection() {
	go func() {
		for {
			allShortestStrings()
			allLongestStrings()
			mostUsedString()
			time.Sleep(sleepTime)
		}
	}()
}

func (str StringInate) Create(w http.ResponseWriter, request *http.Request) (interface{}, int, error) {

	if request.Header.Get("Content-Type") != constants.HTTPMediaTypeJson {
		return nil, http.StatusUnsupportedMediaType, &models.HandledError{Message: "Invalid Content-Type"}
	}

	if request.ContentLength == 0 {
		defaultLog.Error("controllers/stringinate:Create() The request body was not provided")
		return nil, http.StatusBadRequest, &models.HandledError{Message: "The request body was not provided"}
	}

	var requestString StringRequest
	// Decode the incoming json data to note struct
	dec := json.NewDecoder(request.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&requestString)
	if err != nil {
		return nil, http.StatusBadRequest, &models.HandledError{Message: "Unable to decode JSON request body"}
	}

	if requestString.Input == "" {
		return nil, http.StatusBadRequest, &models.HandledError{Message: "Empty string provided"}
	}

	if property, seen := collection[requestString.Input]; seen {
		property.Count = property.Count + 1
		collection[requestString.Input] = property
		return StringResponse{Value: requestString.Input, Property: &property}, http.StatusOK, nil
	}

	property := &Property{
		Count:      1,
		Length:     len(requestString.Input),
		Palindrome: Palindrome(requestString.Input),
	}

	response := StringResponse{
		Value:    requestString.Input,
		Property: property,
	}

	strList = append(strList, requestString.Input)
	collection[requestString.Input] = *property

	return response, http.StatusCreated, nil
}

func (str StringInate) Get(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	defaultLog.Trace("controllers/stringinate:Get() Entering")
	defer defaultLog.Trace("controllers/stringinate:Get() Leaving")

	input := r.URL.Query().Get("input")
	if input == "" {
		return nil, http.StatusBadRequest, errors.New("input is missing")
	}

	var response StringResponse

	if property, value := collection[input]; value {
		response.Value = input
		response.Property = &property
	}

	return response, http.StatusOK, nil
}

func (str StringInate) GetAll(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	defaultLog.Trace("controllers/stringinate:GetAll() Entering")
	defer defaultLog.Trace("controllers/stringinate:GetAll() Leaving")

	var response []StringResponse

	for _, value := range strList {
		prp := collection[value]
		resp := StringResponse{
			Value:    value,
			Property: &prp,
		}
		response = append(response, resp)
	}

	return response, http.StatusOK, nil
}

func Palindrome(str string) bool {
	reversedStr := ""
	for i := len(str) - 1; i >= 0; i-- {
		reversedStr += string(str[i])
	}
	for i := range str {
		if str[i] != reversedStr[i] {
			return false
		}
	}
	return true
}
