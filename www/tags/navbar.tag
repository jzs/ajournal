<navbar>
<nav class="nav">
	<div class="nav-left">
	</div>
	<div class="nav-center">
		<a class="nav-item" href="#/">a-Journal</a>
	</div>
	<div class="nav-right">
		<a class="nav-item" href="#" onclick={logout} if={user}>Log out</a>
	</div>
</nav>
	<script>
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
	</script>

</navbar>
