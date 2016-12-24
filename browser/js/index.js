(function () {
	'use strict';
	require('knockout-mapping');
	var ko       = require('knockout');
	var domready = require('domready');

	var alert     = require('./user-alert');
	var common    = require('./common');
	var state     = require('./game-state');
	var hex_map   = require('./hex-map');
	var players   = require('./players');
	var companies = require('./companies');

	var input_ctrl = {
		description:    ko.observable("Initializing"),
		market_phase:   state.market_phase,

		market: {
			type:      ko.observable("pass"),
			sales:     ko.observableArray([]),
			buy_cnt:   ko.observable(0),
			buy_price: ko.observable(0),
		},
	};
	input_ctrl.confirm = function () {
		if (this.market_phase()) {
			take_market_turn(this.market);
		} else if (this.business_part1()) {
			console.log("inventory control not implemented");
		} else {
			console.log("income control not implemented");
		}
	};

	function update_options() {
		var turn   = state.turn();
		// make sure the state information has been initialized.
		if (!turn) {
			return;
		}

		if (state.market_phase()) {
			var market = input_ctrl.market;
			var held_shares = players.held_shares(turn);
			if (!held_shares) {
				console.warn('held shares for "'+turn+'" not present yet');
				setTimeout(update_options, 1000);
				return;
			}
			input_ctrl.description(turn + "'s Turn");
			market.type("player");
			market.buy_cnt(0);
			market.sales(held_shares.map(function (info) {
				info.count = ko.observable(0);
				return info;
			}));
		}
	}
	state.market_phase.subscribe(update_options);
	state.turn.subscribe(update_options);
	state.refresh();

	function take_market_turn(market) {
		if (market.type() === "pass") {
			market.sales([]);
			market.buy_cnt(0);
		}

		var data = {};
		data.sales = market.sales().filter(function (sale) {
			return sale.count() > 0;
		}).map(function (sale) {
			return {
				company: sale.name,
				count:   parseInt(sale.count(), 10),
			};
		});
		if (market.buy_cnt() > 0) {
			data.purchase = {
				company: companies.selected(),
				count:   parseInt(market.buy_cnt(), 10),
				price:   parseInt(market.buy_price(), 10),
			};
		}
		data.player_name = state.turn();

		var req_obj = {};
		req_obj.url = '/game/market_turn';
		req_obj.method = 'POST';
		req_obj.data = data;
		common.request(req_obj, function (errs) {
			if (errs) {
				alert('Market Turn Failed', errs);
			}
			state.refresh();
			players.refresh();
			companies.refresh();
		});
	}

	domready(function () {
		ko.applyBindings(input_ctrl, document.getElementById("user-input"));
	});
})();
