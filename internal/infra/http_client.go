package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPClient struct {
	client  *http.Client
	headers map[string]string
}

func NewHTTPClient(headers map[string]string, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client:  &http.Client{Timeout: timeout},
		headers: headers,
	}
}

func (c *HTTPClient) Get(url string, response interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	c.setHeaders(req)
	return c.doRequest(req, response)
}

func (c *HTTPClient) Post(url string, body interface{}, response interface{}, customHeader map[string]string) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	if len(customHeader) > 0 {
		c.setCustomHeader(req, customHeader)
	} else {
		c.setHeaders(req)
	}

	req.Header.Set("Content-Type", "application/json")
	return c.doRequest(req, response)
}

func (c *HTTPClient) Put(url string, body interface{}, response interface{}) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	return c.doRequest(req, response)
}

func (c *HTTPClient) Delete(url string, response interface{}) error {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	c.setHeaders(req)
	return c.doRequest(req, response)
}

func (c *HTTPClient) setHeaders(req *http.Request) {
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
}

func (c *HTTPClient) setCustomHeader(req *http.Request, customMap map[string]string) {
	headers := mergeMaps(customMap, c.headers)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

func (c *HTTPClient) doRequest(req *http.Request, response interface{}) error {
	fmt.Printf("request entrnado -> %v \n", req)

	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Printf("error in request -> %v", err)
		return err
	}
	fmt.Printf("request saindo -> %v \n", resp)

	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 202 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if response != nil {
		if err := json.Unmarshal(bodyBytes, response); err != nil {
			return err
		}
	}

	return nil
}

func mergeMaps(m1, m2 map[string]string) map[string]string {
	result := make(map[string]string)

	for k, v := range m1 {
		result[k] = v
	}

	for k, v := range m2 {
		result[k] = v
	}

	return result
}
