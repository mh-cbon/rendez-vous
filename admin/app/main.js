
var $ = window.jQuery || window.$ || require('jquery');
var vdom = require('./vdom')
var hyperx = require('hyperx')
var hx = hyperx(vdom.h);
var main = require('main-loop')

var router = require("./router/router")
var r = new router()
window.addEventListener("hashchange", function(){r.resolve(location.hash.slice(1));}, false);

var webUtils = require('./web-utils');
var appsMainView = require('./apps/main');
var webMainView = require('./web/main');

(new mainView(r)).install(document.body.querySelector(".main"));
r.resolve(location.hash.slice(1));

function mainView (router) {

  var state = {
    view: "",
  };

  function render (state) {
    return hx`
<div class="app">
  <div class="pusher">
    <div class="app-view apps-container ${state.view=='apps'?'active':'hide'}"></div>
    <div class="app-view web-container ${state.view=='web'?'active':'hide'}"></div>
  </div>
  <div class="ui vertical compact menu visible fixmenu">
    <a class="item ${state.view=='apps'?'active':''}"
      name=apps title="Manage apps" href="#/apps/">
      Apps <i class="hashtag icon"></i>
    </a>
    <a class="item ${state.view=='web'?'active':''}"
      name=web title="Browse the web" href="#/web/">
      Web <i class="cloud icon"></i>
    </a>
  </div>
</div>
`
  }

  var that = this;
  var loop = main(state, render, vdom);
  this.enable = function(params){
    if (state.view!=params.view) {
      state.view = params.view;
    }
    loop.update(state);
  }
  router.on('/:view', that.enable)

  var appsMain = new appsMainView(router);
  var webMain = new webMainView(router);

  this.install = function(to){
    appsMain.install(loop.target.querySelector(".apps-container"));
    webMain.install(loop.target.querySelector(".web-container"));
    loop.update(state);
    to.appendChild(loop.target);
  }
  this.uninstall = function(from){
    appsMain.uninstall(loop.target.querySelector(".apps-container"));
    webMain.uninstall(loop.target.querySelector(".web-container"));
    from.removeChild(loop.target)
  }
}
