<content>
	<page-dash if={loggedin && dash}></page-dash>
	<page-login if={loggedout}></page-login>
	<page-journal-create if={loggedin && journalcreate}></page-journal-create>
	<page-journal if={loggedin && journal} journalid={journalid}></page-journal>
	<page-entryeditor if={loggedin && entry} journalid={journalid} entryid={entryid}></page-entryeditor>
	<page-viewjournalentry if={viewjournalentry} journalid={journalid} entryid={entryid}></page-viewjournalentry>
	<page-viewjournal if={viewjournal} journalid={journalid}></page-viewjournal>
	<page-viewjournals if={viewjournals} userid={userid}></page-viewjournals>
	<page-profile if={profile} userid={userid}></page-profile>
	<script>
var self = this;
self.loggedin = false;
self.loggedout = !self.loggedin;
self.journalcreate = false;
self.dash = false;
self.entry = false
self.viewjournals = false;
self.viewjournal = false;
self.viewjournalentry = false;
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
	self.viewjournalentry = false;
	self.profile = false;
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
	</script>
</content>
