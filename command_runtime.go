package runtime

import (
	"encoding/json"
	event "github.com/swift9/ares-event"
	"io"
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

func (c *Command) toJSON() string {
	bs, err := json.Marshal(c)
	if err != nil {
		return "{\"cmd\": \"" + c.Cmd + "\" }"
	}
	return string(bs)
}

type Runtime struct {
	event.Emitter
	Meta        map[string]interface{}
	Command     Command
	Cmd         *exec.Cmd
	EnableEvent bool
	LogWriter   io.Writer
}

type emptyLogWriter struct {
}

func (w *emptyLogWriter) Write(bytes []byte) (n int, err error) {
	return len(bytes), nil
}

type runtimeLogWriterProxy struct {
	r *Runtime
}

func (w *runtimeLogWriterProxy) Write(bytes []byte) (n int, err error) {
	if w.r.EnableEvent {
		w.r.Emit("log", string(bytes))
	}
	return w.r.LogWriter.Write(bytes)
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

	enhanceCmd(r.Cmd)

	r.Cmd.Stdout = &runtimeLogWriterProxy{
		r: r,
	}
	r.Cmd.Stderr = &runtimeLogWriterProxy{
		r: r,
	}
	r.getLogWriter().Write([]byte("runtime starting command:" + r.Command.toJSON()))
	err := r.Cmd.Start()
	if err != nil {
		r.getLogWriter().Write([]byte("runtime start error " + err.Error() + " command:" + r.Command.toJSON()))
		if r.EnableEvent {
			r.Emit("exit", 1)
		}
		return 1
	}
	if r.EnableEvent {
		r.Emit("ready")
	}
	err = r.Cmd.Wait()
	if err != nil {
		r.getLogWriter().Write([]byte("runtime wait error " + err.Error() + " command:" + r.Command.toJSON()))
		if r.EnableEvent {
			r.Emit("exit", 1)
		}
		return 1
	} else {
		if r.EnableEvent {
			r.Emit("exit", 0)
		}
		return 0
	}
}

func (r *Runtime) Stop() int {
	defer func() {
		if e := recover(); e != nil {
			r.getLogWriter().Write([]byte("runtime stop unknown error command:" + r.Command.toJSON()))
		}
	}()
	r.killGroup(r.Cmd.Process.Pid)
	r.kill()
	return 0
}

func (r *Runtime) kill() error {
	defer func() {
		if e := recover(); e != nil {
			r.getLogWriter().Write([]byte("runtime kill unknown error command:" + r.Command.toJSON()))
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

func NewRuntime(command Command, log io.Writer, enableEvent bool) IRuntime {
	if log == nil {
		log = &emptyLogWriter{}
	}
	var r IRuntime = &Runtime{
		Command:     command,
		LogWriter:   log,
		EnableEvent: enableEvent,
	}
	r.Init()
	return r
}
