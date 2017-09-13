
var $ = window.jQuery || window.$;
var vdom = require('virtual-dom');
var hyperx = require('hyperx');
var hx = hyperx(vdom.h);
var main = require('main-loop');
var webUtils = require('../web-utils');

require('util').inherits(AppsList, require('events').EventEmitter);
module.exports = AppsList;

function AppsList(){

  var that = this;
  var state = {
    apps: [],
    loading:false,
  };

  function render (state) {
      return hx`
<div style="position:absolute;width:100%;" class="">
  <div style="position:absolute;right:0;top: -35px;">
    <button class="small ui button green compact  labeled icon button" onclick=${refreshClick}>
      <i class="refresh icon"></i> Refresh
    </button>
    <a class="small ui button blue compact  labeled icon button" href="#/apps/new">
      <i class="plus icon"></i> New
    </a>
  </div>

  <h2 class="ui medium  header " style="margin-top:0">
    <br>Apps list
  </h2>

  <table class="ui celled striped table">
    <tbody>
      ${state.apps.map(function (w, i) {
        return hx`<tr>
          <td class="collapsing">${w.URL || w.Name}</td>
          <td>${w.UpdatedAt.substring(0,10)}</td>
          <td class="right aligned collapsing">
            <a disabled="${w.IsSystem ? 'disabled' : ''}" href="#/apps/edit/${w.ID}">
              <i class="write icon ${w.IsSystem ? 'grey' : ''}"></i>
            </a>
          </td>
          <td class="right aligned collapsing">
            <a disabled="${w.IsSystem ? 'disabled' : ''}" id="${w.ID}" onclick=${deleteClick}  href="#/apps/delete/${w.ID}">
              <i class="trash icon ${w.IsSystem ? 'grey' : ''}"></i>
            </a>
          </td>
        </tr>`
      })}
    </tbody>
  </table>
  <div align="center">
    <button class="bt-loadmore" onclick=${loadMore}>Load more</button>
    <div class="ui text loader ${state.loading?'active':''}">loading...</div>
  </div>
</div>
`
  }

  function getApps(start,limit,done) {
    state.loading = true
    loop.update(state);
    $.get("/list/"+start+"/"+limit, function(res) {
      res = JSON.parse(res)
      if (res) {
        done(res)
      }
      state.loading = false
      loop.update(state);
    }).fail(function(res){
      state.loading = false
      loop.update(state);
    })
  }
  var start = 0;
  var limit = 10;
  function loadMore() {
    getApps(start,limit,function(res) {
      state.apps = state.apps.concat(res);
      start += res.length;
      loop.update(state);
    })
  }
  function refreshClick() {
    getApps(0,start+limit,function(res) {
      state.apps = res;
      loop.update(state);
    })
  }
  function deleteClick() {
    webUtils.post("/delete/"+this.id, {}, function(res){
      refreshClick();
    }).catch(console.log)
    return false;
  }

  var loop = main(state, render, vdom);
  this.install = function(to){
    to.appendChild(loop.target);
    loop.update(state);
    loadMore();
  }
  this.uninstall = function(from){
    from.removeChild(loop.target)
  }
}
