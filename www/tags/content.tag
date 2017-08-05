<content>
	<page-dash if={loggedin && dash}></page-dash>
	<page-login if={login}></page-login>
	<page-register if={register}></page-register>
	<page-journal-create if={loggedin && journalcreate}></page-journal-create>
	<page-journal if={loggedin && journal} journalid={journalid}></page-journal>
	<page-entryeditor if={loggedin && entry} journalid={journalid} entryid={entryid}></page-entryeditor>
	<page-viewjournalentry if={viewjournalentry} journalid={journalid} entryid={entryid}></page-viewjournalentry>
	<page-viewjournal if={viewjournal} username={username} journalid={journalid}></page-viewjournal>
	<page-viewjournals if={viewjournals} userid={userid}></page-viewjournals>
	<page-viewuser if={viewuser} username={username}></page-viewuser>
	<page-profile if={profile} userid={userid}></page-profile>
	<script>
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
			// TODO Handle views...
			// id == userid,
			// method == journal
			// 
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
				// show create dialog
				self.journalcreate = true;
			} else {
				if(method == 'entries') {
					// Show entry.
					self.journalid = id;
					self.entryid = mid;
					self.entry = true;
					// show entries
				} else {
					// show journal
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
	</script>
</content>
