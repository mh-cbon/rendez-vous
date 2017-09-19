
function load(prm, update, states){
  var states = Array.prototype.slice.call(arguments).slice(2);
  states.map(function(s){s.loading=true;})
  update()
  return prm.then(function(res){
    states.map(function(s){s.loading=false;})
    update()
    return res
  })
}

function run(prm, update, states){
  var states = Array.prototype.slice.call(arguments).slice(2);
  states.map(function(s){
    s.loading=true;
    s.failure = "";
    s.failed = true;
  })
  update()
  var oks = [];
  var fails = [];
  var u = {ok:function(f){oks.push(f); return u;}, fail:function(f){fails.push(f); return u;}}
  prm.then(function(res){
    oks.map(function(ok){
      ok(res)
    })
    states.map(function(s){
      s.loading=false;
      s.failed = false;
    })
    update()
    return res
  }).catch(function(error){
    fails.map(function(fail){
      fail(error)
    })
    states.map(function(s){
      s.loading=false;
      s.failed=true;
      s.failure = error.message;
    })
    update()
    return error
  })
  return u
}

module.exports = {
  load: load,
  run: run,
}
