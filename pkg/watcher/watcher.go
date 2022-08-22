package watcher

import (
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog"
)

type Watcher struct {
	lastChangeIndex uint64
	needInit        bool
	queryOptions    *api.QueryOptions
	client          *api.Client
}

type AllocData struct {
	ID        string
	Node      string
	Job       string
	Task      string
	TaskGroup string
	Message   string
	EventAt   time.Time
}

func (ad AllocData) MarshalZerologObject(e *zerolog.Event) {
	e.Time("eventAt", ad.EventAt).Str("node", ad.Node).Str("job", ad.Job).Str("taskGroup", ad.TaskGroup).Str("task", ad.Task).Str("alloc", ad.ID).Msg(ad.Message)
}

func changed(new, old uint64) bool {
	if new < old {
		return false
	}
	return true
}
