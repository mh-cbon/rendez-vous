
var $ = window.jQuery || window.$;
var vdom = require('virtual-dom')
var hyperx = require('hyperx')
var hx = hyperx(vdom.h)
var main = require('main-loop')
var webUtils = require('../web-utils');

require('util').inherits(WebMain, require('events').EventEmitter);
module.exports = WebMain;

function WebMain () {

  var state = {
    view: "",
    port: 53241,
    status: {
      text:"",
      icon:"",
      color:"",
    },
    peers: [],
    loadingPort: false,
  };
  function render (state) {
    return hx`
<div style="position:relative;">
  <h1 class="ui medium header">
    <br>Browse the web
  </h1>

  <div class="web-view ${state.view=='connect'?'visible':'invisible'}">

    <div class="ui form"  style="position:absolute;top:16px;right:16px;">
      <div class="field">
        <div class="ui mini left fluid action input">
          <input type="text" id="port" placeholder="Port to listen" value=${state.port} disabled="${state.loadingPort?'disabled':''}"/>
          <div class="ui button">
            <i class="icon settings"></i>Change port
          </div>
        </div>
      </div>
      <div class="field">
        <button class="ui compact right labeled icon button basic ${state.status.color}" style="white-space: nowrap;">
          Port status: ${state.status.text}
          <i class="${state.status.icon} icon"></i>
        </button>
      </div>
    </div>

    <div align="center" style="margin-top:160px;">
      <div class="field" style="width:60%;">
        <div class="ui  fluid action input">
          <input type="text" id="open" placeholder="#public key" />
          <div class="ui button" onclick=${openClick}>
            <i class="icon cloud"></i>Open
          </div>
        </div>
      </div>
    </div>

    <div align="center" style="margin-top:60px">
      <div class="ui button" onclick=${findClick}>
        <i class="icon comments"></i>Find peers
      </div>
      <div>
        ${state.peers.map(function(p){
          return hx`<div> <b>${p.Name}</b>: ${p.Pbk} </div>`
        })}
      </div>
    </div>

  </div>
</div>
`
  }

  var portStatuses = {
    0: {
      text:"unknown",
      icon:"help red",
      color:"grey",
    },
    1: {
      text:"open",
      icon:"green pause",
      color:"grey",
    },
    2: {
      text:"close",
      icon:"warning orange",
      color:"grey",
    }
  }

  function changeClick() {
    state.loadingPort = true;
    loop.update(state);
    var port = loop.target.querySelector("#port").value;
    webUtils.post("/change_port/"+port, null, function(res){
      loop.target.querySelector("#port").value = res.Port;
      if (portStatuses[res.Status]!==null) {
        state.status = portStatuses[res.Status];
      }
      loop.update(state);
    }).always(function(){
      state.status = portStatuses[0];
      state.loadingPort = false;
      loop.update(state);
    })
  }

  function openClick() {
    var open = loop.target.querySelector("#open")
    var v = open.value
    open.classList.remove("error")
    if (v == "") {
      open.classList.add("error")
    } else {
      window.open("http://"+v+".me.com/")
    }
  }

  function findClick() {
    webUtils.post("/list/0/30/", null, function(res){
      state.peers = res
      loop.update(state)
    })
  }

  function testPort(){
    state.loadingPort = true;
    loop.update(state);
    webUtils.post("/test_port/", null, function(res){
      loop.target.querySelector("#port").value = res.Port;
      if (portStatuses[res.Status]!==null) {
        state.status = portStatuses[res.Status];
      }
      loop.update(state);
    }).always(function(){
      state.status = portStatuses[0];
      state.loadingPort = false;
      loop.update(state);
    })
  }

  var that = this;
  var loop = main(state, render, vdom);
  this.enable = function(_, view){
    if(view==undefined) {view = "connect"}
    if (state.view!=view) {
      state.view = view;
      loop.update(state);
    }
  }
  var watch = webUtils.hashChanged(
    webUtils.matchURL(/^\/web(\/(connect)\/)?/, that.enable)
  );

  var loop = main(state, render, vdom);
  this.install = function(to){
    watch.begin();
    to.appendChild(loop.target);
    loop.update(state);
    testPort();
  }
  this.uninstall = function(from){
    watch.close();
    from.removeChild(loop.target)
  }
}
