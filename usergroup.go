package hdsstorage

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/tidwall/gjson"
)

type Usergroup struct {
	UserGroupObjectId   string   `json:"userGroupObjectId"`
	UserGroupId         string   `json:"userGroupId"`
	RoleNames           []string `json:"roleNames"`
	ResourceGroupIds    []int64  `json:"resourceGroupIds"`
	IsBuiltin           bool     `json:"isBuiltin"`
	HasAllResourceGroup bool     `json:"hasAllResourceGroup"`
}

func (s *Storage) GetAllUsergroups() []Usergroup {
	url := BaseURL + "/v1/objects/storages/" + s.StorageDeviceId + "/user-groups"
	resp, err := GetRequest(url, s.Token)
	if err != nil {
		log.Printf("Get Usergroup error: %v\n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Read body error %v\n", err)
	}
	var usergroups []Usergroup
	for _, usergroup := range gjson.Get(string(body), "data").Array() {
		tmpUsergroup := Usergroup{}
		json.Unmarshal([]byte(usergroup.String()), &tmpUsergroup)
		usergroups = append(usergroups, tmpUsergroup)
	}
	return usergroups
}
