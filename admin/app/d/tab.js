
var vdom = require('./vdom');
var hyperx = require('hyperx')
var hx = hyperx(vdom.h);

var emit = require('./emit')
var ui = require('./ui')
var ev = require('./ev')

module.exports = Tab;

function Tab(update){

  var state = this;
  state.tabs = [];

  this.add = function(id, opts){
    opts = Object.assign({
      active: state.tabs.length===0,
      id: id,
    }, opts)
    state.tabs.push(opts)
    return state.tabs[state.tabs.length-1]
  }
  this.setActive = function(id){
    state.tabs.map(function(s){
      s.active=s.id==id;
    })
    update()
  }
  this.DisableTab = function(id){
    state.tabs.filter(function(s){
      return s.id==id;
    }).map(function(t){
      t.disabled = true;
    })
    update()
  }
  this.EnableTab = function(id){
    state.tabs.filter(function(s){
      return s.id==id;
    }).map(function(t){
      t.disabled = false;
    })
    update()
  }
  this.GetTab = function(id){
    return state.tabs.find(function(s){
      return s.id==id;
    })
  }
  this.getActive = function(id){
    return state.tabs.find(function(s){
      return s.active;
    })
  }

  this.render = function(){
    return hx`<div>
      <div class="ui top attached tabular menu">
        ${state.tabs.map(renderTitle)}
      </div>
      ${state.tabs.map(renderContent)}
    </div>`
  }
  function renderTitle(s){
    return hx`
    <div data-tab="${s.id}" onclick=${ev.cancel('disabled', titleClick)}
      class="item
      ${ui.cl(s, 'active','disabled')}
      "
    >
      ${s.title}
    </div>`
  }
  function renderContent(s){
    return hx`
    <div class="ui bottom attached tab segment
      ${s.loading?"loading":""} ${s.active?"active":""}"
    >
      ${s.render()}
    </div>`
  }
  function titleClick(ev) {
    if (this.hasAttribute("data-tab")) {
      state.setActive(
        this.getAttribute("data-tab")
      )
    }
    return true
  }
}
