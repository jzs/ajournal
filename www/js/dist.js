riot.tag2('content', '<page-dash if="{loggedin}"></page-dash> <page-login if="{loggedout}"></page-login>', '', '', function(opts) {
var self = this;
self.loggedin = false;
self.loggedout = !self.loggedin;

RiotControl.on('logout', function() {
	self.loggedin = false;
	self.loggedout = !self.loggedin;

	self.update();
});
RiotControl.on('login', function() {
	self.loggedin = true;
	self.loggedout = !self.loggedin;
	self.update();
});

});

riot.tag2('navbar', '<nav class="nav"> <div class="nav-left"> </div> <div class="nav-center"> <a class="nav-item" href="#/">Journal</a> </div> <div class="nav-right"> <a class="nav-item" href="#" onclick="{logout}" if="{user}">Log out</a> </div> </nav>', '', '', function(opts) {
var self = this;
self.user = null;

self.logout = function(e) {
	e.preventDefault();
	RiotControl.trigger('perform-logout', null);
};

RiotControl.on('login', function(user) {
	self.user = user;
	self.update();
});
RiotControl.on('logout', function() {
	self.user = null;
	self.update();
});
});


riot.tag2('page-dash', '<div class="section"> <div class="container"> <div class="columns"> <div class="column"> <h3 class="title">My Journals</h3> <p> Welcome to journal... </p> <div class="box" each="{journal in journals}"> <article class="media"> <div class="media-left"> <figure class="image is-64x64"> <img src="images/128x128.png" alt="Image"> </figure> </div> <div class="media-content"> <div class="content" onclick="{}" style="cursor:pointer;"> <p> <strong>{journal.Title}</strong> <small>@jzs</small> <small>31m</small> <br> Description here... </p> </div> <nav class="level"> <div class="level-left"> <a class="level-item"> <span class="icon is-small"><i class="fa fa-plus"></i></span> </a> </div> </nav> </div> </article> </div> <button class="button">New Journal</button> </div> </div> </div> </div>', '', '', function(opts) {
var self = this;
self.journals = [{Title: "A new beginning"},{Title: "A journey abroad"}];
});



riot.tag2('page-login', '<section class="hero is-fullheight is-primary"> <div class="hero-head"></div> <div class="hero-body"> <div class="container"> <div class="columns"> <div class="column is-half is-offset-one-quarter has-text-left"> <h3 class="title">Journal</h3> <div class="card"> <div class="card-content is-clearfix"> <label class="label">Username</label> <p class="control has-icon has-icon-right"> <input class="input" type="text" placeholder="Username" onkeyup="{onusername}" value=""> <span class="icon is-small hidden"> <i class="fa fa-check"></i> </span> <span class="hidden help is-success">This username is available</span> </p> <label class="label">Password</label> <p class="control has-icon has-icon-right"> <input class="input" type="password" placeholder="Password" value="" onkeyup="{onpassword}"> <span class="icon is-small hidden"> <i class="fa fa-warning"></i> </span> <span if="{loginerr}" class="help is-danger">{errmsg}</span> </p> <div class="control is-grouped is-pulled-right"> <p class control> <button class="button is-link {is-disabled : loggingin}" onclick="{register}">Register</button> </p> <p class="control"> <button class="button is-success {is-loading : loggingin}" onclick="{login}">Login</button> </p> </div> </div> </div> </div> </div> </div> </div> <div class="hero-foot"></div> </section>', '', '', function(opts) {
var self = this;

self.loggingin = false;
self.loginerr = false;
self.errmsg = "";

self.username = "";
self.onusername = function(e) {
	self.username = e.target.value;
};

self.password = "";
self.onpassword = function(e) {
	self.password = e.target.value;
};

RiotControl.on('logout', function() {
	self.update();
});
RiotControl.on('login', function() {
	self.update();
});

self.login = function() {

	var http = new XMLHttpRequest();
	http.open("POST", "/api/users/login", true);

	http.addEventListener("load", function(e) {

		var data = JSON.parse(http.response);
		if(data.Status != 200) {
			self.errmsg = data.Error;
			self.loginerr = true;
		} else {
			RiotControl.trigger('perform-login', data.Data);
		}

		self.loggingin = false;
		self.update();
	});
	http.addEventListener("error", function(e) {
	});
	http.addEventListener("abort", function(e) {
	});

	self.loggingin = true;
	self.update();

	var user = {
		Username: self.username,
		Password: self.password
	};
	http.send(JSON.stringify(user));
};

self.register = function(e) {
	var http = new XMLHttpRequest();
	http.open("POST", "/api/users", true);

	http.addEventListener("load", function(e) {

		var data = JSON.parse(http.response);
		if (data.Status != 200) {
			self.errmsg = data.Error;
			self.loginerr = true;
		} else {
			self.login();
		}
		self.update();
	});
	http.addEventListener("error", function(e) {
	});
	http.addEventListener("abort", function(e) {
	});

	var user = {
		Username: self.username,
		Password: self.password
	};
	http.send(JSON.stringify(user));
};

});
