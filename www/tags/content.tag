<content>
	<page-dash if={loggedin}></page-dash>
	<page-login if={loggedout}></page-login>
	<script>
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

	</script>
</content>
