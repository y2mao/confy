package confy

import "fmt"

func logf(f string, v ...interface{}) {
	s := fmt.Sprintf("[confy] "+f+"\n", v...)

	if Logger != nil {
		Logger(s)
	}
}

func panicf(f string, v ...interface{}) {
	logf(f, v...)
	panic(fmt.Sprintf(f, v...))
}
