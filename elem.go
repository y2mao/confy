package confy

import (
	"strconv"
	"sync"
	"time"
)

type elem struct {
	name string
	val  string

	defined    bool
	assigned   bool
	defineLock sync.Mutex
}

func newElem(name string, val string) *elem {
	e := elem{
		name:     name,
		val:      val,
		assigned: true,
	}

	logf("NEW %-10s -> %s", e.name, e.val)
	return &e
}

func (e *elem) Define(v string) {
	e.defineLock.Lock()
	defer e.defineLock.Unlock()

	if e.defined {
		panicf("duplicate config definition: %s", e.name)
	}

	e.defined = true

	if !e.assigned {
		e.val = v
		logf("DEF %-10s -> %s", e.name, e.val)
	} else {
		logf("IGN %-10s != %s", e.name, e.val)
	}
}

func (e *elem) Set(v string) {
	e.assigned, e.val = true, v
	logf("SET %-10s -> %s", e.name, e.val)
}

func (e *elem) GetString() string {
	return e.val
}

func (e *elem) GetBool() bool {
	v, err := strconv.ParseBool(e.val)
	if err != nil {
		panicf("%v", err)
	}

	return v
}

func (e *elem) GetInt64() int64 {
	v, err := strconv.ParseInt(e.val, 10, 64)
	if err != nil {
		panicf("%v", err)
	}

	return v
}

func (e *elem) GetFloat64() float64 {
	v, err := strconv.ParseFloat(e.val, 64)
	if err != nil {
		panicf("%v", err)
	}

	return v
}

func (e *elem) GetTime() time.Time {
	v, err := time.Parse(TIME_FORMAT, e.val)
	if err != nil {
		panicf("%v", err)
	}

	return v
}

func (e *elem) GetDuration() time.Duration {
	v, err := time.ParseDuration(e.val)
	if err != nil {
		panicf("%v", err)
	}

	return v
}
