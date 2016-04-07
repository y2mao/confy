package confy

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var (
	signature string
	memMap    = make(map[string]*elem)

	reloadLock   sync.Mutex
	reloadTicker *time.Ticker
)

func reloadTimely() {
	for _ = range reloadTicker.C {
		reload()
	}
}

func getElemFromMap(name string) *elem {
	e, ok := memMap[name]
	if !ok {
		panicf("missing config name: %s", name)
	}

	return e
}

func define(name string, defaultValue interface{}) {
	reloadLock.Lock()
	defer reloadLock.Unlock()

	v := castToString(defaultValue)

	if e, ok := memMap[name]; ok {
		e.Define(v)
	} else {
		memMap[name] = newElem(name, v)
	}
}

func reload() {
	reloadLock.Lock()
	defer reloadLock.Unlock()

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
		logf("same signature - ignored (%s)", signature)
		return
	} else {
		signature = sum
		logf("new signature (%s)", signature)
	}

	// 3. parse data to m
	if err := json.Unmarshal(data, &m); err != nil {
		panicf("parsing config data failed. %v", err)
	}

	// 4. validate kind and assign m to memMap
	for name, newV := range m {
		v := castToString(newV)
		if oe, ok := memMap[name]; ok {
			oe.Set(v)
		} else {
			memMap[name] = newElem(name, v)
		}
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

func castToString(v interface{}) string {
	if t, ok := v.(time.Time); ok {
		return t.Format(TIME_FORMAT)
	} else {
		return fmt.Sprintf("%v", v)
	}
}
