package ratelimit

import (
	"container/list"
	"log"
	"sync"
	"time"
)

type rateLimit struct {
	cfg   *RtFailedLoginConfig
	mtx   sync.Mutex
	times *list.List
}

type RtFailedLogin struct {
	cfg *RtFailedLoginConfig
	mtx sync.RWMutex
	rts map[string]*rateLimit
}

type RtFailedLoginConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Limit    int           `yaml:"burst"`
	Interval time.Duration `yaml:"interval"`
}

// New creates a new rate limiter for the limit and interval.
func NewRtFailedLogin(cfg *RtFailedLoginConfig) *RtFailedLogin {
	rt := &RtFailedLogin{
		cfg: cfg,
		mtx: sync.RWMutex{},
		rts: make(map[string]*rateLimit),
	}
	if cfg.Enabled {
		go rt.clean()
	}
	return rt
}

// CanAct is expected to be called before Act
func (r *RtFailedLogin) CanAct(key string) (ok bool, remaining time.Duration) {
	if !r.cfg.Enabled {
		return true, 0
	}
	r.mtx.RLock()
	rt, ok := r.rts[key]
	r.mtx.RUnlock()
	if !ok {
		return true, 0
	}
	return rt.canAct()
}

// Act can be called after CanAct returns true.
func (r *RtFailedLogin) Act(key string) {
	if !r.cfg.Enabled {
		return
	}
	r.mtx.RLock()
	rt, ok := r.rts[key]
	r.mtx.RUnlock()
	if !ok {
		rt = &rateLimit{
			cfg:   r.cfg,
			mtx:   sync.Mutex{},
			times: list.New(),
		}
		r.mtx.Lock()
		r.rts[key] = rt
		r.mtx.Unlock()
	}
	rt.act()
}

func (r *RtFailedLogin) clean() {
	for {
		r.mtx.Lock()
		for k, v := range r.rts {
			if v.empty() {
				delete(r.rts, k)
			}
		}
		r.mtx.Unlock()
		time.Sleep(time.Hour)
	}
}

func (r *rateLimit) empty() bool {
	v := r.times.Back().Value
	b := v.(time.Time)
	now := time.Now()
	return now.Sub(b) > r.cfg.Interval
}

func (r *rateLimit) canAct() (ok bool, remaining time.Duration) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	now := time.Now()
	l := r.times.Len()
	if l < r.cfg.Limit {
		return true, 0
	}
	frnt := r.times.Front()
	t := frnt.Value.(time.Time)
	diff := now.Sub(t)
	if diff < r.cfg.Interval {
		return false, r.cfg.Interval - diff
	}
	return true, 0
}

func (r *rateLimit) act() {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	now := time.Now()
	l := r.times.Len()
	log.Print("len act: ", l)
	if l < r.cfg.Limit {
		r.times.PushBack(now)
		return
	}
	frnt := r.times.Front()
	frnt.Value = now
	r.times.MoveToBack(frnt)
}
