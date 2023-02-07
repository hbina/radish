package util

import (
	"os"
)

type ILogger interface {
	Fatal(v ...any)
	Fatalf(format string, v ...any)
	Fatalln(v ...any)
	Print(v ...any)
	Printf(format string, v ...any)
	Println(v ...any)
}

var Logger ILogger

var _ ILogger = &StubLogger{}

type StubLogger struct {
}

func (l *StubLogger) Fatal(v ...any) {
	os.Exit(1)
}
func (l *StubLogger) Fatalf(format string, v ...any) {
	os.Exit(1)
}
func (l *StubLogger) Fatalln(v ...any) {
	os.Exit(1)
}
func (l *StubLogger) Print(v ...any)                 {}
func (l *StubLogger) Printf(format string, v ...any) {}
func (l *StubLogger) Println(v ...any)               {}
