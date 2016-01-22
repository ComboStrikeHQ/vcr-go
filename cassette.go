package vcr

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
)

type cassette struct {
	name     string
	Episodes []episode
}

type episode struct {
	Request  *vcrRequest
	Response *vcrResponse
}

func (c *cassette) fileName() string {
	return "fixtures/vcr/" + c.name + ".json"
}

func (c *cassette) exists() bool {
	_, err := os.Stat(c.fileName())
	return err == nil
}

func (c *cassette) read() {
	jsonData, _ := ioutil.ReadFile(c.fileName())
	err := json.Unmarshal(jsonData, c)
	if err != nil {
		panic("VCR: Cannot parse JSON!")
	}
}

func (c *cassette) write() {
	jsonData, err := json.Marshal(currentCassette)

	var jsonOut bytes.Buffer
	json.Indent(&jsonOut, jsonData, "", "  ")

	os.MkdirAll("fixtures/vcr", 0755)
	err = ioutil.WriteFile(c.fileName(), jsonOut.Bytes(), 0644)
	if err != nil {
		panic("VCR: Cannot write cassette file!")
	}
}
func (c *cassette) matchEpisode(request *vcrRequest) *episode {
	if len(c.Episodes) == 0 {
		panic("VCR: No more episodes!")
	}

	episode := c.Episodes[0]
	if !reflect.DeepEqual(episode.Request, request) {
		panic("VCR: Episodes do not match!")
	}

	c.Episodes = c.Episodes[1:]
	return &episode
}
