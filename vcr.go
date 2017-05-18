package vcr

import (
	"net/http"
)

const (
	modeRecord = iota
	modeReplay
)

// RequestModifierFunc is a function that can be used to manipulate HTTP requests
// before they are sent to the server.
// Useful for adding row-limits in integration tests.
type RequestModifierFunc func(request *http.Request)

var currentFilterMap map[string]string
var currentRequestModifier RequestModifierFunc
var currentMode = modeRecord
var currentCassette *cassette

type roundTripper struct {
	originalRT http.RoundTripper
}

func init() {
	http.DefaultTransport = &roundTripper{originalRT: http.DefaultTransport}
}

// Start starts a VCR session with the given cassette name.
// Records episodes if the cassette file does not exists.
// Otherwise plays back recorded episodes.
func Start(cassetteName string, modFunc RequestModifierFunc) {
	if currentCassette != nil {
		panic("VCR: Session already started!")
	}

	currentCassette = &cassette{name: cassetteName}
	currentRequestModifier = modFunc
	currentFilterMap = make(map[string]string)

	if currentCassette.exists() {
		currentMode = modeReplay
		currentCassette.read()
	} else {
		currentMode = modeRecord
	}
}

// FilterData allows replacement of sensitive data with a dummy-string
func FilterData(original string, replacement string) {
	currentFilterMap[original] = replacement
}

// Stop stops the VCR session and writes the cassette file (when recording)
func Stop() {
	if currentMode == modeRecord {
		currentCassette.write()
	}

	currentCassette = nil
}

func (rt *roundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	vcrReq := newVCRRequest(request, currentFilterMap)
	var vcrRes *vcrResponse

	if currentCassette == nil {
		return rt.originalRT.RoundTrip(request)
	}

	if currentRequestModifier != nil {
		currentRequestModifier(request)
	}

	if currentMode == modeRecord {
		response, err := rt.originalRT.RoundTrip(request)
		if err != nil {
			return nil, err
		}
		vcrRes = newVCRResponse(response)

		e := episode{Request: vcrReq, Response: vcrRes}
		currentCassette.Episodes = append(currentCassette.Episodes, e)
	} else {
		e := currentCassette.matchEpisode(vcrReq)
		vcrRes = e.Response
	}

	return vcrRes.httpResponse(), nil
}
