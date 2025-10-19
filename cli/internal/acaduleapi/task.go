package acaduleapi

import (
	"acadule-cli/internal/easyhttp"
	"acadule-cli/internal/simplejson"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type TaskProgress string

var (
	NOT_STARTED TaskProgress = "NOT_STARTED"
	IN_PROGRESS TaskProgress = "IN_PROGRESS"
	COMPLETE    TaskProgress = "COMPLETE"
	SUSPENDED   TaskProgress = "SUSPENDED"
)

type TaskResponse struct {
	Id          string       `json:"id,omitempty"`
	OwnerId     *int         `json:"ownerId,omitempty"`
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	Progress    TaskProgress `json:"progress,omitempty"`
	Deadline    *time.Time   `json:"deadline,omitempty"`
	LastUpdated CustomTime   `json:"lastUpdated,omitempty"`
	HasDone     bool         `json:"hasDone,omitempty"`
}

type RequestFailError struct {
	Status string `json:"status,omitempty"`
	Reason string `json:"reason,omitempty"`
}

func (e *RequestFailError) Error() string {
	return fmt.Sprintf("failed to get task status: %s, reason: %s", e.Status, e.Reason)
}

func GetAll(apiUrl, token string) (data *[]TaskResponse, err error) {
	res, err := easyhttp.GetJsonWithBearer(apiUrl+"/task", token)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		result, err := simplejson.UnmarshalResponse[RequestFailError](res)
		if err != nil {
			return nil, err
		}
		return nil, &result
	}

	result, err := simplejson.UnmarshalResponse[[]TaskResponse](res)
	if err != nil {
		return
	}

	return &result, nil
}

type TaskAddRequest struct {
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	Progress    TaskProgress `json:"progress,omitempty"`
	Deadline    *CustomTime  `json:"deadline,omitempty"`
}

type TaskAddResponse struct {
	Status string
	Id     string
}

func Add(apiUrl, token string, request TaskAddRequest) (data *TaskAddResponse, err error) {
	postData, err := json.Marshal(request)
	if err != nil {
		return
	}
	res, err := easyhttp.PostJsonWithBearer(apiUrl+"/task", token, postData)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf("token is not valid")
		}
		errorData, err := simplejson.UnmarshalResponse[RequestFailError](res)
		if err != nil {
			return nil, err
		}
		return nil, &errorData
	}

	response, err := simplejson.UnmarshalResponse[TaskAddResponse](res)
	return &response, err
}

func View(apiUrl, token, id string) (data *TaskResponse, err error) {
	res, err := easyhttp.GetJsonWithBearer(apiUrl+"/task/"+id, token)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		errorData, err := simplejson.UnmarshalResponse[RequestFailError](res)
		if err != nil {
			return nil, err
		}
		return nil, &errorData
	}

	response, err := simplejson.UnmarshalResponse[TaskResponse](res)
	return &response, err
}

type UpdateRequest struct {
	Id          string       `json:"id,omitempty"`
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	Progress    TaskProgress `json:"progress,omitempty"`
	Deadline    *CustomTime  `json:"deadline,omitempty"`
	HasDone     bool         `json:"hasDone,omitempty"`
}

func Update(apiUrl, token string, request UpdateRequest) (response *TaskResponse, err error) {
	patchData, err := json.Marshal(request)
	if err != nil {
		return
	}
	slog.Debug("Send update", slog.String("body", string(patchData)))
	res, err := easyhttp.PatchJsonWithBearer(apiUrl+"/task", token, patchData)
	if err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		errorData, err := simplejson.UnmarshalResponse[RequestFailError](res)
		if err != nil {
			return nil, err
		}
		return nil, &errorData
	}
	taskData, err := simplejson.UnmarshalResponse[TaskResponse](res)
	return &taskData, err
}
