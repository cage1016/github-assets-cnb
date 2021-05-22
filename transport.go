package cnbtest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/gabriel-vasile/mimetype"
)

type Transport struct {
	Token string
}

func NewTransport(options ...TransportOption) *Transport {
	t := &Transport{}

	for _, option := range options {
		option(t)
	}

	return t
}

// TransportOption sets an optional parameter for Transport.
type TransportOption func(*Transport)

func SetToken(token string) TransportOption {
	return func(t *Transport) {
		t.Token = token
	}
}

func (t *Transport) prepareRequest(url string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)
	if t.Token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("token %s", t.Token))
	}
	req.Header.Add("User-Agent", "metal3d-go-client")
	return req
}

func (t *Transport) Fetch(repo, asset, tag string) (id float64, err error) {
	if len(repo) == 0 {
		return 0, errors.New("No repository provided")
	}

	// command to call
	command := "releases/latest"
	if len(tag) > 0 {
		command = fmt.Sprintf("releases/tags/%s", tag)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", repo, command)

	// create a request with basic-auth
	req := t.prepareRequest(url)

	// Add required headers
	req.Header.Add("Accept", "application/vnd.github.v3.text-match+json")
	req.Header.Add("Accept", "application/vnd.github.moondragon+json")

	// call github
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return 0, fmt.Errorf("Repo: %s, Asset: %s Error while making request: %v", repo, asset, err)
	}

	// status in <200 or >299
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return 0, fmt.Errorf("Repo: %s, Asset: %s Error: %d, %v", repo, asset, resp.StatusCode, resp.Status)
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("Repo: %s, Asset: %s Error reading response: %v", repo, asset, err)
	}

	// prepare result
	result := make(map[string]interface{})
	json.Unmarshal(bodyText, &result)

	// filter
	result["assets"] = filter(result["assets"].([]interface{}), func(in interface{}) bool {
		if in.(map[string]interface{})["name"] == asset {
			return true
		}
		return false
	})

	if len(result["assets"].([]interface{})) == 0 {
		return 0, fmt.Errorf("Repo: %s, Asset: %s Asset not found", repo, asset)
	} else if len(result["assets"].([]interface{})) > 1 {
		return 0, fmt.Errorf("Repo: %s, Asset: %s Asset found more than one item", repo, asset)
	} else {
		return result["assets"].([]interface{})[0].(map[string]interface{})["id"].(float64), nil
	}
}

func (t *Transport) Drop(repo string, id float64) (io.ReadCloser, string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/assets/%.0f", repo, id)
	fmt.Printf("Download: %s\n", url)
	req := t.prepareRequest(url)

	req.Header.Add("Accept", "application/octet-stream")

	client := http.Client{}
	resp, _ := client.Do(req)

	disp := resp.Header.Get("Content-disposition")
	re := regexp.MustCompile(`filename=(.+)`)
	matches := re.FindAllStringSubmatch(disp, -1)

	if len(matches) == 0 || len(matches[0]) == 0 {
		log.Println("WTF: ", matches)
		log.Println(resp.Header)
		log.Println(req)
		return nil, "", errors.New("asset not found")
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close() //  must close
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	mtype, err := mimetype.DetectReader(bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, "", err
	}

	return resp.Body, mtype.String(), nil
}

func filter(arr []interface{}, predicate func(interface{}) bool) []interface{} {
	out := make([]interface{}, 0)

	for _, e := range arr {
		if predicate(e) {
			out = append(out, e)
		}
	}

	return out
}
