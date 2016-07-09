(function () {
	'use strict';
	var fabric = require('fabric').fabric;

	var initial_state = {
		'A30': {price: 40, coal: false, city: null},
		'B27': {price: 30, coal: false, city: null},
		'B29': {price: 30, coal: false, city: null},
		'C22': {price: 10, coal: false, city: null},
		'C24': {price: 10, coal: false, city: null},
		'C26': {price: 40, coal: false, city: null},
		'C28': {price: 20, coal: false, city: null},
		'C30': {price: 20, coal: false, city: null},
		'D19': {price: 10, coal: false, city: null},
		'D21': {price: 10, coal: false, city: null},
		'D23': {price: 10, coal: false, city: null},
		'D25': {price: 40, coal: false, city: null},
		'D27': {price: 40, coal: false, city: null},
		'D29': {price: 10, coal: false, city: null},
		'E4':  {price: 10, coal: false, city: null},
		'E8':  {price: 10, coal: false, city: null},
		'E10': {price: 10, coal: false, city: null},
		'E12': {price: 10, coal: false, city: null},
		'E16': {price: 10, coal: false, city: null},
		'E18': {price: 10, coal: false, city: null},
		'E20': {price: 10, coal: false, city: null},
		'E22': {price: 10, coal: false, city: null},
		'E24': {price: 10, coal: false, city: null},
		'E26': {price: 10, coal: false, city: null},
		'E28': {price: 10, coal: false, city: null},
		'E30': {price: 10, coal: false, city: null},
		'F3':  {price: 10, coal: false, city: null},
		'F5':  {price: 10, coal: false, city: null},
		'F7':  {price: 10, coal: false, city: null},
		'F9':  {price: 10, coal: false, city: null},
		'F11': {price: 10, coal: false, city: null},
		'F13': {price: 10, coal: false, city: null},
		'F15': {price: 10, coal: false, city: null},
		'F17': {price: 10, coal: false, city: null},
		'F19': {price: 10, coal: false, city: null},
		'F21': {price: 10, coal: false, city: null},
		'F23': {price: 10, coal: false, city: null},
		'F25': {price: 10, coal: false, city: null},
		'F27': {price: 10, coal: false, city: null},
		'G2':  {price: 10, coal: false, city: null},
		'G4':  {price: 10, coal: false, city: null},
		'G6':  {price: 10, coal: false, city: null},
		'G8':  {price: 10, coal: false, city: null},
		'G10': {price: 10, coal: false, city: null},
		'G12': {price: 10, coal: false, city: null},
		'G14': {price: 10, coal: false, city: null},
		'G16': {price: 10, coal: false, city: null},
		'G18': {price: 10, coal: true,  city: null},
		'G20': {price: 10, coal: false, city: null},
		'G22': {price: 10, coal: false, city: null},
		'G24': {price: 10, coal: false, city: null},
		'H1':  {price: 10, coal: false, city: null},
		'H3':  {price: 10, coal: false, city: null},
		'H5':  {price: 10, coal: false, city: null},
		'H7':  {price: 10, coal: false, city: null},
		'H9':  {price: 10, coal: false, city: null},
		'H11': {price: 10, coal: false, city: null},
		'H13': {price: 10, coal: false, city: null},
		'H15': {price: 10, coal: false, city: null},
		'H17': {price: 10, coal: true,  city: null},
		'H19': {price: 10, coal: false, city: null},
		'H21': {price: 10, coal: false, city: null},
		'H23': {price: 10, coal: false, city: null},
		'I0':  {price: 10, coal: false, city: null},
		'I2':  {price: 10, coal: false, city: null},
		'I4':  {price: 10, coal: false, city: null},
		'I6':  {price: 10, coal: false, city: null},
		'I8':  {price: 10, coal: false, city: null},
		'I10': {price: 10, coal: false, city: null},
		'I12': {price: 10, coal: false, city: null},
		'I14': {price: 10, coal: false, city: null},
		'I16': {price: 10, coal: true,  city: null},
		'I18': {price: 10, coal: false, city: null},
		'I20': {price: 10, coal: false, city: null},
		'I22': {price: 10, coal: false, city: null},
		'I24': {price: 10, coal: false, city: null},
		'J1':  {price: 10, coal: false, city: null},
		'J3':  {price: 10, coal: false, city: null},
		'J5':  {price: 10, coal: false, city: null},
		'J7':  {price: 10, coal: false, city: null},
		'J9':  {price: 10, coal: false, city: null},
		'J11': {price: 10, coal: false, city: null},
		'J13': {price: 10, coal: false, city: null},
		'J15': {price: 10, coal: true,  city: null},
		'J17': {price: 10, coal: false, city: null},
		'J19': {price: 10, coal: false, city: null},
		'J21': {price: 10, coal: false, city: null},
		'K2':  {price: 10, coal: false, city: null},
		'K4':  {price: 10, coal: false, city: null},
		'K6':  {price: 10, coal: false, city: null},
		'K8':  {price: 10, coal: false, city: null},
		'K10': {price: 10, coal: false, city: null},
		'K12': {price: 10, coal: false, city: null},
		'K14': {price: 10, coal: true,  city: null},
		'K16': {price: 10, coal: false, city: null},
		'K18': {price: 10, coal: false, city: null},
		'K20': {price: 10, coal: false, city: null},
		'K22': {price: 10, coal: false, city: null},
	};

	function create() {
		var radius = 50;
		var pnts = [0, 1, 2, 3, 4, 5].map(function (index) {
			var angle = (2*index + 1)*Math.PI/6;
			return {
				x: Math.round(1000*radius*(1 + Math.cos(angle)))/1000,
				y: Math.round(1000*radius*(1 + Math.sin(angle)))/1000,
			};
		});
		var hex_opts = {
			fill:    'green',
			stroke:  'black',
			strokeWidth: 1,
		};
		var price_opts = {
			fontSize:  15,
			textAlign: 'center',
			originX:   'center',
			originY:   'bottom',
			top:       1.87*radius,
			left:      radius,
		};
		var coal_opts = {
			fontSize:  15,
			textAlign: 'center',
			originX:   'center',
			originY:   'center',
			top:       radius,
			left:      radius,
			fill:      'white',
			backgroundColor: 'black',
		};

		var x_sep = (pnts[0].x - pnts[3].x)/2;
		var y_sep = 1.5 * radius;
		var hex_items = Object.keys(initial_state).map(function (id) {
			var details = initial_state[id];
			var items = [];

			items.push(new fabric.Polygon(pnts, hex_opts));
			items.push(new fabric.Text('$'+details.price, price_opts));
			if (details.coal) {
				items.push(new fabric.Text('COAL', coal_opts));
			}

			return new fabric.Group(items, {
				top:  y_sep*(id.charCodeAt(0) - 65),
				left: x_sep*parseInt(id.slice(1)),
			});
		});
		return new fabric.Group(hex_items, {selectable: false});
	}

	module.exports.create = create;
})();
