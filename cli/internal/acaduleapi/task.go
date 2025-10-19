package acaduleapi

import (
	"acadule-cli/internal/easyhttp"
	"acadule-cli/internal/simplejson"
	"fmt"
	"net/http"
	"time"
)

type TaskProgress string

var (
	NOT_STARTED TaskProgress = "not_started"
	IN_PROGRESS TaskProgress = "in_progress"
	COMPLETE    TaskProgress = "complete"
	SUSPENDED   TaskProgress = "suspended"
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

type TaskFailError struct {
	Status string `json:"status,omitempty"`
	Reason string `json:"reason,omitempty"`
}

func (e *TaskFailError) Error() string {
	return fmt.Sprintf("failed to get task status: %s, reason: %s", e.Status, e.Reason)
}

func GetAll(apiUrl, token string) (data *[]TaskResponse, err error) {
	res, err := easyhttp.GetJsonWithBearer(apiUrl+"/task", token)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		result, err := simplejson.UnmarshalResponse[TaskFailError](res)
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
