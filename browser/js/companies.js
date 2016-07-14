(function () {
	'use strict';
	var ko       = require('knockout');
	var domready = require('domready');
	var request  = require('./json-request');

	var company_colors = {
		'pennsylvania':                 'rgb(255,  48,  48)',
		'boston_maine':                 'rgb(255, 128, 192)',
		'illinois_central':             'rgb(255, 128,   0)',
		'chesapeake_ohio':              'rgb(248, 248,   7)',
		'new_york_central':             'rgb( 64, 224,  96)',
		'baltimore_ohio':               'rgb(  0, 128, 255)',
		'new_york_chicago_saint_louis': 'rgb(128,  64, 255)',
		'erie':                         'rgb(160,  80,  0)',
		'wabash':                       'rgb(100, 100, 100)',
		'new_york_new_haven_hartford':  'rgb(255, 255, 255)',
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

			company_list.push(view_model);
		});
	});
	domready(function () {
		ko.applyBindings({companies: company_list}, document.getElementById("company-list"));
	});
})();
