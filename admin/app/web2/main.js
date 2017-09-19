
var vdom = require('../vdom')
var hyperx = require('hyperx')
var hx = hyperx(vdom.h);

var emit = require('../d/emit')
var ui = require('../d/ui')
var ev = require('../d/ev')
var form = require('../d/form')
var vs = require('../d/vs')

var api = require('../node_api')
api = new api().wrapped();

module.exports = WebMain;

function WebMain (update, router, opts) {

  var state = this;
  Object.assign(state, {
    port: 53241,

    change:{},
    test:{icon:"ban", text:"close"},
    find:{},

    openURL: "",
    URLerror: "",
    peers: [],

    loadingPort: false,
  }, opts)

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

  state.status = portStatuses[0];

  this.render = function () {
    return hx`
      <div class="ui container fluid">
        <div>

          <div style="float: right;">
            <div style="display:inline-block;" style="text-align:left">
              <button class="ui small labeled icon fluid button
                ${ui.cl(state.test, 'compact', 'small', 'loading', 'disabled')}"
                style="white-space: nowrap;"
                onclick=${getPort}
                >
                <i class="${state.test.icon}  ${ui.print(state.test.color)} icon"></i>
                Port is ${state.test.text}
              </button>
            </div>
            <div style="display:inline-block;width:180px;">
              <div class="ui small  mini left fluid action input
                ${ui.cl(state.change, 'compact', 'small', 'loading', 'disabled')}">
                <input type="text" id="port" placeholder="Port to listen"
                  value=${state.port}
                  onkeydown=${setPortValue}
                  disabled="${ui.cl(state, 'disabled')}"
                  />
                <div class="ui button
                  ${ui.print(state.change.color)}
                  ${ui.cl(state.change, 'compact', 'small', 'loading', 'disabled')}
                  "
                  onclick=${changePort}
                  >
                  <i class="icon settings"></i>Change
                </div>
              </div>
            </div>
          </div>

          <h2 class="ui medium header">
            <br>Browse the web
          </h2>

          <div class="ui red floating icon message ${ui.myb('hidden', !state.failure)}">
            <i class="warning icon"></i>
            <i class="close icon" onclick=${dismiss}></i>
            <div class="content">
              <p>${state.failure}</p>
            </div>
          </div>

          <div align="center">
            <div class="field" style="width:60%;">
              <div class="ui  fluid action input">
                <input type="text" placeholder="#public key" onkeydown=${setOpenURL} class=${state.URLerror} />
                <div class="ui button" onclick=${openClick}>
                  <i class="icon cloud"></i>Browse
                </div>
              </div>
            </div>
          </div>
          <br>
          <div align="center">
            <div class="ui button
              ${ui.print(state.find.color)}
              ${ui.cl(state.find, 'compact', 'small', 'loading', 'disabled')}
              " onclick=${findClick}>
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

  function dismiss() {
    state.failure = "";
    update()
  }

  function setPortValue(){
    state.port = this.value;
  }
  function updatePort(res){
    state.port = res.Port;
    if (portStatuses[res.Status]!==null) {
      state.status = portStatuses[res.Status];
    }
    update();
  }

  function testPort(){
    var p = api.testPort()
    p.then(updatePort)
    vs(state.test).load(p, update)
      .begin(vs.loading, vs.disable).then(vs.loaded, vs.undisable).catch(vs.loaded, vs.undisable)
    vs(state).load(p, update)
      .begin(vs.unfail, vs.unfailure).catch(vs.fail, vs.failure())
  }
  function changePort() {
    var p = api.changePort(state.port)
    p.then(updatePort)
    vs(state.change).load(p, update)
      .begin(vs.loading, vs.disable).then(vs.loaded, vs.undisable).catch(vs.loaded, vs.undisable)
    vs(state).load(p, update)
      .begin(vs.unfail, vs.unfailure).catch(vs.fail, vs.failure())
    update();
  }

  function setOpenURL(){
    state.openURL = this.value;
  }
  function openClick() {
    var v = state.openURL
    state.URLerror = ""
    if (v == "") {
      state.URLerror = "red"
    } else {
      window.open("http://"+v+".me.com/")
    }
    update()
  }

  function findClick() {
    var p = api.list(0,30)
    p.then(function(res){
      state.peers = res
      update();
    })
    vs(state.find).load(p, update)
      .begin(vs.goBlue).then(vs.goGreen).catch(vs.goRed)
  }

  function getPort() {
    var p = api.getPort()
    p.then(updatePort)
    vs(state.change, state.test).load(p, update)
      .begin(vs.loading, vs.disable).then(vs.loaded, vs.undisable).catch(vs.loaded, vs.undisable)
    vs(state).load(p, update)
      .begin(vs.unfail, vs.unfailure).catch(vs.fail, vs.failure())
  }
  getPort();
}
