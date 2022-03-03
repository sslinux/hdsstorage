package hdsstorage

import (
	"fmt"
	"io/ioutil"
	"log"
)

// 定义存储结构体
type Storage struct {
	StorageDeviceId       string              `json:"storageDeviceId"`
	Model                 string              `json:"model"`
	SerialNumber          int                 `json:"serialNumber"`
	SvpIp                 string              `json:"svpIp"`
	Ctl1Ip                string              `json:"ctl1Ip"`
	Ctl2Ip                string              `json:"ctl2Ip"`
	DkcMicroVersion       string              `json:"dkcMicroVersion"`
	DetailDkcMicroVersion string              `json:"detailDkcMicroVersion"`
	Ctl1MicroVersion      string              `json:"ctl1MicroVersion"`
	Ctl2MicroVersion      string              `json:"ctl2MicroVersion"`
	CommunicationModes    []map[string]string `json:"communicationModes"`
	IsSecure              bool                `json:"isSecure"`
	TargetCtl             string              `json:"targetCtl"`
	UseSvp                bool                `json:"usesSvp"`
	Token                 string
}

func (s *Storage) Summary(detailInfoType string) {
	/*
			"message": "This API request can be executed only for configurations linked to an SVP.",
		    "solution": "Make sure that the configuration is linked to an SVP, and then try again.",
	*/
	var url string
	if detailInfoType == "parityGroupCapacity" {
		url = BaseURL + "/v1/objects/storages/" + s.StorageDeviceId + "/storage-summaries/instance" + "?detailInfoType=" + detailInfoType
	} else {
		url = BaseURL + "/v1/objects/storages/" + s.StorageDeviceId + "/storage-summaries/instance"
	}

	res, err := GetRequest(url, s.Token)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

// 删除已注册的存储设备
func (s *Storage) DeleteStorage() {
	storage := GetSpecificStoages(s.SerialNumber, "")
	url := BaseURL + "/v1/objects/storages/" + s.StorageDeviceId
	response, err := DeleteRequest(url, storage.Token, nil)
	if err != nil {
		log.Fatalf("DeleteStorage error: %v\n", err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("DeleteStorage error: %v\n", err)
	}
	fmt.Println(string(body))
}
