
var $ = window.jQuery || window.$;
var vdom = require('../vdom')
var hyperx = require('hyperx')
var hx = hyperx(vdom.h);
var main = require('main-loop')

var serialize = require('form-serialize');
var webUtils = require('../web-utils')

module.exports = AppsItem;

function AppsItem(router){

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
      ID:0,
      URL:"",
      Name:"",
      Type:"",
      StartCommand:"",
      StartPattern:"",
      KillCommand:"",
      KillPattern:"",
      AnnouncementName:"",
      Credentials: "",
      AnnouncementName: "",
      RequireCredentials: false,
      Announce: false,
      ExtraKill: false,
    },
    main: new runState(),
    add: new runState(),
    install: new runState(),
    check: new runState(),
    configData: new runState(),
    save: new runState(),
    start: new runState(),
    stop: new runState(),
    tab: "install",
  };
  function setApp(app){
    state.App = Object.assign(state.App, app);
    if(app.IsInstalled) {
      state.install.success()
    }
    if(app.PassTest) {
      state.check.success()
    }
    loop.update(state)
  }
  function setRequireCredentials(){
    state.App.RequireCredentials = !!this.checked;
    loop.update(state)
  }
  function setAnnounce(){
    state.App.Announce = !!this.checked;
    loop.update(state)
  }
  function setExtraKill(){
    state.App.ExtraKill = !!this.checked;
    loop.update(state)
  }
  function setURL(){
    state.App.URL = this.value;
    loop.update(state)
  }
  function setBinaryName(){
    state.App.BinaryName = this.value;
    loop.update(state)
  }
  function setType(){
    state.App.Type = this.value;
    loop.update(state)
  }
  function setStartCommand(){
    state.App.StartCommand = this.value;
    loop.update(state)
  }
  function setKillCommand(){
    state.App.KillCommand = this.value;
    loop.update(state)
  }
  function setAnnouncementName(){
    state.App.AnnouncementName = this.value;
    loop.update(state)
  }
  function setCredentials(){
    state.App.Credentials = this.value;
    loop.update(state)
  }
  function setAnnouncementName(){
    state.App.AnnouncementName = this.value;
    loop.update(state)
  }

  function render (state) {
    return hx`
    <div style="position:absolute;width:100%;">

      <div style="position:absolute;right:0;top: -35px;">
        <a class="small ui button green compact labeled icon button" href="#/apps">
          <i class="chevron left icon"></i> Return
        </a>
      </div>

      <h2 class="ui medium  header ${show(isFormNew())}" style="margin-top:0">
        <br>Create new app
      </h2>
      <h2 class="ui medium  header ${show(isFormEdit())}" style="margin-top:0">
        <br>Edit app ${state.App.URL}
      </h2>

      <div class="${show(state.main.failed)}">
        ${state.main.failure}
      </div>

      <form class="ui form ${hide(state.main.failed)}">
        <input type=hidden name=ID value=${state.App.ID} data-value-type="number"/>
        <div class="ui top attached tabular menu" onclick=${tabClick}>
          <div class="item ${active(isTab("install"))}" data-tab="install">Install</div>
          <div class="item ${active(isTab("controls"))} ${enabled(isFormEdit() && state.App.IsInstalled)}"
            onclick=${configClick} data-tab="controls">Controls</div>
          <div class="item ${active(isTab("announce"))} ${enabled(isFormEdit() && state.App.IsInstalled)}" data-tab="announce">Announce</div>
        </div>
        ${renderInstall(state, "install")}
        ${renderConfig(state, "controls")}
        ${renderAnnounce(state, "announce")}
      </form>
    </div>`
  }
  function isTab(a){return state.tab==a};
  function tabClick(ev) {
    if (ev.target.hasAttribute("data-tab")) {
      if (ev.target.classList.contains("disabled")) {
        ev.stopImmediatePropagation()
        ev.stopPropagation()
        return false
      }
      state.tab=ev.target.getAttribute("data-tab")
      loop.update(state)
    }
    return true
  }
  function configClick(ev) {
    if (ev.target.classList.contains("disabled")==false) {
      if(state.configData.ran==false){
        var p = webUtils.postJSON("/config/"+state.App.ID, {}, setApp)
        state.configData.follow(p, function(){loop.update(state)});
      }
    }
  }


  function renderInstall (state, tabName) {
    return hx`
    <div class="ui bottom attached tab segment ${active(isTab(tabName))} ${loading(state.main.loading)}" data-tab="${tabName}">
      <div class="field">
        <div class="two fields">
           <div class="field ${enabled(isFormNew())}">
             <input placeholder="The url location of the app..."
                type="text" name="URL" value="${state.App.URL}" onchange=${setURL} />
           </div>
           <div class="field">
             <input placeholder="The name of the binary..."
                type="text" name="BinaryName" value="${state.App.BinaryName}" onchange=${setBinaryName} />
           </div>
        </div>
        <div class="field">
           <select class="ui fluid search dropdown ${enabled(isFormNew())}" onchange=${setType} name=Type>
             ${state.okTypes.map(function(w){
               return hx`<option value="${w.Name}" ${state.App.Type==w.Name?'selected':''}>${w.Name}</option>`
             })}
           </select>
        </div>
      </div>
      <div class="${show(isFormNew())}" align="right">
        <a class="small ui button  compact right labeled icon ${state.add.class()}" onclick=${addClick}>
          <i class="plus left icon"></i> Add
        </a>
        ${state.add.failure}
      </div>
      <div class="${show(isFormEdit())}" align="right">
        <div class="ui divider"></div>
        <div class="${hide(state.App.IsInstalled && state.App.PassTest)}">
          The application is not yet installed. Use below buttons to install then configure it.
        </div>
        <a class="small ui button compact  labeled icon button
          ${state.install.class()} "
          onclick=${installClick}>
          <i class="power left icon"></i> Install
        </a>
        <a class="small ui button compact labeled icon button
          ${state.check.class()} ${enabled(state.App.IsInstalled)}"
          onclick=${checkClick}>
          <i class="plug left icon"></i> Check
        </a>
        <span class="${show(state.install.failed)}">
          <br>Install failed: ${state.install.failure}
        </span>
        <span class="${show(state.check.failed)}">
          <br>Configuration failed: ${state.check.failure}
        </span>
        <div class="${show(state.App.IsInstalled && state.App.PassTest)}" align="right">
          <div class="ui divider"></div>
          <br><br>
          The application is installed and pass the tests!<br><br>
          <a class="small ui button right blue compact labeled icon button" onclick=${continueClick}>
            <i class="chevron right icon"></i> Continue
          </a>
        </div>
      </div>
    </div>
    `
  }
  function renderConfig (state, tabName) {
    return hx`
    <div class="ui bottom attached tab segment ${active(isTab(tabName))} ${loading(state.main.loading)}" data-tab="${tabName}">
      <div class="field">
        <label>Start command</label>
        <div class="ui action input">
          <input placeholder="The command to start the app"
           type="text" name="StartCommand" value="${state.App.StartCommand}"
           onchange=${setStartCommand} />
          <div type="submit" class="ui button icon ${state.start.class()}"
            onclick="${startClick}"><i class="icon plus"></i>Check</div>
         </div>
         <br>
         Pattern: "${state.App.StartPattern}"
      </div>
      <div class="field">
        <label>Kill command</label>
        <div class="ui action input">
          <input placeholder="The command to kill the app"
           type="text" name="KillCommand" value="${state.App.KillCommand}"
           onchange=${setKillCommand} />
          <div type="submit" class="ui button icon ${state.stop.class()} ${state.start.success&&state.start.ran?"":"disabled"}"
            onclick="${stopClick}"><i class="icon plus"></i>Check</div>
        </div>
         Pattern: "${state.App.KillPattern}"
      </div>
      <div class="field">
        <label>Extra kill command</label>
        <div class="ui toggle checkbox ${checked(state.App.ExtraKill)}" onclick=${radioClick}>
          <input class="hidden" type="checkbox" name=ExtraKill value=true
          checked="${checked(state.App.ExtraKill)}"
          onchange=${setExtraKill} />
          <label>Apply a command to force kill the application</label>
        </div>
      </div>
      <div align="right">
        <div class="ui divider"></div>
        <br><br>
        <a class="small ui button compact labeled icon button ${state.save.class()}"
            onclick=${saveClick}>
          <i class="save right icon"></i> Save
        </a>
      </div>
    </div>
    `
  }
  function renderAnnounce (state, tabName) {
    return hx`
    <div class="ui bottom attached tab segment ${active(isTab(tabName))} ${loading(state.main.loading)}" data-tab="${tabName}">
      <div class="field">
        <label>Announce the application</label>
        <div class="ui toggle checkbox ${checked(state.App.Announce)}" onclick=${radioClick}>
          <input tabindex="0" class="hidden" type="radio"
             name=Announce value=true
             checked="${checked(state.App.Announce)}"
             onchange=${setAnnounce} />
          <label>Enable the application announcement</label>
        </div>
         <div class="field ${enabled(state.App.Announce)}">
           <input placeholder="The announcement name"
              type="text" name="AnnouncementName" value="${state.App.AnnouncementName}" onchange=${setAnnouncementName} />
         </div>
      </div>
      <div class="field">
        <label>Require credentials</label>
        <div class="ui toggle checkbox ${checked(state.App.RequireCredentials)}" onclick=${radioClick}>
          <input class="hidden" type="radio"
            checked="${checked(state.App.RequireCredentials)}"
            name=RequireCredentials value=true
             onchange=${setRequireCredentials} />
          <label>Set credentials on the application</label>
        </div>
        <div class="field ${enabled(state.App.RequireCredentials)}">
          <input placeholder="The credentials such as user:pwd"
             type="text" name="Credentials" value="${state.App.Credentials}" onchange=${setCredentials} />
        </div>
      </div>
      <div align="right">
        <div class="ui divider"></div>
        <br><br>
        <a class="small ui button compact labeled icon button ${state.save.class()}"
            onclick=${saveClick}>
          <i class="save right icon"></i> Save
        </a>
      </div>
    </div>
    `
  }
  function radioClick(ev){
    if (this.classList.contains("checked")) {
      this.classList.remove("checked")
      this.querySelector(".hidden").checked = false
    } else {
      this.classList.add("checked")
      this.querySelector(".hidden").checked = true;
    }
    this.querySelector(".hidden").onchange();
  }

  function isFormNew() { return state.formType==="new" }
  function isFormEdit() { return state.formType==="edit" }
  function show(when) { return when ? "visible" : "invisible" }
  function hide(when) { return when ? "invisible" : "visible" }
  function active(when) { return when ? "active" : "" }
  function enabled(when) { return when ? "" : "disabled" }
  function loading(when) { return when ? "loading" : "" }
  function checked(when) { return when ? "checked" : "" }

  function validateForm(){
    loop.target.querySelector("[name=URL]").parentNode.classList.remove("error")
    var form = loop.target.querySelector(".form")
    // var obj = serialize(form, { hash: true,empty:true });
    var obj = $(form).serializeJSON({checkboxUncheckedValue: "false", parseBooleans:"true"});
    obj.URL = obj.URL && obj.URL.trim();
    if(!obj.URL) {
      loop.target.querySelector("[name=URL]").parentNode.classList.add("error")
      return null
    }
    return obj
  }

  function returnClick() {
    that.emit("click-return");
  }
  function addClick() {
    var obj = validateForm();
    if(!obj) return
    setApp(obj);
    var p = webUtils.postJSON("/add", obj, setApp).done(function(){
      router.goto("/apps/edit/"+state.App.ID);
    })
    state.add.follow(p, function(){loop.update(state)});
    return false
  }
  function installClick() {
    var id = state.App.ID
    var p = webUtils.postJSON("/install/"+id, {}, setApp)
    state.install.follow(p, function(){loop.update(state)});
    return false
  }
  function checkClick() {
    var id = state.App.ID
    var p = webUtils.postJSON("/config/"+id, {}, setApp)
    state.check.follow(p, function(){loop.update(state)});
    return false
  }
  function continueClick() {
    var title = $(loop.target).find(".tabular .item[data-tab='configure']");
    title.click()
  }
  function saveClick() {
    var obj = validateForm();
    if(!obj) return
    var p = webUtils.postJSON("/update", obj, setApp)
    state.save.follow(p, function(){loop.update(state)});
    return false
  }
  function startClick() {
    return false
  }
  function stopClick() {
    return false
  }

  var loop = main(state, render, vdom);
  function enter(view){
    state.main.reset();
    state.add.reset();
    state.install.reset();
    state.check.reset();
    state.configData.reset();
    state.save.reset();
    state.start.reset();
    state.stop.reset();
    if (view=="new") {
      state.formType = view;
      var searchParams = new URLSearchParams(window.location.search);
      setApp({
        ID:-1,
        BinaryName: searchParams.get("BinaryName") || "",
        URL: searchParams.get("URL") || "",
        Type: searchParams.get("Type") || state.okTypes[0].Name,
        KillCommand: "",
        StartCommand: "",
        StartPattern: "",
        KillPattern: "",
        Credentials: "",
        AnnouncementName: "",
        RequireCredentials: false,
        Announce: false,
        ExtraKill: false,
      })
      var title = $(loop.target).find(".tabular .item[data-tab='install']");
      title.click()
      loop.update(state)
    } else if (view=="edit") {
      state.tab="install";
      state.formType = view;
      loop.update(state);
      var id = location.hash.slice(1).replace("/apps/edit/","");
      var p = webUtils.postJSON("/status/"+id,{}, setApp)
      state.main.follow(p, function(){loop.update(state)});
    }
  }
  this.enable = function(params){
    if (state.view!=params.view) {
      enter(params.view)
    }
  }
  router.on('/apps/:view', that.enable)

  this.install = function(to){
    loop.update(state);
    to.appendChild(loop.target);
  }
  this.uninstall = function(from){
    from.removeChild(loop.target)
  }
}

function runState(){
  var that = this;
  that.reset = function(){
    that.loading = false;
    that.failed = false;
    that.ran = false;
    that.failure = "";
  }
  that.fail = function(failure){
    that.ran = true;
    that.failed = true;
    that.loading = false;
    that.failure = failure;
  }
  that.success = function(){
    that.ran = true;
    that.failed = false;
    that.loading = false;
    that.failure = "";
  }
  that.class = function(){
    var ret = "blue"
    if(that.ran===true) {
      ret = "green"
      if(that.failed) {
        ret = "red"
      }
    }
    if (that.loading===true) {
      ret += " loading"
    }
    return ret;
  }
  that.follow = function(promise, update){
    that.reset();
    that.loading = true;
    update()
    return promise
      .fail(function(res){
        that.fail(res.responseText)
        update()
      })
      .done(function(res){
        that.success()
        update()
      })
  }
  this.reset();
}
