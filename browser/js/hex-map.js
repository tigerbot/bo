(function () {
	'use strict';
	var fabric = require('fabric').fabric;

	var initial_state = {
		'A30': {price: 40, coal: false, city: {name: 'Augusta', color: 'black', revenue: [20, 20, 20, 20, 30, 40]}},
		'B27': {price: 30, coal: false, city: {name: 'Burlington', color: 'black', revenue: [10, 20, 20, 20, 30, 30]}},
		'B29': {price: 30, coal: false, city: null},
		'C22': {price: 10, coal: false, city: null},
		'C24': {price: 10, coal: false, city: null},
		'C26': {price: 40, coal: false, city: null},
		'C28': {price: 20, coal: false, city: {name: 'Concord', color: 'black', revenue: [20, 20, 20, 20, 20, 30]}},
		'C30': {price: 20, coal: false, city: {name: 'Portsmouth', color: 'black', revenue: [20, 20, 20, 20, 20, 30]}},
		'D19': {price: 10, coal: false, city: {name: 'Buffalo', color: '#109B1C', revenue: [20, 30, 30, 40, 50, 60]}},
		'D21': {price: 10, coal: false, city: {name: 'Syracuse', color: '#109B1C', revenue: [10, 20, 20, 30, 30, 40]}},
		'D23': {price: 10, coal: false, city: {name: 'Utica', color: '#109B1C', revenue: [10, 10, 10, 20, 20, 20]}},
		'D25': {price: 40, coal: false, city: {name: 'Albany', color: '#109B1C', revenue: [30, 30, 40, 40, 40, 50]}},
		'D27': {price: 40, coal: false, city: null},
		'D29': {price: 10, coal: false, city: {name: 'Boston', color: 'black', revenue: [30, 30, 40, 40, 40, 50]}},
		'E4':  {price: 20, coal: false, city: {name: 'Chicago', color: 'blue', revenue: [20, 30, 50, 70, 90, 100]}},
		'E8':  {price: 10, coal: false, city: null},
		'E10': {price: 10, coal: false, city: null},
		'E12': {price: 10, coal: false, city: {name: 'Detroit', color: 'black', revenue: [20, 30, 40, 60, 80, 90]}},
		'E16': {price: 10, coal: false, city: null},
		'E18': {price: 10, coal: false, city: null},
		'E20': {price: 40, coal: false, city: null},
		'E22': {price: 40, coal: false, city: null},
		'E24': {price: 40, coal: false, city: null},
		'E26': {price: 40, coal: false, city: null},
		'E28': {price: 10, coal: false, city: {name: 'Hartford', color: 'black', revenue: [20, 20, 20, 30, 30, 30]}},
		'E30': {price: 20, coal: false, city: {name: 'Providence', color: 'black', revenue: [20, 30, 30, 30, 30, 30]}},
		'F3':  {price: 10, coal: false, city: null},
		'F5':  {price: 10, coal: false, city: null},
		'F7':  {price: 10, coal: false, city: null},
		'F9':  {price: 10, coal: false, city: null},
		'F11': {price: 10, coal: false, city: null},
		'F13': {price: 10, coal: false, city: {name: 'Cleveland', color: 'black', revenue: [20, 30, 40, 50, 60, 60]}},
		'F15': {price: 10, coal: false, city: null},
		'F17': {price: 80, coal: false, city: null},
		'F19': {price: 80, coal: false, city: null},
		'F21': {price: 60, coal: false, city: null},
		'F23': {price: 40, coal: false, city: null},
		'F25': {price: 60, coal: false, city: {name: 'New York', color: '#109B1C', revenue: [30, 40, 50, 60, 70, 80]}},
		'F27': {price: 10, coal: false, city: {name: 'New Haven', color: 'black', revenue: [20, 20, 30, 30, 30, 40]}},
		'G2':  {price: 10, coal: false, city: {name: 'Springfield', color: 'black', revenue: [10, 10, 20, 20, 20, 30]}},
		'G4':  {price: 10, coal: false, city: null},
		'G6':  {price: 10, coal: false, city: null},
		'G8':  {price: 10, coal: false, city: {name: 'Fort Wayne', color: 'black', revenue: [10, 20, 20, 30, 40, 50]}},
		'G10': {price: 10, coal: false, city: null},
		'G12': {price: 10, coal: false, city: null},
		'G14': {price: 20, coal: false, city: null},
		'G16': {price: 60, coal: false, city: {name: 'Pittsburgh', color: 'red', revenue: [20, 30, 40, 60, 70, 80]}},
		'G18': {price: 80, coal: true,  city: null},
		'G20': {price: 20, coal: false, city: {name: 'Harrisburg', color: 'red', revenue: [10, 10, 20, 20, 20, 20]}},
		'G22': {price: 20, coal: false, city: null},
		'G24': {price: 20, coal: false, city: {name: 'Philadelphia', color: 'red', revenue: [30, 40, 40, 40, 50, 60]}},
		'H1':  {price: 10, coal: false, city: null},
		'H3':  {price: 10, coal: false, city: null},
		'H5':  {price: 10, coal: false, city: null},
		'H7':  {price: 10, coal: false, city: {name: 'Indianapolis', color: 'black', revenue: [20, 30, 30, 40, 50, 60]}},
		'H9':  {price: 10, coal: false, city: null},
		'H11': {price: 10, coal: false, city: null},
		'H13': {price: 20, coal: false, city: null},
		'H15': {price: 60, coal: false, city: {name: 'Wheeling', color: 'black', revenue: [20, 20, 30, 40, 50, 60]}},
		'H17': {price: 60, coal: true,  city: null},
		'H19': {price: 20, coal: false, city: null},
		'H21': {price: 10, coal: false, city: null},
		'H23': {price: 10, coal: false, city: {name: 'Baltimore', color: 'black', revenue: [20, 30, 30, 40, 40, 50]}},
		'I0':  {price: 40, coal: false, city: {name: 'Saint Louis', color: 'black', revenue: [30, 40, 50, 60, 70, 90]}},
		'I2':  {price: 10, coal: false, city: null},
		'I4':  {price: 10, coal: false, city: null},
		'I6':  {price: 10, coal: false, city: null},
		'I8':  {price: 10, coal: false, city: null},
		'I10': {price: 20, coal: false, city: {name: 'Cincinnati', color: 'black', revenue: [30, 40, 50, 50, 60, 70]}},
		'I12': {price: 40, coal: false, city: null},
		'I14': {price: 80, coal: false, city: null},
		'I16': {price:100, coal: true,  city: null},
		'I18': {price: 80, coal: false, city: null},
		'I20': {price: 10, coal: false, city: null},
		'I22': {price: 20, coal: false, city: {name: 'Washington', color: 'black', revenue: [20, 20, 30, 30, 30, 30]}},
		'I24': {price: 10, coal: false, city: {name: 'Dover', color: 'black', revenue: [10, 10, 10, 20, 20, 20]}},
		'J1':  {price: 10, coal: false, city: null},
		'J3':  {price: 10, coal: false, city: null},
		'J5':  {price: 10, coal: false, city: null},
		'J7':  {price: 20, coal: false, city: {name: 'Louisville', color: 'black', revenue: [20, 30, 30, 40, 40, 50]}},
		'J9':  {price: 10, coal: false, city: null},
		'J11': {price: 40, coal: false, city: null},
		'J13': {price: 40, coal: false, city: {name: 'Huntington', color: 'black', revenue: [10, 10, 20, 30, 30, 40]}},
		'J15': {price:100, coal: true,  city: null},
		'J17': {price: 80, coal: false, city: null},
		'J19': {price: 20, coal: false, city: null},
		'J21': {price: 20, coal: false, city: {name: 'Richmond', color: 'black', revenue: [30, 30, 20, 20, 20, 30]}},
		'K2':  {price: 20, coal: false, city: {name: 'Cairo', color: 'black', revenue: [10, 20, 20, 20, 20, 20]}},
		'K4':  {price: 10, coal: false, city: null},
		'K6':  {price: 10, coal: false, city: null},
		'K8':  {price: 10, coal: false, city: null},
		'K10': {price: 20, coal: false, city: {name: 'Lexington', color: 'black', revenue: [10, 20, 20, 30, 30, 30]}},
		'K12': {price:100, coal: false, city: null},
		'K14': {price: 80, coal: true,  city: null},
		'K16': {price: 20, coal: false, city: {name: 'Roanoke', color: 'black', revenue: [20, 20, 20, 20, 20, 20]}},
		'K18': {price: 10, coal: false, city: null},
		'K20': {price: 10, coal: false, city: null},
		'K22': {price: 20, coal: false, city: {name: 'Norfolk', color: 'black', revenue: [20, 20, 30, 30, 30, 40]}},
	};

	function create_rev_bar(revenue) {
		var rev_opts = {
			fontSize:  12,
			textAlign: 'center',
			originY:   'center',
			originX:   'center',
			height:    15,
			width:     17,
		};
		var clrs = [{fill: 'white'}, {fill: 'black'}];

		var x_offset = 0;
		return new fabric.Group(revenue.map(function (value, index) {
			var txt = new fabric.Text(value.toString(), rev_opts);
			var bg  = new fabric.Rect(rev_opts);

			txt.set(clrs[(index+0) % 2]);
			bg.set( clrs[(index+1) % 2]);
			x_offset += bg.getWidth();
			return (new fabric.Group([bg, txt], {originY: 'bottom', originX: 'right', left: x_offset}));
		}));
	}

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
			fill:    '#688E45',
			stroke:  'black',
			strokeWidth: 1,
		};
		var price_opts = {
			fontSize:  12,
			textAlign: 'center',
			originX:   'center',
			originY:   'bottom',
			top:       1.82*radius,
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
		var city_opts = {
			fontSize:    15,
			fontWeight:  'bold',
			shadow:      'white 0 0 2px',
			originX:     'center',
			originY:     'bottom',
			top:         radius*1.2,
			left:        radius,
		};
		var rev_bar_opts = {
			originX: 'center',
			originY: 'top',
			top:     radius*1.25,
			left:    radius,
		};

		// we divide the width of the pnts by 2 because the IDs only use every other number.
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
			if (details.city) {
				items.push((new fabric.Text(details.city.name, city_opts)).set({fill: details.city.color}));
				items.push(create_rev_bar(details.city.revenue).scaleToWidth(1.85*x_sep).set(rev_bar_opts));
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
