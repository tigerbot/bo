(function () {
	'use strict';
	var reqwest = require('reqwest');

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

	module.exports = make_req;
})();
