
const util = require('util');
const EventEmitter = require('events');

function emit(){
  var em = new EventEmitter();
  this.on = function(e, a){
    a.__removable__ = function(ev, that){
      var args = Array.prototype.slice.call(arguments).slice(2);
      return a.apply(that, [ev].concat(args))
    };
    return em.on(e, a.__removable__)
  }
  this.off = function(e, a){
    if (a && a.__removable__) return em.removeListener(e, a.__removable__)
    else return em.removeListener(e, a)
  }
  this.emit = function(e, a){
    var ret = true;
    var args = Array.prototype.slice.call(arguments).slice(1);
    var listeners = em.listeners(e)
    listeners.map(function(l){
      if(l.apply(null, args)===false){ret=false}
    })
    return ret;
  }
}

module.exports = emit
