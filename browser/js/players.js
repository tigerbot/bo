(function () {
	'use strict';
	var ko       = require('knockout');
	var domready = require('domready');
	var common   = require('./common');

	var player_list = ko.observableArray();

	function select_player(name) {
		player_list().forEach(function (view_model) {
			view_model.selected(name === view_model.name);
		});
	}

	function sort_players() {
		player_list.sort(function (left, right) {
			var cash_diff = right.cash() - left.cash();
			if (cash_diff !== 0) {
				return cash_diff;
			} else {
				return right.net_worth() - left.net_worth();
			}
		});
	}

	common.request('player_info', function (err, result) {
		if (err) {
			console.error('failed to get initial player info', err);
			return;
		}

		Object.keys(result).forEach(function (name) {
			var view_model = ko.mapping.fromJS(result[name]);

			view_model.name = name;
			view_model.stock_list = ko.observableArray(Object.keys(view_model.stocks).map(function (name) {
				return {
					name:      name,
					color:     common.get_company_color(name),
					count:     view_model.stocks[name].count,
					president: view_model.stocks[name].president,
				};
			}));
			view_model.selected = ko.observable(false);
			view_model.select = select_player.bind(null, name);

			player_list.push(view_model);
		});

		sort_players();
		player_list()[0].selected(true);
	});
	domready(function () {
		ko.applyBindings({players: player_list}, document.getElementById("player-list"));
	});
})();
