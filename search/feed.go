package search

import (
	"encoding/json"
	"os"
)

type Feed struct {
	Name string `json:"site"`
	URI  string `json:"link"`
	Type string `json:"type"`
}

const dataFile="data/data.json"
func RetrieveFeeds() ([]*Feed, error) {

	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	//Schedule the file to be closed once the function returns
	defer file.Close()
	var feeds []*Feed
	err = json.NewDecoder(file).Decode(&feeds)

	return feeds, err
}
