package hdsstorage

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/tidwall/gjson"
)

//锁定分配给用户的资源组
func (s *Storage) LockResourceGroup(waitTime int) Job {
	/*
		请求参数：
		  waitTime int, lock timeout in seconds; range: 0-7200; default: 0
	*/
	url := BaseURL + "/v1" + s.StorageDeviceId + "/services/resource-group-service/actions/lock/invoke"

	tmpMap := make(map[string]interface{})
	tmpMap["parameters"] = map[string]int{"waitTime": waitTime}

	content, _ := json.Marshal(tmpMap)
	postContent := bytes.NewBuffer(content)

	response, err := PostRequest(url, s.Token, postContent)
	if err != nil {
		if response.StatusCode == 503 {
			log.Println("LockResourceGroup error:", err)
			return Job{}
		}
		log.Fatalf("LockResourceGroup error: %v\n", err)
		// return Job{}
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	job := Job{}
	json.Unmarshal(body, &job)

	return job
}

// 解锁resource group
func (s *Storage) UnlockResourceGroup() Job {
	url := BaseURL + "/v1" + s.StorageDeviceId + "/services/resource-group-service/actions/unlock/invoke"
	response, err := PostRequest(url, s.Token, nil)
	if err != nil {
		if response.StatusCode == 503 {
			log.Println("UnlockResourceGroup error:", err)
			return Job{}
		}
		log.Fatalf("UnlockResourceGroup error: %v\n", err)
		// return Job{}
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	job := Job{}
	json.Unmarshal(body, &job)

	return job
}

// 获取存储的所有resource group
func (s *Storage) GetAllResourceGroups() []ResourceGroup {
	url := BaseURL + "/v1/objects/storages" + s.StorageDeviceId + "/resource-groups"
	resp, err := GetRequest(url, s.Token)
	if err != nil {
		log.Printf("Get Resource Groups error: %v\n", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var rgroups []ResourceGroup
	for _, r := range gjson.Get(string(body), "data").Array() {
		rgroup := ResourceGroup{}
		json.Unmarshal([]byte(r.String()), &rgroup)
		rgroups = append(rgroups, rgroup)
	}
	return rgroups
}

// 定义ResourGroup结构体
type ResourceGroup struct {
	ResourceGroupId   int64    `json:"resourceGroupId"`
	ResourceGroupName string   `json:"resourceGroupName"`
	LockStatus        string   `json:"lockStatus"`
	VirtualStorageId  int64    `json:"virtualStorageId"`
	LdevIds           []int64  `json:"ldevIds"`
	ParityGroupIds    []string `json:"parityGroupIds"`
	PortIds           []string `json:"portIds"`
	HostGroupIds      []string `json:"hostGroupIds"`
	LockOwner         string   `json:"lockOwner"`
	LockHost          string   `json:"lockHost"`
}
