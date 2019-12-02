package runtime_test

import (
	runtime "github.com/swift9/ares-runtime-command-plugin"
	"log"
	"testing"
	"time"
)

func restart() {
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()
	r := runtime.NewRuntime(runtime.Command{Cmd: "tail", Args: []string{"-f", "go.mod"}})
	r.Start()
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
