
var $ = window.jQuery || window.$;
var vdom = require('virtual-dom')
var hyperx = require('hyperx')
var hx = hyperx(vdom.h)
var main = require('main-loop')
var webUtils = require('../web-utils');

require('util').inherits(AppsMain, require('events').EventEmitter);
module.exports = AppsMain;

var appsListView = require('./list');
var appsItemView = require('./item');

function AppsMain () {

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
  this.enable = function(_, _, view){
    if(view===undefined) {
      view = "list"
    } else if(view=="new" || view=="edit") {
      view = "item"
    }
    if (state.view!=view) {
      state.view = view;
      loop.update(state);
    }
  }
  var watch = webUtils.hashChanged(
    webUtils.matchURL(/^\/apps(\/(new|edit)\/?)?/, that.enable)
  );

  var appsList = new appsListView();
  var appsItem = new appsItemView();

  var loop = main(state, render, vdom);
  this.install = function(to){
    appsList.install(loop.target.querySelector(".apps-list-view"));
    appsItem.install(loop.target.querySelector(".apps-item-view"));
    watch.begin();
    to.appendChild(loop.target);
    loop.update(state);
  }
  this.uninstall = function(from){
    appsList.uninstall(loop.target.querySelector(".apps-list-view"));
    appsItem.uninstall(loop.target.querySelector(".apps-item-view"));
    watch.close();
    from.removeChild(loop.target)
  }
}
