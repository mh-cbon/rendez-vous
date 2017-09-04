# rendez-vous

UDP meeting point server


# cli

#### $ go run main.go -h
```sh
rendez-vous - noversion

	A server to expose your endpoints with a public key.
	A client to find/register endpoints for a public key.
	A website to expose website.
	A browser to visit remote website.
	An http client to test a remote website..

Usage
	rendez-vous [server|client|website|browser|http] <options>

Usage of main:
  -h	show help
  -help
    	show help
  -v	show version
  -version
    	show version
```

# tests

#### $ go test -v
```sh
=== RUN   Test1
=== RUN   Test1/1
2017/09/05 00:01:58 go run main.go serve -listen 8070
2017/09/05 00:01:58 Listening... :8070
2017/09/05 00:01:59 go run main.go client -query ping -remote :8070
2017/09/05 00:01:59 [::1]:60326 <- {ping 0 [] []   }
model.Message{Query:"", Code:200, Pbk:[]uint8(nil), Sign:[]uint8(nil), Value:"", Address:"[::1]:60326", Response:""}
=== RUN   Test1/2
2017/09/05 00:01:59 go run main.go serve -listen 8090
2017/09/05 00:01:59 Listening... :8090
2017/09/05 00:02:00 go run main.go website -remote :8090 -listen 8091 -local 8092 -pvk 202d229c0f09f41c858066496b21c27e59266ec7c5b0933275518b351da5e92e -static demows
2017/09/05 00:02:00 Public Website listening on  [::]:8091
2017/09/05 00:02:00 Local Website listening on  127.0.0.1:8092
pvk= 202d229c0f09f41c858066496b21c27e59266ec7c5b0933275518b351da5e92e
pbk= b6b8113748fe0795658fa9d6ab3e36d27d72e97b7df407e7a8080d61ec405d74
sig= cdd8ea95c3f2957edf80d6a77d7efd93c5678453637a6fdfd6b2c8da286a93d3a36fa67dd246d598752bcd1c0f29ff542a37d853db98747f7d667c9307bc190a
2017/09/05 00:02:00 [::1]:8091 <- {reg 0 [182 184 17 55 72 254 7 149 101 143 169 214 171 62 54 210 125 114 233 123 125 244 7 231 168 8 13 97 236 64 93 116] [205 216 234 149 195 242 149 126 223 128 214 167 125 126 253 147 197 103 132 83 99 122 111 223 214 178 200 218 40 106 147 211 163 111 166 125 210 70 213 152 117 43 205 28 15 41 255 84 42 55 216 83 219 152 116 127 125 102 124 147 7 188 25 10] website  }
2017/09/05 00:02:00 registration  { 200 [] []  [::1]:8091 }
2017/09/05 00:02:01 go run main.go http -url http://127.0.0.1:8091/index.html
A demo website.
<img src="tumblr_nosz96g2GT1t0pdxwo1_1280.jpg" />
=== RUN   Test1/3
2017/09/05 00:02:01 go run main.go serve -listen 8080
2017/09/05 00:02:01 Listening... :8080
2017/09/05 00:02:02 go run main.go website -remote :8080 -listen 8081 -local 8082 -pvk 202d229c0f09f41c858066496b21c27e59266ec7c5b0933275518b351da5e92e -static demows
2017/09/05 00:02:02 Public Website listening on  [::]:8081
2017/09/05 00:02:02 Local Website listening on  127.0.0.1:8082
pvk= 202d229c0f09f41c858066496b21c27e59266ec7c5b0933275518b351da5e92e
pbk= b6b8113748fe0795658fa9d6ab3e36d27d72e97b7df407e7a8080d61ec405d74
sig= cdd8ea95c3f2957edf80d6a77d7efd93c5678453637a6fdfd6b2c8da286a93d3a36fa67dd246d598752bcd1c0f29ff542a37d853db98747f7d667c9307bc190a
2017/09/05 00:02:02 [::1]:8081 <- {reg 0 [182 184 17 55 72 254 7 149 101 143 169 214 171 62 54 210 125 114 233 123 125 244 7 231 168 8 13 97 236 64 93 116] [205 216 234 149 195 242 149 126 223 128 214 167 125 126 253 147 197 103 132 83 99 122 111 223 214 178 200 218 40 106 147 211 163 111 166 125 210 70 213 152 117 43 205 28 15 41 255 84 42 55 216 83 219 152 116 127 125 102 124 147 7 188 25 10] website  }
2017/09/05 00:02:02 registration  { 200 [] []  [::1]:8081 }
2017/09/05 00:02:03 go run main.go browser -remote :8080 -listen 8083 -ws 8085 -proxy 8084 -headless
2017/09/05 00:02:03 me.com server listening on 127.0.0.1:8085
<html>
<head>
  <title>me.com: myself on the internet</title>
  <script src="jquery-3.2.1.min.js"></script>
</head>
<body>

  <div align="center">
    Find: <input type="text" id="search" style="width:60%" /><input type="button" value="find" id="bt" />
  </div>
  <script>
    $("#search").css("border", "solid 1px gray")
    $("#bt").click(function(){
      var v = $("#search").val();
      if (v == "") {
        $("#search").css("border", "solid 1px red")
      } else {
        $("#search").css("border", "solid 1px gray")
        window.open("http://"+v+".me.com/")
      }
    })
  </script>
</body>
</html>
<html>
<head>
  <title>me.com: myself on the internet</title>
  <script src="jquery-3.2.1.min.js"></script>
</head>
<body>

  <div align="center">
    Find: <input type="text" id="search" style="width:60%" /><input type="button" value="find" id="bt" />
  </div>
  <script>
    $("#search").css("border", "solid 1px gray")
    $("#bt").click(function(){
      var v = $("#search").val();
      if (v == "") {
        $("#search").css("border", "solid 1px red")
      } else {
        $("#search").css("border", "solid 1px gray")
        window.open("http://"+v+".me.com/")
      }
    })
  </script>
</body>
</html>
--- PASS: Test1 (7.11s)
    --- PASS: Test1/1 (1.02s)
    --- PASS: Test1/2 (2.02s)
    --- PASS: Test1/3 (3.02s)
PASS
ok  	github.com/mh-cbon/rendez-vous	7.116s
```

# todos

#### $ grep --include='*go' -r todo -B 1 -A 1 -n
```sh
server/server.go-67-	case model.Register:
server/server.go:68:		//todo: rendez-vous server should implement a write token
server/server.go-69-
--
server/server.go-94-	case model.Unregister:
server/server.go:95:		//todo: unregister should accept/verify a pbk/sig/value with a special value to identify the query issuer.
server/server.go-96-		if len(v.Pbk) == 0 {
--
server/server.go-116-	case model.Join:
server/server.go:117:		//todo: Join the swarm
server/server.go-118-	case model.Leave:
server/server.go:119:		//todo: leave the swarm
server/server.go-120-	}
--
server/registration.go-2-
server/registration.go:3://todo: add storage clean up with ttl on entry
--
main.go-66-
main.go:67://todo: rendez-vous server should check ttl registrations
main.go:68://todo: rendez-vous server should impelment a write token concept to register
main.go:69://todo: rendez-vous server unregister should accept/verify a pbk/sig/value with a special value to identify the query issuer.
main.go-70-
--
socket/tx.go-66-	t.id++
socket/tx.go:67:	//todo: find a better way
socket/tx.go-68-	if t.id > 10000 {
--
socket/tx.go-84-	_, err := t.UDP.Write(data, remote)
socket/tx.go:85:	//todo: handle _
socket/tx.go-86-	return err
--
socket/tx.go-95-	_, err := t.UDP.Write(data, remote)
socket/tx.go:96:	//todo: handle _
socket/tx.go-97-	return err
```
