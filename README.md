# rendez-vous

UDP meeting point server


# cli

#### $ go run main.go -h
```sh
Usage:
  rendez-vous [OPTIONS] <command>

Help Options:
  -h, --help  Show this help message

Available commands:
  browser  Run a browser to visit websites within a rendez-vous network.
  client   Run rendez-vous client
  http     Run an http request using a rendez-vous client.
  serve    Run rendez-vous server
  website  Run and announce a website on given rendez-vous remote.
  wsadmin  Run the backend server without node
```

# tests

#### $ go test -v
```sh
=== RUN   Test1
=== RUN   Test1/1
go run main.go serve -l 0.0.0.0:8070
2017/09/10 23:51:07 Listening... [::]:8070
go run main.go client -q ping -r 127.0.0.1:8070
2017/09/10 23:51:08 [::]:47340       snd 127.0.0.1:8070  : txid=1      len=16   "{\"q\":\"ping\"}"
2017/09/10 23:51:08 [::]:8070        rcv 127.0.0.1:47340 : txid=1      len=16   "{\"q\":\"ping\"}"
2017/09/10 23:51:08 [::]:8070        snd 127.0.0.1:47340 : txid=1      len=35   "{\"c\":200,\"a\":\"127.0.0.1:47340\"}"
2017/09/10 23:51:08 [::]:47340       rcv 127.0.0.1:8070  : txid=1      len=35   "{\"c\":200,\"a\":\"127.0.0.1:47340\"}"
model.Message{Query:"", Code:200, Pbk:[]uint8(nil), Sign:[]uint8(nil), Value:"", Address:"127.0.0.1:47340", Data:"", Token:"", PortStatus:0, Peers:[]*model.Peer(nil), Start:0, Limit:0}
=== RUN   Test1/2
go run main.go serve -l 0.0.0.0:8090
2017/09/10 23:51:08 Listening... [::]:8090
go run main.go website -r 127.0.0.1:8090 -l 0.0.0.0:8091 --local 0.0.0.0:8092 --pvk 504bc61393e5d7ea991dbfad4d5bb98093562d472fa22d425a35bcd46341d8f678e7d4c5aa13e3d9a538a5aa2a027cb5343931a48a6fd7b7b1ae699ec8125f12 --dir demows
2017/09/10 23:51:09 [::]:8091        snd 127.0.0.1:8090  : txid=1      len=175  "{\"q\":\"reg\",\"p\":\"eOfUxaoT49mlOKWqKgJ8tTQ5MaSKb9e3sa5pnsgSXxI=\",\"s\":\"pQ7EVE7h8c3Z6wRHoVGsDiMr51gzwG8w3aQot62sp9oKX9ej8aEJ+bm901IugJMHs1xWXW0KQ+LqwnvIhgIXCw==\",\"v\":\"website\"}"
2017/09/10 23:51:09 [::]:8090        rcv 127.0.0.1:8091  : txid=1      len=175  "{\"q\":\"reg\",\"p\":\"eOfUxaoT49mlOKWqKgJ8tTQ5MaSKb9e3sa5pnsgSXxI=\",\"s\":\"pQ7EVE7h8c3Z6wRHoVGsDiMr51gzwG8w3aQot62sp9oKX9ej8aEJ+bm901IugJMHs1xWXW0KQ+LqwnvIhgIXCw==\",\"v\":\"website\"}"
2017/09/10 23:51:09 [::]:8090        snd 127.0.0.1:8091  : txid=1      len=34   "{\"c\":200,\"a\":\"127.0.0.1:8091\"}"
2017/09/10 23:51:09 [::]:8091        rcv 127.0.0.1:8090  : txid=1      len=34   "{\"c\":200,\"a\":\"127.0.0.1:8091\"}"
2017/09/10 23:51:09 Public Website listening on  [::]:8091
2017/09/10 23:51:09 Local Website listening on  0.0.0.0:8092
pvk= 504bc61393e5d7ea991dbfad4d5bb98093562d472fa22d425a35bcd46341d8f678e7d4c5aa13e3d9a538a5aa2a027cb5343931a48a6fd7b7b1ae699ec8125f12
pbk= 78e7d4c5aa13e3d9a538a5aa2a027cb5343931a48a6fd7b7b1ae699ec8125f12
sig= a50ec4544ee1f1cdd9eb0447a151ac0e232be75833c06f30dda428b7adaca7da0a5fd7a3f1a109f9b9bdd3522e809307b35c565d6d0a43e2eac27bc88602170b
2017/09/10 23:51:09 [::]:8091        snd 127.0.0.1:8090  : txid=2      len=33   "{\"q\":\"testport\",\"t\":\"random\"}"
2017/09/10 23:51:09 [::]:8090        rcv 127.0.0.1:8091  : txid=2      len=33   "{\"q\":\"testport\",\"t\":\"random\"}"
2017/09/10 23:51:09 [::]:43885       snd 127.0.0.1:8091  : txid=1      len=54   "{\"q\":\"porttest\",\"a\":\"127.0.0.1:8091\",\"t\":\"random\"}"
2017/09/10 23:51:09 [::]:8091        rcv 127.0.0.1:43885 : txid=1      len=54   "{\"q\":\"porttest\",\"a\":\"127.0.0.1:8091\",\"t\":\"random\"}"
2017/09/10 23:51:09 porttest q:  127.0.0.1:43885 random
2017/09/10 23:51:09 porttest success
2017/09/10 23:51:09 [::]:8091        snd 127.0.0.1:43885 : txid=1      len=48   "{\"c\":200,\"a\":\"127.0.0.1:43885\",\"d\":\"random\"}"
2017/09/10 23:51:09 [::]:43885       rcv 127.0.0.1:8091  : txid=1      len=48   "{\"c\":200,\"a\":\"127.0.0.1:43885\",\"d\":\"random\"}"
2017/09/10 23:51:09 { 200 [] []  127.0.0.1:43885 random  0 [] 0 0}
go run main.go http --url http://127.0.0.1:8091/index.html --remote 127.0.0.1:8090
A demo website.
<img src="tumblr_nosz96g2GT1t0pdxwo1_1280.jpg" />
<img src="tumblr_nosz96g2GT1t0pdxwo1_1280.jpg?1" />
<img src="tumblr_nosz96g2GT1t0pdxwo1_1280.jpg?2" />
<img src="tumblr_nosz96g2GT1t0pdxwo1_1280.jpg?3" />
go run main.go http --url http://78e7d4c5aa13e3d9a538a5aa2a027cb5343931a48a6fd7b7b1ae699ec8125f12.me.com/index.html --remote 127.0.0.1:8090
2017/09/10 23:51:10 [::]:52516       snd 127.0.0.1:8090  : txid=1      len=81   "{\"q\":\"find\",\"p\":\"eOfUxaoT49mlOKWqKgJ8tTQ5MaSKb9e3sa5pnsgSXxI=\",\"v\":\"website\"}"
2017/09/10 23:51:10 [::]:8090        rcv 127.0.0.1:52516 : txid=1      len=81   "{\"q\":\"find\",\"p\":\"eOfUxaoT49mlOKWqKgJ8tTQ5MaSKb9e3sa5pnsgSXxI=\",\"v\":\"website\"}"
2017/09/10 23:51:10 [::]:8090        snd 127.0.0.1:52516 : txid=1      len=222  "{\"c\":200,\"p\":\"eOfUxaoT49mlOKWqKgJ8tTQ5MaSKb9e3sa5pnsgSXxI=\",\"s\":\"pQ7EVE7h8c3Z6wRHoVGsDiMr51gzwG8w3aQot62sp9oKX9ej8aEJ+bm901IugJMHs1xWXW0KQ+LqwnvIhgIXCw==\",\"v\":\"website\",\"a\":\"127.0.0.1:52516\",\"d\":\"127.0.0.1:8091\",\"u\":1}"
2017/09/10 23:51:10 [::]:52516       rcv 127.0.0.1:8090  : txid=1      len=222  "{\"c\":200,\"p\":\"eOfUxaoT49mlOKWqKgJ8tTQ5MaSKb9e3sa5pnsgSXxI=\",\"s\":\"pQ7EVE7h8c3Z6wRHoVGsDiMr51gzwG8w3aQot62sp9oKX9ej8aEJ+bm901IugJMHs1xWXW0KQ+LqwnvIhgIXCw==\",\"v\":\"website\",\"a\":\"127.0.0.1:52516\",\"d\":\"127.0.0.1:8091\",\"u\":1}"
A demo website.
<img src="tumblr_nosz96g2GT1t0pdxwo1_1280.jpg" />
<img src="tumblr_nosz96g2GT1t0pdxwo1_1280.jpg?1" />
<img src="tumblr_nosz96g2GT1t0pdxwo1_1280.jpg?2" />
<img src="tumblr_nosz96g2GT1t0pdxwo1_1280.jpg?3" />
=== RUN   Test1/3
go run main.go serve -l 0.0.0.0:8080
2017/09/10 23:51:10 Listening... [::]:8080
go run main.go website -r 127.0.0.1:8080 -l 0.0.0.0:8081 --local 0.0.0.0:8082 --pvk 504bc61393e5d7ea991dbfad4d5bb98093562d472fa22d425a35bcd46341d8f678e7d4c5aa13e3d9a538a5aa2a027cb5343931a48a6fd7b7b1ae699ec8125f12 --dir demows
2017/09/10 23:51:11 [::]:8081        snd 127.0.0.1:8080  : txid=1      len=175  "{\"q\":\"reg\",\"p\":\"eOfUxaoT49mlOKWqKgJ8tTQ5MaSKb9e3sa5pnsgSXxI=\",\"s\":\"pQ7EVE7h8c3Z6wRHoVGsDiMr51gzwG8w3aQot62sp9oKX9ej8aEJ+bm901IugJMHs1xWXW0KQ+LqwnvIhgIXCw==\",\"v\":\"website\"}"
2017/09/10 23:51:11 [::]:8080        rcv 127.0.0.1:8081  : txid=1      len=175  "{\"q\":\"reg\",\"p\":\"eOfUxaoT49mlOKWqKgJ8tTQ5MaSKb9e3sa5pnsgSXxI=\",\"s\":\"pQ7EVE7h8c3Z6wRHoVGsDiMr51gzwG8w3aQot62sp9oKX9ej8aEJ+bm901IugJMHs1xWXW0KQ+LqwnvIhgIXCw==\",\"v\":\"website\"}"
2017/09/10 23:51:11 [::]:8080        snd 127.0.0.1:8081  : txid=1      len=34   "{\"c\":200,\"a\":\"127.0.0.1:8081\"}"
2017/09/10 23:51:11 [::]:8081        rcv 127.0.0.1:8080  : txid=1      len=34   "{\"c\":200,\"a\":\"127.0.0.1:8081\"}"
2017/09/10 23:51:11 Public Website listening on  [::]:8081
2017/09/10 23:51:11 Local Website listening on  0.0.0.0:8082
pvk= 504bc61393e5d7ea991dbfad4d5bb98093562d472fa22d425a35bcd46341d8f678e7d4c5aa13e3d9a538a5aa2a027cb5343931a48a6fd7b7b1ae699ec8125f12
pbk= 78e7d4c5aa13e3d9a538a5aa2a027cb5343931a48a6fd7b7b1ae699ec8125f12
sig= a50ec4544ee1f1cdd9eb0447a151ac0e232be75833c06f30dda428b7adaca7da0a5fd7a3f1a109f9b9bdd3522e809307b35c565d6d0a43e2eac27bc88602170b
2017/09/10 23:51:11 [::]:8081        snd 127.0.0.1:8080  : txid=2      len=33   "{\"q\":\"testport\",\"t\":\"random\"}"
2017/09/10 23:51:11 [::]:8080        rcv 127.0.0.1:8081  : txid=2      len=33   "{\"q\":\"testport\",\"t\":\"random\"}"
2017/09/10 23:51:11 [::]:58283       snd 127.0.0.1:8081  : txid=1      len=54   "{\"q\":\"porttest\",\"a\":\"127.0.0.1:8081\",\"t\":\"random\"}"
2017/09/10 23:51:11 [::]:8081        rcv 127.0.0.1:58283 : txid=1      len=54   "{\"q\":\"porttest\",\"a\":\"127.0.0.1:8081\",\"t\":\"random\"}"
2017/09/10 23:51:11 porttest q:  127.0.0.1:58283 random
2017/09/10 23:51:11 porttest success
2017/09/10 23:51:11 [::]:8081        snd 127.0.0.1:58283 : txid=1      len=48   "{\"c\":200,\"a\":\"127.0.0.1:58283\",\"d\":\"random\"}"
2017/09/10 23:51:11 [::]:58283       rcv 127.0.0.1:8081  : txid=1      len=48   "{\"c\":200,\"a\":\"127.0.0.1:58283\",\"d\":\"random\"}"
2017/09/10 23:51:11 { 200 [] []  127.0.0.1:58283 random  0 [] 0 0}
go run main.go browser -r 127.0.0.1:8080 -l 0.0.0.0:8083 --ws 0.0.0.0:8085 --proxy 0.0.0.0:8084 --headless
2017/09/10 23:51:12 me.com server listening on 0.0.0.0:8085
2017/09/10 23:51:12 browser proxy listening on 0.0.0.0:8084
2017/09/10 23:51:12 [::]:8083        snd 127.0.0.1:8080  : txid=1      len=33   "{\"q\":\"testport\",\"t\":\"random\"}"
2017/09/10 23:51:12 [::]:8080        rcv 127.0.0.1:8083  : txid=1      len=33   "{\"q\":\"testport\",\"t\":\"random\"}"
2017/09/10 23:51:12 [::]:58283       snd 127.0.0.1:8083  : txid=2      len=54   "{\"q\":\"porttest\",\"a\":\"127.0.0.1:8083\",\"t\":\"random\"}"
2017/09/10 23:51:12 [::]:8083        rcv 127.0.0.1:58283 : txid=2      len=54   "{\"q\":\"porttest\",\"a\":\"127.0.0.1:8083\",\"t\":\"random\"}"
2017/09/10 23:51:12 porttest q:  127.0.0.1:58283 random
2017/09/10 23:51:12 porttest success
2017/09/10 23:51:12 [::]:8083        snd 127.0.0.1:58283 : txid=2      len=48   "{\"c\":200,\"a\":\"127.0.0.1:58283\",\"d\":\"random\"}"
2017/09/10 23:51:12 port 8083  is open: true
2017/09/10 23:51:12 [::]:58283       rcv 127.0.0.1:8083  : txid=2      len=48   "{\"c\":200,\"a\":\"127.0.0.1:58283\",\"d\":\"random\"}"
2017/09/10 23:51:12 { 200 [] []  127.0.0.1:58283 random  0 [] 0 0}
HTTP GET  http://127.0.0.1:8085/index.html
<html>
<head>
  <title>me.com: myself on the internet</title>
  <script src="jquery.min.js"></script>
</head>
<body>

  <div style="position:absolute;top:16px;right:16px;border:solid 1px gray;padding:4px;">
    Port status: <input type="text" id="port" style="width:80px" />
    <input type="button" value="Change" id="change" />
    <span id="status">unknown</span>
  </div>

  <div align="center" style="margin-top:60px">
    Find: <input type="text" id="search" style="width:60%" /><input type="button" value="find" id="find" />
  </div>

  <div align="center" style="margin-top:60px">
    List: <input type="button" id="list" value="refresh" />
    <div id="peers"></div>
  </div>
  <script>
    $("#list").click(function(){
      $("#peers").html("")
      $.post("/list/0/30", function(res) {
        res.forEach(function(item){
          $("#peers").append($("<div><b>"+item.Name+"</b>: "+item.Pbk+"</div>"))
        })
      })
    })
  </script>
  <script>
    $("#list").attr("disabled", "disabled")
    $("#change").attr("disabled", "disabled")
    $("#port").val("searching").attr("disabled", "disabled")
    $("#change").click(function(){
      var port = $("#port").val()
      $.post("/change_port/"+port, function(res) {
        $("#port").val(res.Port).attr("disabled", false)
        if(res.Status==0){ $("#status").html("unknown") }
        if(res.Status==1){ $("#status").html("open") }
        if(res.Status==2){ $("#status").html("close") }
        $("#change").attr("disabled", false)
      })
    })
    $.post("/test_port", function(res) {
      $("#list").attr("disabled", false)
      $("#port").val(res.Port).attr("disabled", false)
      if(res.Status==0){ $("#status").html("unknown") }
      if(res.Status==1){ $("#status").html("open") }
      if(res.Status==2){ $("#status").html("close") }
      $("#change").attr("disabled", false)
    })
  </script>
  <script>
    $("#search").css("border", "solid 1px gray")
    $("#find").click(function(){
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
HTTP GET  http://78e7d4c5aa13e3d9a538a5aa2a027cb5343931a48a6fd7b7b1ae699ec8125f12.me.com/index.html
2017/09/10 23:51:13 [001] INFO: Got request /index.html 78e7d4c5aa13e3d9a538a5aa2a027cb5343931a48a6fd7b7b1ae699ec8125f12.me.com GET http://78e7d4c5aa13e3d9a538a5aa2a027cb5343931a48a6fd7b7b1ae699ec8125f12.me.com/index.html
2017/09/10 23:51:13 [::]:8083        snd 127.0.0.1:8080  : txid=2      len=81   "{\"q\":\"find\",\"p\":\"eOfUxaoT49mlOKWqKgJ8tTQ5MaSKb9e3sa5pnsgSXxI=\",\"v\":\"website\"}"
2017/09/10 23:51:13 [::]:8080        rcv 127.0.0.1:8083  : txid=2      len=81   "{\"q\":\"find\",\"p\":\"eOfUxaoT49mlOKWqKgJ8tTQ5MaSKb9e3sa5pnsgSXxI=\",\"v\":\"website\"}"
2017/09/10 23:51:13 [::]:8080        snd 127.0.0.1:8083  : txid=2      len=221  "{\"c\":200,\"p\":\"eOfUxaoT49mlOKWqKgJ8tTQ5MaSKb9e3sa5pnsgSXxI=\",\"s\":\"pQ7EVE7h8c3Z6wRHoVGsDiMr51gzwG8w3aQot62sp9oKX9ej8aEJ+bm901IugJMHs1xWXW0KQ+LqwnvIhgIXCw==\",\"v\":\"website\",\"a\":\"127.0.0.1:8083\",\"d\":\"127.0.0.1:8081\",\"u\":1}"
2017/09/10 23:51:13 [::]:8083        rcv 127.0.0.1:8080  : txid=2      len=221  "{\"c\":200,\"p\":\"eOfUxaoT49mlOKWqKgJ8tTQ5MaSKb9e3sa5pnsgSXxI=\",\"s\":\"pQ7EVE7h8c3Z6wRHoVGsDiMr51gzwG8w3aQot62sp9oKX9ej8aEJ+bm901IugJMHs1xWXW0KQ+LqwnvIhgIXCw==\",\"v\":\"website\",\"a\":\"127.0.0.1:8083\",\"d\":\"127.0.0.1:8081\",\"u\":1}"
2017/09/10 23:51:13 [001] INFO: Copying response to client 200 OK [200]
2017/09/10 23:51:13 [001] INFO: Copied 222 bytes to client error=<nil>
A demo website.
<img src="tumblr_nosz96g2GT1t0pdxwo1_1280.jpg" />
<img src="tumblr_nosz96g2GT1t0pdxwo1_1280.jpg?1" />
<img src="tumblr_nosz96g2GT1t0pdxwo1_1280.jpg?2" />
<img src="tumblr_nosz96g2GT1t0pdxwo1_1280.jpg?3" />
--- PASS: Test1 (7.18s)
    --- PASS: Test1/1 (1.02s)
    --- PASS: Test1/2 (2.06s)
    --- PASS: Test1/3 (3.02s)
PASS
ok  	github.com/mh-cbon/rendez-vous	7.185s
```

# todos

#### $ grep --include='*go' -r todo -B 1 -A 1 -n
```sh
client/token.go-11-	tokens map[string]string
client/token.go:12:} //todos: add a limit on maximum number of tokens.
client/token.go-13-
--
server/handler.go-74-		var res *model.Message
server/handler.go:75:		//todo: rendez-vous server should implement a write token
server/handler.go-76-		if len(m.Pbk) == 0 {
--
server/handler.go-106-		var res *model.Message
server/handler.go:107:		//todo: unregister should accept/verify a pbk/sig/value with a special value to identify the query issuer.
server/handler.go-108-		if len(m.Pbk) == 0 {
--
server/handler.go-184-	return QueryHandler(model.DoKnock, func(remote net.Addr, m model.Message, reply MessageResponseWriter) error {
server/handler.go:185:		//todo: protect from undesired usage.
server/handler.go-186-		addrToKnock := m.Data
--
server/handler.go-225-	return QueryHandler(model.TestPort, func(remote net.Addr, m model.Message, reply MessageResponseWriter) error {
server/handler.go:226:		//todo: protect from undesired usage.
server/handler.go-227-		go func(remote string, token string) {
--
main.go-30-
main.go:31://todo: rendez-vous server should implement a write token concept to register
main.go:32://todo: rendez-vous server unregister should accept/verify a pbk/sig/value with a special value to identify the query issuer.
main.go-33-
--
main.go-302-	handler := http.FileServer(http.Dir(opts.Dir))
main.go:303:	public := utils.ServeHTTPFromListener(ln, httpServer(handler, "")) //todo: replace with a transparent proxy, so the website can live into another process
main.go-304-	local := httpServer(handler, opts.Local)
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
--
store/pendingops.go-26-	if token == "" {
store/pendingops.go:27:		token = "random" //todo: random token
store/pendingops.go-28-	}
--
node/peer.go-140-	}
node/peer.go:141:	// todo: implement challenge of the remote here. need to receieve the extra pbk key associated with the dial request
node/peer.go-142-	_, err = conn.Write(id)
--
node/peer.go-174-	if err == nil {
node/peer.go:175:		// todo: implement challenge of the remote here. need to receieve the extra pbk key of the node receiving incoming conn
node/peer.go-176-		b := make([]byte, len(s.id))
```
