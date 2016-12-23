(function () {
	'use strict';
	require('knockout-mapping');
	var ko       = require('knockout');
	var domready = require('domready');

	var common    = require('./common');
	var hex_map   = require('./hex-map');
	var players   = require('./players');
	var companies = require('./companies');

	var game_state = {
		round:      ko.observable(0),
		phase:      ko.observable(""),

		tech_level:    ko.observable(1),
		trains_bought: ko.observable(0),

		orphan_stocks: ko.observableArray([]),
	};

	domready(function () {
		ko.applyBindings(game_state, document.getElementById('game-state'));
	});
	common.request('/game/state',function (err, repsonse) {
		if (err) {
			console.error('failed to get game global state');
			return;
		}

		game_state.round(repsonse.round);
		game_state.phase(repsonse.phase);
		game_state.tech_level(repsonse.tech_level);
		game_state.trains_bought(repsonse.trains_bought);
		Object.keys(repsonse.orphan_stocks).forEach(function (name) {
			game_state.orphan_stocks.push({
				name:  name,
				color: common.get_company_color(name),
				count: repsonse.orphan_stocks[name],
			});
		});

		hex_map.set_coal(repsonse.unmined_coal);
	});

	window.hex_map   = hex_map;
	window.players   = players;
	window.companies = companies;
})();
