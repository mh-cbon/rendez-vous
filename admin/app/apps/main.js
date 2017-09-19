
var $ = window.jQuery || window.$;
var vdom = require('../vdom')
var hyperx = require('hyperx')
var hx = hyperx(vdom.h);
var main = require('main-loop')
var webUtils = require('../web-utils');

var appsListView = require('./list');
var appsItemView = require('./item');

module.exports = AppsMain;

function AppsMain (router) {

  var state = {
    view: "",
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
  this.enable = function(params){
    if(!params || !params.view) {
      params = {view:"list"}
    } else if(params.view=="new" || params.view=="edit") {
      params.view = "item"
    }
    if (state.view!=params.view) {
      state.view = params.view;
    }
    loop.update(state);
  }
  router.on('/apps', that.enable)
  router.on('/apps/:view', that.enable)

  var appsList = new appsListView(router);
  var appsItem = new appsItemView(router);

  var loop = main(state, render, vdom);
  this.install = function(to){
    appsList.install(loop.target.querySelector(".apps-list-view"));
    appsItem.install(loop.target.querySelector(".apps-item-view"));
    to.appendChild(loop.target);
    loop.update(state);
  }
  this.uninstall = function(from){
    appsList.uninstall(loop.target.querySelector(".apps-list-view"));
    appsItem.uninstall(loop.target.querySelector(".apps-item-view"));
    from.removeChild(loop.target)
  }
}
