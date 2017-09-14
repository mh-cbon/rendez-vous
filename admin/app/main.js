
var $ = window.jQuery || window.$ || require('jquery');
var vdom = require('virtual-dom')
var hyperx = require('hyperx')
var hx = hyperx(vdom.h)
var main = require('main-loop')

var webUtils = require('./web-utils');
var appsMainView = require('./apps/main');
var webMainView = require('./web/main');

(new mainView()).install(document.body.querySelector(".main"));

webUtils.triggerHashChange()

function mainView () {

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
  this.enable = function(_, view){
    if (state.view!=view) {
      state.view = view;
      loop.update(state);
    }
  }
  var watch = webUtils.hashChanged(
    webUtils.matchURL(/^\/(apps|web)\/?/, that.enable)
  );

  var appsMain = new appsMainView();
  var webMain = new webMainView();

  this.install = function(to){
    appsMain.install(loop.target.querySelector(".apps-container"));
    webMain.install(loop.target.querySelector(".web-container"));
    watch.begin();
    to.appendChild(loop.target);
    loop.update(state);
  }
  this.uninstall = function(from){
    watch.close();
    appsMain.uninstall(loop.target.querySelector(".apps-container"));
    webMain.uninstall(loop.target.querySelector(".web-container"));
    from.removeChild(loop.target)
  }
}
