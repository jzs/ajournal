<navbar>
<nav class="nav">
	<div class="nav-left">
	</div>
	<div class="nav-center">
		<a class="nav-item" href="#/">a-Journal</a>
	</div>
	<div class="nav-right">
		<a class="nav-item" href="#/profile">{user.Username}</a>
		<a class="nav-item" href="#" onclick={logout} if={user.Username}>Log out</a>
	</div>
</nav>
	<script>
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
	</script>

</navbar>
