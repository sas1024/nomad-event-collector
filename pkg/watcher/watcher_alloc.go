package watcher

import (
	"strconv"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog/log"
)

func NewAllocWatcher(client *api.Client) Watcher {
	return Watcher{
		lastChangeIndex: 0,
		needInit:        true,
		queryOptions: &api.QueryOptions{
			Namespace:  "*",
			AllowStale: true,
			WaitIndex:  0,
			WaitTime:   5 * time.Minute,
		},
		client: client,
	}
}

func (w *Watcher) Run(msgChan chan AllocData) {
	allocCache, _ := lru.New(10240)
	for {
		var lastIdx uint64
		allocations, meta, err := w.client.Allocations().List(w.queryOptions)
		if err != nil {
			log.Error().Err(err).Msg("unable to fetch allocations list from Nomad")
			time.Sleep(10 * time.Second)
			continue
		}

		if w.needInit {
			log.Info().Msg("running initial allocation watcher index update")
			w.queryOptions.WaitIndex = meta.LastIndex
			w.lastChangeIndex = meta.LastIndex
			w.needInit = false
			continue
		}

		if !changed(meta.LastIndex, w.queryOptions.WaitIndex) {
			continue
		}

		lastIdx = meta.LastIndex

		for _, alloc := range allocations {
			last5min := time.Now().Add(-5 * time.Minute)

			if !changed(alloc.ModifyIndex, lastIdx) {
				continue
			}

			if len(alloc.TaskGroup) == 0 {
				continue
			}

			lastIdx = alloc.ModifyIndex

			for task, v := range alloc.TaskStates {
				for _, e := range v.Events {
					if e.Type != api.TaskTerminated {
						continue
					}
					eventAt := time.Unix(0, e.Time)

					if eventAt.Before(last5min) {
						continue
					}

					if v, exist := e.Details["oom_killed"]; exist && v == "true" {
						adId := strconv.Itoa(int(e.Time)) + alloc.ID
						if !allocCache.Contains(adId) {
							allocCache.Add(adId, nil)
							ad := AllocData{
								ID:        alloc.ID,
								Node:      alloc.NodeName,
								Job:       alloc.JobID,
								Task:      task,
								TaskGroup: alloc.TaskGroup,
								Message:   e.Message,
								EventAt:   eventAt,
							}
							msgChan <- ad
						}
					}
				}
			}
		}
		w.queryOptions.WaitIndex = lastIdx
		w.lastChangeIndex = lastIdx
	}
}
