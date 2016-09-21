# draftTerm

Concept for a Unix Terminal emulator (frontend javascript, channel websocket, backend go)

The program has a webserver that exposes the javascript client and the websocket to transfer data from/to the terminal



# GOAL OF THE PROJECT

Expose by web terminal/based linux command (interactive or not)



# Status of the project

The initial commit is a draft, a work in progress


# Docker container

Is it possible to test it in a docker container

Some examples

- Run a contaniner that show a terminal with a top running. To connect to the container open a browser to https://{dockerhost}:8443
```
	docker run -d -p 8443:8443  mysinmyc/draftterm top
```


- Run a container with a login session. To connect to the container open a browser to https://{dockerhost}:8443. User id is guest, password is password01 (forced by GUEST_PASSWORD environment variable
```
	docker run -d -p 8443:8443 -e GUEST_PASSWORD=password01 mysinmyc/draftterm
```
 

# Direct Build info

this project depends on "github.com/kr/pty"

to build it execute `go get -t github.com/mysinmyc/draftTerm/cmd/draftTermd`



# Command parameters

```
Usage of draftTermd:
  -cert string
    	Public certificate file
  -cmd string
    	Initial command
  -debug
    	Enable debug
  -key string
    	Private key file
  -listen string
    	Listening address (default "0.0.0.0:8080")
  -secure
    	Enable protocol encryption
```



# Terminal implementation

In the initial implementation only some escape sequences have been implemented

Regarding the input check inside the source code input.go what has been implmented



# Some Samples

The following command start a daemon that expose "top" command to clients. Try to run then open a browser to http://{serverName}:8080/

`draftTermd -listen 0.0.0.0:8080 -cmd top`

The following command start a daemon that show the server date then exit

`draftTermd -listen 0.0.0.0:8080 -cmd date`



# Know issues and limitations

Not all the escape sequences has been tested.

The daemon doesn't passed all the vttest command tests. Will be improved in the next releases

At the moment the buffer of the terminal doens't maintains the history (for scrolling)

