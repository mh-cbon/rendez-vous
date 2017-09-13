
var $ = window.jQuery || window.$;
var vdom = require('virtual-dom')
var hyperx = require('hyperx')
var hx = hyperx(vdom.h)
var main = require('main-loop')

require('util').inherits(AppsMain, require('events').EventEmitter);
module.exports = AppsMain;

var appsListView = require('./list');
var appsItemView = require('./item');

function AppsMain () {

  var state = {
    view: "list",
  };
  function render (state) {
    return hx`
<div style="position:relative;">
  <h1 class="ui medium header">
    <br>Manage apps
  </h1>

  <div class="apps-list-view ${state.view=='list'?'visible':'invisible'}"></div>
  <div class="apps-item-view ${state.view=='item'?'visible':'invisible'}"></div>
</div>
`
  }

  var that = this;
  var loop = main(state, render, vdom);
  function hashChanged(){
    var url = location.hash.substring(1).split("/").splice(1)
    if (url.length>0 && url[0]=="apps") {
      var view = "list"
      if (url.length>1 && url[1]=="new") {
        view = "item"
      } else if (url.length>2 && url[1]=="edit" && url[2]!="") {
        view = "item"
      }
      if (state.view!=view) {
        state.view = view
        loop.update(state);
      }
    }
  }

  var appsList = new appsListView();
  var appsItem = new appsItemView();

  // appsList.on("click-new", function(){
  //   state.view = "item"
  //   loop.update(state);
  // });
  //
  // appsList.on("click-edit", function(){
  //   state.view = "item"
  //   loop.update(state)
  // });
  //
  // appsItem.on("click-return", function(){
  //   state.view = "list"
  //   loop.update(state)
  // });

  var loop = main(state, render, vdom);
  this.install = function(to){
    appsList.install(loop.target.querySelector(".apps-list-view"));
    appsItem.install(loop.target.querySelector(".apps-item-view"));
    window.addEventListener("hashchange", hashChanged, false);
    hashChanged()
    to.appendChild(loop.target);
    loop.update(state);
  }
  this.uninstall = function(from){
    appsList.uninstall(loop.target.querySelector(".apps-list-view"));
    appsItem.uninstall(loop.target.querySelector(".apps-item-view"));
    window.removeEventListener("hashchange", hashChanged, false);
    from.removeChild(loop.target)
  }
}
