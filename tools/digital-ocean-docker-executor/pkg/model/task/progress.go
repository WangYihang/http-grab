package http_grab_task

import (
	"encoding/json"
	"fmt"

	"github.com/WangYihang/digital-ocean-docker-executor/pkg/model/task"
	"github.com/WangYihang/gojob"
)

type GoJobProgress struct {
	GoJobStatus *gojob.Status
}

var PendingProgress *GoJobProgress
var DoneProgress *GoJobProgress

func init() {
	PendingProgress, _ = NewHTTPGrabProgress(`{"timestamp":"2024-03-18T10:46:40Z","num_done":0,"num_total":100}`)
	DoneProgress, _ = NewHTTPGrabProgress(`{"timestamp":"2024-03-18T10:46:40Z","num_done":100,"num_total":100}`)
}

func NewHTTPGrabProgress(message string) (*GoJobProgress, error) {
	progress := GoJobProgress{
		GoJobStatus: &gojob.Status{},
	}
	err := json.Unmarshal([]byte(message), progress.GoJobStatus)
	if err != nil {
		return nil, err
	}
	return &progress, nil
}

func (h *GoJobProgress) GetStatus() task.TaskStatus {
	if h.NumDoneWithSuccess() == 0 && h.NumDoneWithError() == 0 {
		return task.PENDING
	}
	if h.NumDoneWithSuccess()+h.NumDoneWithError() >= h.NumTotal() {
		return task.FINISHED
	}
	return task.RUNNING
}

func (h *GoJobProgress) NumTotal() int64 {
	return h.GoJobStatus.NumTotal
}

func (h *GoJobProgress) NumDoneWithSuccess() int64 {
	return h.GoJobStatus.NumSucceed
}

func (h *GoJobProgress) NumDoneWithError() int64 {
	return h.GoJobStatus.NumFailed
}

func (h *GoJobProgress) String() string {
	return fmt.Sprintf("%f%%", float64(h.NumDoneWithSuccess()+h.NumDoneWithError())/float64(h.NumTotal())*100)
}
