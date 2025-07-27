package cli

import (
	"bytes"
	"encoding/json"
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
	fmt.Println("headers", headers)
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

func (c *Cli) putGenericWithHeaders(kind string, body []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/api/v1/%s", c.baseURL, kind)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
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

// Session management methods
func (c *Cli) GetAllSessions() ([]types.Session, error) {
	body, err := c.getGeneric("session")
	if err != nil {
		return nil, err
	}
	var sessions []types.Session
	err = json.Unmarshal(body, &sessions)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (c *Cli) GetSessionEvents() ([]types.SessionEvent, error) {
	body, err := c.getGeneric("session/events")
	if err != nil {
		return nil, err
	}
	var events []types.SessionEvent
	err = json.Unmarshal(body, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (c *Cli) DeleteExpiredSessions() error {
	_, err := c.putEmptyBodyGeneric("session/cleanup")
	return err
}

func (c *Cli) DeleteSessionByID(sessionID int64) error {
	url := fmt.Sprintf("%s/api/v1/session/%d", c.baseURL, sessionID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
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
		return err
	}
	if resp.StatusCode != http.StatusOK {
		reqBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("return code %v: %s", resp.StatusCode, string(reqBody))
	}
	return nil
}
