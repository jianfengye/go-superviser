package main

import(
    "bytes"
    "os/exec"
    "sync"
    "fmt"
    "go/build"
)

func restart(root string) {
    stopRun()

    // TODO: change to getGopath
    context := build.Context.Default
    splits := strings.Split(root, "/")
    bin := context.GOPATH + "/bin/" + splits[len(splits) - 1]
    out, err := run(context.GOPATH + "/bin/", "go", "build", "-o", bin, splits[len(splits) - 1])
    defer os.Remove(bin)
    if err != nil {
        return
    }
    return run("", bin)
}

func stopRun() {
    running.Lock()
    if running.cmd != nil {
        running.cmd.Process.Kill()
        running.cmd = nil
    }
    running.Unlock()
}


var running struct {
    sync.Mutex
    cmd *exec.Cmd
}

func run(dir string, args ...string) ([]byte, error) {
    var buf bytes.Buffer
    cmd := exec.Command(args[0], args[1:]...)
    cmd.Dir = dir
    cmd.Stdout = &buf
    cmd.Stderr = cmd.Stdout

    // Start command and leave in 'running'.
    running.Lock()
    if running.cmd != nil {
        defer running.Unlock()
        return nil, fmt.Errorf("already running %s", running.cmd.Path)
    }
    if err := cmd.Start(); err != nil {
        running.Unlock()
        return nil, err
    }
    running.cmd = cmd
    running.Unlock()

    // Wait for the command.  Clean up,
    err := cmd.Wait()
    running.Lock()
    if running.cmd == cmd {
        running.cmd = nil
    }
    running.Unlock()
    return buf.Bytes(), err
}