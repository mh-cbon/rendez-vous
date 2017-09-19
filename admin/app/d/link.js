
var vdom = require('./vdom');
var hyperx = require('hyperx')
var hx = hyperx(vdom.h);

var emit = require('./emit')
var ui = require('./ui')
var ev = require('./ev')

module.exports = Link;

function Link(update, opts){

  var state = this;
  Object.assign(state, {
    color:"",
    label:"",
    href:"",
  }, opts)

  this.render = function(){
    return hx`
    <a class="ui
      ${ui.cl(state, 'right', 'floated', 'labeled', 'icon')}
      button
      ${ui.print(state.color)}
      ${ui.cl(state, 'compact', 'small', 'loading', 'disabled')}
      "
      href="${ui.print(state.href)}"
      onclick=${ev.cancel('disabled', click)}
    >
      ${ui.icon(state.icon)}
      ${ui.print(state.label)}
    </a>`
  }

  emit.call(state);
  function click(ev){
    return state.emit("click", ev, this);
  }
}
