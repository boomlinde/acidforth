window.onload = function () {
	var chk = function (d, code) { console.log("ERROR " + code + ": " + d); }
	var createController = function (name, min, max, step) {
		var row = document.createElement("tr");
		var nameColumn = document.createElement("td");
		var inputColumn = document.createElement("td");
		var input = document.createElement("input");
		input.setAttribute("type", "range");
		input.setAttribute("min", min);
		input.setAttribute("max", max);
		if (step !== undefined) input.setAttribute("step", step);
		$GET("/controls/" + name).success(function (d) {
			input.value = d;
		});
		input.oninput = function (e) {
			$POST("/controls/" + name, input.value).fail(chk);
		}
		inputColumn.appendChild(input);
		nameColumn.innerHTML = name;
		row.appendChild(nameColumn);
		row.appendChild(inputColumn);
		$('#controllers').appendChild(row);
	}

	createController('volume', 0, 1, 0.1);
	createController('pitch', 0, 10000);

	var container = T('ul', {}, [
		T('li', {}, []),
		T('li', {}, []),
		T('li', {}, [])
	]);

	document.body.appendChild(container.el);
}
