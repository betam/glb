package cache

import (
	"cmp"
	"fmt"
	"sync"
	"time"
)

type Stored interface {
	cmp.Ordered
}

type cached[T Stored] struct {
	validUntil int64
	key        string
	item       T
}

type Changed[T Stored] struct {
	Changed bool
	Old     T
	New     T
}

type Ocache[T Stored] struct {
	cache     map[string]cached[T]
	m         sync.Mutex
	so        sync.Once
	ttl       int64
	gcRan     bool
	getCnt    int64
	getTryCnt int64
	setCnt    int64
	unsetCnt  int64
}

func NewOCache[T Stored](ttl int64) *Ocache[T] {
	c := Ocache[T]{}
	return c.Init(ttl)
}

func (e *Ocache[T]) Init(ttl int64) *Ocache[T] {
	e.so.Do(func() {
		e.ttl = ttl
		e.cache = make(map[string]cached[T])
	})
	return e
}

func (e *Ocache[T]) Set(item T, key string) {

	e.m.Lock()
	defer e.m.Unlock()

	e.gcrun()
	e.setCnt++

	e.cache[key] = cached[T]{
		validUntil: time.Now().Unix() + e.ttl,
		key:        key,
		item:       item,
	}
}

func (e *Ocache[T]) SetIfChanged(item T, key string) Changed[T] {

	e.m.Lock()
	defer e.m.Unlock()

	e.gcrun()

	exist, ok := e.cache[key]
	e.getTryCnt++

	if !ok || (ok && cmp.Compare(item, exist.item) != 0) {
		e.getCnt++
		e.setCnt++

		e.cache[key] = cached[T]{
			validUntil: time.Now().Unix() + e.ttl,
			key:        key,
			item:       item,
		}
		return Changed[T]{
			Changed: true,
			Old:     exist.item,
			New:     item,
		}
	}
	return Changed[T]{}
}

//func (e *Ocache[T]) equal(t1, t2 T) bool {
//
//	switch v1 := t1.(type) {
//	case string:
//		return strings.Compare(v1, t2.(string)) < 0 // import "strings"
//		// return v1 < t2.(string)
//	case int:
//		return v1 < t2.(int)
//	}
//
//	return false
//}

func (e *Ocache[T]) Get(key string) *T {
	e.m.Lock()
	defer e.m.Unlock()

	e.getTryCnt++

	exist, ok := e.cache[key]
	if !ok {
		return nil
	}

	if e.isOverdue(exist) {
		return nil
	}

	e.getCnt++

	return &exist.item
}

func (e *Ocache[T]) GetAll() map[string]T {
	e.m.Lock()
	defer e.m.Unlock()

	all := map[string]T{}
	for key, v := range e.cache {
		all[key] = v.item
	}
	return all
}

func (e *Ocache[T]) isOverdue(c cached[T]) bool {

	if c.validUntil < time.Now().Unix() {
		return true
	}
	return false
}

func (e *Ocache[T]) Report() string {
	return fmt.Sprintf("Try/Get/Set/Unset cnt: %v/%v/%v/%v", e.getTryCnt, e.getCnt, e.setCnt, e.unsetCnt)
}

func (e *Ocache[T]) gcrun() {
	if !e.gcRan {
		e.gcRan = true
		go func(e *Ocache[T]) {
			defer func() { e.gcRan = false }()
			for {
				e.m.Lock()

				if len(e.cache) == 0 {
					e.m.Unlock()
					return
				}

				for key, cached := range e.cache {
					if e.isOverdue(cached) {
						delete(e.cache, key)
						e.unsetCnt++
					}
				}

				e.m.Unlock()
				time.Sleep(3 * time.Second)
			}
		}(e)
	}
}
