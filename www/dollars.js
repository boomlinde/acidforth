$ = function (s) {
	var mangle = function (e) {e.on = e.addEventListener; return e}
	if (s[0] == '#') return mangle(document.querySelector(s));
	var e;
	var a = [];
	var es = document.querySelectorAll(s);
	for(var i = 0, e; e = es[i]; ++i) a.push(mangle(e));
	return a;
}
$WORD = function () {
	return {
		s: [], f: [], d: [],
		success: function (f) {this.s.push(f); return this},
		fail: function (f) {this.f.push(f); return this},
		done: function (f) {this.d.push(f); return this}
	}
}
$ALL = function () {
	var rp = $WORD();
	var re = [];
	var t = 0;
	var args = arguments;
	for (var _ = 0; _ < args.length; _++)
		re.push(undefined);
	for (var i = 0; i < args.length; i++) {
		args[i].done((function (i) {
			return function (data, s) {
				re[i] = {data: data, status: s};
				if (++t == args.length) {
					var failed = false;
					re.forEach(function (r) {
						if (r.status < 200 || r.status > 299) failed = true;
					});
					if (failed) rp.f.forEach(function (f) {f(re, 0)});
					else rp.s.forEach(function (f) {f(re, 200)});
					rp.d.forEach(function (f) {f(re, 0)});
				}
			}
		})(i));
	}
	return rp;
}
$REQ = function (u, h, d, t) {
	if (h === undefined) h = [];
	var r = new XMLHttpRequest();
	var p = $WORD();
	if (t !== undefined) r.timeout = t;
	r.onreadystatechange = function () {
		if (r.readyState == 4) {
			var s = r.status;
			if (s > 199 && s < 300) p.s.forEach(function (f) {f(r.responseText, s)});
			else p.f.forEach(function (f) {f(r.responseText, s)});
			p.d.forEach(function (f) {f(r.responseText, s)});
		}
	}
	if (d === undefined) r.open("GET", u, true);
	else r.open("POST", u, true);
	h.forEach(function (pair) {
		xmlhttp.setRequestHeader(pair[0], pair[1]);
	});
	r.send(d);
	return p;
}
$GET = function (u, h, t) {return $REQ(u, h, undefined, t)}
$POST = function (u, d, h, t) {return $REQ(u, h, d, t)}
$DELAY = function (t) {
	var p = $WORD();
	setTimeout(function () {p.d.forEach(function (f) {f("DELAY", 200)})}, t);
	return p;
}
