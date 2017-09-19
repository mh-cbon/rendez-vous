
var vdom = require('./vdom');
var hyperx = require('hyperx')
var hx = hyperx(vdom.h);

var emit = require('./emit')
var ui = require('./ui')

module.exports = Toggle;

function Toggle(update, opts){

  var state = this;
  Object.assign(state, {
    value: true,
    type:"checkbox",
    label:"",
    name:"",
    checked:false,
    disabled:false,
  }, opts)

  this.render = function(){
    return hx`
    <div class="ui toggle ${state.type} ${ui.cl(state, 'checked', 'disabled')}" onclick=${click}>
      <input class="hidden"
      type="${state.type}"
      name="${state.name}"
      value=${state.value}
      checked="${ui.cl(state, 'checked')}"
      onchange=${change} />
      <label>${state.label}</label>
    </div>`
  }

  emit.call(state);
  function click(ev){
    if (this.classList.contains("checked")) {
      this.classList.remove("checked")
      this.querySelector(".hidden").checked = false;
      state.checked = false;
    } else {
      this.classList.add("checked")
      this.querySelector(".hidden").checked = true;
      state.checked = true;
    }
    this.querySelector(".hidden").onchange();
    if(!this.classList.contains("disabled")) {
      state.emit("click", ev, this)
    }
  }
  function change(ev){
    if(!this.classList.contains("disabled")) {
      state.emit("change", ev, this)
    }
  }
}
