(function () {
	'use strict';
	var ko       = require('knockout');
	var domready = require('domready');

	var common  = require('./common');
	var hex_map = require('./hex-map');

	var company_list = ko.observableArray([]);

	function select_company(name) {
		company_list().forEach(function (view_model) {
			view_model.selected(name === view_model.name);
		});
	}
	function get_selected() {
		var result = '';
		company_list().some(function (company) {
			if (company.selected()) {
				result = company.name;
				return true;
			}
		});
		return result;
	}

	function sort_companies() {
		company_list.sort(function (left, right) {
			var price_diff = right.stock_price() - left.stock_price();
			if (price_diff !== 0) {
				return price_diff;
			}
			if (right.price_changed() < left.price_changed()) {
				return 1;
			} else {
				return -1;
			}
		});
	}

	function get_president(name) {
		var result = '';
		company_list().some(function (company) {
			if (company.name === name) {
				result = company.president();
				return true;
			}
		});
		return result;
	}

	function convert_equipment() {
		// jshint validthis:true
		var result = this.equipment();

		result = result.map(function (count, index) {
			var level = index + 1;
			return {
				name:     'Tech ' + level,
				count:    count,
				capacity: count * level,
			};
		});

		result = result.filter(function (obj) {
			return obj.count > 0;
		});

		result.push(result.reduce(function (total, next) {
			total.count    += next.count;
			total.capacity += next.capacity;
			return total;
		}, {name: 'Total', count: 0, capacity: 0}));

		return result;
	}

	domready(function () {
		ko.applyBindings({companies: company_list}, document.getElementById("company-list"));
	});

	function refresh(first_time) {
		common.request('/game/companies', function (err, result) {
			if (err) {
				console.error('failed to get initial company info', err);
				return;
			}

			company_list().forEach(function (company) {
				var already_built = company.built_track();
				result[company.name].built_track.forEach(function (id) {
					if (already_built.indexOf(id) < 0) {
						hex_map.build_rail(id, company.color);
					}
				});
				ko.mapping.fromJS(result[company.name], company);
				delete result[company.name];
			});
			Object.keys(result).forEach(function (name) {
				var view_model = ko.mapping.fromJS(result[name]);

				view_model.name = name;
				view_model.icon = common.get_company_logo(name);
				view_model.color = common.get_company_color(name);
				view_model.selected = ko.observable(false);
				view_model.select = select_company.bind(null, name);
				view_model.equipment_list = ko.computed(convert_equipment, view_model);
				view_model.built_track().forEach(function (id) {
					hex_map.build_rail(id, view_model.color);
				});

				company_list.push(view_model);
			});
			sort_companies();
			if (first_time) {
				company_list()[0].selected(true);
			}
		});
	}
	refresh(true);

	module.exports.refresh   = refresh;
	module.exports.selected  = get_selected;
	module.exports.president = get_president;
})();
