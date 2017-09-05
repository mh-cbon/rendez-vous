# rendez-vous

UDP meeting point server


# cli

#### $ go run main.go -h
```sh
Usage:
  commands [OPTIONS] <command>

Help Options:
  -h, --help  Show this help message

Available commands:
  browser  Run a browser to visit websites within a rendez-vous network.
  client   Run rendez-vous client
  http     Run an http request using a rendez-vous client.
  serve    Run rendez-vous server
  website  Run and announce a website on given rendez-vous remote.
```

# tests

#### $ go test -v
```sh
=== RUN   Test1
=== RUN   Test1/1
go run main.go serve -l 8070
2017/09/05 17:14:41 Listening... :8070
go run main.go client -q ping -r :8070
2017/09/05 17:14:42 [::]:53347       snd :8070           : txid=1      len=10   "\n\x04ping"
2017/09/05 17:14:42 [::]:8070        rcv [::1]:53347     : txid=1      len=10   "\n\x04ping"
2017/09/05 17:14:42 [::]:8070        snd [::1]:53347     : txid=1      len=20   "\x10\xc8\x012\v[::1]:53347"
2017/09/05 17:14:42 [::]:53347       rcv [::1]:8070      : txid=1      len=20   "\x10\xc8\x012\v[::1]:53347"
model.Message{Query:"", Code:200, Pbk:[]uint8(nil), Sign:[]uint8(nil), Value:"", Address:"[::1]:53347", Response:""}
=== RUN   Test1/2
go run main.go serve -l 8090
2017/09/05 17:14:42 Listening... :8090
go run main.go website -r :8090 -l 8091 --local 8092 --pvk 202d229c0f09f41c858066496b21c27e59266ec7c5b0933275518b351da5e92e --dir demows
2017/09/05 17:14:43 Public Website listening on  [::]:8091
2017/09/05 17:14:43 Local Website listening on  127.0.0.1:8092
pvk= 202d229c0f09f41c858066496b21c27e59266ec7c5b0933275518b351da5e92e
pbk= b6b8113748fe0795658fa9d6ab3e36d27d72e97b7df407e7a8080d61ec405d74
sig= cc273ddc49feed05551cc891c1b6bbc86974e63a4974812525929453beb7e917fed3d1b583b9d81bfe7cab9e1d15c092c24ea7764f18d1bef1e048744ed6770b
2017/09/05 17:14:43 [::]:8091        snd :8090           : txid=1      len=109  "\n\x03reg\x1a \xb6\xb8\x117H\xfe\a\x95e\x8f\xa9֫>6\xd2}r\xe9{}\xf4\a\xe7\xa8\b\ra\xec@]t\"@\xcc'=\xdcI\xfe\xed\x05U\x1cȑ\xc1\xb6\xbb\xc8it\xe6:It\x81%%\x92\x94S\xbe\xb7\xe9\x17\xfe\xd3ѵ\x83\xb9\xd8\x1b\xfe|\xab\x9e\x1d\x15\xc0\x92\xc2N\xa7vO\x18Ѿ\xf1\xe0HtN\xd6w\v"
2017/09/05 17:14:43 [::]:8090        rcv [::1]:8091      : txid=1      len=109  "\n\x03reg\x1a \xb6\xb8\x117H\xfe\a\x95e\x8f\xa9֫>6\xd2}r\xe9{}\xf4\a\xe7\xa8\b\ra\xec@]t\"@\xcc'=\xdcI\xfe\xed\x05U\x1cȑ\xc1\xb6\xbb\xc8it\xe6:It\x81%%\x92\x94S\xbe\xb7\xe9\x17\xfe\xd3ѵ\x83\xb9\xd8\x1b\xfe|\xab\x9e\x1d\x15\xc0\x92\xc2N\xa7vO\x18Ѿ\xf1\xe0HtN\xd6w\v"
2017/09/05 17:14:43 [::]:8090        snd [::1]:8091      : txid=1      len=19   "\x10\xc8\x012\n[::1]:8091"
2017/09/05 17:14:43 [::]:8091        rcv [::1]:8090      : txid=1      len=19   "\x10\xc8\x012\n[::1]:8091"
2017/09/05 17:14:43 registration  { 200 [] []  [::1]:8091 }
go run main.go http --url http://127.0.0.1:8091/index.html
A demo website.
<img src="tumblr_nosz96g2GT1t0pdxwo1_1280.jpg" />
=== RUN   Test1/3
go run main.go serve -l 8080
2017/09/05 17:14:44 Listening... :8080
go run main.go website -r :8080 -l 8081 --local 8082 --pvk 202d229c0f09f41c858066496b21c27e59266ec7c5b0933275518b351da5e92e --dir demows
2017/09/05 17:14:45 Public Website listening on  [::]:8081
2017/09/05 17:14:45 Local Website listening on  127.0.0.1:8082
pvk= 202d229c0f09f41c858066496b21c27e59266ec7c5b0933275518b351da5e92e
pbk= b6b8113748fe0795658fa9d6ab3e36d27d72e97b7df407e7a8080d61ec405d74
sig= cc273ddc49feed05551cc891c1b6bbc86974e63a4974812525929453beb7e917fed3d1b583b9d81bfe7cab9e1d15c092c24ea7764f18d1bef1e048744ed6770b
2017/09/05 17:14:45 [::]:8081        snd :8080           : txid=1      len=109  "\n\x03reg\x1a \xb6\xb8\x117H\xfe\a\x95e\x8f\xa9֫>6\xd2}r\xe9{}\xf4\a\xe7\xa8\b\ra\xec@]t\"@\xcc'=\xdcI\xfe\xed\x05U\x1cȑ\xc1\xb6\xbb\xc8it\xe6:It\x81%%\x92\x94S\xbe\xb7\xe9\x17\xfe\xd3ѵ\x83\xb9\xd8\x1b\xfe|\xab\x9e\x1d\x15\xc0\x92\xc2N\xa7vO\x18Ѿ\xf1\xe0HtN\xd6w\v"
2017/09/05 17:14:45 [::]:8080        rcv [::1]:8081      : txid=1      len=109  "\n\x03reg\x1a \xb6\xb8\x117H\xfe\a\x95e\x8f\xa9֫>6\xd2}r\xe9{}\xf4\a\xe7\xa8\b\ra\xec@]t\"@\xcc'=\xdcI\xfe\xed\x05U\x1cȑ\xc1\xb6\xbb\xc8it\xe6:It\x81%%\x92\x94S\xbe\xb7\xe9\x17\xfe\xd3ѵ\x83\xb9\xd8\x1b\xfe|\xab\x9e\x1d\x15\xc0\x92\xc2N\xa7vO\x18Ѿ\xf1\xe0HtN\xd6w\v"
2017/09/05 17:14:45 [::]:8080        snd [::1]:8081      : txid=1      len=19   "\x10\xc8\x012\n[::1]:8081"
2017/09/05 17:14:45 [::]:8081        rcv [::1]:8080      : txid=1      len=19   "\x10\xc8\x012\n[::1]:8081"
2017/09/05 17:14:45 registration  { 200 [] []  [::1]:8081 }
go run main.go browser -r :8080 -l 8083 --ws 8085 --proxy 8084 --headless
2017/09/05 17:14:46 me.com server listening on 127.0.0.1:8085
HTTP GET  http://127.0.0.1:8085/index.html
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
HTTP GET  http://b6b8113748fe0795658fa9d6ab3e36d27d72e97b7df407e7a8080d61ec405d74.me.com/index.html
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
--- PASS: Test1 (7.32s)
    --- PASS: Test1/1 (1.02s)
    --- PASS: Test1/2 (2.01s)
    --- PASS: Test1/3 (3.02s)
PASS
ok  	github.com/mh-cbon/rendez-vous	7.329s
```

# todos

#### $ grep --include='*go' -r todo -B 1 -A 1 -n
```sh
server/server.go-41-		case model.Register:
server/server.go:42:			//todo: rendez-vous server should implement a write token
server/server.go-43-
--
server/server.go-68-		case model.Unregister:
server/server.go:69:			//todo: unregister should accept/verify a pbk/sig/value with a special value to identify the query issuer.
server/server.go-70-			if len(v.Pbk) == 0 {
--
server/server.go-90-		case model.Join:
server/server.go:91:			//todo: Join the swarm
server/server.go-92-		case model.Leave:
server/server.go:93:			//todo: leave the swarm
server/server.go-94-		}
--
server/registration.go-2-
server/registration.go:3://todo: add storage clean up with ttl on entry
--
main.go-31-
main.go:32://todo: rendez-vous server should check ttl registrations
main.go:33://todo: rendez-vous server should implement a write token concept to register
main.go:34://todo: rendez-vous server unregister should accept/verify a pbk/sig/value with a special value to identify the query issuer.
main.go-35-
--
main.go-235-	handler := http.FileServer(http.Dir(opts.Dir))
main.go:236:	public := utils.ServeHTTPFromListener(ln, httpServer(handler, "")) //todo: replace with a transparent proxy, so the website can live into another process
main.go-237-	local := httpServer(handler, "127.0.0.1:"+opts.Local)
--
socket/tx.go-66-	t.id++
socket/tx.go:67:	//todo: find a better way
socket/tx.go-68-	if t.id > 10000 {
--
socket/tx.go-85-	_, err := t.UDP.Write(data, remote)
socket/tx.go:86:	//todo: handle _
socket/tx.go-87-	return err
--
socket/tx.go-97-	_, err := t.UDP.Write(data, remote)
socket/tx.go:98:	//todo: handle _
socket/tx.go-99-	return err
```
