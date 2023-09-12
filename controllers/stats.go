package controllers

import "net/http"

type Stats struct{}

type StatsResponse struct {
	Total       int      `json:"num-of-strings"`
	Shortest    string   `json:"shortest-string"`
	Longest     string   `json:"longest-string"`
	Palindromes []string `json:"palidromes"`
	MostUsed    string   `json:"most-used"`
}

func (st Stats) GetStats(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var response StatsResponse
	response.Total = len(strList)

	allShortestStrings()
	allLongestStrings()
	mostUsedString()

	for _, value := range strList {

		property := collection[value]
		if property.Shortest {
			response.Shortest = value
		}
		if property.Longest {
			response.Longest = value
		}
		if property.Palindrome {
			response.Palindromes = append(response.Palindromes, value)
		}

		response.MostUsed = mostUsed
	}

	return response, http.StatusOK, nil
}
