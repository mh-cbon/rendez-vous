
var vdom = require('./vdom');
var hyperx = require('hyperx')
var hx = hyperx(vdom.h);

var emit = require('./emit')
var ui = require('./ui')
var ev = require('./ev')

module.exports = Select;

function Select(update, name, options){

  var state = this;
  state.name = name || "";
  state.options = options || [];

  this.add = function(value, text, selected) {
    var option = {
      value: value,
      text: text,
      selected: !!selected,
    }
    state.options.push(option)
    return option;
  }

  this.selectValue = function(v){
    state.options.forEach(function(o){
      o.selected = o.value===v
    })
  }
  this.selectAtIndex = function(index){
    state.options.forEach(function(o, i){
      state.options[index].selected = index === i
    })
  }
  this.selectedIndex = function(){
    var ret = -1;
    state.options.forEach(function(o, i){
      if (o.selected) ret = i
    })
    return ret
  }
  this.getValue = function(){
    var ret = null;
    var index = this.selectedIndex()
    if (index==-1 && state.options.length>0){
      ret = state.options[0].value
    } else {
      ret = state.options[index].value
    }
    return ret
  }

  this.render = function(){
    return hx`
    <select class="ui fluid search dropdown
      ${ui.cl(state,'disabled')} ${ui.myb('error', state.failed)}"
      onchange=${ev.cancel('disabled', change)}
      name=${ui.print(state.name)}
      >
      ${state.options.map(renderOption)}
    </select>`
  }
  function renderOption(s){
    return hx`
    <option value="${s.value}" selected="${ui.cl(s,'selected')}">
      ${s.text}
    </option>`
  }

  emit.call(state);
  function change(ev){
    var opt = state.options[this.selectedIndex];
    if (opt) {
      state.options.map(function(o, i){ o.selected = false; })
      opt.selected = true;
      state.value = opt.value;
    }
    state.emit("change", ev, this)
    return true
  }
}
