package i18nPkg

import (
	"fmt"
	"os"
)

type Log interface {
	//Warnf(string, ...any)
	Errorf(string, ...any)
}

type StdLogger struct{}

func (l *StdLogger) Errorf(s string, args ...any) {
	fmt.Fprintf(os.Stdout, "[Error] "+s, args...)
}
