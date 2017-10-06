package log

import (
	"fmt"
	"os"
)

var Debug = (os.Getenv("DEBUG") == "1")

type Exiter interface {
	Exit(status int)
}

type DefaultExiter struct {
}

func (e DefaultExiter) Exit(status int) {
	os.Exit(status)
}

var exiter Exiter = DefaultExiter{}

// TODO Handle concurrent access correctly
func GetExiter() Exiter {
	return exiter
}

// TODO Handle concurrent access correctly
func SetExiter(e Exiter) {
	exiter = e
}

func Exit(status int) {
	exiter.Exit(status)
}

func Fatal(v ...interface{}) {
	Print(v...)
	exiter.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	Printf(format, v...)
	exiter.Exit(1)
}

func Fatalln(v ...interface{}) {
	Println(v...)
	exiter.Exit(1)
}

func Print(v ...interface{}) {
	fmt.Print(v...)
}

func Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func Println(v ...interface{}) {
	fmt.Println(v...)
}

func Trace(v ...interface{}) {
	if Debug {
		Print(v...)
	}
}

func Tracef(format string, v ...interface{}) {
	if Debug {
		Printf(format, v...)
	}
}

func Traceln(v ...interface{}) {
	if Debug {
		Println(v...)
	}
}
