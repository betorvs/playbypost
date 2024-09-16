package cli

import (
	"fmt"
	"net/http"
)

func (c *Cli) putGeneric(kind string) error {
	url := fmt.Sprintf("%s/api/v1/%s", c.baseURL, kind)
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("status code not expected %d", resp.StatusCode)
}

func (c *Cli) PutReload() error {
	err := c.putGeneric("reload")
	return err
}
