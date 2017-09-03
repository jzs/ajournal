<navbar>
<nav class="nav">
	<div class="nav-left">
		<a class="nav-item" href="/"><img src="images/logo.png"/></a>
		<a class="nav-item" href="#/">Home</a>
	</div>
	<div class="nav-center">
		<a class="nav-item" href="#/profile">{user.Username}</a>
	</div>
	<span id="nav-toggle" class="nav-toggle {is-active: isActive}" onclick={toggle}>
    <span></span>
    <span></span>
    <span></span>
  </span>
  <div id="nav-menu" class="nav-right nav-menu {is-active: isActive}">
		<a class="nav-item" href="#" onclick={logout} if={user.Username}>Log out</a>
	</div>
</nav>
	<script>
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
	</script>

</navbar>
