<page-login>
	<section class="hero is-fullheight is-primary">
		<div class="hero-head"></div>
		<div class="hero-body">
			<div class="container">
				<div class="columns">
					<div class="column is-half is-offset-one-quarter has-text-left">
						<h3 class="title">Journal</h3>
						<div class="card">
							<div class="card-content is-clearfix">
								<label class="label">Username</label>
								<p class="control has-icon has-icon-right">
								<input class="input" type="text" placeholder="Username" onkeyup={onusername} value="">
								<span class="icon is-small hidden">
									<i class="fa fa-check"></i>
								</span>
								<span class="hidden help is-success">This username is available</span>
								</p>

								<label class="label">Password</label>
								<p class="control has-icon has-icon-right">
								<input class="input" type="password" placeholder="Password" value="" onkeyup={onpassword}>
								<span class="icon is-small hidden">
									<i class="fa fa-warning"></i>
								</span>
								<span if={loginerr} class="help is-danger">{errmsg}</span>
								</p>

								<div class="control is-grouped is-pulled-right">
									<p class"control">
									<button class="button is-link {is-disabled : loggingin}" onclick={register}>Register</button>
									</p>
									<p class="control">
									<button class="button is-success {is-loading : loggingin}" onclick={login}>Login</button>
									</p>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
		<div class="hero-foot"></div>
	</section>
	<script>
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

	// Perform login
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
			// Do shit!
			self.errmsg = data.Error;
			self.loginerr = true;
			self.update();
			return;
		}
		self.login();
		self.update();
	});
};

	</script>
</page-login>
