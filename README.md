# rendez-vous

UDP meeting point server.

# cli

```sh
# the meeting point
$ go run main.go -op server
# register yourself on the server
go run main.go -query=register -op client -remote "104.197.173.209:8080" -value dd -pbk 16050fff0df375c9390165a2ba0639dde9c14d453fa0067968bfb05382a9cef7 -sign 038b5eb599a6f77c284dcf8b4288da8672d2890d9e6db21f473553ff9a75296dea0cabdac281d74212cb82e2f918bd2b974a13ee5f22221e62db599b2201ce07 -port 8081
$ go run main.go -query=register -op client -remote ":8080" -value dd -pbk 16050fff0df375c9390165a2ba0639dde9c14d453fa0067968bfb05382a9cef7 -sign 038b5eb599a6f77c284dcf8b4288da8672d2890d9e6db21f473553ff9a75296dea0cabdac281d74212cb82e2f918bd2b974a13ee5f22221e62db599b2201ce07 -port 8081
# find a peer with given pbk
$ go run main.go -query=find -op=client -remote=":8080" -pbk 16050fff0df375c9390165a2ba0639dde9c14d453fa0067968bfb05382a9cef7 -port 8082
go run main.go -query=find -op=client -remote="104.197.173.209:8080" -pbk 16050fff0df375c9390165a2ba0639dde9c14d453fa0067968bfb05382a9cef7 -port 8082
# query the peer found
$ go run main.go -query=find -op=client -remote=":8081" -pbk 16050fff0df375c9390165a2ba0639dde9c14d453fa0067968bfb05382a9cef7 -port 8082

go run main.go -query=ping -op=client -remote="23.227.123.8:8081" -port 8082

```
