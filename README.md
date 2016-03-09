# wsc

A simplistic tool for sending and receiving websocket messages from a command line.
Mainly useful to test websocket servers.

Getting started:
```
$ go get github.com/raphael/wsc
$ wsc ws://echo.websocket.org
2016/03/08 22:51:51 connecting to ws://echo.websocket.org...
2016/03/08 22:51:52 ready, exit with CTRL+C.
foo 
>> foo
<< foo
^C
exiting
```
