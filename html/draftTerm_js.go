package html



//
// Javascript client
//
const DRAFTTERM_JS=
`

//
//	Trace show on browser console diagnostic messages
//
const DEBUG_ENABLED=false;

function writeTraceMessage(pMessage) {
	if (DEBUG_ENABLED) {
		console.trace(pMessage);
	}
}

function TerminalBufferLine(pColumnSize) {
	this._empty=true
	this._text = new Array(pColumnSize);
	this._textColor=new Array(pColumnSize);
	this._backgroundColor=new Array(pColumnSize);

	this.getCharAt = function(pColumn) {
		return this._text[pColumn];
	};

	this.setCharAt = function(pColumn, pChar) {
		this._empty=false
		this._text[pColumn] = pChar;
	};

	this.removeCharAt = function (pColumn) {
		for (vCnt=pColumn;vCnt<this._text.length;vCnt++) {
			this._text[vCnt] = vCnt < ( this._text.length-1) ? this._text[vCnt+1] : null;
			this._textColor[vCnt] = vCnt < ( this._textColor.length-1) ? this._textColor[vCnt+1] : null;
			this._backgroundColor[vCnt] = vCnt < ( this._backgroundColor.length-1) ? this._backgroundColor[vCnt+1] : null;
		}
	};

	this.addCharAt = function (pColumn) {

		for (vCnt=pColumn+1;vCnt<this._text.length;vCnt++) {
			this._text[vCnt] = vCnt > 1 ? this._text[vCnt-1] : null;
			this._textColor[vCnt] = vCnt > 1 ? this._textColor[vCnt-1] : null;
			this._backgroundColor[vCnt] = vCnt > 1 ? this._backgroundColor[vCnt-1] : null;
		}
	};

	this.setTextColorAt = function(pColumn,pColor) {
		this._textColor[pColumn]=pColor;
	};

	this.getTextColorAt = function(pColumn) {
		return this._textColor[pColumn];
	};

	this.setBackgroundColorAt = function(pColumn,pColor) {
		this._backgroundColor[pColumn]=pColor;
	};

	this.getBackgroundColorAt = function(pColumn) {
		return this._backgroundColor[pColumn];
	};

	this.asString = function () {
		vRis="";
		for ( var vCnt=0;vCnt < this._text.length ;vCnt++) {
			vCurChar=this._text[vCnt];
			vRis= vRis+(vCurChar==null?'?':vCurChar);
		}
		return vRis;
	};

	this.isEmpty = function() {
		return this._empty=0;
	};


}


const DEFAULT_SIZE_COLUMNS=80;
const DEFAULT_SIZE_ROWS=40;

function TerminalBuffer() {

	this._size_columns=DEFAULT_SIZE_COLUMNS;
	this._size_rows=DEFAULT_SIZE_ROWS;
	this._margin_top=0;
	this._margin_bottom=this._size_rows-1;

	this.clear= function() {
		this._lines = new Array(this._size_columns);
		this._currentLine=0;
		this._currentColumn=0;
		this._currentTextColor="grey";
		this._currentBackgroundColor="black";
	};

	this.clear();

	this.getCharAt = function(pColumn, pRow) {

		//if (pRow > this._currentLine) {
		//	return null;
		//}

		vCurRow= this._lines[pRow]

		if (vCurRow !=null ) {
			return vCurRow.getCharAt(pColumn);
		}

		//return "?";
		return null;
	}

	this.getColorsAt = function(pColumn, pRow) {

		//if (pRow > this._currentLine) {
		//	return [null,null];
		//}

		vCurRow= this._lines[pRow]

		if (vCurRow !=null ) {
			return [ vCurRow.getTextColorAt(pColumn), vCurRow.getBackgroundColorAt(pColumn)];
		}

		//return "?";
		return [null,null];
	}


	this.newLine = function() {
		this._currentColumn=0;
		this._currentLine++;

		while (this._currentLine >this._margin_bottom) {
			this._currentLine--;

			for (vCnt=this._margin_top; vCnt<this._margin_bottom;vCnt++) {
				this._lines[vCnt] = this._lines[vCnt+1];
			}
			this._lines[this._margin_bottom]=new TerminalBufferLine(this._size_columns);
		}

	};

	this.moveBack = function (){
		this._currentColumn--;
		if (this._currentColumn< 0) {
			if (this._currentLine >0) {
				this._currentColumn=this._size_columns-1;
				this._currentLine--
			} else {
				this._currentColumn=0;
			}
		}
	}

	this.writeChar = function(pChar) {
		if (this._currentColumn >= this._size_columns) {
			this.newLine();
		}
		vCurLine =this._lines[this._currentLine];
		if (vCurLine==null){
			vCurLine = new TerminalBufferLine(this._size_columns);
			this._lines[this._currentLine] = vCurLine;
		}
		vCurLine.setCharAt(this._currentColumn,pChar);
		vCurLine.setTextColorAt(this._currentColumn,this._currentTextColor);
		vCurLine.setBackgroundColorAt(this._currentColumn,this._currentBackgroundColor);
		this._currentColumn++;
	};

	this.controlCarachter = function(pCharacterCode) {

		switch(pCharacterCode) {
			case CC_NUL:
				return true;
			case CC_BS:
				this.moveBack();
				return true;
			case CC_LF:
				this.newLine();
				return true;
			case CC_VT:
				this.newLine();
				return true;
			case CC_FF:
				this.newLine();
				return true;
			case CC_CR:
				this._currentColumn=0;
				return true;

		}
		return false;
	};

	this.escapeSequence = function(pSequence) {
		pLastChar=pSequence.charAt(pSequence.length-1);
		switch (pLastChar) {
			case 'm':
				vDirectives=pSequence.substr(1,pSequence.length-2).split(';');
				for (vCnt=0;vCnt<vDirectives.length;vCnt++) {
					switch (vDirectives[vCnt]) {

						case "":
							this._currentTextColor="grey";
							this._currentBackgroundColor="black";
							break;
						case "0":
							this._currentTextColor="grey";
							this._currentBackgroundColor="black";
							break;
						case "1":
							this._currentTextColor="white";
							break;
						case "7":
							vTmp=this._currentTextColor;
							this._currentTextColor=this._currentBackgroundColor;
							this._currentBackgroundColor=vTmp;
							break;
						case "27":
							vTmp=this._currentTextColor;
							this._currentTextColor=this._currentBackgroundColor;
							this._currentBackgroundColor=vTmp;
							break;
						case "30":
							this._currentTextColor="black";
							break;
						case "31":
							this._currentTextColor="red";
							break;
						case "32":
							this._currentTextColor="green";
							break;
						case "33":
							this._currentTextColor="brown";
							break;
						case "34":
							this._currentTextColor="blue";
							break;
						case "35":
							this._currentTextColor="magenta";
							break;
						case "36":
							this._currentTextColor="cyan";
							break;
						case "37":
							this._currentTextColor="white";
							break;
						case "39":
							this._currentTextColor="gray";
							break;
						case "40":
							this._currentBackgroundColor="black";
							break;
						case "41":
							this._currentBackgroundColor="red";
							break;
						case "42":
							this._currentBackgroundColor="green";
							break;
						case "43":
							this._currentBackgroundColor="brown";
							break;
						case "44":
							this._currentBackgroundColor="blue";
							break;
						case "45":
							this._currentBackgroundColor="magenta";
							break;
						case "46":
							this._currentBackgroundColor="cyan";
							break;
						case "47":
							this._currentBackgroundColor="white";
							break;
						case "49":
							this._currentBackgroundColor="black";
							break;
						default:
							console.log("Discarded escape sequence ["+vDirectives[vCnt]+"m");
							break;
					}
				}

				return true;

			case 'K':
				this.eraseToTheEndOfTheLine();
				return true;

			//CLEAR SCREEN
			case 'J':

				switch (pSequence) {
					case "[J":
						this.clearDown();
						return true;
					case "[0J":
						this.clearDown();
						return true;
					case "[1J":
						this.clearUp();
						return true;
					case "[2J":
						this.clear();
						return true;
				}
				break;

			case 'H':
				var vMoveCursorRegEx = /\[(\d+?);(\d+?)H/;
				vMatch=vMoveCursorRegEx.exec(pSequence);
				if (vMatch!=null) {
					writeTraceMessage("Moved cursor to "+vMatch[1]+ " "+vMatch[2]);
					this._currentLine=parseInt(vMatch[1])-1;
					this._currentColumn=parseInt(vMatch[2])-1;
					return true;
				}

				this._currentColumn=0;
				this._currentLine=0;
				return true;

			case 'r':
				var vMoveCursorRegEx = /\[(\d+?);(\d+?)r/;
				vMatch=vMoveCursorRegEx.exec(pSequence);
				if (vMatch!=null) {
					writeTraceMessage("DECSTBM set margin to "+vMatch[1]+ " "+vMatch[2]);
					this._margin_top=parseInt(vMatch[1])-1;
					this._margin_bottom=parseInt(vMatch[2])-1;
					return true;
				}
				break;

			case 'd':
				var vMoveCursorRegEx = /\[(\d+?)d/;
				vMatch=vMoveCursorRegEx.exec(pSequence);
				if (vMatch!=null) {
					writeTraceMessage("VPA set line to "+vMatch[1]);
					this._currentLine=parseInt(vMatch[1])-1;
					return true;
				}
				break;

			case 'G':
				var vMoveCursorRegEx = /\[(\d+?)d/;
				vMatch=vMoveCursorRegEx.exec(pSequence);
				if (vMatch!=null) {
					writeTraceMessage("CHA set column to "+vMatch[1]);
					this._currentColumn=parseInt(vMatch[1])-1;
					return true;
				}
				break;

			case 'A':
				var vMoveCursorRegEx = /\[(\d+?)A/;
				vMatch=vMoveCursorRegEx.exec(pSequence);
				if (vMatch!=null) {
					writeTraceMessage("Move cursor up "+vMatch[1]+ " lines");
					this._currentLine=this._currentLine-parseInt(vMatch[1]);
					if (this._currentLine < 0) {
						this._currentLine=0;
					}
					return true;
				}
				break

			case 'B':
				var vMoveCursorRegEx = /\[(\d+?)B/;
				vMatch=vMoveCursorRegEx.exec(pSequence);
				if (vMatch!=null) {
					writeTraceMessage("Move cursor down "+vMatch[1]+ " lines");
					this._currentLine=this._currentLine+parseInt(vMatch[1]);
					return true;
				}
				break;

			case 'C':
				var vMoveCursorRegEx = /\[(\d+?)C/;
				vMatch=vMoveCursorRegEx.exec(pSequence);
				if (vMatch!=null) {
					writeTraceMessage("Move cursor right "+vMatch[1]+ " cols");
					this._currentColumn=this._currentColumn+parseInt(vMatch[1]);
					return true;
				}
				break;

			case 'D':
				var vMoveCursorRegEx = /\[(\d+?)D/;
				vMatch=vMoveCursorRegEx.exec(pSequence);
				if (vMatch!=null) {
					writeTraceMessage("Move cursor left "+vMatch[1]+ " cols");
					this._currentColumn=this._currentColumn-parseInt(vMatch[1]);
					return true;
				}
				break;

			case 'P':
				var vDeleteChar = /\[(\d+?)P/;
				vMatch=vDeleteChar.exec(pSequence);
				if (vMatch!=null) {
					writeTraceMessage("Delete "+vMatch[1]+ " characters");
					for (vCnt=0; vCnt<parseInt(vMatch[1]); vCnt++) {
						this._lines[this._currentLine].removeCharAt(this._currentColumn);
					}
					return true;
				}
				break;

			case '@':
				var vAddChar = /\[(\d+?)@/;
				vMatch=vAddChar.exec(pSequence);
				if (vMatch!=null) {
					writeTraceMessage("Adding "+vMatch[1]+ " characters");
					for (vCnt=0; vCnt<parseInt(vMatch[1]); vCnt++) {
						this._lines[this._currentLine].addCharAt(this._currentColumn);
					}
					return true;
				}
				break;
		}

		return false;
	};

	this.getY = function() {
		return this._currentLine;
	};

	this.getX = function() {
		return this._currentColumn;
	};

	this.clearUp= function() {
		for (vCntRow=this._currentLine;vCntRow>-1;vCntRow--) {
			 this._lines[vCntRow]=null;
		}
	};

	this.clearDown= function() {
		for (vCntRow=this._currentLine;vCntRow<this._size_rows;vCntRow++) {
			 this._lines[vCntRow]=null;
		}
	};

	this.eraseToTheEndOfTheLine= function () {
		var vCurLine=this._lines[this._currentLine];
		if (vCurLine !=null) {
			for (var vCntColumn=this._currentColumn;vCntColumn<this._size_columns;vCntColumn++) {
				vCurLine.setCharAt(vCntColumn,null);
			}
		}
	};

	this.resize = function(pColumns, pRows) {
		this._size_columns=pColumns;
		this._size_rows=pRows;
		this._margin_bottom=pRows-1;
	}


}

function TerminalVideoCanvas(pCanvas,pBuffer) {

	this._canvas=pCanvas;
	this._buffer=pBuffer;
	this._renderPending=false;
	this._renderDate=new Date();
	this._backgroundColor=null;
	this._title="";
	this._statusBarMessage=["",null,null];
	this._enabled=true;

	this.resize=function(pColumns,pRows) {
		this._canvas.width=pColumns*11;
		if (DEBUG_ENABLED) {
			this._canvas.width+=100;
		}
		this._canvas.height=(pRows+2)*21;
		this.render();
		this.renderTitle();
		this.renderStatusBarMessage();
	};

	this.setEnabled = function(pEnabled) {
		this._enabled=pEnabled;
		this.renderTitle();
	};

	this.render=function() {
		vCurDate=new Date();
		this._renderPending =true;
		if (vCurDate.getTime()-this._renderDate.getTime()>100) {
			this.renderImmediate();
		}
	};

	window.setInterval(function() {
		if (this._renderPending==true)  {
			this.renderImmediate();
		}
	}.bind(this),200);

	this.setBackgroundColor=function(pBackgroundColor) {
		this._backgroundColor=pBackgroundColor;
		this.render();
	}

	this.renderImmediate = function() {
		this._renderDate=new Date();
		this._renderPending=false;
		vContext = this._canvas.getContext("2d");
		vContext.font="20px Courier";
		vContext.textBaseline="top";

		for (var vCurRow=0;vCurRow<this._buffer._size_rows;vCurRow++) {
			//Clear the line
			vContext.fillStyle=(this._backgroundColor == null ? "black" : this._backgroundColor);
			vContext.fillRect(0,(vCurRow+1)*21,this._canvas.width,21);

			for (var vCurCol=0;vCurCol<this._buffer._size_columns;vCurCol++) {
				var vCurChar = this._buffer.getCharAt(vCurCol,vCurRow);
				if (vCurChar ==null) {
					continue;
				}
				var vColors=this._buffer.getColorsAt(vCurCol,vCurRow);
				if (vColors[1] != this._backgroundColor) {
					vContext.fillStyle=vColors[1];
					vContext.fillRect(vCurCol*11,(vCurRow+1)*21,11,21);
				}
				vContext.fillStyle=vColors[0];
				vContext.fillText(vCurChar,vCurCol*11,(vCurRow+1)*21);
			}
			if (DEBUG_ENABLED) {
				vContext.fillStyle="white";
				vContext.fillRect(this._canvas.width-100,(vCurRow+1)*21,150,21);
				vContext.fillStyle="black";
				vContext.fillText(vCurRow+"/"+this._buffer._size_rows,this._canvas.width-100,(vCurRow+1)*21);
			}
		}

		//Cursor
		vContext.fillStyle="green";
		vContext.fillRect( this._buffer.getX()*11,(( this._buffer.getY()+1)*21)+19,11,2);

		if (DEBUG_ENABLED) {
			vContext.fillStyle="white";
			vContext.fillRect(this._canvas.width-100,(this._buffer._size_rows+1)*21,this._canvas.width,21);
			vContext.fillStyle="red";
			vContext.fillText(this._buffer.getX()+","+this._buffer.getY(),this._canvas.width-100,(this._buffer._size_rows+1)*21);
		}
	};

	this.escapeSequence = function(pSequence) {
		//WINDOWS TITLE
		if (pSequence.startsWith("]0;")) {
			this.writeTitle(pSequence.substr(3));
			return true;
		}
	};

	this.writeTitle = function(pTitle) {
		this._title=pTitle;
		this.renderTitle();
	};

	this.renderTitle = function() {
		vContext = this._canvas.getContext("2d");
		vContext.fillStyle= this._enabled ? "blue": "gray";
		vContext.fillRect(0,0,this._canvas.width,21);
		vContext.fillStyle="white";
		vContext.fillText(this._title,0,0);
	};

	this.writeStatusBarMessage = function(pMessage,pTextColor,pBackgroundColor) {
		this._statusBarMessage=[pMessage,pTextColor,pBackgroundColor];
		this.renderStatusBarMessage();
	};

	this.renderStatusBarMessage=function() {
		vContext = this._canvas.getContext("2d");
		vContext.font="20px Courier";
		vContext.textBaseline="top";
		vContext.fillStyle=this._statusBarMessage[2]==null?"darkgray":this._statusBarMessage[2];
		vContext.fillRect(0,(this._buffer._size_rows+1)*21,this._canvas.width,21);
		vContext.fillStyle=this._statusBarMessage[1]==null?"black":this._statusBarMessage[1];
		vContext.fillText(this._statusBarMessage[0],0,(this._buffer._size_rows+1)*21);
	};

	this.getOptimalSize =function() {
		vOptimalX=Math.trunc(this._canvas.clientWidth / 11);
		vOptimalY=Math.trunc(this._canvas.clientHeight / 21)-2;
		return [vOptimalX,vOptimalY];
	};
}


const IGNORED_KEYS=['Control','Alt', 'Shift'];


function CommunicationChannel(pRemoteUrl) {

	this._connected=false;
	this.connect(pRemoteUrl);
}

CommunicationChannel.prototype.connect=function(pUrl) {
	this._url=pUrl;
	this._webSocket=new WebSocket(this._url);
	this._webSocket.addEventListener("open", function() {
		this._connected = true;	
		console.log("[WEBSOCKET "+this._url+"] opened connection");
		this.onConnected(this._url);
	}.bind(this));

	this._webSocket.addEventListener("message", function(pEvent) {
		 writeTraceMessage("[WEBSOCKET "+this._url+"]  received {{{"+pEvent.data+"}}}");
		 this.onDataToDisplay(pEvent.data);
	}.bind(this));

	this._webSocket.addEventListener("error", function(pEvent) {
		console.log("[WEBSOCKET "+this._url+"] error " + JSON.stringify(pEvent));
		this.onConnectionError();
	}.bind(this));

	this._webSocket.addEventListener("close", function() {
	   this._connected=false;
	   console.log("[WEBSOCKET "+this._url+"] connection closed");
	   this.onConnectionClosed();
	}.bind(this));
}



CommunicationChannel.prototype.onDataToDisplay = function(pData) {
}



CommunicationChannel.prototype.onConnected = function() {
}



CommunicationChannel.prototype.onConnectionError = function() {
}



CommunicationChannel.prototype.onConnectionClosed = function() {
}



CommunicationChannel.prototype.isConnected= function () {
		return this._connected;
}



CommunicationChannel.prototype.forwardKeyEvent = function(pEvent) {

	if (!this._connected) {
		this.connect(this._url);
	}

	pEvent.preventDefault();

	if (IGNORED_KEYS.indexOf(pEvent.key) > -1) {
		return;
	}

	if (   this.isConnected()==false ) {
		console.log("[WEBSOCKET "+this._url+"] discarded key event due to disconnected channel");
		return;
	}

	this.send ({ type:"key", code:pEvent.keyCode, key:pEvent.key, ctrl: pEvent.ctrlKey, alt: pEvent.altKey, shift: pEvent.shiftKey });
}



CommunicationChannel.prototype.send = function(pMessage) {
	this._webSocket.send (JSON.stringify(pMessage));
}



const CC_NUL=0, CC_BEL=0x7,CC_BS=0x8, CC_HT=0x9, CC_LF=0xa, CC_VT=0xb, CC_FF=0xc, CC_CR=0xd, CC_SO=0xe, CC_SI=0xf, CC_CAN=0x18, CC_SUB=0x1a, CC_ESC=0x1b, CC_DEL=0x7f;
const CC_ALL=[  CC_NUL, CC_BEL, CC_BS, CC_HT, CC_LF, CC_VT, CC_FF, CC_CR, CC_SO, CC_SI, CC_CAN, CC_SUB, CC_ESC, CC_DEL]



const ESCAPE_PATTERNS =[
/\[\d+@/,
/\[[0-9]{0,1}m/,
/\[.+[a-zA-Z]$/,
/[A-Z=>]{1}$/,
/[\(\)][A-Z0-9]$/,
/\].*\7/,
/#\d$/,
/\?\d+[lh]$/
];



function isControlCharacterByCode(pCode) {
		return CC_ALL.indexOf(pCode) >-1;
}



function getEscapeSequence(pString) {
	for (var vCntPattern =0;vCntPattern<ESCAPE_PATTERNS.length;vCntPattern++) {
		if (pString.search(ESCAPE_PATTERNS[vCntPattern])!=-1){
			return pString;
		}
	}
	//Means not finished
	return null;
}



//http://man7.org/linux/man-pages/man4/console_codes.4.html
function TerminalOutputParser() {

	this._escaping=false;
	this._EscapeSequence="";

	this.onControlCharacter= function(pCharacterCode) {
		return false;
	};

	this.onWriteChar = function(pCharacter) {
		return false;
	};

	this.onEscapeSequence = function(pSequence) {
		return false;
	};

	this.parse=function(pString) {
		for (var vCnt=0;vCnt < pString.length;vCnt++) {
				vChar=pString.charAt(vCnt);
				vCharCode=pString.charCodeAt(vCnt);

				if (this._escaping==false){

					if ( isControlCharacterByCode(vCharCode)) {

						if (vCharCode == CC_ESC) {
							this._escaping=true;
							this._EscapeSequence="";
						} else {
							if (true!=this.onControlCharacter(vCharCode)) {
								console.warn("Ignored control character "+vCharCode+"!!!");
							}
						}
						continue;
					}

					if (true !=this.onWriteChar(vChar)) {
							console.log("IGNORED: on write char "+vChar);
					}
				} else {
					if  (vChar=="[" && this._EscapeSequence.length > 0) {
						console.error("Something wrong on parsing of escape sequence, discarded"+this._EscapeSequence);
						this._EscapeSequence="[";
					}else {
						this._EscapeSequence+=vChar;
					}
					var vEscapeSequence = getEscapeSequence(this._EscapeSequence);
					//PENDING
					if (vEscapeSequence == null ) {
						continue;
					}
					this._escaping=false;
					if (true !=this.onEscapeSequence(vEscapeSequence)) {
						console.warn("Escape sequence discarded "+vEscapeSequence+"!!!");
					}
				}
		}
	};
}



function TerminalEmulator(pCanvas, pUrl) {
		this._terminalBuffer=new TerminalBuffer();
		this._display =  new TerminalVideoCanvas(pCanvas,this._terminalBuffer);

		this._TerminalOutputParser=new TerminalOutputParser();
		this._enabled=true;
		this._resized=0;

		this.isEnabled = function() {
			return this._enabled;
		};

		this.setEnabled = function(pEnabled) {
			this._enabled=pEnabled;
			this._display.setEnabled(pEnabled);
		};

		this.resize = function(pColumns,pRows) {
			this._terminalBuffer.resize(pColumns,pRows);
			this._display.resize(pColumns,pRows);
		};

		this.adjustSize = function() {
			vSize=this._display.getOptimalSize();
			console.log("Resizing to "+vSize[0]+" cols "+vSize[1]+" rows");
			this.resize(vSize[0],vSize[1]);
			if (this._communicationChannel!=null && this._communicationChannel.isConnected) {
				this._communicationChannel.send({ type:"resize", cols:vSize[0], rows:vSize[1]});
			}
			return vSize;
		};

		//this.adjustSize();
		window.addEventListener("resize", function(pEvent) {
			var vResizedByMe=++this._resized;
			setTimeout(function() {
				if (this._resized>vResizedByMe) {
					//console.log("Resize skipped because pending "+vResizedByMe+"/"+this._resized);
					return;
				}
				this.adjustSize();
			}.bind(this), 500);
		}.bind(this));

		pCanvas.addEventListener("mouseover", function() {
			this.setEnabled(true);
		}.bind(this));

		pCanvas.addEventListener("mouseout", function() {
			this.setEnabled(false);
		}.bind(this));

		this._display.render();

		document.addEventListener("keydown", function(pEvent) {
			if (this.isEnabled()) {
				this._communicationChannel.forwardKeyEvent(pEvent);
			}
		}.bind(this));

		this._communicationChannel= new CommunicationChannel(pUrl);
		this._communicationChannel.onConnected = function(pUrl) {
			vSize=this.adjustSize();
			this._terminalBuffer.clear();
			this._display.writeStatusBarMessage("Connected to "+pUrl,"darkgrey","darkgreen");
		}.bind(this);

		this._communicationChannel.onConnectionError = function() {
			this._display.writeStatusBarMessage("connection error","yellow","darkred");
		}.bind(this);

		this._communicationChannel.onConnectionClosed = function() {
			this._display.writeStatusBarMessage("connection closed","white","darkgray");
		}.bind(this);

		this._communicationChannel.onDataToDisplay = function(pData){
			this._TerminalOutputParser.parse(pData);
			this._display.render();
		}.bind(this);

		this._TerminalOutputParser.onWriteChar=function (pChar) {
			this._terminalBuffer.writeChar(pChar);
			return true;
		}.bind(this);

		this._TerminalOutputParser.onControlCharacter=function (pCharCode) {
			return this._terminalBuffer.controlCarachter(pCharCode);
		}.bind(this);

		this._TerminalOutputParser.onEscapeSequence=function (pSequence) {
			//Forward escape sequence first on terminal buffer 
			if (this._terminalBuffer.escapeSequence(pSequence)) {
				return true;
			}
			//then on display
			return this._display.escapeSequence(pSequence);
		}.bind(this);
}
`;




