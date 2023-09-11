package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"stringinator-go/constants"
	"stringinator-go/models"
	"stringinator-go/utils"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type StringInate struct{}

type Property struct {
	Count      int  `json:"count"`
	Length     int  `json:"length"`
	Palindrome bool `json:"palindrome"`
	Shortest   bool `json:"shortest"`
	Longest    bool `json:"longest"`
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

var defaultLog = utils.ConfigureLogs()

var sleepTime = 10 * time.Second

func init() {

	defaultLog.Println("controller/stringinate.go init() entering")
	defer defaultLog.Println("controller/stringinate.go init() Leaving")
	mu = &sync.Mutex{}

	collection = make(map[string]Property)

	monitorStringCollection()
}

func allLongestStrings() {

	defaultLog.Println("controller/stringinate.go allLongestStrings() entering")
	defer defaultLog.Println("controller/stringinate.go allLongestStrings() Leaving")

	// if len(strList) == 0 {
	// 	return
	// }
	defaultLog.Println("controller/stringinate.go allShortestStrings() proceeds calculating longest string")
	mu.Lock()

	var maxLength int
	if len(strList) > 0 {
		maxLength = len(strList[0])
	}

	for _, value := range strList {
		if len(value) > maxLength {
			maxLength = len(value)
			property := collection[value]
			property.Longest = true

			collection[value] = property
		}
	}
	defaultLog.Println("controller/stringinate.go allShortestStrings() going to sleep")
	time.Sleep(sleepTime)

	mu.Unlock()
}

func allShortestStrings() {

	defaultLog.Println("controller/stringinate.go allShortestStrings() entering")
	defer defaultLog.Println("controller/stringinate.go allShortestStrings() Leaving")

	// if len(strList) == 0 {
	// 	return
	// }
	defaultLog.Println("controller/stringinate.go allShortestStrings() proceeds calculating shortest string")

	mu.Lock()
	var minLength int
	if len(strList) > 0 {
		minLength = len(strList[0])
	}

	for _, value := range strList {
		if len(value) < minLength {
			minLength = len(value)
			property := collection[value]
			property.Shortest = true

			collection[value] = property
		}
	}
	time.Sleep(sleepTime)

	mu.Unlock()
}

func monitorStringCollection() {
	go func() {
		for {
			allShortestStrings()
			allLongestStrings()
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
	// defaultLog.Trace("controllers/key_controller:Retrieve() Entering")
	// defer defaultLog.Trace("controllers/key_controller:Retrieve() Leaving")

	input := mux.Vars(r)["input"]
	if input != "" {
		return nil, http.StatusBadRequest, errors.New("input is missing")
	}

	var response StringResponse

	if property, value := collection[input]; value {
		response.Value = input
		response.Property = property
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
