package runtime_test

import (
	"errors"
	runtime "github.com/swift9/ares-runtime-command-plugin"
	"log"
	"strings"
	"testing"
	"time"
)

func restart() {
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()
	r := runtime.NewCommandRuntime("~")
	var a = r.Start("tail", "-f", "1.log")
	println(a)
	r.On("log", func(data string) {
		if strings.Contains(data, "out of memory") ||
			strings.Contains(data, "aborting render") {
			r.Emit("exit", errors.New(data))
		}
	})
	r.On("exit", func(code int) {
		log.Println(code)
		r.Stop()
		go func() {
			time.Sleep(3 * time.Second)
			r.Start("tail", "-f", "1.log")
		}()
	})
}

func TestStart(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()
	restart()
	time.Sleep(1 * time.Hour)
}
