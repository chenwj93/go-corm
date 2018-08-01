package eureka

import (
	"encoding/json"
	"strings"
	"errors"
)

func (c *Client) RegisterInstance(appId string, instanceInfo *InstanceInfo) error {
	values := []string{"apps", appId}
	path := strings.Join(values, "/")
	instance := &Instance{
		Instance: instanceInfo,
	}
	body, err := json.Marshal(instance)
	if err != nil {
		return err
	}

	resp, err := c.Post(path, body)
	if err == nil {
		//fmt.Println(resp.StatusCode, string(resp.Body))
		if resp.StatusCode > 299 || resp.StatusCode < 200 {
			err = errors.New(string(resp.Body))
		}
	}
	return err
}
