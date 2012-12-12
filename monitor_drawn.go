package main

import (
    "syscall"
)

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
    var (
        buf [syscall.SizeofInotifyEvent * 4069]byte
        n  int
    )

    for {
        n, errno = syscall.Read(fd, buf[0:])
        if n > 0 {
            action <- 1
        }
    }
}