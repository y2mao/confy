package confy

import (
	"flag"
	"fmt"
	"time"
)

func init() {
	flag.Parse()

	if !CfgNoLog() {
		Logger = func(s string) {
			fmt.Print(s)
		}
	}

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
