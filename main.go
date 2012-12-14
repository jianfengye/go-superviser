package main

import (
    "flag"
    "go/build"
    "log"
)

var (
    project = flag.String("project", "", "project to be supervised")
    withrun = flag.Bool("run", false, "whether to run after build the project")
    end = make(chan bool)
)

func main() {
    flag.Parse()

    if *project == "" {
        println("Params Error:")
        flag.PrintDefaults()
        return
    }

    p, err := build.Default.Import(*project, "", build.FindOnly)
    if err != nil {
        log.Fatalf("Couldn't find project: %v", err)
    }

    root := p.Dir
    action := make(chan uint32)
    go monitor(root, action)

    restart(root)

    for {
        select {
        case <- action:
            out, err := restart(root)
            if err != nil {
                log.Println("===============notice================")
                log.Println("error out:", string(out))
                log.Println("could not restart the project: %v", err)
                log.Println("======================================")
            }
        case <- end:
            log.Fatalf("superviser end")
        }
    }
}