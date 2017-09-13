

var $ = window.jQuery || window.$;
var vdom = require('virtual-dom');
var hyperx = require('hyperx');
var hx = hyperx(vdom.h);
var main = require('main-loop');
var serialize = require('form-serialize');
var webUtils = require('../web-utils')

require('util').inherits(AppsItem, require('events').EventEmitter);
module.exports = AppsItem;


function AppsItem(){

  var that = this;
  var state = {
    formType: "new",
    okTypes: [{
      Name: 'GO',
      Description: 'Install with go get',
    },{
      Name: 'NPM',
      Description: 'Install with npm i',
    }],
    App: {
      URL:"",
      Type:"GO",
    },
    download: {
      working:false,
      failed:false,
      failure:""
    }
  };

  function render (state) {
      return hx`
<div style="position:absolute;width:100%;">

  <div style="position:absolute;right:0;top: -35px;">
    <a class="small ui button green compact  labeled icon button" href="#/apps">
      <i class="chevron left icon"></i> Return
    </a>
  </div>

  <h2 class="ui medium  header ${state.formType=="new"?'visible':'invisible'}" style="margin-top:0">
    <br>Create new app
  </h2>
  <h2 class="ui medium  header ${state.formType=="edit"?'visible':'invisible'}" style="margin-top:0">
    <br>Edit app ...
  </h2>

  <form class="ui form">
    <h4 class="ui dividing header">Download settings</h4>

    <div class="field">
      <div class="two fields">
         <div class="field">
           <input placeholder="The url location of the app..." type="text" name="URL" value="${state.App.URL}" />
         </div>
         <div class="field">
         <select class="ui fluid search dropdown" name=Type>
            ${state.okTypes.map(function(w){
              return hx`<option value="${w.Name}" ${state.App.Type==w.Name?'selected':''}>${w.Name}</option>`
            })}
          </select>
         </div>
      </div>
    </div>

    <div align="center" class="${state.download.working ? 'active' : ''}">
      <button class="ui  icon basic button ${state.download.failed ? 'red' : 'blue'}" type="button" onclick=${downloadClick}>
        <i class="download icon huge"></i>
        <div class="ui text ${state.download.failed ? 'invisible' : ''}">Click here to download your app</div>
        <div class="ui text ${state.download.failed ? '' : 'invisible red'}">
          ${state.download.failure}
          <br>Click again to retry
        </div>
      </button>
      <div class="ui text loader">Downloading the app...</div>
    </div>

    <div align="center" class="download-log ${state.download.working ? '' : 'invisible'}">
      todo..
    </div>
  </form>

</div>
`
  }

  var loop = main(state, render, vdom);

  function returnClick() {
    that.emit("click-return");
  }
  function downloadClick() {
    state.download.working = true;
    state.download.failed = false;
    state.download.failure = "";
    loop.target.querySelector("[name=URL]").parentNode.classList.remove("error")
    loop.update(state);
    var form = loop.target.querySelector(".form")
    var obj = serialize(form, { hash: true,empty:true });
    obj.URL = obj.URL && obj.URL.trim();
    if(!obj.URL) {
      loop.target.querySelector("[name=URL]").parentNode.classList.add("error")
      return
    }
    state.App.URL = obj.URL;
    state.App.Type = obj.Type;
    webUtils.post("/add", obj, function (res){
      state.download.working = false;
      loop.update(state);
    }).fail(function(res){
      state.download.failure = res.responseText;
      state.download.working = false;
      state.download.failed = true;
      loop.update(state)
    })
    return false
  }

  function semanticInstall(){
    $(loop.target).find(".ui.dropdown").dropdown();
  }
  function semanticUninstall(){
    $(loop.target).find(".ui.dropdown").destroy();
  }

  this.install = function(to){
    to.appendChild(loop.target);
    loop.update(state);
    setTimeout(semanticInstall,0)
  }
  this.uninstall = function(from){
    from.removeChild(loop.target)
    semanticUninstall()
  }
}
