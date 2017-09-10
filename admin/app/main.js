
var $ = window.jQuery || window.$;
  var vdom = require('virtual-dom')
  var hyperx = require('hyperx')
  var hx = hyperx(vdom.h)

  var main = require('main-loop')
  var loop = main({ times: 0 }, render, vdom)
  document.querySelector('.pusher').appendChild(loop.target)
  console.log(loop.target)

  function render (state) {
    return hx`<div>
      <h1>clicked ${state.times} times</h1>
      <button onclick=${onclick}>click me!</button>
    </div>`
    function onclick () {
      loop.update({ times: state.times + 1 })
    }
  }

  $('.sidebar')
    // .sidebar('setting', 'transition', 'overlay')
    .sidebar('show')
  ;

  console.log(Date.now())
