<content>
	<page-dash if={loggedin && dash}></page-dash>
	<page-login if={loggedout}></page-login>
	<page-journal-create if={loggedin && journalcreate}></page-journal-create>
	<page-journal if={loggedin && journal} journalid={journalid}></page-journal>
	<page-entryeditor if={loggedin && entry} entryid={entryid}></page-entryeditor>
	<script>
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
				// show create dialog
				self.journalcreate = true;
			} else {
				if(method == 'entries') {
					if(mid == 'create') {
					} else {
						// Show entry.
						self.entryid = id;
						self.entry = true;
					}
					// show entries
				} else {
					// show journal
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
	</script>
</content>
