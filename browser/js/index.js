(function () {
	'use strict';
	var domready = require('domready');
	var fabric   = require('fabric').fabric;

	var hex_grid = [
		'A30',
		'B27', 'B29',
		'C22', 'C24', 'C26', 'C28', 'C30',
		'D19', 'D21', 'D23', 'D25', 'D27', 'D29',
		'E4',  'E8',  'E10', 'E12', 'E16', 'E18', 'E20', 'E22', 'E24', 'E26', 'E28', 'E30',
		'F3',  'F5',  'F7',  'F9',  'F11', 'F13', 'F15', 'F17', 'F19', 'F21', 'F23', 'F25', 'F27',
		'G2',  'G4',  'G6',  'G8',  'G10', 'G12', 'G14', 'G16', 'G18', 'G20', 'G22', 'G24',
		'H1',  'H3',  'H5',  'H7',  'H9',  'H11', 'H13', 'H15', 'H17', 'H19', 'H21', 'H23',
		'I0',  'I2',  'I4',  'I6',  'I8',  'I10', 'I12', 'I14', 'I16', 'I18', 'I20', 'I22', 'I24',
		'J1',  'J3',  'J5',  'J7',  'J9',  'J11', 'J13', 'J15', 'J17', 'J19', 'J21',
		'K2',  'K4',  'K6',  'K8',  'K10', 'K12', 'K14', 'K16', 'K18', 'K20', 'K22',
	]

	function create_hex_map() {
		var radius = 50;
		var pnts = [0, 1, 2, 3, 4, 5].map(function (index) {
			var angle = (2*index + 1)*Math.PI/6;
			return {
				x: Math.round(1000*radius*(1 + Math.cos(angle)))/1000,
				y: Math.round(1000*radius*(1 + Math.sin(angle)))/1000,
			};
		});
		var opts = {
			fill: 'green',
			stroke: 'black',
			strokeWidth: 1,
		};

		var x_sep = (pnts[0].x - pnts[3].x)/2;
		var y_sep = 1.5 * radius;
		var hex_items = hex_grid.map(function (id) {
			var hex = new fabric.Polygon(pnts, opts);
			hex.setTop(y_sep*(id.charCodeAt(0) - 65));
			hex.setLeft(x_sep*parseInt(id.slice(1)));
			return hex;
		});
		return new fabric.Group(hex_items, {selectable: false});
	}

	domready(function () {
		var canvas = new fabric.Canvas('main', {
			backgroundColor: 'rgb(200, 200, 255)',
			height: 650,
			width: 1050,
		});

		var map = create_hex_map();
		map.scaleToWidth(canvas.getWidth() - 20);
		map.set({top: 10, left: 10});

		canvas.add(map);
		window.canvas = canvas;
	});
})();
