package runtime

import (
	"github.com/codeskyblue/go-sh"
	"github.com/swift9/ares-sdk/runtime"
	"syscall"
	"time"
)

type ShellRuntime struct {
	runtime.Runtime
	ShSession  *sh.Session
	IdleFunc   func() int
	HealthFunc func() int
	Dir        string
}

type logWriter struct {
	runtime *ShellRuntime
}

func (irayLogWriter *logWriter) Write(bytes []byte) (n int, err error) {
	irayLogWriter.runtime.Emit("log", string(bytes))
	return len(bytes), nil
}

func (runtime *ShellRuntime) Start(cmd string, args ...string) int {
	irayLogWriter := logWriter{
		runtime: runtime,
	}
	runtime.ShSession.Stdout = &irayLogWriter

	var err error = nil
	go func() {
		if args == nil || len(args) < 1 {
			err = runtime.ShSession.Command(cmd).Run()
		} else {
			cmdArgs := make([]interface{}, len(args))
			for i := range args {
				cmdArgs[i] = args[i]
			}
			err = runtime.ShSession.Command(cmd, cmdArgs...).Run()
		}
		runtime.Emit("exit", err)
	}()

	if err != nil {
		return 1
	}

	return 0
}

func (runtime *ShellRuntime) Stop() {
	if runtime.ShSession != nil {
		runtime.ShSession.Kill(syscall.SIGKILL)
	}
}

func (runtime *ShellRuntime) ReStart(cmd string, args ...string) int {
	runtime.Stop()
	return runtime.Start(cmd, args...)
}

func (runtime *ShellRuntime) Idle() int {
	return runtime.IdleFunc()
}

func (runtime *ShellRuntime) Health() int {
	return runtime.HealthFunc()
}

func (runtime *ShellRuntime) Init() {
	runtime.CreateTime = time.Now()
	runtime.ShSession = sh.NewSession()
	runtime.ShSession.SetDir(runtime.Dir)
}

func New(idleFunc func() int, healthFunc func() int, dir string) *runtime.IRuntime {
	var runtime runtime.IRuntime = &ShellRuntime{
		IdleFunc:   idleFunc,
		HealthFunc: healthFunc,
		Dir:        dir,
	}
	runtime.Init()
	return &runtime
}
