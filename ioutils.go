package draftTerm

import (
	"io"
	"time"
	"unicode/utf8"
)



//
// Remove invalid bytes that can cause issues
//
func cleanInvalidBytes( pBytesToClean []byte) []byte {

	logMessage("decoding rune...")
	vRis:=make([]byte,len(pBytesToClean))
	vProcessedSize:=0
	for vProcessedSize< len(pBytesToClean) {
		vCurRune, _ := utf8.DecodeRune(pBytesToClean[vProcessedSize:vProcessedSize+1])
		if vCurRune < 2 {
			pBytesToClean[vProcessedSize]=byte('?')	
		}
		vProcessedSize++
	}
	//logDebug(fmt.Sprintf("decoded rune:%s",vRis))
	return vRis
}



//
// dump to the logs a seguence of bytes
//
func dumpBytes(pBytes []byte) {

	vBytes:=make([]byte,20)
	for vCur:=0;vCur< len(pBytes); vCur++ {
		vBytes =append(vBytes, pBytes[vCur])
		if vCur % 20 == 0 {
			logDebug("| %s |  %v \n", string(vBytes), vBytes)
			vBytes=make([]byte,20)
		}
	}
	if (len(vBytes) >0) {
		logDebug("| %s |  %v \n", string(vBytes), vBytes)
	}
}



func pipe(pReader io.ReadCloser, pWriter io.WriteCloser, pSyncChannel chan bool) error {
	defer pReader.Close()
	defer pWriter.Close()
	vBytes := make([]byte, 8000)
	for {
		vLen, vErr := pReader.Read(vBytes)
		if vErr != nil {
			return vErr
		}
		if vLen > 0 {
			if IsDebug() {
				dumpBytes(vBytes[0:vLen])
			}

			/*
			if utf8.Valid(vBytes)==false {
				logMessage(fmt.Sprintf("WARNING invalid sequence %c",vBytes[0:vLen]))
				vBytes = cleanInvalidBytes(vBytes)

			}
			*/
			_,vErr= pWriter.Write(vBytes[0:vLen])
			if vErr != nil {
				return vErr
			}
		}

		//SyncChannel allow to wait a sync message from client before sending other data
		//has been introduced to avoid an overload the client in case fast refresh (example find /)
		//To avoid infinite wait in case of client issues ha been introduced a timeout
		if pSyncChannel != nil {
			vTimeoutChannel:= time.After(10000*time.Millisecond)
			select {
				case <-pSyncChannel:	
				case <-vTimeoutChannel:
					logMessage("Timeout waiting for sync channel")
			
			}
		}
	}
}

