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
	Count      int  `json:"count"`
	Length     int  `json:"length"`
	Palindrome bool `json:"palindrome,omitempty"`
	Shortest   bool `json:"shortest,omitempty"`
	Longest    bool `json:"longest,omitempty"`
}

type StringRequest struct {
	Input string `json:"input"`
}

type StringResponse struct {
	Value    string   `json:"string"`
	Property Property `json:"property"`
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

	defaultLog.Info("controller/stringinate.go init() entering")
	defer defaultLog.Info("controller/stringinate.go init() Leaving")
	mu = &sync.Mutex{}

	collection = make(map[string]Property)

	monitorStringCollection()
}

func allLongestStrings() {

	mu.Lock()
	defer mu.Unlock()

	defaultLog.Println("controller/stringinate.go allLongestStrings() entering")
	defer defaultLog.Println("controller/stringinate.go allLongestStrings() Leaving")

	defaultLog.Println("controller/stringinate.go allShortestStrings() proceeds calculating longest string")

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

	defaultLog.Println("controller/stringinate.go allShortestStrings() going to sleep")
}

func allShortestStrings() {

	defaultLog.Println("controller/stringinate.go allShortestStrings() entering")
	defer defaultLog.Println("controller/stringinate.go allShortestStrings() Leaving")

	mu.Lock()
	defer mu.Unlock()

	defaultLog.Println("controller/stringinate.go allShortestStrings() proceeds calculating shortest string")

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

func monitorStringCollection() {
	go func() {
		for {
			allShortestStrings()
			allLongestStrings()
			time.Sleep(sleepTime)
		}
	}()
}

func (str StringInate) Create(w http.ResponseWriter, request *http.Request) (interface{}, int, error) {

	if request.Header.Get("Content-Type") != constants.HTTPMediaTypeJson {
		return nil, http.StatusUnsupportedMediaType, &models.HandledError{Message: "Invalid Content-Type"}
	}

	if request.ContentLength == 0 {
		// secLog.Error("controllers/key_controller:Create() The request body was not provided")
		return nil, http.StatusBadRequest, &models.HandledError{Message: "The request body was not provided"}
	}

	var requestString StringRequest
	// Decode the incoming json data to note struct
	dec := json.NewDecoder(request.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&requestString)
	if err != nil {
		// secLog.WithError(err).Errorf("controllers/key_controller:Create() %s : Failed to decode request body as KeyRequest", commLogMsg.InvalidInputBadEncoding)
		return nil, http.StatusBadRequest, &models.HandledError{Message: "Unable to decode JSON request body"}
	}

	if requestString.Input == "" {
		return nil, http.StatusBadRequest, &models.HandledError{Message: "Empty string provided"}
	}

	property := Property{
		Length:     len(requestString.Input),
		Palindrome: Palindrome(requestString.Input),
	}

	response := StringResponse{
		Value:    requestString.Input,
		Property: property,
	}

	// if already present increase the count
	if property, seen := collection[requestString.Input]; seen {
		property.Count++
		collection[requestString.Input] = property
		return response, http.StatusOK, nil
	}

	property.Count = 1

	strList = append(strList, requestString.Input)
	collection[requestString.Input] = property

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
		response.Property = property
	}

	return response, http.StatusOK, nil
}

func (str StringInate) GetAll(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	defaultLog.Trace("controllers/stringinate:GetAll() Entering")
	defer defaultLog.Trace("controllers/stringinate:GetAll() Leaving")

	var response []StringResponse

	for _, value := range strList {
		resp := StringResponse{
			Value:    value,
			Property: collection[value],
		}
		response = append(response, resp)
	}

	return response, http.StatusOK, nil
}

func Palindrome(str string) bool {
	lastIdx := len(str) - 1
	// using for loop
	for i := 0; i < lastIdx/2 && i < (lastIdx-i); i++ {
		if str[i] != str[lastIdx-i] {
			return false
		}
	}
	return true
}
