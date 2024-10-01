package cli

import (
	"encoding/json"
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (c *Cli) Ping() error {
	var a types.Msg
	body, err := c.getGeneric("validate")
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &a)
	if err != nil {
		return err
	}
	if a.Msg != "authenticated" {
		return fmt.Errorf("error not authenticated: %s", a.Msg)
	}
	return nil
}
