function vs(){
  var states = Array.prototype.slice.call(arguments);
  var doer = {
    set: function(){
      var modifiers = Array.prototype.slice.call(arguments);
      modifiers.map(function(m){
        states.map(function(s){m(s)})
      })
    },
    load: function(p, update){
      var thens = [];
      var catchs = [];
      var soder = {
        begin: function(f){
          var modifiers = Array.prototype.slice.call(arguments);
          doer.set.apply(null, modifiers);
          update()
          return soder;
        },
        then: function(f){
          var modifiers = Array.prototype.slice.call(arguments);
          thens = thens.concat(modifiers);
          return soder;
        },
        catch: function(f){
          var modifiers = Array.prototype.slice.call(arguments);
          catchs = catchs.concat(modifiers);
          return soder;
        },
      }
      p.then(function(res){
        thens.map(function(m){
          states.map(function(s){m(s, res)})
        })
        update();
        return res
      })
      p.catch(function(err){
        catchs.map(function(m){
          states.map(function(s){m(s, err)})
        })
        update();
        return err
      })
      return soder
    }
  }
  return doer;
}

vs.unfailure = function(state){
  state.failure = null;
}
vs.unfail = function(state){
  state.failed = false;
}
vs.fail = function(state){
  state.failed = true;
}
vs.disable = function(state){
  state.disabled = true;
}
vs.undisable = function(state){
  state.disabled = false;
}
vs.loading = function(state){
  state.loading = true;
}
vs.loaded = function(state){
  state.loading = false;
}
vs.ran = function(state){
  state.ran = true;
}
vs.unran = function(state){
  state.ran = false;
}
vs.color = function(c){
  return function(state){
    state.color = c;
  }
}
vs.failure = function(f){
  var f = Array.prototype.slice.call(arguments);
  return function(state, err){
    f.map(function(ff){err=ff(err)})
    state.failure = err;
  }
}

vs.beginBlue = function(state){
  vs.unfail(state)
  vs.unfailure(state)
  vs.loading(state)
  vs.color("blue")(state)
}
vs.beginGreen = function(state){
  vs.unfail(state)
  vs.unfailure(state)
  vs.loading(state)
  vs.color("green")(state)
}
vs.goBlue = function(state){
  vs.unfail(state)
  vs.unfailure(state)
  vs.loaded(state)
  vs.color("blue")(state)
}
vs.goGreen = function(state){
  vs.unfail(state)
  vs.unfailure(state)
  vs.loaded(state)
  vs.color("green")(state)
}
vs.goRed = function(state){
  vs.fail(state)
  vs.loaded(state)
  vs.color("red")(state)
}
module.exports = vs;
