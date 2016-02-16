package confy

import "flag"

var (
	configFile     = flag.String("confy-file", "", "")
	configURL      = flag.String("confy-url", "", "")
	configInterval = flag.Int("confy-interval", 60, "")
)

func CfgFile() string {
	if s := *configFile; len(s) > 0 {
		return s
	}

	return "./app.confy"
}

func CfgURL() string {
	if s := *configURL; len(s) > 0 {
		return s
	}

	return ""
}

func CfgReloadInterval() int {
	if i := *configInterval; i >= 10 {
		return i
	}

	return 60
}
