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

func (self keyPressedMessage)Bytes() []byte {

        if len(self.Key) > 1 {
                switch self.Key {
			case "Tab":
				return []byte{'\t'}
			case "Backspace":
				return []byte{'\b'}
			case "Escape":
				return []byte{27}
			case "Enter":
				return []byte{10}
			case "ArrowUp":
				return []byte{27, 91, 65}
			case "ArrowDown":
				return []byte{27, 91, 66}
			case "ArrowRight":
				return []byte{27, 91, 67}
			case "ArrowLeft":
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


func NewInputChannel(pInput io.Reader, pPty *os.File) *inputChannel {
	
	return &inputChannel{ input:pInput, pty: pPty }
}	


func (self *inputChannel) ProcessIncomingMessage() error {

	vBytes := make([]byte, 8000)

        vLen, vError:= self.input.Read(vBytes)
        if vError !=nil  {
                return vError
        }
       
	if vLen == 0 {
		return nil
	} 

        fmt.Printf("Received %s \n",vBytes)
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
        		//fmt.Printf("Processing message %s... \n",vKeyPressed)
			_,vError=self.pty.Write(vKeyPressed.Bytes())
			return vError
		
                case "resize":
                        vResize:=&resizeTerminalMessage{}
                        vError=json.Unmarshal(vBytes[0:vLen],vResize)
			if (vError!=nil) {
				return vError
			}
			//fmt.Printf("Processing message %s ...\n",vResize)
			vError=ResizeTerminal(self.pty, vResize.Cols,vResize.Rows)
                        return vError
        }


        return nil


}
