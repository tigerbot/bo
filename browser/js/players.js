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
	function get_selected() {
		var result = '';
		player_list().some(function (player) {
			if (player.selected()) {
				result = player.name;
				return true;
			}
		});
		return result;
	}

	function get_stocks(name) {
		var result = null;
		player_list().some(function (player) {
			if (player.name === name) {
				result = player.stock_list();
				return true;
			}
		});

		if (result) {
			result = result.map(function (info) {
				return {
					company: info.name,
					color:   info.color,
					count:   info.count,
				};
			});
		}
		return result;
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

	domready(function () {
		ko.applyBindings({players: player_list}, document.getElementById("player-list"));
	});
	function refresh(select_first) {
		common.request('/game/players', function (err, result) {
			if (err) {
				console.error('failed to get initial player info', err);
				return;
			}

			player_list().forEach(function (player) {
				var update = result[player.name];
				ko.mapping.fromJS(update, player);
				player.stock_list(Object.keys(update.stocks).map(function (name) {
					return {
						name:      name,
						color:     common.get_company_color(name),
						count:     update.stocks[name],
					};
				}));
				delete result[player.name];
			});
			Object.keys(result).forEach(function (name) {
				var view_model = ko.mapping.fromJS(result[name]);

				view_model.name = name;
				view_model.stock_list = ko.observableArray(Object.keys(view_model.stocks).map(function (name) {
					return {
						name:      name,
						color:     common.get_company_color(name),
						count:     view_model.stocks[name],
					};
				}));
				view_model.selected = ko.observable(false);
				view_model.select = select_player.bind(null, name);

				player_list.push(view_model);
			});

			sort_players();
			if (select_first) {
				player_list()[0].selected(true);
			}
		});
	}
	refresh(true);

	module.exports.refresh     = refresh;
	module.exports.select      = select_player;
	module.exports.selected    = get_selected;
	module.exports.held_shares = get_stocks;
})();
