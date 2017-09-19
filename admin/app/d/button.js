
var vdom = require('./vdom');
var hyperx = require('hyperx')
var hx = hyperx(vdom.h);

var emit = require('./emit')
var ui = require('./ui')

module.exports = Button;

function Button(update, opts){

  var state = this;
  Object.assign(state, {
    color:"",
    label:"",
  }, opts)

  this.render = function(){
    return hx`
    <button class="ui
      ${ui.cl(state, 'right', 'floated', 'labeled', 'icon')}
      button
      ${ui.print(state.color)}
      ${ui.cl(state, 'compact', 'small', 'loading', 'disabled')}
      "
      onclick=${click}
    >
      ${ui.icon(state.icon)}
      ${ui.print(state.label)}
    </button>`
  }

  emit.call(state);
  function click(ev){
    if(!this.classList.contains("disabled")) {
      return state.emit("click", ev, this)
    }
  }
}
