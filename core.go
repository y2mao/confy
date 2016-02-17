package confy

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"
	"time"
)

const (
	kindInt64 = iota
	kindFloat64
	kindString
	kindBool
	kindTime
	kindDuration
)

func kindName(kind int) string {
	switch kind {
	case kindInt64:
		return "int64"
	case kindFloat64:
		return "float64"
	case kindString:
		return "string"
	case kindBool:
		return "bool"
	case kindTime:
		return "time"
	case kindDuration:
		return "duration"
	}

	return "n/a"
}

func cast(originVal interface{}) (v interface{}, kind int) {
	switch originVal.(type) {
	case int:
		v = int64(originVal.(int))
		kind = kindInt64
	case int16:
		v = int64(originVal.(int16))
		kind = kindInt64
	case int32:
		v = int64(originVal.(int32))
		kind = kindInt64
	case int64:
		v = originVal.(int64)
		kind = kindInt64
	case uint:
		v = int64(originVal.(uint))
		kind = kindInt64
	case uint16:
		v = int64(originVal.(uint16))
		kind = kindInt64
	case uint32:
		v = int64(originVal.(uint32))
		kind = kindInt64
	case uint64:
		v = int64(originVal.(uint64))
		kind = kindInt64
	case float32:
		v = float64(originVal.(float32))
		kind = kindFloat64
	case float64:
		v = float64(originVal.(float64))
		kind = kindFloat64
	case string:
		v = originVal.(string)
		kind = kindString
	case bool:
		v = originVal.(bool)
		kind = kindBool
	case time.Time:
		v = originVal.(time.Time)
		kind = kindTime
	case time.Duration:
		v = originVal.(time.Duration)
		kind = kindDuration
	default:
		panicf(
			"invalid type:%s, val:%v",
			reflect.TypeOf(originVal).String(),
			originVal,
		)
	}

	return
}

type elem struct {
	name string
	kind int
	val  interface{}
}

var (
	signature string
	memMap    = make(map[string]elem)

	reloadLock   sync.Mutex
	reloadTicker *time.Ticker
	tickerQuit   = make(chan struct{})
)

func reloadTimely() {
	for t := range reloadTicker.C {
		logf("ticker fired at %v", t)
		reload()
	}
	/*
		for {
			select {
			case <-reloadTicker.C:
				reload()
			case <-tickerQuit:
				reloadTicker.Stop()
				return
			}
		}
	*/
}

func getVal(name string) (elem, bool) {
	v, ok := memMap[name]
	return v, ok
}

func setVal(name string, e elem) {
	memMap[name] = e
}

func reload() {
	reloadLock.Lock()
	defer reloadLock.Unlock()

	logf("start reloading")

	var data []byte
	var m = make(map[string]interface{})

	defer func() {
		// trigger reloadHandler for further process of app, such as logging
		if ReloadHandler != nil {
			ReloadHandler(m)
		}
	}()

	// 1. load data from file and URL
	// file data will be overwritten if URL is valid
	if s := CfgFile(); s != "" {
		if b, err := loadFromFile(s); err == nil {
			data = b
			logf("file loaded[%d] (%s)", len(data), s)
		} else {
			logf("file missing: %v (%s)", err, s)
		}
	}
	if s := CfgURL(); s != "" {
		if b, err := loadFromURL(s); err == nil {
			data = b
			logf("url loaded[%d] (%s)", len(data), s)
		} else {
			logf("url missing: %v (%s)", err, s)
		}
	}

	// 2. calculate signature
	// ignore if signature as same as last result
	if data == nil || len(data) == 0 {
		logf("data is empty")
		return
	}

	if sum := fmt.Sprintf("%x", sha1.Sum(data)); sum == signature {
		logf("same signature - ignored")
		return
	} else {
		signature = sum
		logf("new signature:%s", signature)
	}

	// 3. parse data to m
	if err := json.Unmarshal(data, &m); err != nil {
		panicf("parsing config data failed. %v", err)
	}

	// 4. validate kind and assign m to memMap
	for name, newV := range m {
		castV, kind := cast(newV)

		if oldElem, ok := getVal(name); ok {
			// validate kind and assign
			if oldElem.kind != kind {
				panicf(
					"new config kind doesn't match current one. %s -> %s",
					kindName(kind),
					kindName(oldElem.kind),
				)
			}
		}

		setVal(name, elem{name, kind, castV})
	}
}

func loadFromFile(fn string) (b []byte, err error) {
	b, err = ioutil.ReadFile(fn)
	return
}

func loadFromURL(url string) (b []byte, err error) {
	var resp *http.Response
	resp, err = http.Get(url)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code:%d from %s", resp.StatusCode, url)
		return
	}

	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)

	return
}
