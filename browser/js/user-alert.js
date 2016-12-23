(function () {
	'use strict';
	var ko       = require('knockout');
	var domready = require('domready');

	var dialogue_manager = {
		callback:    null,
		active:      ko.observable(false),
		description: ko.observable(""),
		messages:    ko.observableArray([]),
	};

	dialogue_manager.accept = function () {
		this.description("");
		this.messages([]);
		this.callback();
		this.callback = null;
		this.active(false);
	};

	domready(function () {
		ko.applyBindings(dialogue_manager, document.getElementById('alert-mask'));
	});

	function alert_user(description, messages, cb) {
		if (dialogue_manager.active()) {
			throw 'alert box already in use';
		}

		if (typeof cb !== 'function') {
			cb = function () {};
		}
		if (!messages) {
			messages = [];
		}
		if (!Array.isArray(messages)) {
			messages = [messages];
		}

		dialogue_manager.description(description);
		dialogue_manager.messages(messages);
		dialogue_manager.callback = cb;
		dialogue_manager.active(true);
	}

	module.exports = alert_user;
})();
