package draftTerm

import (
        "os"
        "syscall"
        "unsafe"
)

type winsize struct {
        ws_row    uint16
        ws_col    uint16
        ws_xpixel uint16
        ws_ypixel uint16
}


//
// Resize the terminal
// to evaluate if to suggest as enhanchement of "github.com/kr/pty"
func ResizeTerminal(pTerminalFd *os.File, pCols int, pRows int) (vErr error) {
        vNewSize:=winsize{ws_row:uint16(pRows), ws_col:uint16(pCols)}
	_, _, vErrNo := syscall.Syscall(syscall.SYS_IOCTL, pTerminalFd.Fd(),syscall.TIOCSWINSZ,uintptr(unsafe.Pointer(&vNewSize)))
        if vErrNo != 0 {
		
                return syscall.Errno(vErrNo)
        }
        return nil
}

