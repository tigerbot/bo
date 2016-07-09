(function () {
	'use strict';
	var domready = require('domready');
	var fabric   = require('fabric').fabric;

	var hex_map  = require('./hex-map');

	domready(function () {
		var canvas = new fabric.Canvas('main', {
			backgroundColor: '#DEEFF7',
			height: 650,
			width: 1050,
		});

		var map = hex_map.create();
		map.scaleToWidth(canvas.getWidth() - 20);
		map.set({top: 10, left: 10});

		canvas.add(map);
		window.canvas = canvas;
	});
})();
