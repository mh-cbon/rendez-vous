
var $ = window.jQuery || window.$ || require('jquery');

function post(url, data, cb) {
  return $.ajax({
        url : url,
        type: "POST",
        data: JSON.stringify(data),
        contentType: "application/json; charset=utf-8",
        dataType   : "json",
        success    : cb
    });
}

module.exports = {
  post:post,
}
