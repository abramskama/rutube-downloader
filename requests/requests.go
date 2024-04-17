package requests

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetResponse(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	cl := http.DefaultClient
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http status code != 200, url: %s, code: %d", url, resp.StatusCode)
	}
	return resp, nil
}

func GetResponseBody(url string) ([]byte, error) {
	resp, err := GetResponse(url)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()
	return respBody, nil
}
