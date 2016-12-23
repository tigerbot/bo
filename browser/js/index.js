(function () {
	'use strict';
	require('knockout-mapping');

	var common    = require('./common');
	var hex_map   = require('./hex-map');
	var players   = require('./players');
	var companies = require('./companies');

	common.request('/game/state',function (err, repsonse) {
		if (err) {
			console.error('failed to get game global state');
			return;
		}

		hex_map.set_coal(repsonse.unmined_coal);
	});

	window.hex_map   = hex_map;
	window.players   = players;
	window.companies = companies;
})();
