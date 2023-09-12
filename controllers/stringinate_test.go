package controllers

import "testing"

func Test_allShortestStrings(t *testing.T) {

	strList = []string{"ihh", "hi", "a"}

	tests := []struct {
		name string
	}{
		{name: "valid case"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allShortestStrings()
		})
	}
}
