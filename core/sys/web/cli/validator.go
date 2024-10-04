package cli

import "fmt"

func (c *Cli) ValidatorPut(obj string, id int) ([]byte, error) {
	kind := fmt.Sprintf("validator/%s/%d", obj, id)
	res, err := c.putEmptyBodyGeneric(kind)
	return res, err
}

func (c *Cli) GetValidator() ([]byte, error) {
	kind := "validator"
	res, err := c.getGeneric(kind)
	return res, err
}
