
var vdom = require('./vdom');
var hyperx = require('hyperx')
var hx = hyperx(vdom.h);
var main = require('main-loop');

function loop(el, ready){
  var state = {};
  var render = function(){return hx`<i></i>`}
  var loop = main(state, function(){return render(state)}, vdom);
  var update = function(){
    return loop.update(state);
  }
  state = ready(update, loop)
  render = state.render;
  update();
  el.appendChild(loop.target)
  return state;
}

module.exports = loop
