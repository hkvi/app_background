//go:build !NDEBUG

package hkvilog

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

func Info(a ...any) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("hkvi: %s:%s:%d %s", time.Now().Format("2006-01-02 15:04:05"), file, line, fmt.Sprintln(a...))
}
func Infof(format string, a ...any) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("hkvi: %s:%s:%d %s\n", time.Now().Format("2006-01-02 15:04:05"), file, line, fmt.Sprintf(format, a...))
}

func Error(a ...any) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("hkvi error: %s:%s:%d %s%s\n", time.Now().Format("2006-01-02 15:04:05"), file, line, fmt.Sprintln(a...), debug.Stack())
}

func Errorf(format string, a ...any) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("hkvi error: %s:%s:%d %s\n%s\n", time.Now().Format("2006-01-02 15:04:05"), file, line, fmt.Sprintf(format, a...), debug.Stack())
}
