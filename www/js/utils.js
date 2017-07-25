_aj = new function() {

	var send = function(url, type, data, callback) {
		var http = new XMLHttpRequest();
		http.open(type, url, true);
		http.addEventListener("load", function(e) {
			// Completed...
			try {
				var data = JSON.parse(http.response);
				if (data.Status != 200) {
					callback(data.Data, data.Error);
				} else {
					callback(data.Data, null);
				}
				return
			} catch (e) {
				callback(null, {"error":"boom"});
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

	this.get = function(url, callback) {
		send(url, "GET", null, callback);
	};

	this.post = function(url, data, callback) {
		send(url, "POST", data, callback);
	};

	this.put = function(url, data, callback) {
		send(url, "PUT", data, callback);
	};

	this.delete = function(url, callback) {
		send(url, "DELETE", null, callback);
	};
};

