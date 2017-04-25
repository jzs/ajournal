<page-register>
	<section class="hero is-fullheight is-primary">
		<div class="hero-head"></div>
		<div class="hero-body">
			<div class="container">
				<div class="columns">
					<div class="column is-half is-offset-one-quarter has-text-left">
						<h3 class="title">Register</h3>
						<div class="card">
							<div class="card-content is-clearfix">
								<label class="label">Username*</label>
								<p class="control has-icon has-icon-right">
								<input class="input" type="text" placeholder="Username" onkeyup={onusername} value="">
								<span class="icon is-small hidden">
									<i class="fa fa-check"></i>
								</span>
								<span class="hidden help is-success">This username is available</span>
								</p>

								<label class="label">Password*</label>
								<p class="control has-icon has-icon-right">
								<input class="input" type="password" placeholder="Password" value="" onkeyup={onpassword}>
								<span class="icon is-small hidden">
									<i class="fa fa-warning"></i>
								</span>
								<span if={loginerr} class="help is-danger">{errmsg}</span>
								</p>

								<label class="label">Full Name</label>
								<p class="control has-icon has-icon-right">
								<input class="input" type="text" placeholder="Full Name" value="" onkeyup={onname}>
								</p>

								<label class="label">E-mail</label>
								<p class="control has-icon has-icon-right">
								<input class="input" type="text" placeholder="E-mail" value="" onkeyup={onemail}>
								</p>

								<p class="control">
								<label class="checkbox">
									<input type="checkbox" checked={user.Public} onchange={onlicense}>
									Accept terms*
								</label>
								</p>

								<div class="control is-grouped is-pulled-right">
									<p class"control">
									<button disabled={formdisabled} class="button is-success {is-disabled : loggingin}" onclick={register}>Register</button>
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

// TODO: Move login function out into it's own class.
self.login = function() {
	self.loggingin = true;
	self.update();

	// Perform login
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
			self.user.ID = data.UserID; // Cause profile expects user id set on the post
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
			// Do shit!
			self.errmsg = err;
			self.loginerr = true;
			self.update();
			return;
		}

		self.login();
		self.update();
	});
};

	</script>
</page-register>
