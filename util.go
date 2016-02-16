package confy

import (
	"fmt"
	"os"
)

func logf(f string, v ...interface{}) {
	fmt.Fprintf(os.Stdout, "[confy] "+f+"\n", v...)
}

func panicf(f string, v ...interface{}) {
	logf(f, v...)
	panic(fmt.Sprintf(f, v...))
}
