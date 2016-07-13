(function () {
	'use strict';
	var ko       = require('knockout');
	var domready = require('domready');
	var request  = require('./json-request');

	var company_colors = {
		'pennsylvania':                 'red',
		'boston_maine':                 'pink',
		'illinois_central':             'orange',
		'chesapeake_ohio':              'yellow',
		'new_york_central':             'green',
		'baltimore_ohio':               'blue',
		'new_york_chicago_saint_louis': 'purple',
		'erie':                         'brown',
		'wabash':                       'grey',
		'new_york_new_haven_hartford':  'white',
	};
	var company_list = ko.observableArray();

	request('company_info', function (err, result) {
		if (err) {
			console.error('failed to get initial company info', err);
			return;
		}

		Object.keys(result).forEach(function (name) {
			var view_model = ko.mapping.fromJS(result[name]);
			var safe_name = name.toLowerCase().replace(/[^a-z ]/g, '').replace(/ +/g, '_');

			view_model.name = name;
			view_model.icon = '/logos/'+safe_name+'.png';
			view_model.color = company_colors[safe_name];

			view_model._president = ko.computed(function () {
				return view_model.president.name() + ' (' + view_model.president.shares() + ')';
			}, view_model);

			company_list.push(view_model);
		});
	});
	domready(function () {
		ko.applyBindings({companies: company_list}, document.getElementById("company-list"));
	});
})();
