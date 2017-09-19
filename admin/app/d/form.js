
function test(update, tests){
  var tests = Array.prototype.slice.call(arguments).slice(1);

  var failed = [];
  setTimeout(function(){
    tests.map(function(test){
      var f = test()
      if (f!==true){
        failed.push(f)
      }
    })
    if (failed.length===0) {
      thens.map(function(c){return c(failed)})
    } else {
      catchs.map(function(c){return c(failed)})
    }
    update();
  }, 0)

  var thens = [];
  var catchs = [];
  var continuation = {
    then: function(){
      thens = thens.concat(Array.prototype.slice.call(arguments));
      return continuation
    },
    catch: function(){
      catchs = catchs.concat(Array.prototype.slice.call(arguments));
      return continuation
    }
  }
  return continuation
}

function notEmpty(state, msg){
  return function(){
    var v = state.value || "";
    v = v.trim();
    if (v.length==0) {
      state.failure = msg || "It must not be empty"
      state.failed = true;
      return state.failure
    }
    state.failure = "";
    state.failed = false;
    return true
  }
}
function isSelected(state, options, msg){
  return function(){
    var v = state.getValue();
    var index = state.selectedIndex()
    if (!v && index==-1) {
      state.failure = msg || "Select an option"
      state.failed = true;
      return state.failure
    } else {
      var ok = false;
      options.map(function(o){if(o.value==v){ok=true}})
      if (!ok) {
        state.failure = msg || "Invalid option"
        state.failed = true;
        return state.failure
      }
    }
    state.failure = "";
    state.failed = false;
    return true
  }
}
function match(state, p, msg){
  return function(){
    var v = state.value || "";
    if (v.match(p)==null) {
      state.failure = msg || "Invalid option"
      state.failed = true;
      return state.failure
    }
    state.failure = "";
    state.failed = false;
    return true
  }
}
function ok(){
  return function(){
    return true
  }
}

module.exports = {
  test:test,
  notEmpty:notEmpty,
  isSelected:isSelected,
  match:match,
  ok:ok,
}
