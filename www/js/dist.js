riot.tag2('content', '<page-dash if="{loggedin && dash}"></page-dash> <page-login if="{loggedout}"></page-login> <page-journal-create if="{loggedin && journalcreate}"></page-journal-create> <page-journal if="{loggedin && journal}" journalid="{journalid}"></page-journal> <page-entryeditor if="{loggedin && entry}" journalid="{journalid}" entryid="{entryid}"></page-entryeditor> <page-viewjournal if="{viewjournal}" journalid="{journalid}"></page-viewjournal> <page-viewjournals if="{viewjournals}" userid="{userid}"></page-viewjournals> <page-profile if="{profile}" userid="{userid}"></page-profile>', '', '', function(opts) {
var self = this;
self.loggedin = false;
self.loggedout = !self.loggedin;
self.journalcreate = false;
self.dash = false;
self.entry = false
self.viewjournals = false;
self.viewjournal = false;
self.profile = false;

RiotControl.on('logout', function() {
	self.loggedin = false;
	self.loggedout = !self.loggedin;

	self.update();
});
RiotControl.on('login', function(user) {
	self.user = user;
	self.loggedin = true;
	self.loggedout = !self.loggedin;
	self.update();
});

self.clear = function() {
	self.dash = false;
	self.journalcreate = false;
	self.journal = false;
	self.entry = false
	self.viewjournals = false;
	self.viewjournal = false;
	self.profile = false;
}

route(function(collection, id, method, mid) {
	self.clear()
	switch(collection) {
		case 'view':

			if(method == "journal") {

			}
			break;
		case 'journals':
			if(id == 'create') {

				self.journalcreate = true;
			} else {
				if(method == 'entries') {

					self.journalid = id;
					self.entryid = mid;
					self.entry = true;

				} else {

					self.journalid = id;
					self.journal = true;
				}
			}
			break;
		case 'profile':
			self.userid = self.user.ID;
			self.profile = true;
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

riot.tag2('datepicker', '<input class="input" type="text" placeholder="yyyy/mm/dd" onblur="{onblur}" onfocus="{onfocus}" onkeydown="{onkeydown}" onmouseup="{onmouseup}" riot-value="{inputval}">', '', '', function(opts) {
var self = this;

var year = "yyyy";
var month = "mm";
var day = "dd";

self.inputval = "";

var datestr = "yyyy/mm/dd";

var optsdate = null;
self.on('update', function() {
	if(opts.date != '' && opts.date != optsdate) {
		datestr = moment(opts.date).format('YYYY/MM/DD');
		year = moment(opts.date).format('YYYY');
		month = moment(opts.date).format('MM');
		day = moment(opts.date).format('DD');
		self.inputval = datestr;
		optsdate = opts.date;
		self.update();
	}
});

self.onfocus = function(e) {

	e.target.value = datestr;
	self.update();
};

var KEY_DELETE = 8;
var KEY_0 = 48;
var KEY_1 = 49;
var KEY_3 = 51;
var KEY_9 = 57;
var ARROW_LEFT = 37;
var ARROW_RIGHT = 39;

self.onkeydown = function(e) {
	e.preventDefault();
	var start = e.target.selectionStart;
	var end = e.target.selectionEnd;
	if( start <= 4 ) {
		if(e.keyCode == KEY_DELETE) {
			if(year.length > 0) {
				year = year.substring(0, year.length -1);
				datestr = year + "/" + month + "/" + day;
				e.target.value = datestr;
				e.target.setSelectionRange(start-1,start-1);
			}
		}
		if(e.keyCode >= KEY_0 && e.keyCode <= KEY_9) {
			if(year.length == 4) {
				year = "";
			}
			year = year + e.key;
			datestr = year + "/" + month + "/" + day;
			e.target.value = datestr;
			e.target.setSelectionRange(start+1,start+1);

			if(year.length >= 4) {
				e.target.setSelectionRange(5,7);
			}
		}
	} else if( start <= 7 ) {
		if(e.keyCode == KEY_DELETE) {
			if(month.length > 0) {
				month = month.substring(0, month.length -1);
				datestr = year + "/" + month + "/" + day;
				e.target.value = datestr;
				e.target.setSelectionRange(start-1,start-1);
			}
		}
		if(e.keyCode >= KEY_0 && e.keyCode <= KEY_9) {

			if(month.length == 0 && e.keyCode > KEY_1) {
				return;
			}
			if(month.length == 2) {
				month = "";
			}
			month = month + e.key;
			datestr = year + "/" + month + "/" + day;
			e.target.value = datestr;
			e.target.setSelectionRange(start+1,start+1);

			if(month.length >= 2) {
				e.target.setSelectionRange(8,10);
			}

		}
		if(month.length >= 2) {
			e.target.setSelectionRange(8,10);
		}
	} else if( start <= 10 ) {
		if(e.keyCode == KEY_DELETE) {
			if(day.length > 0) {
				day = day.substring(0, day.length -1);
				datestr = year + "/" + month + "/" + day;
				e.target.value = datestr;
				e.target.setSelectionRange(start-1,start-1);
			}
		}
		if(e.keyCode >= KEY_0 && e.keyCode <= KEY_9) {

			if(day.length == 0 && e.keyCode > KEY_3) {
				return;
			}
			if(day.length >= 2 && start != end) {
				day = "";
			} else if(day.length == 2) {
				return;
			}
			day = day + e.key;
			datestr = year + "/" + month + "/" + day;
			e.target.value = datestr;
			e.target.setSelectionRange(start+1,start+1);
		}

	} else if( start <= 13 ) {
	} else if( start <= 16 ) {
	}
};

self.onblur = function(e) {
};

self.onmouseup = function(e) {
	e.preventDefault();
	var start = e.target.selectionStart;
	var end = e.target.selectionEnd;
	if( start <= 3 ) {
		e.target.setSelectionRange(0,4);
	} else if( start <= 6 ) {
		e.target.setSelectionRange(5,7);
	} else if( start <= 9 ) {
		e.target.setSelectionRange(8,10);
	} else if( start <= 12 ) {
		e.target.setSelectionRange(11,13);
	} else if( start <= 15 ) {
		e.target.setSelectionRange(14,16);
	}
};

self.date = function() {
	return new Date(Date.UTC(year,month-1,day));
};
});

riot.tag2('navbar', '<nav class="nav"> <div class="nav-left"> </div> <div class="nav-center"> <a class="nav-item" href="#/">a-Journal</a> </div> <div class="nav-right"> <a class="nav-item" href="#/profile">{user.Username}</a> <a class="nav-item" href="#" onclick="{logout}" if="{user.Username}">Log out</a> </div> </nav>', '', '', function(opts) {
var self = this;
self.user = {};

self.logout = function(e) {
	e.preventDefault();
	RiotControl.trigger('perform-logout', null);
};

RiotControl.on('login', function(user) {
	if(self.user == null) {
		self.user = {};
	} else {
		self.user = user;
	}
	self.update();
});
RiotControl.on('logout', function() {
	self.user = {};
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

riot.tag2('page-entryeditor', '<div class="section"> <div class="container"> <div class="columns"> <div class="column"> <label class="label">Title</label> <p class="control"> <input class="input" type="text" placeholder="Title" onkeyup="{onTitle}" riot-value="{entry.Title}"> </p> <label class="label">Date</label> <p> <datepicker date="{entry.Date}"></datepicker> </p> <label class="label">Content</label> <p class="control"> <textarea style="min-height: 200px;" class="textarea" placeholder="Textarea" onkeyup="{contentchange}">{entry.Content}</textarea> <span if="{err}" class="help is-danger">{err}</span> </p> <p> <br> <a class="button {is-link : showpreview}" onclick="{togglepreview}">Preview</a> <button class="button is-pulled-right {is-loading : saving}" onclick="{saveEntry}">Save</button> </p> </div> <div class="column" if="{showpreview}"> <label class="label">Preview</label> <raw class="markdown" content="{preview}"></raw> </div> </div> </div> </div>', '', '', function(opts) {
var self = this;
self.showpreview = false;
self.preview = "";
self.saving = false;
var converter = new showdown.Converter();

self.entry = {};

self.on('mount', function() {

	self.entry = {
		JournalID: parseInt(opts.journalid),
		Date: "",
		Title: "",
		Content: "",
		Tags: []
	};

	if( opts.entryid != 'create' ) {

		_aj.get("/api/journals/"+self.entry.JournalID+"/entries/"+opts.entryid, function(data, err) {

			if(err != null) {

				return;
			}
			self.entry = data;
			self.editContent = self.entry.Content;
			self.update();
		});
	}

	self.preview = converter.makeHtml(self.entry.Content);
	self.update();
});

self.editContent = "";

self.contentchange = function(e) {
	self.editContent = e.target.value;
	if(self.showpreview) {
		self.preview = converter.makeHtml(self.editContent);
	}
	self.update();
};

self.onTitle = function(e) {
	self.entry.Title = e.target.value;
};

self.saveEntry = function(e) {
	self.saving = true;
	self.entry.Date = self.tags.datepicker.date().toISOString();
	self.entry.Content = self.editContent;
	if(typeof(self.entry.ID) != 'undefined') {

		_aj.post("/api/journals/"+self.entry.JournalID+"/entries/"+self.entry.ID, self.entry, function(data, err) {
			self.saving = false;

			if( err != null ) {

				self.err = err;
				self.update();
				return
			}

			self.update();
		});

	} else {

		_aj.post("/api/journals/"+self.entry.JournalID+"/entries", self.entry, function(data, err) {
			self.saving = false;

			if( err != null ) {

				self.err = err;
				self.update();
				return
			}
			self.entry.ID = data.ID;

			self.update();
		});
	}
	self.entry.Content = self.editContent;
};

self.togglepreview = function(e) {
	self.showpreview = !self.showpreview;
	if( self.showpreview ) {
		self.preview = converter.makeHtml(self.editContent);
	}
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

riot.tag2('page-journal', '<section class="section"> <div class="container"> <section class="section"> <h3 class="title">Journal: {journal.Title}</h3> <p> {journal.Description} </p> <a class="button" href="#/journals/{opts.journalid}/view">View Journal</a> </section> <section class="section"> <button class="button" onclick="{newentry}">New Entry</button> </section> <section class="section"> <div class="box" each="{entry in journal.Entries}" onclick="{onentry}" style="cursor:pointer;"> <article class="media"> <div class="media-content"> <div class="content"> <p> <strong>{entry.Title}</strong> <small>@jzs</small> <small>31m</small> <br> <pre>{entry.Content.substring(0, 200)}...</pre> <br> <span each="{tag in parent.Tags}">{tag}</span> </p> </div> <nav class="level"> <div class="level-left"> <a class="level-item"> <span class="icon is-small"><i class="fa fa-reply"></i></span> </a> <a class="level-item"> <span class="icon is-small"><i class="fa fa-retweet"></i></span> </a> <a class="level-item"> <span class="icon is-small"><i class="fa fa-heart"></i></span> </a> </div> </nav> </div> </article> </div> </section> </div> </section>', '', '', function(opts) {
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
	}, {
		ID: 2,
		Title: "My second entry",
		Content: "",
		Tags: []
	}],
	Created: ""
};

self.on('mount', function() {
	_aj.get("/api/journals/" + opts.journalid, function(data, err) {
		if(err != null) {
			self.err = err;
			self.update();
			return;
		}
		self.journal = data;
		self.update();
	});
});

self.newentry = function(e) {
	route("/journals/"+opts.journalid+"/entries/create");
}

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
	self.loggingin = true;
	self.update();

	var user = {
		Username: self.username,
		Password: self.password
	};

	_aj.post("/api/users/login", user, function(data, err) {
		if( err != null ) {
			self.errmsg = data.Error;
			self.loginerr = true;
			self.update();
			return;
		}
		self.loggingin = false;
		self.update();
		data.Username = user.Username;
		RiotControl.trigger('perform-login', data);
		route("/");
	});
};

self.register = function(e) {
	var user = {
		Username: self.username,
		Password: self.password
	};
	_aj.post("/api/users", user, function(data, err) {
		if( err != null ) {

			self.errmsg = data.Error;
			self.loginerr = true;
			self.update();
			return;
		}
		self.login();
		self.update();
	});
};

});

riot.tag2('page-profile', '<div> Hello world!! wefewf </div>', '', '', function(opts) {
var self = this;

self.on('mount', function() {

	_aj.get("/api/profile", null, function(data, err) {
		if ( err != null ) {
			return;
		}
		self.profile = data;
	});
});

});

riot.tag2('page-viewjournal', '<section class="section"> <div class="container"> <section class="section"> <h3 class="title">Journal: {journal.Title}</h3> <p> {journal.Description} </p> <a class="button" href="#/journals/{opts.journalid}/view">View Journal</a> </section> <section class="section"> <button class="button" onclick="{newentry}">New Entry</button> </section> <section class="section"> <div class="box" each="{entry in journal.Entries}" onclick="{onentry}" style="cursor:pointer;"> <article class="media"> <div class="media-content"> <div class="content"> <p> <strong>{entry.Title}</strong> <small>@jzs</small> <small>31m</small> <br> <pre>{entry.Content.substring(0, 200)}...</pre> <br> <span each="{tag in parent.Tags}">{tag}</span> </p> </div> <nav class="level"> <div class="level-left"> <a class="level-item"> <span class="icon is-small"><i class="fa fa-reply"></i></span> </a> <a class="level-item"> <span class="icon is-small"><i class="fa fa-retweet"></i></span> </a> <a class="level-item"> <span class="icon is-small"><i class="fa fa-heart"></i></span> </a> </div> </nav> </div> </article> </div> </section> </div> </section>', '', '', function(opts) {
var self = this;
self.journal = {
	Tags: [],
	Entries: []
};

self.on('mount', function() {
	_aj.get("/api/journals/" + opts.journalid, function(data, err) {
		if(err != null) {
			self.err = err;
			self.update();
			return;
		}
		self.journal = data;
		self.update();
	});
});
});


riot.tag2('page-viewjournals', '', '', '', function(opts) {
});

riot.tag2('raw', '<span></span>', '', '', function(opts) {
this.root.innerHTML = opts.content

this.on('update', function() {
	this.root.innerHTML = opts.content
});
});
