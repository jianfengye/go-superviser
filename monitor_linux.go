package main

import (
    "syscall"
    "unsafe"
    "log"
)

func monitor(root string, action chan uint32) {
    fd, err := syscall.InotifyInit()
    if fd == -1 || err != nil {
        end <- true
        return
    }
    
    flags := syscall.IN_MODIFY | syscall.IN_CREATE | syscall.IN_DELETE // 2/128/512
    wd, _ := syscall.InotifyAddWatch(fd, root, uint32(flags))
    if wd == -1 {
        end <- true
        return
    }
    var (
        buf [syscall.SizeofInotifyEvent * 10]byte
        n  int
    )

    for {
        n, _ = syscall.Read(fd, buf[0:])
        if n > syscall.SizeofInotifyEvent {
            var offset = 0
            for offset < n {
                raw := (*syscall.InotifyEvent)(unsafe.Pointer(&buf[offset]))
                mask := uint32(raw.Mask)
                offset = offset + int(raw.Len) + syscall.SizeofInotifyEvent
                action <- mask

                log.Println("action:", mask)
            }
        }
    }
}