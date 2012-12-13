package main

import (
    "flag"
    "go/build"
    "log"
)

var (
    project = flag.String("project", "", "project to be supervised")
    withrun = flag.Bool("withrun", false, "whether to run after build the project")
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
    log.Println("Found the project", *project, ":", root)

    action := make(chan uint32)

    // monitor folder
    log.Println("start project")
    go restart(root) //start prject

    log.Println("start monitor")
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