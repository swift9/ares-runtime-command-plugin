package runtime

import (
	"errors"
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
	r := New("/opt/ihome/iray")
	var a = r.Start("tail", "-f", "1.log")
	println(a)
	r.On("log", func(data string) {
		if strings.Contains(data, "out of memory") ||
			strings.Contains(data, "aborting render") {
			r.Emit("exit", errors.New(data))
		}
	})
	r.On("exit", func(err error) {
		log.Println(err.Error())
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
