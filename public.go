package confy

import (
	"fmt"
	"strconv"
	"time"
)

var (
	Logger        func(s string)
	ReloadHandler func(map[string]interface{})
)

func Ready() {
	// close old ticker if exist
	if reloadTicker != nil {
		reloadTicker.Stop()
	}

	// initial reload
	reload()

	// start new ticker
	d := time.Duration(CfgReloadInterval()) * time.Second
	logf("start ticker and trigger it every %fs", d.Seconds())
	reloadTicker = time.NewTicker(d)
	go reloadTimely()
}

func Define(name string, defaultValue interface{}) {
	// cast value
	v, kind := cast(defaultValue)

	// dupldate check
	if e, ok := getVal(name); ok {
		if e.defined {
			// if element exist in map with same name
			// and defined before, then panic
			panicf("duplicate name: %s", name)
		} else {
			// otherwise, validate the exist element's kind
			if kind == kindInt64 && e.kind == kindFloat64 {
				e.kind = kindInt64
				e.val = int64(e.val.(float64))
			}

			if e.kind != kind {
				panicf(
					`invalid kind for "%s". %v(%s) -> %v(%s)`,
					name,
					defaultValue, kindName(kind),
					e.val, kindName(e.kind),
				)
			} else {
				// if kind is same.
				// just set property "defined" to true
				e.defined = true
				setVal(name, e)
			}
		}
	} else {
		// assignment
		setVal(name, elem{name, kind, v, true})
	}
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
