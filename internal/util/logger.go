package util

import (
	"io"
)

type ILogger interface {
	Fatal(v ...any)
	Fatalf(format string, v ...any)
	Fatalln(v ...any)
	Panic(v ...any)
	Panicf(format string, v ...any)
	Panicln(v ...any)
	Print(v ...any)
	Printf(format string, v ...any)
	Println(v ...any)
	SetFlags(flag int)
	SetOutput(w io.Writer)
	SetPrefix(prefix string)
}

var Logger ILogger

var _ ILogger = &StubLogger{}

type StubLogger struct {
}

func (l *StubLogger) Fatal(v ...any)                 {}
func (l *StubLogger) Fatalf(format string, v ...any) {}
func (l *StubLogger) Fatalln(v ...any)               {}
func (l *StubLogger) Panic(v ...any)                 {}
func (l *StubLogger) Panicf(format string, v ...any) {}
func (l *StubLogger) Panicln(v ...any)               {}
func (l *StubLogger) Print(v ...any)                 {}
func (l *StubLogger) Printf(format string, v ...any) {}
func (l *StubLogger) Println(v ...any)               {}
func (l *StubLogger) SetFlags(flag int)              {}
func (l *StubLogger) SetOutput(w io.Writer)          {}
func (l *StubLogger) SetPrefix(prefix string)        {}
