function fixVdom (h) {
  return function (tagName, props, children) {
    if (!props.attributes) props.attributes = {}
    Object.keys(props).forEach(function (key) {
       if (/^data-/.test(key)) { props.attributes[key] = props[key] }
    })
    return h(tagName, props, children)
  }
}
var vdom = require('virtual-dom')
vdom.h = fixVdom(vdom.h);
module.exports = vdom;
