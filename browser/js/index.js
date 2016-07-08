(function () {
	'use strict';
	var domready = require('domready');
	var fabric   = require('fabric').fabric;

	var hex_radius = 60;
	var hex_pnts = [0, 1, 2, 3, 4, 5].map(function (index) {
		var angle = (2*index + 1)*Math.PI/6;
		return {
			x: hex_radius * (1 + Math.cos(angle)),
			y: hex_radius * (1 + Math.sin(angle)),
		};
	});
	var hex = new fabric.Polygon(hex_pnts, {
		left: 100,
		top:  100,
		fill: 'green',
		stroke: 'black',
		strokeWidth: 5,
	});

	domready(function () {
		var canvas = new fabric.Canvas('main', {
			backgroundColor: 'rgb(200, 200, 255)',
			height: 1000,
			width:  1000,
		});
		canvas.add(hex);

		window.canvas = canvas;
	});
})();