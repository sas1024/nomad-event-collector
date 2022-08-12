package watcher

import (
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
	TaskGroup string
	Message   string
}

func (ad AllocData) MarshalZerologObject(e *zerolog.Event) {
	e.Str("node", ad.Node).Str("job", ad.Job).Str("taskGroup", ad.TaskGroup).Str("alloc", ad.ID).Msg(ad.Message)
}

func changed(new, old uint64) bool {
	if new < old {
		return false
	}
	return true
}
