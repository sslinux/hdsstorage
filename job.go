package hdsstorage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/tidwall/gjson"
)

type Job struct {
	JobId             string            `json:"jobId"`
	Self              string            `json:"self"`
	UserId            string            `json:"userId"`
	Status            string            `json:"status"`
	State             string            `json:"state"`
	CreatedTime       string            `json:"createdTime"`
	UpdatedTime       string            `json:"updatedTime"`
	CompletedTime     string            `json:"completedTime"`
	Request           map[string]string `json:"request"`
	AffectedResources []string          `json:"affectedResources"`
}

// 获取指定存储的所有job信息；
func (s *Storage) GetJobs() []Job {
	// 支持的查询参数：
	/*
		startCreatedTime ISO8601string, 可选参数；
		endCreatedTime ISO8601string, 可选参数；
		count int, range: 1-100,默认为100；
		status string, 可选参数，可以是：Initializing, Running, Completed;
			如果也要指定state，则status必须为：Succeeded,Failed,Unknown其中之一；
		state string, 可选参数,可以是：Queued,Started,Succeeded,Failed,Unknown；

		Example:
			?startCreatedTime=2015-05-01T08:00:00Z&endCreatedTime=2015-05-31T23:59:59Z&count=30&state=Succeeded
	*/
	url := BaseURL + "/v1/objects/storages/" + s.StorageDeviceId + "/jobs"
	res, err := GetRequest(url, s.Token)
	if err != nil {
		log.Fatalf("Get Jobs error: %v\n", err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	var jobs []Job
	for _, job := range gjson.Get(string(body), "data").Array() {
		tmpJob := Job{}
		json.Unmarshal([]byte(job.String()), &tmpJob)
		jobs = append(jobs, tmpJob)
	}
	return jobs
}

// 获取指定job信息；
func (s *Storage) GetJob(jobId int64) Job {
	url := BaseURL + "/v1/objects/storages/" + s.StorageDeviceId + "/jobs/" + fmt.Sprintf("%d", jobId)
	res, err := GetRequest(url, s.Token)
	if err != nil {
		if res.StatusCode != 404 {
			log.Printf("Job: %d not found\n", jobId)
			return Job{}
		} else {
			log.Fatalf("Get Job:%d error: %v\n", jobId, err)
		}
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	var job Job
	json.Unmarshal([]byte(body), &job)
	return job
}
