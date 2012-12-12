package main

import (
    "flag"
    "go/build"
    "log"
)

var (
    project = flag.String("project", "", "project to be supervised")
)

func main() {
    flag.Parse()

    if project == nil {
        log.Fatalf("Could not find supervise project")
    }

    p, err := build.Default.Import(*project, "", build.FindOnly)
    if err != nil {
        log.Fatalf("Couldn't find project: %v", err)
    }

    root := p.Dir
    log.Println("Found the project", project)

    action := make(chan uint32)

    // monitor folder
    go restart(root) //start prject
    go monitor(root, action)

    for {
        a := <- action
        switch {
        case a > 0:
            restart(root)
        default:
            log.Fatalf("monitor internal error")
        }
    }
}