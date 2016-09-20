package html


const INDEX_HTML=
`
<html>
<head>
<script src="draftTerm.js"  type="text/javascript"></script>

<script  type="text/javascript">
	
var vTerminal1;
function initTerminal() {
	vTerminal1=new TerminalEmulator(myCanvas1,"ws://localhost:8080/draftTerm.socket");
}
</script>
</head>
<body onload="initTerminal()">
	<canvas id="myCanvas1" style="width:100%; height:100%">Not supported by your browser</canvas>
</body>
</html>

`
