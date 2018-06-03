# Chat Server
This is a dumb server, which just relays the mwssage from one connecion to another.
just start server by setting environment variable: MY_SERVER_PORT
and you can use telnet as your chat client.
e.g.
## Server Side
```
export MY_SERVER_PORT 6789
go run main.go
```

## Client side
telnet localhost 6789
```
Rajnis-Air:~ rajnikant$ telnet localhost 6789
Trying ::1...
Connected to localhost.
Escape character is '^]'.
You are connected , please enter your user name: **rajni**
rajni: You are conected
```
you can connect multiple users by running multiple telnet client.

### Communicate
For communicating just send:
```
{"From":"rajni","To":"Sharad","Message":"Hi Rajni, how are you"
```
**You can develope clients arounnd it and play.**
