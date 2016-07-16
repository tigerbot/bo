(function () {
	'use strict';
	var reqwest = require('reqwest');


	var company_colors = {
		'pennsylvania':                 'rgb(255,  80,  80)',
		'boston_maine':                 'rgb(255, 128, 192)',
		'illinois_central':             'rgb(255, 128,   0)',
		'chesapeake_ohio':              'rgb(248, 248,   7)',
		'new_york_central':             'rgb( 64, 224,  96)',
		'baltimore_ohio':               'rgb( 64, 160, 255)',
		'new_york_chicago_saint_louis': 'rgb(144,  96, 255)',
		'erie':                         'rgb(160,  80,  0)',
		'wabash':                       'rgb(128, 128, 128)',
		'new_york_new_haven_hartford':  'rgb(255, 255, 255)',
	};

	function convert_name(pretty_name) {
		return pretty_name.toLowerCase().replace(/[^a-z ]/g, '').replace(/ +/g, '_');
	}

	function make_req(options, cb) {
		if (typeof cb !== 'function') {
			console.warn('json request has no callback');
			cb = function() {};
		}

		function callback(err, resp) {
			cb(err, resp);
			cb = null;
		}

		if (typeof options === 'string') {
			options = {url: options};
		}

		// make a default timeout of one minute.
		options.timeout = options.timeout || 60*1000;
		options.processData = false;
		options.type = 'json';
		if (typeof options.data === 'object' && !(options.data instanceof FormData)) {
			options.data = JSON.stringify(options.data);
			options.contentType = 'application/json';
		}

		if (options.success || options.error) {
			throw new Error('json request would overwrite cb functions');
		}
		options.error = function (err) {
			callback(err, null);
		};
		options.success = function (resp) {
			callback(null, resp);
		};

		return reqwest(options);
	}

	module.exports.request = make_req;
	module.exports.get_company_color = function (name) {
		return company_colors[convert_name(name)];
	};
	module.exports.get_company_logo = function (name) {
		return '/logos/'+convert_name(name)+'.png';
	};
})();
