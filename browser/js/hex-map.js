(function () {
	'use strict';
	var domready = require('domready');
	var fabric   = require('fabric').fabric;
	var Join     = require('join').Join;
	var request  = require('./json-request');

	var canvas;
	var x_sep, y_sep;
	var water_clr = '#DEEFF7';
	var grass_clr = '#688E45';
	var hex_elems = {};

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

	function create_map(hex_background, map_state) {
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
			fill:    hex_background,
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
			originY:   'top',
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
		var rails_opts = {
			originX: 'center',
			originY: 'top',
			top:     pnts[1].y - pnts[0].y,
			left:    radius,
		};

		// we divide the width of the pnts by 2 because the IDs only use every other number.
		x_sep = (pnts[0].x - pnts[3].x + border_width)/2;
		y_sep = 1.5 * radius;
		return new fabric.Group(Object.keys(map_state).map(function (id) {
			var details = map_state[id];
			var items = [];
			var coal = null;
			var rails;

			items.push(new fabric.Polygon(pnts, hex_opts));
			items.push(new fabric.Text('$'+details.price, price_opts));
			if (details.coal) {
				coal = new fabric.Text('COAL', coal_opts);
				items.push(coal);
			}
			if (details.city) {
				items.push(create_rev_bar(details.city.revenue).scaleToWidth(1.85*x_sep).set(rev_bar_opts));
				items.push((new fabric.Text(details.city.name, city_opts)).set({fill: details.city.color}));
			}
			rails = new fabric.Group([], rails_opts);
			items.push(rails);

			var group = new fabric.Group(items, {
				top:  y_sep*(id.charCodeAt(0) - 65),
				left: x_sep*parseInt(id.slice(1)),
			});

			hex_elems[id] = {
				group:    group,
				coal:     coal,
				rails:    rails,
				selected: false,
			};

			return group;
		}), {selectable: false, hoverCursor: 'auto'});
	}

	function hex_selected(id) {
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
	function select_hex(id) {
		if (hex_elems.hasOwnProperty(id)) {
			hex_elems[id].selected = true;
			hex_elems[id].group.item(0).set({stroke: 'yellow', opacity: 0.8});
		}
		canvas.renderAll();
	}
	function deselect_hex(id) {
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

	function has_coal(id) {
		if (!hex_elems.hasOwnProperty(id)) {
			return false;
		}
		return hex_elems[id].coal !== null;
	}
	function mine_coal(id) {
		if (!hex_elems.hasOwnProperty(id)) {
			return false;
		} else if (hex_elems[id].coal === null) {
			return false;
		}

		hex_elems[id].group.remove(hex_elems[id].coal);
		hex_elems[id].coal = null;
		canvas.renderAll();
		return true;
	}

	function build_rail(id, color) {
		if (!hex_elems.hasOwnProperty(id)) {
			return;
		}
		var rails = hex_elems[id].rails;
		rails.add(new fabric.Rect({
			height: 15,
			width:  15,
			fill:   color,
		}));

		var total_cnt = rails.size();
		var row_cnt, row_ind, col_ind;
		for (col_ind = 0; total_cnt > 0; col_ind += 1) {
			row_cnt = total_cnt;
			if (row_cnt > 3) {
				row_cnt = 3;
			}
			for (row_ind = 0; row_ind < row_cnt; row_ind += 1) {
				rails.item(3*col_ind + row_ind).set({
					originX: 'center',
					originY: 'top',
					top:  20*col_ind,
					left: 20*row_ind - 10*(row_cnt-1),
				});
			}
			total_cnt -= row_cnt;
		}
		canvas.renderAll();
	}

	function add_listeners() {
		var parent_elem = document.getElementById('hex-map-parent');
		var map = canvas.item(0);
		var max_space = 10;
		var pressed = false;
		var time = 0;

		function verify_map_pos() {
			var extra_height = canvas.getHeight() - map.getHeight();
			var extra_width  = canvas.getWidth()  - map.getWidth();
			if (extra_width > 2*max_space && extra_height > 2*max_space) {
				if (canvas.getWidth()/map.getWidth() < canvas.getHeight()/map.getHeight()) {
					map.scaleToWidth(canvas.getWidth() - 2*max_space);
				} else {
					map.scaleToHeight(canvas.getHeight() - 2*max_space);
				}
			}

			extra_height = canvas.getHeight() - map.getHeight();
			if (extra_height > 0) {
				map.set({ top: extra_height/2 });
			} else if (map.getTop() >= max_space) {
				map.set({ top: max_space });
			} else if (map.getTop() + map.getHeight() < canvas.getHeight() - max_space) {
				map.set({ top: canvas.getHeight() - map.getHeight() - max_space });
			}
			extra_width  = canvas.getWidth()  - map.getWidth();
			if (extra_width > 0) {
				map.set({ left: extra_width/2 });
			} else if (map.getLeft() >= max_space) {
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
			if (event.target !== map || Date.now() - time > 250) {
				return;
			}
			var row = Math.floor(((event.e.layerY - event.target.getTop() ) / event.target.scaleY) / y_sep);
			var col = Math.floor(((event.e.layerX - event.target.getLeft()) / event.target.scaleX) / x_sep);

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
			map.set({
				top:  map.getTop()  + event.e.movementY,
				left: map.getLeft() + event.e.movementX,
			});

			verify_map_pos();
		});

		parent_elem.onwheel = function (event) {
			var scale = map.scaleX || 1;
			if (event.wheelDeltaY < 0) {
				scale /= 1.1;
			} else if (event.wheelDeltaY > 0) {
				scale *= 1.1;
			}
			map.scale(scale);
			verify_map_pos();
		};

		var resize_canvas = (function () {
			var timeout_id;
			function resize() {
				timeout_id = null;
				canvas.setWidth(parent_elem.offsetWidth);
				canvas.setHeight(parent_elem.offsetHeight);
				verify_map_pos();
			}

			return function () {
				if (timeout_id) {
					clearTimeout(timeout_id);
				}
				timeout_id = setTimeout(resize, 50);
			};
		})();
		window.onresize = resize_canvas;

		// Set the scale really small so that after the canvas is sized to be exactly the same as the
		// parent element it will be expanded to fit the entire area by verify_map_pos.
		map.scale(0.01);
		resize_canvas();
	}

	(function () {
		var join = new Join();
		fabric.util.loadImage('hex_background.png', join.add());
		request('map_state', join.add());

		join.then(function (img_args, map_args) {
			if (map_args[0]) {
				console.error('failed to get map state', map_args[0]);
				return;
			}
			var pattern = new fabric.Pattern({ source: img_args[0], repeat: 'repeat' });
			var map = create_map(pattern, map_args[1]);

			domready(function () {
				canvas = new fabric.Canvas('hex-map', {
					backgroundColor: water_clr,
					selection: false,
					height: 650,
					width: 1050,
				});
				canvas.add(map);
				add_listeners();
			});
		});
	})();

	module.exports.hex_selected = hex_selected;
	module.exports.select_hex   = select_hex;
	module.exports.deselect_hex = deselect_hex;

	module.exports.has_coal   = has_coal;
	module.exports.mine_coal  = mine_coal;
	module.exports.build_rail = build_rail;
})();
