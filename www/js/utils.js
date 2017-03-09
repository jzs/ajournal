_aj = new function() {
	this.get = function(url, callback) {
		var http = new XMLHttpRequest();
		http.open("GET", url, true);
		http.addEventListener("load", function(e) {
			// Completed...
			var data = JSON.parse(http.response);
			if (data.Status != 200) {
				callback(data.Data, data.Error);
			} else {
				callback(data.Data, null);
			}
		});
		http.addEventListener("error", function(e) {
			callback(null, "Call aborted");
		});
		http.addEventListener("abort", function(e) {
			callback(null, "Call aborted");
		});
		http.send();
	};

	this.post = function(url, data, callback) {
		var http = new XMLHttpRequest();
		http.open("POST", url, true);
		http.addEventListener("load", function(e) {
			// Completed...
			var data = JSON.parse(http.response);
			if (data.Status != 200) {
				callback(data.Data, data.Error);
			} else {
				callback(data.Data, null);
			}
		});
		http.addEventListener("error", function(e) {
			callback(null, "Call aborted");
		});
		http.addEventListener("abort", function(e) {
			callback(null, "Call aborted");
		});

		if(typeof(data) != 'undefined' && data != null) {
			http.send(JSON.stringify(data));
		} else {
			http.send(null);
		}
	};

	this.delete = function(url, callback) {
		var http = new XMLHttpRequest();
		http.open("DELETE", url, true);
		http.addEventListener("load", function(e) {
			// Completed...
			var data = JSON.parse(http.response);
			if (data.Status != 200) {
				callback(data.Data, data.Error);
			} else {
				callback(data.Data, null);
			}
		});
		http.addEventListener("error", function(e) {
			callback(null, "Call aborted");
		});
		http.addEventListener("abort", function(e) {
			callback(null, "Call aborted");
		});
		http.send();
	};
};

