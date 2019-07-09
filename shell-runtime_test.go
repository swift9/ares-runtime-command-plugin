package runtime

import (
	"github.com/swift9/ares-sdk/runtime"
	"strings"
	"testing"
	"time"
)

func restart() {
	var irayRuntime *runtime.IRuntime = New(func() int {
		return 0
	}, func() int {
		return 0
	}, "~")

	(*irayRuntime).On("log", func(data string) {
		if strings.Contains(data, "out of memory") ||
			strings.Contains(data, "aborting render") {
			(*irayRuntime).Stop()
		}
	})
	(*irayRuntime).Start("tail", "-f", "1.log")

	(*irayRuntime).On("exit", func(err error) {
		println(err.Error())
		go func() {
			time.Sleep(3 * time.Second)
			restart()
		}()
	})
}

func TestStart(t *testing.T) {
	restart()

	time.Sleep(1 * time.Hour)

}
