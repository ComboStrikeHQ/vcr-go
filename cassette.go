package vcr

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	jsonData, _ := json.Marshal(currentCassette)

	var jsonOut bytes.Buffer
	json.Indent(&jsonOut, jsonData, "", "  ")

	os.MkdirAll("fixtures/vcr", 0755)
	err := ioutil.WriteFile(c.fileName(), jsonOut.Bytes(), 0644)
	if err != nil {
		panic("VCR: Cannot write cassette file!")
	}
}

func panicEpisodeMismatch(request *vcrRequest, field string, expected string, actual string) {
	panic(fmt.Sprintf(
		"VCR: Problem with Episode for %s %s\n  Episode %s does not match:\n  expected: %s\n  but got: %s",
		request.Method,
		request.URL,
		field,
		expected,
		actual,
	))
}

func (c *cassette) matchEpisode(request *vcrRequest) *episode {
	if len(c.Episodes) == 0 {
		panic("VCR: No more episodes!")
	}

	e := c.Episodes[0]
	expected := e.Request

	if expected.Method != request.Method {
		panicEpisodeMismatch(request, "Method", expected.Method, request.Method)
	}

	if expected.URL != request.URL {
		panicEpisodeMismatch(request, "URL", expected.URL, request.URL)
	}

	if !reflect.DeepEqual(expected.Body, request.Body) {
		panicEpisodeMismatch(request, "Body", string(expected.Body[:]), string(request.Body[:]))
	}

	c.Episodes = c.Episodes[1:]
	return &e
}
