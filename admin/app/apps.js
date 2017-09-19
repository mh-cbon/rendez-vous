
var vdom = require('./vdom');
var hyperx = require('hyperx')
var hx = hyperx(vdom.h);

var emit = require('./d/emit')
var ui = require('./d/ui')
var ev = require('./d/ev')
var form = require('./d/form')
var vs = require('./d/vs')

var Button = require('./d/button')
var Link = require('./d/link')
var Tab = require('./d/tab')
var Select = require('./d/select')
var Text = require('./d/text')
var Toggle = require('./d/toggle')

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

module.exports = AppsItem;

function joinStr(glue){
  return function(res){
    if (res) {
      return res.join(glue)
    }
    return "";
  }
}

function AppsItem(update, router, opts){

  var state = this;
  Object.assign(state, {
    formType: "new",
    okTypes: [{
      value: 'GO',
      text: 'GO',
      Description: 'Install with go get',
    },{
      value: 'NPM',
      text: 'NPM',
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
  }, opts)

  function isFormNew() { return state.formType==="new" }
  function isFormEdit() { return state.formType==="edit" }
  function isAppReady() { return state.App.IsInstalled && state.App.PassTest }

  state.return = new Link(update,
    smallRightFloatedBtn({icon: "chevron left", color:"green", label:"Return", href:"#/apps"})
  );

  state.install = new AppsItemInstall(update, router, {App:state.App, okTypes:state.okTypes})
  state.controls = new AppsItemControls(update, {App:state.App})
  state.announce = new AppsItemAnnounce(update, {App:state.App})

  state.tab = new Tab(update)

  var installTab = state.tab.add("install", {title: "Install", render: state.install.render})
  var controlsTab = state.tab.add("controls", {title: "Controls", render: state.controls.render})
  var announceTab = state.tab.add("announce", {title: "Announce", render: state.announce.render})

  emit.call(state);
  state.return.on('click',function(e){
    state.emit("click-return",e,this)
  })
  state.install.add.on("click", function(ev){
    var p = state.install.validate()
    var about = vs(state.install.add)
    about.load(p, update)
      .begin(vs.beginBlue).catch(vs.goRed, vs.failure(joinStr(", ")))
      .then(function(){
        var p2 = api.add(state.install.getData());
        p2.then(state.updateApp, enterEdit)
        about.load(p2, update).then(vs.goGreen).catch(vs.goRed);
        vs(state).load(p2, update).begin(vs.unfailure).catch(vs.failure());
      })
    return false
  })
  state.install.install.on("click", function(){
    var p = api.install(state.App.ID);
    p.then(state.updateApp)
    vs(state.install.install).load(p, update)
      .begin(vs.beginBlue).then(vs.goGreen).catch(vs.goRed, vs.failure());
    return false
  })
  state.install.check.on("click", function(){
    var p = api.config(state.App.ID);
    p.then(state.updateApp)
    vs(state.install.check).load(p, update)
      .begin(vs.beginBlue).then(vs.goGreen).catch(vs.goRed, vs.failure());
    return false
  })
  state.install.continue.on("click", function(ev){
    if (isAppReady()) {
      window.location.hash = "#/apps/edit/"+state.App.ID
      state.tab.setActive("controls")
    } else{
      var p = api.config(params.ID);
      vs(state.install.continue).load(p, update).begin(vs.goBlue).then(vs.goGreen).catch(vs.goRed);
      p.then(state.updateApp, enterEdit, function(app){
        window.location.hash = "#/apps/edit/"+app.ID
      }).catch(setNotFound)
    }
    return false
  })
  state.controls.save.on("click", function(){
    var p = state.controls.validate()
    var about = vs(state.controls.save)
    about.load(p, update)
      .begin(vs.beginBlue).catch(vs.goRed, vs.failure(joinStr(", ")))
      .then(function(){
        var p2 = api.update(state.getData());
        p2.then(state.updateApp)
        about.load(p2, update)
          .then(vs.goGreen).catch(vs.goRed, vs.failure());
      })
    return false
  })
  state.announce.save.on("click", function(){
    var p = state.announce.validate()
    var about = vs(state.announce.save)
    about.load(p, update)
      .begin(vs.beginBlue).catch(vs.goRed, vs.failure(joinStr(", ")))
      .then(function(){
        var p2 = api.update(state.getData());
        p2.then(state.updateApp)
        about.load(p2, update)
          .then(vs.goGreen).catch(vs.goRed, vs.failure());
      })
    return false
  })

  state.getData = function(){
    return Object.assign({},state.App,state.install.getData(),state.controls.getData(),state.announce.getData())
  }
  state.updateApp = function(data){
    state.App = data;
    vs(controlsTab, announceTab).set(isAppReady() ? vs.undisable : vs.disable)
    state.install.updateApp(data)
    state.controls.updateApp(data)
    state.announce.updateApp(data)
    update()
  }
  function enterNew(){
    state.formType = "new";

    var searchParams = new URLSearchParams(window.location.search);
    state.updateApp({
      BinaryName: searchParams.get("BinaryName"),
      URL: searchParams.get("URL"),
      Type: searchParams.get("Type"),
    })
    state.install.enterNew()
    state.controls.enterNew()
    state.announce.enterNew()

    vs(state).set(vs.unfail, vs.unfailure)
    vs(controlsTab, announceTab).set(vs.disable)
    state.tab.setActive("install")

    update()
  }
  function enterEdit(params){
    state.formType = "edit";

    state.install.enterEdit()
    state.controls.enterEdit()
    state.announce.enterEdit()

    var p = api.status(params.ID);
    p.then(state.updateApp).catch(setNotFound)
    vs(state).load(p, update).begin(vs.unfail, vs.unfailure).catch(vs.fail, vs.failure())
    vs(state.tab.getActive()).load(p, update).begin(vs.loading).then(vs.loaded).catch(vs.loaded)
    vs(installTab, controlsTab, announceTab ).load(p, update).begin(vs.disable)
    vs(installTab).load(p, update).then(vs.undisable).catch(vs.undisable)

  }
  function setNotFound(){
    vs(controlsTab, announceTab).set(vs.disable)
    update()
  }
  router.on('/apps/new', enterNew)
  router.on('/apps/edit/:ID', enterEdit)

  function dismiss() {
    state.failure = "";
    update()
  }
  function returnClick() {
    state.emit("click-return", ev, this)
  }
  this.render = function () {
    return hx`
    <div class="apps-item">

      <div>
        ${state.return.render()}
      </div>

      <h2 class="ui medium  header ${ui.myb('invisible', isFormEdit())}">
        <br>Create new app
      </h2>
      <h2 class="ui medium  header ${ui.myb('invisible', isFormNew())}">
        <br>Edit app ${state.App.URL}
      </h2>

      <div class="ui red floating icon message ${ui.myb('hidden', !state.failure)}">
        <i class="warning icon"></i>
        <i class="close icon" onclick=${dismiss}></i>
        <div class="content">
          <p>${state.failure}</p>
        </div>
      </div>

      <form class="ui form ${ui.myb('hidden', state.failure)}">
        ${state.tab.render()}
      </form>
    </div>`
  }
}


function AppsItemInstall(update, router, opts){

  var state = this;
  Object.assign(state, {
    formType: "new",
    App: {},
    okTypes: [],
  }, opts)

  state.add = new Button(update,
    smallRightBtn({icon: "plus", color:"blue", label:"Add a new application"})
  );
  state.install = new Button(update,
    smallRightBtn({icon: "power", color:"blue", label:"Install the application"})
  );
  state.check = new Button(update,
    smallRightBtn({icon: "plug", color:"blue", label:"Check the configuration"})
  );
  state.continue = new Button(update,
    smallRightBtn({icon: "chevron right", color:"blue", label:"Continue"})
  );
  state.loc = new Text(update, {
    value:"",
    name:"URL",
    placeholder:"The url location of the app..."
  });
  state.binName = new Text(update, {
    value:"",
    name:"BinaryName",
    placeholder:"The name of the binary..."
  });
  state.selTypes = new Select(update, "Type", state.okTypes);

  emit.call(state);

  state.updateApp = function(app){
    state.loc.value = app.URL || "";
    state.binName.value = app.BinName || "";
    state.selTypes.selectValue(app.Type || state.okTypes[0].value);
    state.App = app;
    if ('IsInstalled' in app && app.IsInstalled) {
      state.install.color="green";
    }
    if ('PassTest' in app && app.PassTest && app.IsInstalled) {
      state.check.color="green";
    }
  }

  state.validate = function(){
    return form.test(update,
      form.notEmpty(state.loc, "Provide the location of the app"),
      form.isSelected(state.selTypes, state.okTypes, "Select the application type"),
      form.match(state.binName, /^[a-z0-9-_]*$/ig, "Invalid binary name")
    )
  }
  state.getData = function(){
    return {
      URL: state.loc.value,
      Type: state.selTypes.getValue(),
      BinName: state.binName.value || "",
    }
  }

  this.enterNew = function() {
    state.formType = "new"
    vs(state.add, state.install, state.check).set(vs.goBlue)
    vs(state.loc, state.binName, state.selTypes).set(vs.undisable)
    update();
  }
  this.enterEdit = function() {
    state.formType = "edit"
    vs(state.add, state.install, state.check).set(vs.goBlue)
    vs(state.loc, state.binName, state.selTypes).set(vs.disable)
    update();
  }

  function isFormNew() { return state.formType==="new" }
  function isFormEdit() { return state.formType==="edit" }
  function isAppReady() { return state.App.IsInstalled && state.App.PassTest }

  this.render = function () {
    return hx`
    <div>
      <div class="field">
        ${ui.two(state.loc, state.binName)}
        ${ui.field(state.selTypes)}
      </div>
      <div class="${ui.myb('invisible', isFormEdit())}" align="right">
        ${state.add.render()}
        <br>
        ${state.add.failure}
      </div>
      <div class="${ui.myb('invisible', isFormNew())}" align="right">
        <div class="ui divider"></div>
        <div class="${ui.myb('invisible', isAppReady())}">
          The application is not yet ready. Use below buttons to install then configure it.
        </div>
        ${state.install.render()}
        ${state.check.render()}
        <span class="${ui.myb('invisible', !state.install.failed)}">
          <br>Install failed: ${state.install.failure}
        </span>
        <span class="${ui.myb('invisible', !state.check.failed)}">
          <br>Configuration failed: ${state.check.failure}
        </span>
        <div class="${ui.myb('invisible', !isAppReady())}" align="right">
          <div class="ui divider"></div>
          <br><br>
          The application is installed and pass the tests!<br><br>
          ${state.continue.render()}
        </div>
      </div>
    </div>`
  }
}

function AppsItemControls(update, opts){

  var state = this;
  Object.assign(state, {
    App:{},
    KillPattern:"",
    StartPattern:"",
  }, opts)

  state.StartCommand = new Text(update, {
    value:"",
    name:"StartCommand",
    // change:setStartCommand,
    placeholder:"The command to start the app",
    action: {
      type: "submit",
      icon: {icon: "plus"},
      label:"Check",
      disabled: true,
      // click: startClick,
    },
  });
  state.KillCommand = new Text(update, {
    value:"",
    name:"KillCommand",
    // change:setKillCommand,
    placeholder:"The command to kill the app",
    action: {
      type: "submit",
      icon: {icon: "plus"},
      label:"Check",
      disabled: true,
      // click: stopClick,
    },
  });
  state.ExtraKill = new Toggle(update, {
    value:true,
    name:"ExtraKill",
    // change:setExtraKillCommand,
    label:"Apply a command to force kill the application",
  });
  state.save = new Button(update,
    smallRightBtn({icon: "save right", color:"blue", label:"Save",
      // click: saveClick,
    })
  );

  state.updateApp = function(app){
    state.StartCommand.value = app.StartCommand || "";
    state.KillCommand.value = app.KillCommand || "";
    state.ExtraKill.checked = !!app.ExtraKill;
    state.KillPattern = app.KillPattern;
    state.StartPattern = app.StartPattern;
  }

  state.validate = function(){
    return form.test(update,
      form.ok()
    )
  }
  state.getData = function(){
    return {
      StartCommand: state.StartCommand.value,
      KillCommand: state.KillCommand.value,
      ExtraKill: !!state.ExtraKill.checked,
    }
  }

  this.enterNew = function() {
    state.formType = "new"
    vs(state.save).set(vs.goBlue)
    vs(state.save, state.KillCommand, state.ExtraKill, state.KillCommand).set(vs.disable)
    update();
  }
  this.enterEdit = function() {
    state.formType = "edit"
    vs(state.save).set(vs.goBlue)
    vs(state.save, state.KillCommand, state.ExtraKill, state.KillCommand).set(vs.undisable)
    update();
  }

  state.StartCommand.on("keyup", function(){
    state.StartCommand.action.disabled = !this.value;
    update()
  })
  state.KillCommand.on("keyup", function(){
    state.KillCommand.action.disabled = !this.value;
    update()
  })

  emit.call(state);

  this.render = function () {
    return hx`
    <div>
      <div class="field">
        <label>Start command</label>
        ${state.StartCommand.render()}
         <br>
         Pattern: "${state.StartPattern}"
      </div>
      <div class="field">
        <label>Kill command</label>
        ${state.KillCommand.render()}
         Pattern: "${state.KillPattern}"
      </div>
      <div class="field">
        <label>Extra kill command</label>
        ${state.ExtraKill.render()}
      </div>
      <div align="right">
        <div class="ui divider"></div>
        <br><br>
        ${state.save.render()}
      </div>
    </div>`
  }
}

function AppsItemAnnounce(update, opts){

  var state = this;
  Object.assign(state, {
    App: {},
  }, opts)

  state.updateApp = function(app){
    state.AnnouncementName.value = app.AnnouncementName || "";
    state.Announce.checked = !!app.Announce;
    state.RequireCredentials.checked = !!app.RequireCredentials;
    state.Credentials.value = app.Credentials || "";
  }

  state.validate = function(){
    return form.test(update,
      form.ok()
    )
  }
  state.getData = function(){
    return {
      AnnouncementName: state.AnnouncementName.value,
      Announce: !!state.Announce.checked,
      RequireCredentials: !!state.RequireCredentials.checked,
      Credentials: state.Credentials.value,
    }
  }

  this.enterNew = function() {
    state.formType = "new"
    vs(state.save).set(vs.goBlue)
    vs(state.save, state.AnnouncementName, state.Announce, state.RequireCredentials, state.Credentials).set(vs.disable)
    update();
  }
  this.enterEdit = function() {
    state.formType = "edit"
    vs(state.save).set(vs.goBlue)
    vs(state.save, state.AnnouncementName, state.Announce, state.RequireCredentials, state.Credentials).set(vs.undisable)
    update();
  }

  state.Announce = new Toggle(update, {
    value:true,
    name:"Announce",
    // change:setAnnounce,
    label:"Enable the application announcement",
  });
  state.Announce.on("change", function(){
    state.AnnouncementName.disabled = !state.Announce.checked;
    update()
  })
  state.AnnouncementName = new Text(update, {
    value:"",
    name:"AnnouncementName",
    // change:setAnnouncementName,
    disabled:true,
    // disabled:state.App.Announce,
    placeholder:"The announcement name",
  });
  state.RequireCredentials = new Toggle(update, {
    value:true,
    name:"RequireCredentials",
    // change:setRequireCredentials,
    label:"Set credentials on the application",
  });
  state.RequireCredentials.on("change", function(){
    state.Credentials.disabled = !state.RequireCredentials.checked;
    update()
  })
  state.Credentials = new Text(update, {
    value:"",
    name:"Credentials",
    // change:setCredentials,
    disabled:true,
    // disabled:state.App.RequireCredentials,
    placeholder:"The credentials such as user:pwd",
  });
  state.save = new Button(update,
    smallRightBtn({icon: "save right", color:"blue", label:"Save",
    // click: saveClick,
    })
  );

  emit.call(state);

  this.render = function () {
    return hx`
    <div >
      <div class="field">
        <label>Announce the application</label>
        ${state.Announce.render()}
        ${state.AnnouncementName.render()}
      </div>
      <div class="field">
        <label>Require credentials</label>
        ${state.RequireCredentials.render()}
        ${state.Credentials.render()}
      </div>
      <div align="right">
        <div class="ui divider"></div>
        <br><br>
        ${state.save.render()}
      </div>
    </div>`
  }
}
