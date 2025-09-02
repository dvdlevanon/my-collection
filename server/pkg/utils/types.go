package utils

import (
	"my-collection/server/pkg/model"
	"time"
)

type NowTimeGetter struct {
}

func (d NowTimeGetter) GetCurrentTime() time.Time {
	return time.Now()
}

type VideoFilter struct {
}

func (f VideoFilter) Filter(path string) bool {
	return IsVideo(true, path)
}

type PushSender struct {
	listeners []model.PushListener
}

func (f *PushSender) AddPushListener(l model.PushListener) {
	f.listeners = append(f.listeners, l)
}

func (f *PushSender) Push(m model.PushMessage) {
	if f.listeners == nil {
		return
	}

	for _, l := range f.listeners {
		l.Push(m)
	}
}
