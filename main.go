package main

import (
    "flag"
    "go/build"
)

var (
    project = flag.String("project", nil, "project to be supervised")
)

func main() {
    flag.Parse()

    if project == nil {
        log.Fatalf("Could not find supervise project")
    }

    p, err := build.Default.Import(project, "", build.FindOnly)
    if err != nil {
        log.Fatalf("Couldn't find project: %v", err)
    }

    root = p.Dir
    log.Println("Found the project", project)

    action := make(chan uint32)

    // monitor folder
    go restart(root) //start prject
    go monitor(root, action)

    for {
        a := <- action
        switch {
        case a:
            restart(root)
        default:
            log.Fatalf("monitor internal error")
        }
    }
}

func monitor(root string, action chan uint32) {
    fd, err := syscall.InotifyInit()
    if fd == -1 {
        action <- 0
        return
    }
    
    flags := syscall.IN_MASK_ADD | syscall.IN_MODIFY | syscall.IN_CREATE | syscall.IN_DELETE
    wd, errno := syscall.InotifyAddWatch(fd, root, flags)
    if wd == -1 {
        action <- 0
        return
    }

    for {
        n, errno = syscall.Read(fd, buf[0:])
        if n > 0 {
            action <- 1
        }
    }
}