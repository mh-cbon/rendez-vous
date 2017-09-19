
var vdom = require('./vdom');
var hyperx = require('hyperx')
var hx = hyperx(vdom.h);

var emit = require('./emit')
var ui = require('./ui')
var ev = require('./ev')
var Button = require('./button')

module.exports = Text;

function Text(update, opts){

  var state = this;
  Object.assign(state, {
    value: "",
    type:"text",
    name:"",
    placeholder:"",
    disabled:false,
    action:null,
  }, opts)

  if (state.action) {
    state.action = new Button(update, state.action)
  }

  this.render = function(){
    return hx`
    <div class="ui field ${ui.cl(state,'action','disabled')} ${ui.myb('error', state.failed)} input">
      <input placeholder="${ui.print(state.placeholder)}"
         type="${ui.print(state.type)}"
         name="${ui.print(state.name)}"
         value="${state.value}"
         onchange=${ev.cancel('disabled', change)}
         onkeydown=${ev.cancel('disabled', keydown)}
         onkeyup=${ev.cancel('disabled', keyup)}
         />
       ${state.action ? state.action.render() : ''}
    </div>`
  }

  emit.call(state);
  function change(ev){
    state.value = this.value;
    return state.emit("change", ev, this)
  }
  function keydown(ev){
    return state.emit("keydown", ev, this)
  }
  function keyup(ev){
    state.value = this.value;
    return state.emit("keyup", ev, this)
  }
}
