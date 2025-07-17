package cli

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	storyteller string = "writer"
)

func (c *Cli) GetWriter() ([]types.Writer, error) {
	var storytellers []types.Writer
	body, err := c.getGeneric(storyteller)
	if err != nil {
		return storytellers, err
	}
	err = json.Unmarshal(body, &storytellers)
	if err != nil {
		return storytellers, err
	}
	return storytellers, nil
}

func (c *Cli) GetWriterByUsername(username string) (types.Writer, error) {
	body, err := c.getGeneric(storyteller)
	if err != nil {
		return types.Writer{}, err
	}
	writers := []types.Writer{}
	err = json.Unmarshal(body, &writers)
	if err != nil {
		return types.Writer{}, err
	}
	for _, writer := range writers {
		if writer.Username == username {
			return writer, nil
		}
	}
	return types.Writer{}, nil
}

func (c *Cli) CreateWriter(username, password string) ([]byte, error) {
	u := types.Writer{
		Username: username,
		Password: password,
	}
	body, err := json.Marshal(u)
	if err != nil {
		return []byte{}, err
	}
	res, err := c.postGeneric(storyteller, body)
	return res, err
}

func (c *Cli) GetUserByUserID(id string) (types.User, error) {
	body, err := c.getGeneric("user/" + id)
	if err != nil {
		return types.User{}, err
	}
	users := types.User{}
	err = json.Unmarshal(body, &users)
	if err != nil {
		return types.User{}, err
	}
	return users, nil
}

func (c *Cli) CreateWriterUserAssociation(writerID, userID int) (int, error) {
	obj := types.WriterUserAssociation{WriterID: writerID, UserID: userID}
	body, err := json.Marshal(obj)
	if err != nil {
		return 0, err
	}
	resp, err := c.postGeneric("writer/user", body)
	if err != nil {
		return 0, err
	}
	var msg types.Msg
	err = json.Unmarshal(resp, &msg)
	if err != nil {
		return 0, err
	}
	var id int
	_, err = fmt.Sscanf(msg.Msg, "writer user association id %d", &id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (c *Cli) DeleteWriterUserAssociation(id int) error {
	url := fmt.Sprintf("%s/api/v1/writer/association/%d", c.baseURL, id)
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
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Println("error closing response body", err.Error())
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code not expected %d", resp.StatusCode)
	}
	return nil
}

func (c *Cli) CheckWriterUserAssociationExists(writerID, userID int) (bool, error) {
	// This function is not exposed via an API endpoint, so we'll simulate it by trying to create an association
	// and checking the error message. This is a workaround for testing purposes.
	associations, err := c.GetWriterUsersAssociation()
	if err != nil {
		return false, err
	}
	for _, a := range associations {
		if a.WriterID == writerID && a.UserID == userID {
			return true, nil
		}
	}
	return false, nil
}

func (c *Cli) GetWriterUsersAssociation() ([]types.WriterUserAssociation, error) {
	body, err := c.getGeneric("writer/association")
	if err != nil {
		return []types.WriterUserAssociation{}, err
	}
	associations := []types.WriterUserAssociation{}
	err = json.Unmarshal(body, &associations)
	if err != nil {
		return []types.WriterUserAssociation{}, err
	}
	return associations, nil
}
