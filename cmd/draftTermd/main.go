package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"flag"
	"fmt"
	"net/http"
	"github.com/mysinmyc/draftTerm"
	"github.com/mysinmyc/draftTerm/html"
	"golang.org/x/net/websocket"
)



var (
	flag_Secure=flag.Bool("secure", false, "Enable protocol encryption")
	flag_CertFile=flag.String("cert", "", "Public certificate file")
	flag_KeyFile=flag.String("key", "", "Private key file")
	flag_Listen=flag.String("listen", "0.0.0.0:8080", "Listening address")
	flag_Command=flag.String("cmd", "", "Initial command")
)



func createWebSocketUrlFromRequest(pWebSocketPath string, pRequest *http.Request) string {
	var vProtocol string

	//Removed test of ssl due to ssl in case of ssl offloading
	//if pRequest.TLS == nil {
	if *flag_Secure {
		vProtocol = "wss"
	}else {
		vProtocol = "ws"
	}
	return fmt.Sprintf("%s://%s/%s", vProtocol,pRequest.Host,pWebSocketPath)
}



func getUserNameById(pId int)(vName string, vError error) {
	vCommand:=exec.Command("id","-u","-n",fmt.Sprintf("%d",pId))
	vOutput,vError:=vCommand.Output()
	return strings.TrimRight(string(vOutput),"\n"),vError
}



func evaluateInitialCommand() (vCommand string, vArguments []string,vError error) {

	if *flag_Command == "" {
		vUid:=os.Getuid()
		if vUid==0 {
			vCommand="/bin/login"
		} else {

			vUserName,vError:=getUserNameById(vUid)
			if vError==nil {
				vCommand="/bin/su"
				vArguments=[]string{"-",vUserName}
			}
		}

	} else {
		vCommandSplitted:=strings.Split(*flag_Command," ")
		vCommand=vCommandSplitted[0]
		if len(vCommandSplitted) > 0 {
			vArguments= vCommandSplitted[1:]
		}
	}
	return
}

func main() {

	flag.Parse()

	vCommand,vArguments,vError:=evaluateInitialCommand()

	if vError!= nil  {
		log.Fatalf("Error evaluating command to execute %v",vError)
	}

	log.Printf("Initial command for terminal: %s %s",vCommand,strings.Join(vArguments," "))

	http.Handle("/draftTerm.socket",websocket.Handler(draftTerm.NewTerminalServer(vCommand, vArguments...).Handler))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w,strings.Replace(html.INDEX_HTML,"ws://localhost:8080/draftTerm.socket",createWebSocketUrlFromRequest("draftTerm.socket", r),-1))
	})

	http.HandleFunc("/draftTerm.js", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w,html.DRAFTTERM_JS)
	})

	var vErr error
	if *flag_Secure {
		log.Printf("Listening TLS on %s",*flag_Listen)
		vErr=http.ListenAndServeTLS(*flag_Listen, *flag_CertFile, *flag_KeyFile, nil)
	} else {
		log.Printf("Listening on %s",*flag_Listen)
		vErr=http.ListenAndServe(*flag_Listen, nil)
	}
	log.Fatal(vErr)
}
