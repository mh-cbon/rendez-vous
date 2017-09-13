
var $ = window.jQuery || window.$ || require('jquery');
var vdom = require('virtual-dom')
var hyperx = require('hyperx')
var hx = hyperx(vdom.h)
var main = require('main-loop')

var appsMainView = require('./apps/main');

(new mainView()).install(document.body.querySelector(".main"))

function mainView () {

  var state = {
    view: "apps",
  };

  function render (state) {
    return hx`
<div class="app">
  <div class="pusher">
    <div class="app-view apps-container ${state.view=='apps'?'active':'hide'}"></div>
    <div class="app-view web-container ${state.view=='web'?'active':'hide'}"></div>
  </div>
  <div class="ui vertical compact menu visible fixmenu">
    <a class="item ${state.view=='apps'?'active':''}" onclick=${onclick}
      name=apps title="Manage apps" href="#/apps/">
      Apps <i class="hashtag icon"></i>
    </a>
    <a class="item ${state.view=='web'?'active':''}" onclick=${onclick}
      name=web title="Browse the web" href="#/web/">
      Web <i class="cloud icon"></i>
    </a>
  </div>
</div>
`
  }

  var that = this;
  var loop = main(state, render, vdom);
  function onclick (ev) {
    state.view = this.name;
    loop.update(state);
  }
  function hashChanged(){
    var url = location.hash.substring(1).split("/").splice(1)
    if (url.length>0 && url[0]=="apps" || url[0]=="web") {
      if (state.view!=url[0]) {
        state.view = url[0];
        loop.update(state);
      }
    }
  }

  var appsMain = new appsMainView();

  this.install = function(to){
    appsMain.install(loop.target.querySelector(".apps-container"));
    window.addEventListener("hashchange", hashChanged, false);
    hashChanged();
    to.appendChild(loop.target);
    loop.update(state);
  }
  this.uninstall = function(from){
    window.removeEventListener("hashchange", hashChanged, false);
    appsMain.uninstall(loop.target.querySelector(".apps-container"));
    from.removeChild(loop.target)
  }
}

// return ;
//
// var fixmenu = document.querySelector('.fixmenu');
// var pusher = document.querySelector('.pusher');
//
// var webApp = {
//   Name:"web",
//   URL:"",
//   Icon:"cloud",
//   MenuText:"Web",
//   Enabled:false,
//   IsSystem:true,
//   element: document.createElement("div"),
// };
// var adminAppsList = {
//   Name:"apps",
//   URL:"",
//   Icon:"hashtag",
//   MenuText:"Apps",
//   Enabled:true,
//   IsSystem:true,
//   element: document.createElement("div")
// };
// var loadedApps = [webApp, adminAppsList]
//
// var adminApps = new appsMainView();
// // var webApps = new apps(webApp.element);
// adminApps.install();
// // webApps.install();
//
// var views = new views(pusher);
// var menus = new menus(fixmenu, views.enable);
//
// views.concat(loadedApps);
// menus.concat(loadedApps);
// adminApps.concat(loadedApps);
//
// views.install();
// menus.install();
//
// function views (target) {
//
//   var state = { views: [] };
//   var loop = main(state, render, vdom);
//
//   function render (state) {
//     return hx`<div>
//       ${state.views.map(function (w, i) {
//         return hx`<div class="app-view ${w.Name}-container ${w.Enabled ? 'active' : 'hide'}">
//         </div>`
//       })}
//     </div>`
//   }
//
//
//   this.install = function(){
//     target.appendChild(loop.target);
//     loop.update(state);
//     setTimeout(function(){
//       state.views.map(function(w){
//         target.querySelector("."+w.Name+"-container").appendChild(w.element)
//       })
//     },0)
//   }
//   this.uninstall = function(){
//     target.removeChild(loop.target)
//     target.querySelector("."+w.Name+"-container").removeChild(w.element)
//   }
//
//   this.concat = function(apps){
//     state.views = state.views.concat(apps)
//     loop.update(state)
//   }
//   this.enable = function(name){
//     state.views.map(function(w){
//       w.Enabled = w.Name==name
//     })
//     loop.update(state)
//   }
//   this.add = function(app){
//     state.views.push(app)
//     loop.update(state)
//   }
//   this.remove = function(name){
//     state.views = state.views.filter(function(w){
//       return w.Name==name
//     })
//     loop.update(state)
//   }
// }
//
// function menus (target, menuClick) {
//
//   var that = this;
//   var state = { views: [] };
//   var loop = main(state, render, vdom);
//
//   function render (state) {
//     return hx`<div>
//       ${state.views.map(function (w, i) {
//         return hx`<a class="item ${w.Name}-menu ${w.Enabled ? 'active' :''}" onclick=${onclick} name=${w.Name}>
//          ${w.MenuText} <i class="${w.Icon} icon"></i>
//         </a>`
//       })}
//     </div>`
//   }
//
//   function onclick (ev) {
//     if(menuClick!==null){
//       menuClick(this.name);
//       that.enable(this.name);
//     }
//   }
//
//   this.install = function(){
//     target.appendChild(loop.target)
//     loop.update(state)
//   }
//   this.uninstall = function(){
//     target.appendChild(loop.target)
//   }
//
//   this.concat = function(apps){
//     state.views = state.views.concat(apps)
//     loop.update(state)
//   }
//   this.add = function(app){
//     state.views.push(app)
//     loop.update(state)
//   }
//   this.remove = function(name){
//     state.views = state.views.filter(function(w){
//       return w.Name==name
//     })
//     loop.update(state)
//   }
//   this.enable = function(name){
//     state.views.map(function(w){
//       w.Enabled = w.Name==name
//     })
//     loop.update(state)
//   }
// }

//
// function post(url, data, cb) {
//   return $.ajax({
//         url : url,
//         type: "POST",
//         data: JSON.stringify(data),
//         contentType: "application/json; charset=utf-8",
//         dataType   : "json",
//         success    : cb
//     });
// }
