package runtime

import (
	event "github.com/swift9/ares-event"
	"io"
	"log"
	"os"
	"os/exec"
)

type Env struct {
	Name  string
	Value string
}

type Command struct {
	Envs     []Env
	Cmd      string
	Args     []string
	Dir      string
	Addition map[string]string
}

type Runtime struct {
	event.Emitter
	Meta          map[string]interface{}
	Command       Command
	Cmd           *exec.Cmd
	LogEverything bool
	LogWriter     io.Writer
}

type logWriter struct {
	r *Runtime
}

func (w *logWriter) Write(bytes []byte) (n int, err error) {
	w.r.Emit("log", string(bytes))
	if w.r.LogEverything && w.r.LogWriter != nil {
		w.r.LogWriter.Write(bytes)
	}
	return len(bytes), nil
}

func (r *Runtime) getLogWriter() io.Writer {
	if r.LogWriter != nil {
		return r.LogWriter
	}
	return os.Stdout
}

func (r *Runtime) Start() int {
	command := r.Command
	r.Cmd = exec.Command(command.Cmd, command.Args...)
	r.Cmd.Dir = command.Dir

	if len(r.Cmd.Env) == 0 {
		r.Cmd.Env = []string{}
	}

	for _, environ := range os.Environ() {
		r.Cmd.Env = append(r.Cmd.Env, environ)
	}

	for _, env := range command.Envs {
		r.Cmd.Env = append(r.Cmd.Env, env.Name+"="+env.Value)
	}

	enableGroupKill(r.Cmd)

	r.Cmd.Stdout = &logWriter{
		r: r,
	}
	r.Cmd.Stderr = &logWriter{
		r: r,
	}
	err := r.Cmd.Start()
	if err != nil {
		r.getLogWriter().Write([]byte("start error " + err.Error()))
		r.Emit("exit", 1)
		return 1
	}
	r.Emit("ready")
	err = r.Cmd.Wait()
	if err != nil {
		r.getLogWriter().Write([]byte("wait error " + err.Error()))
		r.Emit("exit", 1)
		return 1
	} else {
		log.Println("exit")
		r.getLogWriter().Write([]byte("exit"))
		r.Emit("exit", 0)
		return 0
	}
}

func (r *Runtime) Stop() int {
	defer func() {
		if e := recover(); e != nil {
			r.getLogWriter().Write([]byte("unknown"))
		}
	}()
	r.killGroup(r.Cmd.Process.Pid)
	r.kill()
	return 0
}

func (r *Runtime) kill() error {
	defer func() {
		if e := recover(); e != nil {
			r.getLogWriter().Write([]byte("unknown"))
		}
	}()
	return r.Cmd.Process.Kill()
}

func (r *Runtime) Idle() int {
	return -1
}

func (r *Runtime) Health() Status {
	return NewStatusUp()
}

func (r *Runtime) Init() {
}

func NewRuntime(command Command, log io.Writer, logEverything bool) IRuntime {
	var r IRuntime = &Runtime{
		Command:       command,
		LogWriter:     log,
		LogEverything: logEverything,
	}
	r.Init()
	return r
}
