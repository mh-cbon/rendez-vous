
function cancel(cl, v){
  var args = Array.prototype.slice.call(arguments).slice(1);
  return function(ev){
    var that = this;
    if (that.classList && that.classList.contains(cl)) {
      ev.stopPropagation()
      ev.stopImmediatePropagation()
      return false;
    }
    var ret = true
    args.forEach(function(k){
      if (k.apply(that, ev)===false) ret = false
    })
    return ret
  }
}

function emit(w, state) {
  return function(ev){
    var args = Array.prototype.slice.call(arguments).slice(1);
    state.emit && state.emit(w, [ev, this].concat(args))
  }
}

module.exports = {
  cancel: cancel,
  emit: emit,
}
