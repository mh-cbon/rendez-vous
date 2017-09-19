
(function(){
  var failure = function(){
    // Do something with response error
    return Promise.reject(error);
  }
  var success = function(){
    // Do something with response data
    return response;
  }
  // axios.interceptors.response.use(success, null);
  // axios.interceptors.response.use(null, failure);
})();

function api(opts){
  var axios = require('axios').create({
    baseURL: '/',
    paramsSerializer: require('querystring')
  });

  this.list = function(start, limit){
    return axios.get("/list/"+start+"/"+limit)
  }
  this.deleteByID = function(id){
    return axios.post("/delete/"+id, {})
  }
  this.add = function(data){
    return axios.post("/add", data)
  }
  this.install = function(id){
    return axios.post("/install/"+id, {})
  }
  this.config = function(id){
    return axios.post("/config/"+id, {})
  }
  this.status = function(id){
    return axios.post("/status/"+id, {})
  }
  this.update = function(data){
    return axios.post("/update", data)
  }
  this.readData = function (res) {
    if (res.status && res.statusText && res.data) {
      res = res.data;
    }
    return res
  }
  this.readErr = function (error) {
    if (error.response) {
      // The request was made and the server responded with a status code
      // that falls out of the range of 2xx
      error = error.response.data || error.message;
    } else if (error.request) {
      // The request was made but no response was received
      // `error.request` is an instance of XMLHttpRequest in the browser and an instance of
      // http.ClientRequest in node.js
    } else {
      // Something happened in setting up the request that triggered an Error
      error = error.message;
    }
    return error
  }

  var that = this;
  this.wrapped = function(readData, readErr){
    var src = that;
    var ret = {};
    if (!readData) readData = that.readData;
    if (!readErr) readErr = that.readErr;
    Object.keys(src).forEach(function(m){
      ret[m] = function(){
        var args = Array.prototype.slice.call(arguments);
        return wrap(src[m].apply(src, args), readData, readErr)
      }
    })
    return ret;
  }
}

function wrap(p, readData, readErr){
  var thens = [];
  var catchs = [];
  var soder = {
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
    thens.map(function(m){ m(readData(res)) })
    return res
  })
  p.catch(function(err){
    catchs.map(function(m){m(readErr(err)) })
    return err
  })
  return soder
}

module.exports = api;
