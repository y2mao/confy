package confy

import (
	"fmt"
	"strconv"
	"time"
)

var (
	ReloadHandler func(map[string]interface{})
)

func Ready() {
	// close old ticker if exist
	if reloadTicker != nil {
		reloadTicker.Stop()
		tickerQuit <- struct{}{}
	}

	// initial reload
	reload()

	// start new ticker
	reloadTicker = time.NewTicker(time.Duration(CfgReloadInterval()) * time.Second)
	go reloadTimely()
}

func Define(name string, defaultValue interface{}) {
	// cast value
	v, kind := cast(defaultValue)

	// dupldate check
	if _, ok := getVal(name); ok {
		panicf("duplicate name: %s", name)
	}

	// assignment
	setVal(name, elem{name, kind, v})
}

func Int(name string) int64 {
	e, ok := getVal(name)

	if !ok {
		panicf("config missing with given name: %s", name)
	} else if e.kind != kindInt64 {
		panicf("unexpect type: %s -> int64", kindName(e.kind))
	}

	return e.val.(int64)
}

func Float(name string) float64 {
	e, ok := getVal(name)

	if !ok {
		panicf("config missing with given name: %s", name)
	} else if e.kind != kindFloat64 {
		panicf("unexpect type: %s -> float64", kindName(e.kind))
	}

	return e.val.(float64)
}

func Bool(name string) bool {
	e, ok := getVal(name)

	if !ok {
		panicf("config missing with given name: %s", name)
	} else if e.kind != kindBool {
		panicf("unexpect type: %s -> bool", kindName(e.kind))
	}

	return e.val.(bool)
}

func Text(name string) string {
	e, ok := getVal(name)

	if !ok {
		panicf("config missing with given name: %s", name)
	} else if e.kind != kindString {
		panicf("unexpect type: %s -> string", kindName(e.kind))
	}

	return e.val.(string)
}

func Time(name string) time.Time {
	e, ok := getVal(name)

	if !ok {
		panicf("config missing with given name: %s", name)
	} else if e.kind != kindTime {
		panicf("unexpect type: %s -> time.Time", kindName(e.kind))
	}

	return e.val.(time.Time)
}

func Duration(name string) time.Duration {
	e, ok := getVal(name)

	if !ok {
		panicf("config missing with given name: %s", name)
	} else if e.kind != kindDuration {
		panicf("unexpect type: %s -> time.Duration", kindName(e.kind))
	}

	return e.val.(time.Duration)
}

func String(name string) string {
	e, ok := getVal(name)

	if !ok {
		panicf("config missing with given name: %s", name)
	}

	switch e.kind {
	case kindString:
		return e.val.(string)
	case kindInt64:
		return strconv.FormatInt(e.val.(int64), 10)
	case kindFloat64:
		return strconv.FormatFloat(e.val.(float64), 'f', 6, 64)
	case kindBool:
		return strconv.FormatBool(e.val.(bool))
	case kindTime:
		return e.val.(time.Time).Format(time.RFC3339)
	case kindDuration:
		return e.val.(time.Duration).String()
	}

	return fmt.Sprintf("%v", e.val)
}
