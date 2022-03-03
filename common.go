package hdsstorage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/tidwall/gjson"
)

var (
	BaseURL string = "http://192.204.1.91:23450/ConfigurationManager"
	// UserName string = ""
	// Password   string = "hiaa"
	BasicToken string = "Basic aGlhYTpQQHNzdzByZA=="
)

// GET请求的简单封装
func GetRequest(url, token string) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("New Get Request %s, error: %v\n", url, err)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Session "+token)

	response, err := client.Do(request)
	if err != nil {
		log.Printf("Get Request: %s, Err: %v\n", url, err)
		return response, err
	}
	return response, nil
}

// POST请求的简单封装；
func PostRequest(url, token string, postContent *bytes.Buffer) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, postContent)
	if err != nil {
		log.Fatalf("Post Request: %s, Err: %v\n", url, err)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", BasicToken)

	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("Post Request %s error: %v\n", url, err)
		return response, err
	}
	return response, nil
}

// DELETE请求的简单封装；
func DeleteRequest(url, token string, postContent *bytes.Buffer) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("DELETE", url, postContent)
	if err != nil {
		log.Fatalf("DELETE Request %s error: %v\n", url, err)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", BasicToken)

	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("DELETE Request %s error: %v\n", url, err)
		return response, err
	}
	return response, nil
}

// 获取所有已注册到HCM上的存储设备；
func GetAllStorages() []Storage {
	var storages []Storage
	url := BaseURL + "/v1/objects/storages"
	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", url, nil)

	response, _ := client.Do(reqest)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("获取存储列表失败: %v\n", err)
		os.Exit(1)
	}

	for _, strStorage := range gjson.Get(string(body), "data").Array() {
		storage := Storage{}
		json.Unmarshal([]byte(strStorage.String()), &storage)
		storages = append(storages, storage)
	}
	return storages
}

// 获取指定SN的存储设备；
func GetSpecificStoages(sn int, detailInfoType string) Storage {
	url := BaseURL + "/v1/objects/storages"
	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", url, nil)

	response, _ := client.Do(reqest)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("获取存储列表失败: %v\n", err)
		os.Exit(1)
	}
	var targetStorage Storage
	for _, strStorage := range gjson.Get(string(body), "data").Array() {
		storage := Storage{}
		json.Unmarshal([]byte(strStorage.String()), &storage)
		if storage.SerialNumber == sn {
			targetStorage = storage
		}
	}

	if detailInfoType == "version" {
		url = BaseURL + "/storages/" + targetStorage.StorageDeviceId + "?detailInfoType=versions"
	} else {
		url = BaseURL + "/storages/" + targetStorage.StorageDeviceId
	}
	response, err = GetRequest(url, "")
	if err != nil {
		log.Fatalf("获取存储设备详细信息失败: %v\n", err)
	}

	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("获取存储设备详细信息失败: %v\n", err)
	}
	var storage Storage

	json.Unmarshal([]byte(string(body)), &storage)
	return storage
}

// 获取HCM REST API的版本信息；
func GetAPIVersion() (string, string) {
	url := BaseURL + "/configuration/version"
	response, err := GetRequest(url, "")
	if err != nil {
		fmt.Println(err)
		log.Fatalf("GetAPIVersion error: %v\n", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("GetAPIVersion error: %v\n", err)
	}
	ProductName := gjson.Get(string(body), "productName").String()
	APIVersion := gjson.Get(string(body), "apiVersion").String()
	fmt.Printf("ProductName: %s, APIVersion: %s \n", ProductName, APIVersion)
	return ProductName, APIVersion
}

// 注册存储设备到HCM
// func RegisterStorage(model string) {
// 	url := BaseURL + "/v1/objects/storages"
// 	storage := Storage{}
// 	switch model {
// 		case "VSP G350","VSP G370","VSP G700","VSP G900","VSP F350","VSP F370","VSP F700","VSP F900":
// 			storage.Model = model
// 			storage.SerialNumber = "1234567890"
// 	}
// }

// 修改已注册存储的信息；
// func ModifyStorage(sn int, model string) {}
