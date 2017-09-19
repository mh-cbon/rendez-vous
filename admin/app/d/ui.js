
var vdom = require('./vdom');
var hyperx = require('hyperx')
var hx = hyperx(vdom.h);

function print(v){
  if(v) return v
  return ""
}
function cl(s, v){
  var ret = "";
  var args = Array.prototype.slice.call(arguments).slice(1);
  args.forEach(function(k){
    if((k in s) && !!s[k]) {ret+=" "+k}
  })
  return ret
}
function myb(p, v){
  if (v) {
    return p
  }
  return ""
}

function icon(v){
  if (v && v.icon ) {
    return hx`<i class="${print(v.icon)} ${print(v.color)} icon"></i>`
  } else if (v) {
    return hx`<i class="${print(v)} icon"></i>`
  }
  return ""
}
function field() {
  var args = Array.prototype.slice.call(arguments);
  return hx`<div class="field">${args.map(function(a){return a && (a.render ? a.render() : a())})}</div>`
}
function two() {
  var args = Array.prototype.slice.call(arguments);
  return hx`<div class="field"><div class="two fields">${args.map(function(a){return a && (a.render ? a.render() : a())})}</div></div>`
}

module.exports = {
  print: print,
  cl: cl,
  myb: myb,
  icon: icon,
  field: field,
  two: two,
}
