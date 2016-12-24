(function () {
	'use strict';
	var ko       = require('knockout');
	var domready = require('domready');

	var common    = require('./common');
	var hex_map   = require('./hex-map');

	var market_regexp   = /^market/i;
	var game_state = {
		round: ko.observable(0),
		phase: ko.observable(""),
		turn:  ko.observable(""),
		stage: ko.observable(""),

		tech_level:    ko.observable(1),
		trains_bought: ko.observable(0),
		all_costs:     ko.observableArray([]),

		orphan_stocks: ko.observableArray([]),
	};
	game_state.train_costs = ko.computed(function () {
		var tech_level = this.tech_level();
		var all_costs  = this.all_costs();
		var result = [];

		if (tech_level <= all_costs.length) {
			result.push(all_costs[tech_level-1]);
		}
		if (tech_level < all_costs.length) {
			result.push(all_costs[tech_level]);
		}

		return result;
	}, game_state);

	domready(function () {
		ko.applyBindings(game_state, document.getElementById('game-state'));
	});

	common.request('/train_costs', function (err, repsonse) {
		if (err) {
			console.error('failed to get equipment costs', err);
			return;
		}
		game_state.all_costs(repsonse.map(function (costs, lvl) {
			return {
				name:  'Tech ' + (lvl+1),
				costs: costs.map(function (value, num) {
					return {
						value:         value,
						number:        lvl*costs.length + num,
						trains_bought: game_state.trains_bought,
					};
				}),
			};
		}));
	});

	function log_error(err) {
		if (err) {
			console.error('failed to get game state', err);
		}
	}

	function refresh(cb) {
		if (typeof cb !== 'function') {
			cb = log_error;
		}

		common.request('/game/state', function (err, repsonse) {
			if (err) {
				return cb(err);
			}

			game_state.round(repsonse.round);
			game_state.phase(repsonse.phase);
			game_state.turn(repsonse.turn);
			game_state.stage(repsonse.stage);

			game_state.tech_level(repsonse.tech_level);
			game_state.trains_bought(repsonse.trains_bought);

			game_state.orphan_stocks(Object.keys(repsonse.orphan_stocks).map(function (name) {
				return {
					name:  name,
					color: common.get_company_color(name),
					count: repsonse.orphan_stocks[name],
				};
			}));

			hex_map.set_coal(repsonse.unmined_coal);

			cb();
		});
	}

	module.exports.refresh = refresh;
	module.exports.turn = ko.computed(function () {
		return this.turn();
	}, game_state);
	module.exports.market_phase = ko.computed(function () {
		return market_regexp.test(this.phase());
	}, game_state);
	module.exports.business_part1 = ko.computed(function () {
		return this.stage() === "inventory";
	}, game_state);
})();
