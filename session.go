package hdsstorage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/tidwall/gjson"
)

// 定义Session结构体；
type Session struct {
	SessionId        int    `json:"sessionId"`
	UserId           string `json:"userId"`
	IpAddress        string `json:"ipAddress"`
	CreatedTime      string `json:"createdTime"`
	LastAccessedTime string `json:"lastAccessedTime"`
}

// 根据提供的BasicToken生成存储的 session token
func (s *Storage) GenerateSession() {
	// 每台存储最多生成64个session，当会话数量超过最大会话数时，HTTP状态码将返回503；
	// 上述情况可以等待一段时间后重新请求生成会话；
	url := BaseURL + "/v1/objects/storages/" + s.StorageDeviceId + "/sessions/"
	tmpMap := make(map[string]interface{})
	tmpMap["aliveTime"] = 60              // 可选参数，range：1-300，如果省略，默认是300s
	tmpMap["authenticationTimeout"] = 900 // 可选参数，range：1-900，如果省略，默认是120s

	content, _ := json.Marshal(tmpMap)
	postContent := bytes.NewBuffer(content)

	response, err := PostRequest(url, BasicToken, postContent)
	if err != nil {
		log.Fatalf("Generate Token error: %v\n", err)
	}
	body, _ := ioutil.ReadAll(response.Body)
	s.Token = gjson.Get(string(body), "token").String()
	log.Printf("Storage:%d,Token: %s\n", s.SerialNumber, s.Token)
}

// 获取指定存储的所有session信息；
func (s *Storage) GetSessions() ([]Session, error) {
	url := BaseURL + "/v1/objects/storages/" + s.StorageDeviceId + "/sessions"
	res, err := GetRequest(url, s.Token)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var sessions []Session
	for idx, session := range gjson.Get(string(body), "data").Array() {
		tmpSession := Session{}
		json.Unmarshal([]byte(session.String()), &tmpSession)
		sessions = append(sessions, tmpSession)
		log.Printf("第%d个session: %v\n", idx, tmpSession)
	}
	return sessions, nil
}

// 获取指定存储的特定session信息；
func (s *Storage) GetSpecificSession(sessionId int) (token string, err error) {
	url := BaseURL + "/v1/objects/storages/" + s.StorageDeviceId + "/sessions/" + fmt.Sprintf("%d", sessionId)
	res, err := GetRequest(url, s.Token)
	if err != nil {
		log.Fatalf("获取指定存储的特定session信息失败: %v\n", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("获取指定存储的特定session信息失败: %v\n", err)
	}

	token = gjson.Get(string(body), "token").String()
	return token, nil
}

// 抛弃指定session
func (s *Storage) DiscardSession(session Session, force bool) (err error) {
	/*
		执行该操作需要特定的用户权限：
		A user who belongs to the Administrator user group (built-in usergroup) for the VSP 5000 series, VSP E series, VSP Gx00 models, VSP G1000, VSPG1500, VSP Fx00 models, or VSP F1500,
		or a maintenance user for Virtual Storage Platform or Unified Storage VM can specify the value of sessionId that was obtained by the processing to get information about sessions.
	*/

	url := BaseURL + "/v1/objects/storages/" + s.StorageDeviceId + "/sessions/" + fmt.Sprintf("%d", session.SessionId)
	tmpMap := make(map[string]interface{})
	tmpMap["force"] = force // 是否强制抛弃会话；

	content, _ := json.Marshal(tmpMap)
	postContent := bytes.NewBuffer(content)
	res, err := DeleteRequest(url, s.Token, postContent)
	if err != nil {
		log.Fatalf("抛弃指定session失败: %v\n", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("抛弃指定session失败: %v\n", err)
	}

	log.Printf("抛弃指定session成功: %s\n", string(body))
	return nil
}
