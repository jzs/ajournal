riot.tag2('content', '<page-dash if="{loggedin && dash}"></page-dash> <page-login if="{login}"></page-login> <page-register if="{register}"></page-register> <page-journal-create if="{loggedin && journalcreate}"></page-journal-create> <page-journal if="{loggedin && journal}" journalid="{journalid}"></page-journal> <page-entryeditor if="{loggedin && entry}" journalid="{journalid}" entryid="{entryid}"></page-entryeditor> <page-viewjournalentry if="{viewjournalentry}" journalid="{journalid}" entryid="{entryid}"></page-viewjournalentry> <page-viewjournal if="{viewjournal}" username="{username}" journalid="{journalid}"></page-viewjournal> <page-viewjournals if="{viewjournals}" userid="{userid}"></page-viewjournals> <page-viewuser if="{viewuser}" username="{username}"></page-viewuser> <page-profile if="{profile}" userid="{userid}"></page-profile>', '', '', function(opts) {
var self = this;
self.loggedin = false;
self.journalcreate = false;
self.dash = false;
self.entry = false
self.viewjournals = false;
self.viewjournal = false;
self.viewjournalentry = false;
self.viewuser = false;
self.profile = false;
self.register = false;
self.login = false;

RiotControl.on('logout', function() {
	self.loggedin = false;
	self.clear();
	self.login = true;
	self.update();
});
RiotControl.on('login', function(user) {
	self.user = user;
	self.loggedin = true;
	self.clear();
	self.dash = true;
	self.update();
});

self.clear = function() {
	self.dash = false;
	self.journalcreate = false;
	self.journal = false;
	self.entry = false
	self.viewjournals = false;
	self.viewjournal = false;
	self.viewjournalentry = false;
	self.profile = false;
	self.viewuser = false;
	self.register = false;
	self.login = false;
}

route(function(collection, id, method, mid) {
	self.clear()
	switch(collection) {
		case 'view':
			if(method == 'entries') {
				self.entryid = mid;
				self.viewjournalentry = true;
				break;
			}

			self.viewjournal = true;
			self.journalid = id;
			break;
		case 'users':
			if(method == 'journals') {
				self.viewjournal = true;
				self.journalid = mid;
				self.username = id;
				break;
			}
			self.viewuser = true;
			self.username = id;
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
			if(!self.loggedin) {
				route("/");
				return;
			}
			self.userid = self.user.ID;
			self.profile = true;
			break;
		case 'login':
			self.login = true;
			break;
		case 'register':
			self.register = true;
			break;
		default:
			if(self.loggedin) {
				self.dash = true;
			} else {
				self.login = true;
			}
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

riot.tag2('latestjournals', '<article class="media" each="{j in journals}"> <div class="media-content"> <div class="content"> <p> <strong>{j.Entry.Title}</strong> <br> <small>{j.Title}</small> <small>{moment(j.Entry.Date).format(\'YYYY/MM/DD\')}</small> <br> <small> <a href="/app#view/{j.ID}/entries/{j.Entry.ID}">Read more</a> </small> </p> </div> </div> </article>', '', '', function(opts) {
var self = this;
self.entries = {Entries: []};
self.journals = [];
_aj.get("/api/journals/latest?limit=3", function(data, err) {
	if(err != null) {
		self.err = err;
		self.update();
		return;
	}
	self.journals = data.Journals;
	self.update();
});

});

riot.tag2('navbar', '<nav class="nav"> <div class="nav-left"> <a class="nav-item" href="/"><img src="images/logo.png"></a> <a class="nav-item" href="#/">Home</a> </div> <div class="nav-center"> <a class="nav-item" href="#/profile">{user.Username}</a> </div> <span id="nav-toggle" class="nav-toggle {is-active: isActive}" onclick="{toggle}"> <span></span> <span></span> <span></span> </span> <div id="nav-menu" class="nav-right nav-menu {is-active: isActive}"> <a class="nav-item" href="#" onclick="{logout}" if="{user.Username}">Log out</a> </div> </nav>', '', '', function(opts) {
var self = this;
self.user = {};

self.isActive = false;
self.toggle = function(e) {
	self.isActive = !self.isActive;
	self.update();
}

self.logout = function(e) {
	e.preventDefault();
	self.isActive = false;
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


riot.tag2('page-dash', '<div class="section"> <div class="container"> <div class="columns"> <div class="column"> <h3 class="title">My Journals</h3> <section class="section" if="{journals.length < 1 && !loading}"> <p> <span>Welcome to a-Journal. It looks like you haven\'t created any journals yet. Let me help you get started.</span> </p> </section> <span class="help is-danger" if="{err}">{err}</span> <div class="box" each="{journal in journals}"> <article class="media"> <div class="media-left"> <figure class="image is-64x64"> <img src="images/128x128.png" alt="Image"> </figure> </div> <div class="media-content"> <div class="content" onclick="{tojournal}" style="cursor:pointer;"> <p> <strong>{journal.Title}</strong> <small>@jzs</small> <small>31m</small> <br> {journal.Description} </p> </div> <nav class="level"> <div class="level-left"> <a class="level-item" onclick="{createentry}"> <span class="icon is-small"><i class="fa fa-plus"></i></span> </a> </div> <div class="level-right"> <span class="level-item" if="{!journal.Public}"> Private </span> </div> </nav> </div> </article> </div> <button class="button" onclick="{newjournal}">New Journal</button> </div> </div> </div> </div>', '', '', function(opts) {
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

riot.tag2('page-entryeditor', '<div class="section"> <div class="container"> <div class="columns"> <div class="column"> <label class="label">Title</label> <p class="control"> <input class="input" type="text" placeholder="Title" onkeyup="{onTitle}" riot-value="{entry.Title}"> </p> <label class="label">Date</label> <p> <datepicker date="{entry.Date}"></datepicker> </p> <label class="label">Content</label> <p class="control"> <textarea style="min-height: 200px;" class="textarea" placeholder="Textarea" onkeyup="{contentchange}" ondrop="{drop}">{entry.Content}</textarea> <span if="{err}" class="help is-danger">{err}</span> </p> <p> <p class="control"> <ul> <li each="{blobs}"> <img class="thumb160" riot-src="{Links.Orig}" dragable="true" ondragstart="{drag}"> </li> </ul> </p> <p class="control"> <input id="blob" name="blob" type="file" accept="image/png, image/jpeg" onsubmit="{uploadfile}" onchange="{selectfile}"> <span if="{uploadingfile}">Uploading</span> </p> <br> <a class="button {is-link : showpreview}" onclick="{togglepreview}">Preview</a> <button class="button is-pulled-right {is-loading : saving}" onclick="{saveEntry}">Save</button> </p> </div> <div class="column" if="{showpreview}"> <label class="label">Preview</label> <raw class="markdown" content="{preview}"></raw> </div> </div> </div> </div>', '', '', function(opts) {
var self = this;
self.showpreview = false;
self.preview = "";
self.saving = false;
var converter = new showdown.Converter();
converter.setFlavor('github');

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

		_aj.get("/api/journals/"+self.entry.JournalID+"/blobs", function(data, err) {
			if(err != null) {
				return;
			}
			self.blobs = data;
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

self.selectfile = function(e) {
	console.log("hello");

	self.uploadfile(e);
}

self.uploadfile = function(e) {
	e.preventDefault();

	var files = e.target.files;

	var formData = new FormData();

	for (var i = 0; i < files.length; i++) {
		var file = files[i];

		if (!file.type.match('image.*')) {
			continue;
		}

		formData.append('blobs', file, file.name);
	}
	var xhr = new XMLHttpRequest();
	xhr.open('POST', '/api/journals/'+ self.entry.JournalID +'/blobs', true);
	xhr.onload = function () {
		if (xhr.status === 200) {

			self.uploadingfile = false;
			var data = JSON.parse(xhr.response);
			self.blobs.push(data);
			self.update();
		} else {
			alert('An error occurred!');
		}
	};

	self.uploadingfile = true;
	self.update();

	xhr.send(formData);
};

self.dragItem = null;
self.drag = function(e) {
	self.dragItem = e.item;
};
self.drop = function(e) {
	console.log(e);
	e.preventDefault();
	console.log(self.dragItem);
	console.log(e.target.selectionStart);
	e.toElement.value += "![](" + self.dragItem.Links.Orig + ")"
};

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

riot.tag2('page-journal', '<section class="section"> <div class="container"> <section class="section"> <h3 class="title">Journal: {journal.Title}</h3> <p> {journal.Description} </p> <a class="button" href="#/users/{opts.username}/journals/{opts.journalid}">View Journal</a> </section> <section class="section"> <button class="button" onclick="{newentry}">New Entry</button> </section> <section class="section"> <div class="box" each="{entry in entries.Entries}" onclick="{onentry}" style="cursor:pointer;"> <article class="media"> <div class="media-content"> <div class="content"> <p> <strong>{entry.Title}</strong> <small>@jzs</small> <small>31m</small> <br> <pre>{entry.Content.substring(0, 200)}...</pre> <br> <span each="{tag in parent.Tags}">{tag}</span> </p> </div> <nav class="level"> <div class="level-left"> <a class="level-item"> <span class="icon is-small"><i class="fa fa-reply"></i></span> </a> <a class="level-item"> <span class="icon is-small"><i class="fa fa-retweet"></i></span> </a> <a class="level-item"> <span class="icon is-small"><i class="fa fa-heart"></i></span> </a> </div> </nav> </div> </article> </div> <button class="button is-primary" if="{entries.HasNext}" onclick="{loadMore}">Load more</button> </section> </div> </section>', '', '', function(opts) {
var self = this;
self.entries = {Entries: []};
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

		getEntries(null);
	});
});

self.newentry = function(e) {
	route("/journals/"+opts.journalid+"/entries/create");
}

self.onentry = function(e) {
	route("/journals/"+opts.journalid+"/entries/"+e.item.entry.ID);
};

self.loadMore = function(e) {
	getEntries(self.entries.Next);
}

var getEntries = function(from) {
	var fromstr = "?limit=10";
	if(from != null) {
		fromstr += "&from=" +from;
	}
	_aj.get("/api/journals/" + opts.journalid + "/entries" + fromstr, function(data, err) {
		if(err != null) {
			self.err = err;
			self.update();
			return;
		}
		Array.prototype.push.apply(self.entries.Entries, data.Entries);
		self.entries.HasNext = data.HasNext;
		self.entries.Next = data.Next;
		self.update();
	});
};

});

riot.tag2('page-login', '<section class="hero is-fullheight"> <div class="hero-head"></div> <div class="hero-body"> <div class="container"> <div class="columns"> <div class="column is-half is-offset-one-quarter has-text-left"> <h3 class="title">Sign in</h3> <div class="card"> <div class="card-content is-clearfix"> <p class="subtitle has-text-black">Sign in to aJournal to create/edit your own journals</p> <div if="{!loginuserpass}"> <br> <div class="control"> <a class="button is-large is-full-width is-primary" href="/api/oauth/google">Continue with google</a> </div> <div class="control"> <a class="button is-large is-full-width is-primary" onclick="{clickloginuser}">Sign in with email</a> </div> <div class="control is-pulled-right"> <a class="is-large is-full-width" href="#/register">Sign up with email</a> </div> </div> <div if="{loginuserpass}"> <label class="label">Username</label> <p class="control has-icon has-icon-right"> <input class="input" type="text" placeholder="Username" onkeyup="{onusername}" value=""> <span class="icon is-small hidden"> <i class="fa fa-check"></i> </span> <span class="hidden help is-success">This username is available</span> </p> <label class="label">Password</label> <p class="control has-icon has-icon-right"> <input class="input" type="password" placeholder="Password" value="" onkeyup="{onpassword}"> <span class="icon is-small hidden"> <i class="fa fa-warning"></i> </span> <span if="{loginerr}" class="help is-danger">{errmsg}</span> </p> <div class="control is-grouped is-pulled-right"> <p class="control"> <button disabled="{isdisabled}" class="button is-success {is-loading : loggingin}" onclick="{login}">Login</button> </p> </div> </div> </div> </div> </div> </div> </div> </div> <div class="hero-foot"></div> </section>', '', '', function(opts) {
var self = this;

self.loginuserpass = false;
self.clickloginuser = function(e) {
	self.loginuserpass = true;
}

self.isdisabled = true;

self.loggingin = false;
self.loginerr = false;
self.errmsg = "";

self.username = "";
self.onusername = function(e) {
	self.username = e.target.value;
	self.verifylogin();
};

self.password = "";
self.onpassword = function(e) {
	self.password = e.target.value;
	self.verifylogin();
};

self.verifylogin = function() {
	if(self.password != "" && self.username != "") {
		self.isdisabled = false;
	} else {
		self.isdisabled = true;
	}
	self.update();
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
		self.loggingin = false;
		self.update();
		if( err != null ) {
			self.errmsg = err;
			self.loginerr = true;
			self.update();
			return;
		}
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

			self.errmsg = err;
			self.loginerr = true;
			self.update();
			return;
		}
		self.login();
		self.update();
	});
};

});

riot.tag2('page-profile', '<div class="container"> <section class="section"> <h3 class="title">Hi {profile.Name}</h3> <div class="container"> <label class="label">Full Name</label> <p class="control"> <input class="input" type="text" placeholder="Full Name" onkeyup="{onFullName}" riot-value="{profile.Name}"> </p> <label class="label">E-mail</label> <p class="control"> <input class="input" type="text" placeholder="E-mail" onkeyup="{onEmail}" riot-value="{profile.Email}"> </p> <p class="label">Short name</label> <input class="input {is-danger: !shortnameValid}" type="text" placeholder="Short name" onkeyup="{onShortName}" riot-value="{profile.ShortName}"> </p> <label class="label">Public profile description</label> <p class="control"> <textarea class="textarea" type="text" placeholder="Description" onkeyup="{onDesc}" riot-value="{profile.Description}"> </textarea> </p> <button class="button is-primary is-medium" onclick="{save}">Save</button> </div> </section> <section class="section"> <h3 class="title">Memberships</h3> <hr> <p> You are currently subscribed to the basic plan. You can upgrade your subscription below. </p> <div class="columns"> <div class="column"> <div class="box"> <article class="media"> <div class="media-content"> <h3 class="title">Basic Plan</h3> <hr> <ul> <li>5 journals</li> <li>100 posts per journal</li> </ul> <p class="has-text-centered"> <h3 class="title has-text-centered">Free</h3> </p> <p class="hero-buttons"> <button class="button is-large is-primary">Free</button> </p> </div> </article> </div> </div> <div class="column"> <div class="box"> <article class="media"> <div class="media-content"> <h3 class="title">Full Plan</h3> <hr> <ul> <li>Unlimited journals</li> <li>Unlimited posts per journal</li> </ul> <p class="has-text-centered"> <h3 class="title has-text-centered">100 dkk/year</h3> </p> <p class="hero-buttons"> <button class="button is-large is-primary" onclick="{upgrade}">Upgrade</button> </p> </div> </article> </div> </div> </div> </section> <div class="modal {is-active : showmodal}"> <div class="modal-background"></div> <div class="modal-card"> <form action="/charge" method="post" id="payment-form"> <header class="modal-card-head"> <p class="modal-card-title">Upgrade to Paid plan</p> <button class="delete" onclick="{closemodal}"></button> </header> <section class="modal-card-body"> <p class="control"> <label for="email-element" class="label"> E-mail </label> <input class="input" placeholder="E-mail" type="text" riot-value="{profile.Email}" onkeyup="{onemail}"> <div id="email-errors">{emailerr}</div> </p> <p class="control"> <label class="label"> Name on debit or credit card </label> <input class="input" name="cardholder-name" placeholder="Name on debit or credit card" type="text"> </p> <p class="control"> <label for="card-element" class="label"> Credit or debit card </label> <div id="card-element"> </div> <div id="card-errors">{carderr}</div> </p> </section> <footer class="modal-card-foot"> <a class="button is-success {is-loading: upgrading}" onclick="{performUpgrade}">Pay 100 dkk</a> <a class="button" onclick="{closemodal}">Cancel</a> </footer> </form> </div> </div> </div>', '', '', function(opts) {
var self = this;
self.showmodal = false;
self.stripe = null;
self.card = null;
self.carderr = null;
self.shortnameValid = true;

self.profile = {};

self.on('mount', function() {

	self.stripe = Stripe('pk_test_4XUbWX7yh2AAiIsDCktzIRPE');
	var elements = self.stripe.elements();

	var classes = {
		base: "stripe-cardelem"
	};
	var style = {
		base: {
			lineHeight: "2"
		}
	};

	self.card = elements.create('card', {style: style, classes: classes});
	self.card.addEventListener('change', function(event) {
		if(event.error) {
			self.carderr = event.error.message;
		} else {
			self.carderr = null;
		}
		self.update();
	});

	self.card.mount('#card-element');

	_aj.get("/api/profile", function(data, err) {
		if ( err != null ) {
			return;
		}
		self.profile = data;
		self.update();
	});
});

self.save = function(e) {
	_aj.post("/api/profile", self.profile, function(data, err) {
		if ( err != null ) {

			return;
		}
		self.profile = data;
		self.update();
	});
}

self.onFullName = function(e) {
	self.profile.Name = e.target.value;
	self.update();
};
self.onEmail = function(e) {
	self.profile.Email = e.target.value;
};
self.onShortName= function(e) {
	self.profile.ShortName = e.target.value;
	_aj.get("/api/profile/"+self.profile.ID+"/shortname/"+e.target.value, function(data, err) {
		if ( err != null ) {
			return;
		}
		self.shortnameValid = data;
		self.update();
	});
};

self.onDesc = function(e) {
	self.profile.Description = e.target.value;
};

self.upgrade = function(e) {
	self.showmodal = true;
	self.update();
}
self.closemodal = function(e) {
	e.preventDefault();
	self.showmodal = false;
	self.update();
}

self.performUpgrade = function(e) {
	e.preventDefault();
	self.upgrading = true;
	self.update();

	self.stripe.createToken(self.card).then(function(result) {
		self.upgrading = false;
		if (result.error) {

			self.carderr = result.error.message;
			self.update();
		} else {

			var args = {Profile: self.profile, Token: result.token.id, Plan: 2};
			_aj.post("/api/profile/signup", args, function(data, err) {
				if( err != null ) {
					self.carderr = err;
					self.update();

					return;
				}
				console.log(data);
			});
		}
	});
}

});

riot.tag2('page-register', '<section class="hero is-fullheight is-primary"> <div class="hero-head"></div> <div class="hero-body"> <div class="container"> <div class="columns"> <div class="column is-half is-offset-one-quarter has-text-left"> <h3 class="title">Register</h3> <div class="card"> <div class="card-content is-clearfix"> <label class="label">Username*</label> <p class="control has-icon has-icon-right"> <input class="input" type="text" placeholder="Username" onkeyup="{onusername}" value=""> <span class="icon is-small hidden"> <i class="fa fa-check"></i> </span> <span class="hidden help is-success">This username is available</span> </p> <label class="label">Password*</label> <p class="control has-icon has-icon-right"> <input class="input" type="password" placeholder="Password" value="" onkeyup="{onpassword}"> <span class="icon is-small hidden"> <i class="fa fa-warning"></i> </span> <span if="{loginerr}" class="help is-danger">{errmsg}</span> </p> <label class="label">Full Name</label> <p class="control has-icon has-icon-right"> <input class="input" type="text" placeholder="Full Name" value="" onkeyup="{onname}"> </p> <label class="label">E-mail</label> <p class="control has-icon has-icon-right"> <input class="input" type="text" placeholder="E-mail" value="" onkeyup="{onemail}"> </p> <p class="control"> <label class="checkbox"> <input type="checkbox" checked="{user.Public}" onchange="{onlicense}"> Accept terms* </label> </p> <div class="control is-grouped is-pulled-right"> <p class control> <button disabled="{formdisabled}" class="button is-success {is-disabled : loggingin}" onclick="{register}">Register</button> </p> </div> </div> </div> </div> </div> </div> </div> <div class="hero-foot"></div> </section>', '', '', function(opts) {
var self = this;

self.user = {
	Username: "",
	Password: "",
	Email: "",
	License: false
};

self.formdisabled = true;
self.submitting = false;

self.loggingin = false;
self.loginerr = false;
self.errmsg = "";

self.onusername = function(e) {
	self.user.Username = e.target.value;
	self.validateform();
};

self.onpassword = function(e) {
	self.user.Password = e.target.value;
	self.validateform();
};
self.onemail = function(e) {
	self.user.Email = e.target.value;
	self.validateform();
};
self.onlicense = function(e) {
	self.user.License = e.target.checked;
	self.validateform();
};
self.onname = function(e) {
	self.user.Name = e.target.value;
	self.validateform();
};

self.validateform = function() {
	if(self.user.Username != "" && self.user.Password != "" && self.user.License != "") {
		self.formdisabled = false;
		self.update();
		return
	}
	self.formdisabled = true;
	self.update();
};

self.login = function() {
	self.loggingin = true;
	self.update();

	_aj.post("/api/users/login", self.user, function(data, err) {
		self.loggingin = false;
		self.update();
		if( err != null ) {
			self.errmsg = err;
			self.loginerr = true;
			self.update();
			return;
		}

		if(self.user.Name != "" || self.user.Email != "") {
			self.user.ID = data.UserID;
			_aj.put("/api/profile", self.user, function(data, err) {});
		}

		data.Username = self.user.Username;
		RiotControl.trigger('perform-login', data);
		self.user = {};
		route("/");
	});
};

self.register = function(e) {
	_aj.post("/api/users", self.user, function(data, err) {
		if( err != null ) {

			self.errmsg = err;
			self.loginerr = true;
			self.update();
			return;
		}

		self.login();
		self.update();
	});
};

});

riot.tag2('page-viewjournal', '<section class="section"> <div class="container"> <section class="section"> <h3 class="title">Journal: {journal.Title}</h3> <p> {journal.Description} </p> </section> <section class="section"> <div class="box" each="{entry in entries.Entries}" onclick="{onentry}" style="cursor:pointer;"> <article class="media"> <div class="media-content"> <div class="content"> <p> <h3 class="title">{entry.Title}</h3> <small>@{user.Username}</small> <small>{entry.Date}</small> <raw class="markdown" content="{entry.preview}"></raw> <br> <p if="{entry.Content.length > 2000}"><span>Read more</span></p> <span each="{tag in parent.Tags}">{tag}</span> </p> </div> </div> </article> </div> <button class="button is-primary" if="{entries.HasNext}" onclick="{loadMore}">Load more</button> </section> </div> </section>', '', '', function(opts) {
var self = this;

self.user = {};
self.entries = {Entries: []};
self.journal = {
	Tags: [],
	Entries: []
};

var converter = new showdown.Converter();
converter.setFlavor('github');

self.on('mount', function() {
	_aj.get("/api/journals/" + opts.journalid, function(data, err) {
		if(err != null) {
			self.err = err;
			self.update();
			return;
		}
		self.journal = data;

		self.update();

		getEntries(null);

		_aj.get("/api/users/" + opts.username, function(data, err) {
			if(err != null) {
				self.err = err;
				self.update();
				return;
			}
			self.user = data;
			self.update();
		});
	});
});

self.onentry = function(e) {
	route("/view/"+opts.journalid+"/entries/"+e.item.entry.ID);
};

self.loadMore = function(e) {
	getEntries(self.entries.Next);
}

var getEntries = function(from) {
	var fromstr = "?limit=10";
	if(from != null) {
		fromstr += "&from=" +from;
	}
	_aj.get("/api/journals/" + opts.journalid + "/entries" + fromstr, function(data, err) {
		if(err != null) {
			self.err = err;
			self.update();
			return;
		}

		for(var i = 0; i < data.Entries.length; i++) {
			var j = data.Entries[i];
			j.preview = converter.makeHtml(j.Content.substring(0, 2000));

		}

		Array.prototype.push.apply(self.entries.Entries, data.Entries);
		self.entries.HasNext = data.HasNext;
		self.entries.Next = data.Next;
		self.update();
	});
};

});

riot.tag2('page-viewjournalentry', '<div class="section"> <div class="container"> <div class="columns"> <div class="column"> <h1 class="title">{entry.Title}</h1> <h2 class="subtitle">{entry.Date}</h2> <hr> <div class="content"> <raw class="markdown" content="{entry.HtmlContent}"></raw> </div> </div> </div> </div> </div>', '', '', function(opts) {
var self = this;
self.entry = {};

self.on('mount', function() {
	self.entry = {
		JournalID: parseInt(opts.journalid),
		Date: "",
		Title: "",
		Content: "",
		Tags: []
	};

	_aj.get("/api/journals/"+self.entry.JournalID+"/entries/"+opts.entryid, function(data, err) {

		if(err != null) {

			return;
		}
		self.entry = data;
		self.update();
	});

});

});

riot.tag2('page-viewjournals', '', '', '', function(opts) {
});

riot.tag2('page-viewuser', '<section class="section"> <div class="container"> <section class="section"> <div class="columns"> <div class="column is-one-quarter"> <img src="images/profile-placeholder.png"> </div> <div class="column"> <h3 class="title">{profile.Name} ({user.Username})</h3> {profile.Description} </div> </div> <p> </p> <h3 class="title">Journals</h3> <div class="box" each="{journal in journals}"> <article class="media"> <div class="media-content"> <div class="content" onclick="{tojournal}" style="cursor:pointer;"> <p> <strong>{journal.Title}</strong> <br> {journal.Description} </p> </div> </div> </article> </div> </section> </div> </section>', '', '', function(opts) {
var self = this;
self.journals = [];
self.user = {Username: opts.username};
self.profile = {};

self.on('mount', function() {

	_aj.get("/api/users/" + opts.username, function(data, err) {
		if(err != null) {
			self.err = err;
			self.update();
			return;
		}
		self.user = data;
		self.update();

		_aj.get("/api/users/" + data.ID + "/profile", function(data, err) {
			if(err != null) {
				return;
			}
			self.profile = data;
			self.update();
		});

		_aj.get("/api/users/" + data.ID + "/journals", function(data, err) {
			if(err != null) {
				self.err = err;
				self.update();
				return;
			}
			self.journals = data;
			self.update();
		});
	});
});

self.tojournal = function(e) {
	route("/users/"+opts.username+"/journals/" + e.item.journal.ID);
};
});

riot.tag2('raw', '<span></span>', '', '', function(opts) {
this.root.innerHTML = opts.content

this.on('update', function() {
	this.root.innerHTML = opts.content
});
});
