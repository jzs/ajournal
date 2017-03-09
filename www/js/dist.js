riot.tag2('content', '<page-dash if="{loggedin && dash}"></page-dash> <page-login if="{loggedout}"></page-login> <page-journal-create if="{loggedin && journalcreate}"></page-journal-create> <page-journal if="{loggedin && journal}" journalid="{journalid}"></page-journal> <page-entryeditor if="{loggedin && entry}" entryid="{entryid}"></page-entryeditor>', '', '', function(opts) {
var self = this;
self.loggedin = false;
self.loggedout = !self.loggedin;
self.journalcreate = false;
self.dash = false;
self.entry = false

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

self.clear = function() {
	self.dash = false;
	self.journalcreate = false;
	self.journal = false;
	self.entry = false
}

route(function(collection, id, method, mid) {
	self.clear()
	switch(collection) {
		case 'journals':
			if(id == 'create') {

				self.journalcreate = true;
			} else {
				if(method == 'entries') {
					if(mid == 'create') {
					} else {

						self.entryid = id;
						self.entry = true;
					}

				} else {

					self.journalid = id;
					self.journal = true;
				}
			}
			break;
		default:
			self.clear()
			self.dash = true;
			self.update();
			break;
	}
	self.update();
});
});

riot.tag2('navbar', '<nav class="nav"> <div class="nav-left"> </div> <div class="nav-center"> <a class="nav-item" href="#/">a-Journal</a> </div> <div class="nav-right"> <a class="nav-item" href="#" onclick="{logout}" if="{user}">Log out</a> </div> </nav>', '', '', function(opts) {
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


riot.tag2('page-dash', '<div class="section"> <div class="container"> <div class="columns"> <div class="column"> <h3 class="title">My Journals</h3> <section class="section"> <p> Welcome to a-Journal. <span if="{journals.length < 1 && !loading}">It looks like you haven\'t created any journals yet. Let me help you get started.</span> </p> </section> <span class="help is-danger" if="{err}">{err}</span> <div class="box" each="{journal in journals}"> <article class="media"> <div class="media-left"> <figure class="image is-64x64"> <img src="images/128x128.png" alt="Image"> </figure> </div> <div class="media-content"> <div class="content" onclick="{tojournal}" style="cursor:pointer;"> <p> <strong>{journal.Title}</strong> <small>@jzs</small> <small>31m</small> <br> {journal.Description} </p> </div> <nav class="level"> <div class="level-left"> <a class="level-item" onclick="{createentry}"> <span class="icon is-small"><i class="fa fa-plus"></i></span> </a> </div> <div class="level-right"> <span class="level-item" if="{!journal.Public}"> Private </span> </div> </nav> </div> </article> </div> <button class="button" onclick="{newjournal}">New Journal</button> </div> </div> </div> </div>', '', '', function(opts) {
var self = this;
self.journals = [];
self.loading = true;
self.err = null;

self.on('mount', function() {

	self.loading = true;
	_aj.get("/api/journals", function(data, err) {
		if(err != null) {
			self.err = err;
			self.loading = false;
		} else {
			self.journals = data;
			self.loading = false;
		}
		self.update();
	});
});

self.newjournal = function() {
	route("/journals/create");
};

self.tojournal = function(e) {
	route("/journals/" + e.item.journal.ID);
};
self.createentry = function(e) {
	e.preventDefault();
	route("/journals/" + e.item.journal.ID + "/entries/create");
}
});

riot.tag2('page-entryeditor', '<div class="section"> <div class="container"> <div class="columns"> <div class="column"> <pre style="min-height:200px" contenteditable="true" onkeyup="{contentchange}">{entry.Content}</pre> <p> <button class="button">Publish</button> <button class="button">Save</button> </p> </div> <div class="column"> <raw class="markdown" content="{preview}"></raw> </div> </div> </div> </div>', '', '', function(opts) {
var self = this;
self.preview = "";
var converter = new showdown.Converter();

self.editContent = "";
self.entry = {
	Content: "#Title\n\nThis is a paragraph"
}

self.on('mount', function() {
	self.preview = converter.makeHtml(self.entry.Content);
	self.update();
});

self.contentchange = function(e) {
	self.editContent = e.target.innerText;
	self.preview = converter.makeHtml(self.editContent);
	self.update();
}
});

riot.tag2('page-journal-create', '<section class="section"> <div class="container"> <h3 class="title">New Journal</h3> <label class="label">Title</label> <p class="control"> <input class="input" type="text" placeholder="Title" onkeyup="{onTitle}" value=""> </p> <label class="label">Description</label> <p class="control"> <textarea class="textarea" placeholder="Description" onkeyup="{onDescription}"></textarea> <span if="{errmsg}" class="help is-danger">{errmsg}</span> </p> <p class="control"> <label class="checkbox"> <input type="checkbox" checked="{journal.Public}" onchange="{onPublic}"> Public </label> </p> <p class="control"> <label class="label">Tags</label> <p class="control has-addons"> <input class="input" type="text" placeholder="Tagname" onkeyup="{onjournaltag}"> <a class="button is-info" onclick="{addJournalTag}"> Add </a> </p> <span class="tag is-large" each="{t in journal.Tags}"> {t} <button class="delete" onclick="{deleteTag}"></button> </span> </p> <button class="button is-success is-pulled-right" onclick="{create}">Create</button> </div> </section>', '', '', function(opts) {
var self = this;
self.errmsg = null;

self.journal = {
	Tags: [],
	Title: "",
	Description: "",
	Public: false
};

self.onTitle = function(e) {
	self.journal.Title = e.target.value;
};
self.onDescription = function(e) {
	self.journal.Description = e.target.value;
};
self.onPublic = function(e) {
	self.journal.Public = e.target.checked;
};

self.journaltag = "";
self.onjournaltag = function(e) {
	e.preventDefault();
	self.journaltag = e.target.value;
};
self.addJournalTag = function() {
	self.journal.Tags.push(self.journaltag);
	self.journaltag = "";
	self.update();
};
self.deleteTag = function(e) {
	var index = self.journal.Tags.indexOf(e.item.t);
	self.journal.Tags.splice(index, 1);
}

self.create = function() {
	_aj.post("/api/journals", self.journal, function(data,err) {
		if( err != null ) {
			self.errmsg = err;
			self.update();
			return;
		}
		route("/journals/" + data.ID);
	});
};
});

riot.tag2('page-journal', '<section class="section"> <div class="container"> <section class="section"> <h3 class="title">Journal: {journal.Title}</h3> <p> {journal.Description} </p> </section> <section class="section"> <div class="box" each="{entry in journal.Entries}" onclick="{onentry}" style="cursor:pointer;"> <article class="media"> <div class="media-content"> <div class="content"> <p> <strong>{entry.Title}</strong> <small>@jzs</small> <small>31m</small> <br> <pre>{entry.Content.substring(0, 200)}...</pre> <br> <span each="{tag in parent.Tags}">{tag}</span> </p> </div> <nav class="level"> <div class="level-left"> <a class="level-item"> <span class="icon is-small"><i class="fa fa-reply"></i></span> </a> <a class="level-item"> <span class="icon is-small"><i class="fa fa-retweet"></i></span> </a> <a class="level-item"> <span class="icon is-small"><i class="fa fa-heart"></i></span> </a> </div> </nav> </div> </article> </div> </section> </div> </section>', '', '', function(opts) {
var self = this;
self.journal = {
	Title: "Journal title",
	Description: "A fine description of a fine wine",
	Tags: [],
	Entries: [{
		ID: 1,
		Date: "",
		Title: "My first entry",
		Content: "#My first content\n\nHello world",
		Tags: ["diary"],
		Created: "",
		IsPublised: false
	}],
	Created: ""
};

self.on('mount', function() {
	_aj.get("/api/journals/" + opts.journalid, function(data, err) {
		if(err != nil) {
			self.err = err;
			self.update();
			return;
		}
		self.journal = data;
		self.update();
	});
});

self.onentry = function(e) {
	route("/journals/"+opts.journalid+"/entries/"+e.item.entry.ID);
};
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
			self.update();
			return;
		} else {
			RiotControl.trigger('perform-login', data.Data);
		}
		self.loggingin = false;
		self.update();
		route("/");
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

riot.tag2('raw', '<span></span>', '', '', function(opts) {
this.root.innerHTML = opts.content

this.on('update', function() {
	this.root.innerHTML = opts.content
});
});
