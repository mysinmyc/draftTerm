package draftTerm

import (
	"github.com/kr/pty"
	"golang.org/x/net/websocket"
	"os/exec"
	"fmt"
)




type terminalServer struct {
	command string
	arguments []string
}



//
// Create new instance of TerminalServer
//
func NewTerminalServer(pCommand string, pArguments ...string) *terminalServer {

	return &terminalServer{ command: pCommand,arguments: []string(pArguments)}
}



//
// WebSocket Handler 
// 	it execute the command and manage i/o trough the pty for each websocket connection
func (self *terminalServer) Handler(ws *websocket.Conn) {

	logMessage(fmt.Sprintf("Opened connection from %v ",ws.Request().RemoteAddr))
	defer ws.Close()

	vCmd := exec.Command(self.command, self.arguments...)
	vStdinPipe, err := pty.Start(vCmd)
	defer vStdinPipe.Close()
	if err != nil {
		logError("Failed to start command in the pty",err)
		return
	}

	vSyncChannel:=make(chan bool,1)
	go pipe(vStdinPipe, ws, vSyncChannel)


	vInputChannel:= NewInputChannel(ws,vStdinPipe)
	for {
		vErr:=vInputChannel.ProcessIncomingMessage(vSyncChannel)
		if vErr != nil {
			return
		}
	}

}

