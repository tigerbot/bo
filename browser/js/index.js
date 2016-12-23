(function () {
	'use strict';
	require('knockout-mapping');

	var alert     = require('./user-alert');
	var state     = require('./game-state');
	var hex_map   = require('./hex-map');
	var players   = require('./players');
	var companies = require('./companies');

	var market_regexp   = /^market/i;
	var business_regexp = /^business/i;

	function check_action(err) {
		if (err) {
			console.error('failed to determine what turn needs to be taken', err);
		}

		if (market_regexp.test(state.phase)) {
			alert("Market phase on "+state.turn+"'s turn");
		} else if (business_regexp.test(state.phase)) {
			alert("Business phase for company "+state.turn);
		}
	}

	function next_turn() {
		state.refresh(check_action);
	}

	next_turn();
})();
