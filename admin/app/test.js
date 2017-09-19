
var vdom = require('./vdom');
var hyperx = require('hyperx')
var hx = hyperx(vdom.h);


var emit = require('./d/emit')
var ui = require('./d/ui')
var ev = require('./d/ev')
var vs = require('./d/vs')

var Button = require('./d/button')
var Link = require('./d/link')
var Tab = require('./d/tab')
var Select = require('./d/select')

var AppsItem = require('./apps.js')
var WebView = require('./web2/main.js')

var api = require('./api')
api = new api().wrapped();

function smallBtn(opts){
  return Object.assign({labeled: true, small:true}, opts)
}
function smallRightBtn(opts){
  return Object.assign(smallBtn({right:true}), opts)
}
function smallRightFloatedBtn(opts){
  return Object.assign(smallBtn({right:true,floated: true}), opts)
}

var loop = require('./d/loop')

var router = require("./router/router.js")
router = new router()
window.addEventListener("hashchange", function(){router.resolve(location.hash.slice(1));}, false);

loop(document.body.querySelector(".main"), function(update){
  var app = new App(update, router);
  app.add("apps", "Apps", new AppsMain(update, router), {icon:{icon:"hashtag", color:"red"}, href:"#/apps"})
  app.add("web", "Web", new WebView(update, router), {icon:{icon:"cloud", color:"blue"}, href:"#/web"})
  return app;
})
router.resolve(location.hash.slice(1));

function format(f){
  return function() {
    var b = Array.prototype.slice.call(arguments);
    var ret = f
    b.map(function(v){ ret = ret.replace(/%[a-z]/, v); })
    return ret;
  }
};

function AppsList(update, router){

  var state = this;
  state.apps = [];
  state.refresh = new Button(update,
    smallRightFloatedBtn({icon: "refresh", color:"green", label:"Refresh"})
  );
  state.newbt = new Link(update,
    smallRightFloatedBtn({icon: "plus", color:"blue", label:"New",href:"#/apps/new"})
  );
  state.loadmore = new Button(update,
    smallRightBtn({icon: "plus", color:"blue", label:"Load more"})
  );

  emit.call(state);

  var start = 0;
  var limit = 10;
  function setApps(apps) {
    if (apps) {
      state.apps = apps
      start = apps.length;
      update()
    }
  }
  function appendApps(apps) {
    if (apps) {
      state.apps = state.apps.concat(apps);
      start += apps.length;
      update()
    }
  }
  state.loadmore.on('click',function(e){
    var p = api.list(start,limit);
    p.then(appendApps)
    vs(state.loadmore).load(p, update)
      .begin(vs.beginBlue).then(vs.goBlue).catch(vs.goRed)
    vs(state).load(p, update)
      .begin(vs.unfail, vs.unfailure).catch(vs.fail, vs.failure())
    return false
  })
  state.refresh.on('click',function(e){
    var p = api.list(0,start+limit)
    p.then(setApps)

    vs(state.refresh).load(p, update)
      .begin(vs.beginGreen).then(vs.goGreen).catch(vs.goRed)
    vs(state).load(p, update)
      .begin(vs.unfail, vs.unfailure).catch(vs.fail, vs.failure())
  })

  state.newbt.on('click', ev.emit("click-new", state))

  function deleteClick(s){
    return function deleteClick(e){
      var p = api.deleteByID(s.ID)
      p.then(ev.emit("click", state.refresh, e, this))

      vs(state).load(p, update)
        .begin(vs.unfail, vs.unfailure)
        .catch(vs.fail, vs.failure(format("Failed to delete app("+(s.Name || s.URL)+"): %s")))

      return false;
    }
  }
  router.on("/apps$", function(){
    state.refresh.emit("click")
  })

  function dismiss() {
    state.failure = "";
    state.failed = false;
    update()
  }

  this.render = function() {
    return hx`
    <div class="apps-list">

      <div>
        ${state.refresh.render()}
        ${state.newbt.render()}
      </div>

      <h2 class="ui medium header">
        Apps list
      </h2>

      <div class="ui red floating icon message ${ui.myb('hidden', !state.failed)}">
        <i class="warning icon"></i>
        <i class="close icon" onclick=${dismiss}></i>
        <div class="content">
          <p>${state.failure}</p>
        </div>
      </div>

      <table class="ui selectable celled striped blue table">
        <thead>
          <tr>
            <th> </th>
            <th>App</th>
            <th>Status</th>
            <th>Controls</th>
            <th>Update date</th>
            <th> </th>
          </tr>
        </thead>
        <tbody>
          ${state.apps.map(renderTr)}
        </tbody>
      </table>

      <div align="center">
        <br>${state.loadmore.render()}
      </div>
    </div>`
  }

  function renderTr(s){
    return hx`
    <tr>
      <td class="right aligned collapsing">
        <a class="${ui.myb('disabled', s.IsSystem)}" href="#/apps/edit/${s.ID}" onclick=${ev.cancel('disabled')}>
         ${ui.icon({icon:'write', color: ui.myb('grey', s.IsSystem)})}
        </a>
      </td>
      <td class="collapsing">${s.URL || s.Name}</td>
      <td>${s.Status}</td>
      <td class="right aligned collapsing">
        <a class="${ui.myb('disabled', s.IsSystem)}">
         ${ui.icon({icon:'play', color: ui.myb('grey', s.IsSystem)})}
        </a>
        <a class="${ui.myb('disabled', s.IsSystem)}">
         ${ui.icon({icon:'pause', color: ui.myb('grey', s.IsSystem)})}
        </a>
        <a class="${ui.myb('disabled', s.IsSystem)}">
         ${ui.icon({icon:'stop', color: ui.myb('grey', s.IsSystem)})}
        </a>
      </td>
      <td>${ui.print(s.UpdatedAt).substring(0,10)}</td>
      <td class="right aligned collapsing">
        <a class="${ui.myb('', s.IsSystem)}" onclick=${ev.cancel('disabled', deleteClick(s))} href="#/apps/delete/${s.ID}">
         ${ui.icon({icon:'trash', color: ui.myb('grey', s.IsSystem)})}
        </a>
        <div class="ui ${ui.myb('active', s.loading)} loader"></div>
      </td>
    </tr>`
  }

}

function AppsMain(update, router){

  var state = this;
  state.list = new AppsList(update, router)
  state.item = new AppsItem(update, router)
  state.current = "list"

  state.list.on("click-new", function(){
    state.current="item"
    update()
  })
  state.item.on("click-return", function(){
    state.current="list"
    update()
  })

  function is(view) {
    return state.current===view
  }

  function enterList(){
    if (state.current!="list") {
      state.current = "list";
      update()
    }
  }
  function enterItem(){
    if (state.current!="item") {
      state.current = "item";
      update()
    }
  }
  function enterEdit(){
    if (state.current!="item") {
      state.current = "item";
      update()
    }
  }
  router.on('/apps$', enterList)
  router.on('/apps/new$', enterItem)
  router.on('/apps/edit', enterItem)

  this.render = function(){
    return hx`
    <div class="ui container fluid">
      <div class="${ui.myb('invisible', !is("list"))}">${state.list.render()}</div>
      <div class="${ui.myb('invisible', !is("item"))}">${state.item.render()}</div>
    </div>`
  }
}


function App(update, router){

  var state = this;
  state.views = [];

  this.setActive = function(id){
    state.views.map(function(s){
      s.active=s.id==id;
      s.hidden=!s.active;
    })
    update()
  }
  this.add = function(id, title, render, opts) {
    render = render || function(){return hx`<div>missing render ${id}</div>`}
    state.views.push(Object.assign({
      id: id,
      title: title,
      href: "#",
      active: state.views.length==0,
      render: render.render || render
    }, opts || {}))
    return render
  }

  this.enable = function(params){
    if((!params || !params.view) && state.views.length){params.view=state.views[0].id}
    if (params.view && state.view!=params.view) {
      state.setActive(params.view)
    }
    update()
  }
  router.on('/:view$', state.enable)
  router.on('/$', state.enable)

  this.render = function(){
    return hx`
    <div class="app">
      <div>
        ${state.views.map(renderContent)}
      </div>
      <div class="ui vertical compact menu visible vertical-menu">
        ${state.views.map(renderTitle)}
      </div>
    </div>`
  }
  function renderTitle(s){
    return hx`
    <a class="item ${ui.cl(s,'active')} ${ui.print(s.color)}"
      data-tab="${ui.print(s.id)}"
      title="${ui.print(s.description)}"
      href="${ui.print(s.href)}"
      onclick=${menuClick}
    >
      ${s.title} ${ui.icon(s.icon)}
    </a>`
  }
  function renderContent(s){
    return hx`
    <div class="app-view ${ui.cl(s,'active', 'hidden')}">
      ${s.render()}
    </div>`
  }
  function menuClick(ev) {
    if (this.hasAttribute("data-tab")) {
      if (this.classList.contains("disabled")) {
        ev.stopImmediatePropagation()
        ev.stopPropagation()
        return false
      }
      state.setActive(
        this.getAttribute("data-tab")
      )
    }
    return true
  }
}
