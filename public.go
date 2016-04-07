package confy

import "time"

var (
	Logger        func(s string)
	ReloadHandler func(map[string]interface{})
)

func Ready() {

}

func Define(name string, defaultValue interface{}) {
	define(name, defaultValue)
}

func Int(name string) int64 {
	e := getElemFromMap(name)
	return e.GetInt64()
}

func Float(name string) float64 {
	e := getElemFromMap(name)
	return e.GetFloat64()
}

func Bool(name string) bool {
	e := getElemFromMap(name)
	return e.GetBool()
}

func Text(name string) string {
	e := getElemFromMap(name)
	return e.GetString()
}

func Time(name string) time.Time {
	e := getElemFromMap(name)
	return e.GetTime()
}

func Duration(name string) time.Duration {
	e := getElemFromMap(name)
	return e.GetDuration()
}

func String(name string) string {
	e := getElemFromMap(name)
	return e.GetString()
}
