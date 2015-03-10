T = function (tag, attribs, children) {
	var el = document.createElement("tag");
	for (var key in attribs) {
		el.setAttribute(key, attribs[key]);
	}
	children.forEach(function (child) {
		el.appendChild(child.el);
	});
	return {
		el: el,
		children: children,
		clear: function () {
			this.children.forEach(function (child) {
				child.remove();
				this.el.removeChild(child.el);
			});
			this.children = [];
		}
	}
}
