
var $ = window.jQuery || window.$ || require('jquery');

function post(url, data, cb) {
  return $.ajax({
        url : url,
        type: "POST",
        data: JSON.stringify(data),
        contentType: "application/json; charset=utf-8",
        dataType   : "json",
        success    : cb
    });
}

var handler = function(){
  var that = this;
  that.hashChanged = function(ev){
    return that.emit("hashchange", ev)
  }
};
require('util').inherits(handler, require('events').EventEmitter);
handler = new handler();
window.addEventListener("hashchange", handler.hashChanged, false);

function hashChanged(cb) {
  var hashchanged = function(ev){
    cb(location.hash.slice(1))
  };
  var begin = function(){
    handler.on("hashchange", hashchanged, false);
  };
  var close = function(){
    handler.removeListener("hashchange", hashchanged, false);
  };
  return {begin: begin, close: close}
}

function matchURL(pattern, cb) {
  return function(url){
    var r = url.match(pattern)
    if (r!==null) {
      cb.apply(null, r);
    }
  }
}
module.exports = {
  post: post,
  triggerHashChange: function(){handler.hashChanged(null)},
  hashChanged: hashChanged,
  matchURL: matchURL,
}
