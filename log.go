package draftTerm

import (
	"fmt"
	"log"
)

var (
        OUTPUT = make(chan string,1000)
)


func outputHandler() {
        for {
                log.Printf(<-OUTPUT)
        }
}

func logMessage(pMessage string) {
	OUTPUT<-pMessage;
}

func logError(pMessage string,pError error) {
	OUTPUT<-fmt.Sprintf("ERROR:%s (%v)!!!\n ",pMessage,pError);
}

func init() {
	go outputHandler();
}
