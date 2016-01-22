# VCR

[![Circle CI](https://circleci.com/gh/ad2games/vcr-go.svg?style=svg)](https://circleci.com/gh/ad2games/vcr-go)

But for golang.

## Features

- Easy to use
- No flags (records when no cassete file is present, otherwise replays)
- Filter sensitive data (such as API keys)
- Modify requests before they are sent (e.g. to add row limits)
- Ignores request headers (so a timestamp or crypto signature does not break your cassette files)
- No dependencies

## Examples

Simple:
```go
package examples

import (
  "github.com/ad2games/vcr-go"
  "net/http"
  "testing"
)

func TestSomeFeature(t *testing.T) {
  // This will automatically create a cassette file at ./fixtures/vcr/my_feature.json
  vcr.Start("my_feature", nil)
  defer vcr.Stop()

  resp, _ := http.Get("http://google.com/")
  if resp.StatusCode != 200 {
    t.Fatalf("Status code is not 200!")
  }
}
```


Advanced:
```go
func modifier(req *http.Request) {
  req.Header["Content-Type"] = []string{"text/xml"}

  vcr.ModifyHTTPRequestBody(req, func(body string) string {
    return body + "<row_limit>1</row_limit>"
  })
}

func TestSomeOtherFeature(t *testing.T) {
  vcr.Start("my_other_feature", modifier)
  vcr.FilterData(os.Getenv("SOME_API_KEY"), "some-api-key")
  defer vcr.Stop()

  http.Get("http://example.com/")
}
```

## License

MIT, see LICENSE.txt

## Contributing

1. Fork it (https://github.com/ad2games/vcr-go/fork)
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request
