(function () {
	'use strict';
	var ko       = require('knockout');
	var domready = require('domready');
	var common   = require('./common');

	var company_list = ko.observableArray();

	function select_company(name) {
		company_list().forEach(function (view_model) {
			view_model.selected(name === view_model.name);
		});
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

	common.request('company_info', function (err, result) {
		if (err) {
			console.error('failed to get initial company info', err);
			return;
		}

		Object.keys(result).forEach(function (name) {
			var view_model = ko.mapping.fromJS(result[name]);

			view_model.name = name;
			view_model.icon = common.get_company_logo(name);
			view_model.color = common.get_company_color(name);
			view_model.selected = ko.observable(false);
			view_model.select = select_company.bind(null, name);

			company_list.push(view_model);
		});

		sort_companies();
		company_list()[0].selected(true);
	});
	domready(function () {
		ko.applyBindings({companies: company_list}, document.getElementById("company-list"));
	});
})();
