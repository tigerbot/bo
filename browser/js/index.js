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
		business_part1: state.business_part1,

		market: {
			sales:     ko.observableArray([]),
			buy_cnt:   ko.observable(0),
			buy_price: ko.observable(0),
		},
		inventory: {
			scrap:  ko.observableArray([]),
			buy:    ko.observable(0),
			action: ko.observable("build"),
		},
		earnings: {
			dividends: ko.observable(false),
		},
	};
	input_ctrl.confirm = function () {
		if (this.market_phase()) {
			take_market_turn(this.market);
		} else if (this.business_part1()) {
			take_inventory_turn(this.inventory);
		} else {
			take_earnings_turn(this.earnings);
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
			players.select(turn);
			input_ctrl.description(turn + "'s Turn");
			market.buy_cnt(0);
			market.buy_price(0);
			market.sales(held_shares.map(function (info) {
				info.count = ko.observable(0);
				return info;
			}));
		} else {
			hex_map.deselect_hex('*');
			input_ctrl.description(turn +"'s ("+companies.president(turn)+") Turn");
			companies.select(turn);
			if (input_ctrl.business_part1()) {
				var inventory = input_ctrl.inventory;
				inventory.scrap([0,0,0,0,0,0].map(function (count, ind) {
					return {
						level: 'Tech ' + (ind+1),
						count: ko.observable(count),
					};
				}));
				inventory.buy(0);
				inventory.action("build");
			} else {
				input_ctrl.earnings.dividends(false);
			}
		}
	}
	state.market_phase.subscribe(update_options);
	state.business_part1.subscribe(update_options);
	state.turn.subscribe(update_options);
	state.refresh();

	function take_market_turn(market) {
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

	function take_inventory_turn(inventory) {
		var data = {};
		var selected = hex_map.selected();

		if (inventory.action() === "mine") {
			if (selected.length !== 1) {
				alert("Invalid input", ["Exactly one hex can be mined per turn"]);
				return;
			}
			data.mine_coal = selected[0];
		} else {
			data.build_track = selected;
		}
		data.scrap_equipment = inventory.scrap().map(function (info) {
			return parseInt(info.count(), 10);
		});
		data.buy_equipment = parseInt(inventory.buy(), 10);
		data.player_name = companies.president(state.turn());

		var req_obj = {};
		req_obj.url = '/game/business_turn_one';
		req_obj.method = 'POST';
		req_obj.data = data;
		common.request(req_obj, function (errs) {
			if (errs) {
				alert('Inventory Update Failed', errs);
			} else {
				hex_map.deselect_hex('*');
			}
			state.refresh();
			players.refresh();
			companies.refresh();
		});
	}

	function take_earnings_turn(earnings) {
		var data = {
			serviced_cities: hex_map.selected(),
			pay_dividends:   earnings.dividends(),
		};
		data.player_name = companies.president(state.turn());

		var req_obj = {};
		req_obj.url = '/game/business_turn_two';
		req_obj.method = 'POST';
		req_obj.data = data;
		common.request(req_obj, function (errs) {
			if (errs) {
				alert('Earnings Stage Failed', errs);
			} else {
				hex_map.deselect_hex('*');
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
