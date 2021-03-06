package draftTerm

import (
	"encoding/json"
	"os"
	"io"
	"fmt"
)


type Byted interface {
	Bytes() []byte
}



type inputMessage struct {
        Type  string `json:"type"`
}



func (self inputMessage) String() string {
        return fmt.Sprintf("Generic message of type %s", self.Type)
}



func (self inputMessage) Bytes() [] byte {
	return []byte(self.String())
}



//
// Message received when a key has been pressed 
//   Data are obtained from KeyBoardEvent 
type keyPressedMessage struct {
	inputMessage
        Key   string `json:"key"`	
        Code  byte   `json:"code"`
        Shift bool   `json:"shift"`
        Ctrl  bool   `json:"ctrl"`
        Alt   bool   `json:"alt"`
}



func (self keyPressedMessage) String() string {
        vJson, _ := json.Marshal(self)
        return string(vJson)
}


//
// Convert input keyboard data into data to send to the pty
// 
func (self keyPressedMessage)Bytes() []byte {

        if len(self.Key) > 1 {
                switch self.Key {
			case "Spacebar":
				return []byte{' '}
			case "Tab":
				return []byte{'\t'}
			case "Backspace":
				return []byte{'\b'}
			case "Escape":
				return []byte{27}
			case "Enter":
				return []byte{10}
			case "ArrowUp" , "Up":
				return []byte{27, 91, 65}
			case "ArrowDown" , "Down":
				return []byte{27, 91, 66}
			case "ArrowRight" , "Right":
				return []byte{27, 91, 67}
			case "ArrowLeft" , "Left":
				return []byte{27, 91, 68}
                }
        } else {
                if self.Ctrl {
                        if self.Code > 64 && self.Code < 91 {
                                //OUTPUT <- fmt.Sprintf("Pressed ctrl with %s %v ",self.Key,self.Code)
                                return []byte{self.Code-64}
                        }
                }
                vCode:= []byte(self.Key)
                return vCode
        }
        /*
                if self.Code <32 || self.Code >136 {
                        return []byte {}
                }
                if  self.Shift ==false && self.Code >=65 && self.Code <=90 {
                        return   []byte {self.Code+32}
                }else {
                        return   []byte {self.Code}
                }
        */
        return [] byte {}
}



//
// Message received when the client ask for terminal resize
//
type resizeTerminalMessage struct {
	inputMessage
        Cols  int `json:"cols"`
        Rows  int   `json:"rows"`
}

func (self resizeTerminalMessage) String() string {
        return fmt.Sprintf("Resize terminal to %d cols %d rows",self.Cols,self.Rows)
}



type inputChannel struct {
	input io.Reader
	pty *os.File
}



//
// 
//
func NewInputChannel(pInput io.Reader, pPty *os.File) *inputChannel {
	
	return &inputChannel{ input:pInput, pty: pPty }
}	


func (self *inputChannel) ProcessIncomingMessage(pSyncChannel chan bool) error {

	vBytes := make([]byte, 8000)

        vLen, vError:= self.input.Read(vBytes)
        if vError !=nil  {
                return vError
        }
      
	if vLen ==1 {
		pSyncChannel<-true	
		return nil
	} 

	if vLen ==0 {
		return nil
	} 

        //fmt.Printf("Received %s \n",vBytes)


        vTmp:=&inputMessage{}
        vError=json.Unmarshal(vBytes[0:vLen],vTmp)
        if vError !=nil {
                return vError;
        }

        switch vTmp.Type {
                case "key":
                        vKeyPressed:=&keyPressedMessage{}
                        vError=json.Unmarshal(vBytes[0:vLen],vKeyPressed)
			if (vError!=nil) {
				return vError
			}
			if IsDebug() {
	        		logDebug("Processing keyboard message %v... \n",vKeyPressed)
			}
			_,vError=self.pty.Write(vKeyPressed.Bytes())
			return vError
		
                case "resize":
                        vResize:=&resizeTerminalMessage{}
                        vError=json.Unmarshal(vBytes[0:vLen],vResize)
			if (vError!=nil) {
				return vError
			}
			if IsDebug() {
	        		logDebug("Processing resize message %v... \n",vResize)
			}
			vError=ResizeTerminal(self.pty, vResize.Cols,vResize.Rows)
                        return vError
        }


        return nil


}
