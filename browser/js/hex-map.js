(function () {
	'use strict';
	var domready = require('domready');
	var fabric   = require('fabric').fabric;
	var canvas;

	var x_sep, y_sep;
	var water_clr = '#DEEFF7';
	var grass_clr = '#688E45';
	var hex_elems = {};

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

		return new fabric.Group(revenue.map(function (value, index) {
			var txt = new fabric.Text(value.toString(), rev_opts);
			var bg  = new fabric.Rect(rev_opts);

			txt.set(clrs[(index+0) % 2]);
			bg.set( clrs[(index+1) % 2]);
			return new fabric.Group([bg, txt], {originY: 'bottom', left: index*rev_opts.width});
		}));
	}

	function create_map() {
		var radius = 50;
		var border_width = 1.5;
		var pnts = [0, 1, 2, 3, 4, 5].map(function (index) {
			var angle = (2*index + 1)*Math.PI/6;
			return {
				x: Math.round(1000*(radius - border_width/2)*(1 + Math.cos(angle)))/1000,
				y: Math.round(1000*(radius - border_width/2)*(1 + Math.sin(angle)))/1000,
			};
		});

		var hex_opts = {
			fill:    grass_clr,
			stroke:  'black',
			strokeWidth: border_width,
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
			fontSize:  13,
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
		x_sep = (pnts[0].x - pnts[3].x + border_width)/2;
		y_sep = 1.5 * radius;
		return new fabric.Group(Object.keys(initial_state).map(function (id) {
			var details = initial_state[id];
			var items = [];

			items.push(new fabric.Polygon(pnts, hex_opts));
			items.push(new fabric.Text('$'+details.price, price_opts));
			if (details.coal) {
				items.push(new fabric.Text('COAL', coal_opts));
			}
			if (details.city) {
				items.push(create_rev_bar(details.city.revenue).scaleToWidth(1.85*x_sep).set(rev_bar_opts));
				items.push((new fabric.Text(details.city.name, city_opts)).set({fill: details.city.color}));
			}

			var group = new fabric.Group(items, {
				top:  y_sep*(id.charCodeAt(0) - 65),
				left: x_sep*parseInt(id.slice(1)),
			});

			hex_elems[id] = {
				group:    group,
				selected: false,
			};

			return group;
		}));
	}

	function selected(id) {
		if (typeof id === 'undefined') {
			return Object.keys(hex_elems).filter(function (id) {
				return hex_elems[id].selected;
			});
		}
		if (hex_elems.hasOwnProperty(id)) {
			return hex_elems[id].selected;
		}
		return false;
	}
	function select(id) {
		if (hex_elems.hasOwnProperty(id)) {
			hex_elems[id].selected = true;
			hex_elems[id].group.item(0).set({stroke: 'yellow', opacity: 0.8});
		}
		canvas.renderAll();
	}
	function deselect(id) {
		var keys;
		if (id === '*') {
			keys = selected();
		} else if (hex_elems.hasOwnProperty(id)) {
			keys = [ id ];
		} else {
			keys = [];
		}
		keys.forEach(function (id) {
			hex_elems[id].selected = false;
			hex_elems[id].group.item(0).set({stroke: 'black', opacity: 1});
		});
		canvas.renderAll();
	}

	function add_listeners() {
		var max_space = 10;
		var pressed = false;
		var time = 0;

		function verify_map_pos() {
			var map = canvas.item(0);
			if (map.getTop() > max_space) {
				map.set({ top: max_space });
			} else if (map.getTop() + map.getHeight() < canvas.getHeight() - max_space) {
				map.set({ top: canvas.getHeight() - map.getHeight() - max_space });
			}
			if (map.getLeft() > max_space) {
				map.set({ left: max_space });
			} else if (map.getLeft() + map.getWidth() < canvas.getWidth() - max_space) {
				map.set({ left: canvas.getWidth() - map.getWidth() - max_space });
			}
			canvas.renderAll();
		}

		canvas.on('mouse:down', function (event) {
			time = Date.now();
			pressed = true;
		});
		canvas.on('mouse:up', function (event) {
			pressed = false;
			if (!event.target || Date.now() - time > 250) {
				return;
			}
			var row = Math.floor(((event.e.clientY - event.target.getTop() ) / event.target.scaleY) / y_sep);
			var col = Math.floor(((event.e.clientX - event.target.getLeft()) / event.target.scaleX) / x_sep);

			if (row % 2 === 0) {
				col = 2*Math.floor(col/2);
			} else {
				col = 2*Math.floor((col-1)/2) + 1;
			}

			var hex_id = String.fromCharCode(row+65) + col;
			if (hex_selected(hex_id)) {
				deselect_hex(hex_id);
			} else {
				select_hex(hex_id);
			}
		});
		canvas.on('mouse:move', function (event) {
			if (!pressed) {
				return;
			}
			var map = canvas.item(0);
			map.set({
				top:  map.getTop()  + event.e.movementY,
				left: map.getLeft() + event.e.movementX,
			});

			verify_map_pos();
		});

		document.getElementById('hex-map-parent').addEventListener('wheel', function board_zoomed(event) {
			var map = canvas.item(0);
			var scale = map.scaleX || 1;
			if (event.wheelDeltaY < 0) {
				scale /= 1.1;
			} else if (event.wheelDeltaY > 0) {
				scale *= 1.1;
			}
			map.scale(scale);
			if (map.getHeight() < canvas.getHeight() - 2*max_space) {
				map.scaleToHeight(canvas.getHeight() - 2*max_space);
			}
			if (map.getWidth() < canvas.getWidth() - 2*max_space) {
				map.scaleToWidth(canvas.getWidth() - 2*max_space);
			}

			verify_map_pos();
		});
	}

	domready(function () {
		canvas = new fabric.Canvas('hex-map', {
			backgroundColor: water_clr,
			selection: false,
			height: 650,
			width: 1050,
		});
		var map = create_map(initial_state);

		map.scaleToWidth(canvas.getWidth() - 20);
		map.set({top: 10, left: 10, selectable: false});

		canvas.add(map);
		add_listeners();
	});
})();
