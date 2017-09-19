const EventEmitter = require('events');

function quote(str) {
    return (str+'').replace(/[.?*+^$[\]\\(){}|-]/g, "\\$&");
};

function router(){
  var e = new EventEmitter();
  this.on = function(path, handler){
    var fix = "";
    if (path[path.length-1]==="$") {
      fix = "$"
      path = path.substring(0,path.length-1);
    }
    var n = quote(path).replace(/(:[a-z]+)/i, "([^/]+)")
    n = new RegExp("^"+(n)+""+fix);
    e.on("change", function(url){
      var match = url.match(n)
      if (match) {
        var params = {}
        var k = path.match(/:([^/]+)/);
        if (k){
          for(var i=1;i<k.length;i++){
            var p = k[i]
            var v = match[i]
            params[p] = v
          }
        }
        handler(params)
      }
    })
  }
  this.resolve = function(path){
    path = path===""?"/":path;
    e.emit("change", path)
  }
  this.goto = function(path){
    window.location.hash = "#"+path;
  }
}

module.exports = router;
