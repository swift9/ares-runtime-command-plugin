package runtime

import (
	"github.com/codeskyblue/go-sh"
	"github.com/swift9/ares-sdk/runtime"
	"syscall"
	"time"
)

type CommandRuntime struct {
	runtime.Runtime
	ShSession  *sh.Session
	IdleFunc   func() int
	HealthFunc func() int
	Dir        string
}

type logWriter struct {
	r *CommandRuntime
}

func (w *logWriter) Write(bytes []byte) (n int, err error) {
	w.r.Emit("log", string(bytes))
	return len(bytes), nil
}

func (r *CommandRuntime) Start(cmd string, args ...string) int {
	irayLogWriter := logWriter{
		r: r,
	}
	r.ShSession.Stdout = &irayLogWriter

	var err error = nil
	go func() {
		if args == nil || len(args) < 1 {
			err = r.ShSession.Command(cmd).Run()
		} else {
			cmdArgs := make([]interface{}, len(args))
			for i := range args {
				cmdArgs[i] = args[i]
			}
			err = r.ShSession.Command(cmd, cmdArgs...).Run()
		}
		r.Emit("exit", err)
	}()

	if err != nil {
		return 1
	}

	return 0
}

func (r *CommandRuntime) Stop() {
	if r.ShSession != nil {
		r.ShSession.Kill(syscall.SIGKILL)
	}
}

func (r *CommandRuntime) Idle() int {
	return r.IdleFunc()
}

func (r *CommandRuntime) Health() int {
	return r.HealthFunc()
}

func (r *CommandRuntime) Init() {
	r.CreateTime = time.Now()
	r.ShSession = sh.NewSession()
	r.ShSession.SetDir(r.Dir)
}

func New(idleFunc func() int, healthFunc func() int, dir string) runtime.IRuntime {
	var r runtime.IRuntime = &CommandRuntime{
		IdleFunc:   idleFunc,
		HealthFunc: healthFunc,
		Dir:        dir,
	}
	r.Init()
	return r
}
