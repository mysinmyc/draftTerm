package draftTerm

import (
	"fmt"
	"log"
	"os"
)

var (
        _debug_Channel = make(chan string,1000)
	_debug_Logger = log.New(os.Stderr,"DEBUG",0)
	_debug  = false
)


func debugWriter() {
        for {
                log.Printf(<-_debug_Channel)
        }
}

func logDebug(pFormat string, pParm ...interface{}) {
	if _debug == false {
		return;
	}
	_debug_Channel<-fmt.Sprintf(pFormat, pParm...);
}

func logMessage(pMessage string) {
	log.Println(pMessage);
}

func logError(pMessage string,pError error) {
	log.Printf("ERROR:%s (%v)!!!\n",pMessage,pError);
}

func init() {
	go debugWriter();
}

func SetDebug(pDebug bool) {
	_debug = pDebug;
}

func IsDebug() bool {
	return _debug;
}
