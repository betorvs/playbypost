package cli

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

type Cli struct {
	baseURL string
	headers map[string]string
	*http.Client
}

func New(base string) *Cli {
	return &Cli{
		baseURL: base,
		headers: make(map[string]string),
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func NewHeaders(base, user, token string) *Cli {
	return &Cli{
		baseURL: base,
		headers: makeHeaders(user, token, "", "", ""),
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c *Cli) UpdateURL(s string) {
	if s != "" {
		c.baseURL = s
	}
}

// user, token : HeaderUsername, HeaderToken
// id, story, playerID: HeaderUserID, HeaderStory, HeaderPlayerID
func makeHeaders(user, token, id, story, channel string) map[string]string {
	headers := make(map[string]string)
	if id != "" {
		headers[types.HeaderUserID] = id
	}
	if story != "" {
		headers[types.HeaderStory] = story
	}
	if channel != "" {
		headers[types.HeaderStoryChannel] = channel
	}
	// if playerID != 0 {
	// 	headers[types.HeaderPlayerID] = fmt.Sprintf("%d", playerID)
	// }
	if user != "" {
		headers[types.HeaderUsername] = user
	}
	if token != "" {
		headers[types.HeaderToken] = token
	}
	return headers
}

func (c *Cli) getGeneric(kind string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/v1/%s", c.baseURL, kind)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.headers != nil {
		for k, v := range c.headers {
			if k != "" && v != "" {
				req.Header.Set(k, v)
			}
		}
	}
	resp, err := c.Do(req)
	if err != nil {
		return []byte{}, err
	}
	reqBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return reqBody, fmt.Errorf("return code %v", resp.StatusCode)
	}
	return reqBody, nil
}

func (c *Cli) postGeneric(kind string, body []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/api/v1/%s", c.baseURL, kind)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.headers != nil {
		for k, v := range c.headers {
			if k != "" && v != "" {
				req.Header.Set(k, v)
			}
		}
	}
	resp, err := c.Do(req)
	if err != nil {
		return []byte{}, err
	}
	respBody, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return []byte{}, err2
	}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
		return respBody, nil
	}
	return respBody, fmt.Errorf("status code not expected %d", resp.StatusCode)
}

func (c *Cli) postGenericWithHeaders(kind string, body []byte, headers map[string]string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/v1/%s", c.baseURL, kind)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		if k != "" && v != "" {
			req.Header.Set(k, v)
		}
	}

	resp, err := c.Do(req)
	if err != nil {
		return []byte{}, err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	if resp.StatusCode == http.StatusOK {
		return respBody, nil
	}
	return respBody, fmt.Errorf("status code not expected %d", resp.StatusCode)
}

func (c *Cli) putEmptyBodyGeneric(kind string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/v1/%s", c.baseURL, kind)
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.headers != nil {
		for k, v := range c.headers {
			if k != "" && v != "" {
				req.Header.Set(k, v)
			}
		}
	}
	resp, err := c.Do(req)
	if err != nil {
		return []byte{}, err
	}
	respBody, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return []byte{}, err2
	}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
		return respBody, nil
	}
	return respBody, fmt.Errorf("status code not expected %d", resp.StatusCode)
}
