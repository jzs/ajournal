<page-profile>
	<div class="container">
		<section class="section">
			<h3 class="title">Hi {profile.Name}</h3>
			<div class="container">
				<div class="avatar" style="background-image:url('{profile.Picture.Links.Orig}');"></div>
				<label class="label">Full Name</label>
				<p class="control">
				<input class="input" type="text" placeholder="Full Name" onkeyup={onFullName} value={profile.Name} />
				</p>
				<label class="label">E-mail</label>
				<p class="control">
				<input class="input" type="text" placeholder="E-mail" onkeyup={onEmail} value={profile.Email} />
				</p>
				<p class="label">Short name</label>
			<input class="input {is-danger: !shortnameValid}" type="text" placeholder="Short name" onkeyup={onShortName} value={profile.ShortName} />
				</p>
				<label class="label">Public profile description</label>
				<p class="control">
				<textarea class="textarea" type="text" placeholder="Description" onkeyup={onDesc} value={profile.Description}> </textarea>
				</p>
				<button class="button is-primary is-medium" onclick={save}>Save</button>
			</div>

		</section>
		<section class="section">
			<h3 class="title">Memberships</h3>
			<hr/>
			<p>
			You are currently subscribed to the basic plan. You can upgrade your subscription below.
			</p>

			<div class="columns">
				<div class="column">
					<div class="box">
						<article class="media">
							<div class="media-content">
								<h3 class="title">Basic Plan</h3>
								<hr/>
								<ul>
									<li>5 journals</li>
									<li>100 posts per journal</li>
								</ul>
								<p class="has-text-centered">
								<h3 class="title has-text-centered">Free</h3>
								</p>
								<p class="hero-buttons">
								<button class="button is-large is-primary">Free</button>
								</p>
							</div>
						</article>
					</div>
				</div>
				<div class="column">
					<div class="box">
						<article class="media">
							<div class="media-content">
								<h3 class="title">Full Plan</h3>
								<hr/>
								<ul>
									<li>Unlimited journals</li>
									<li>Unlimited posts per journal</li>
								</ul>
								<p class="has-text-centered">
								<h3 class="title has-text-centered">100 dkk/year</h3>
								</p>
								<p class="hero-buttons">
								<button class="button is-large is-primary" onclick={upgrade}>Upgrade</button>
								</p>
							</div>
						</article>
					</div>
				</div>
			</div>
		</section>


		<div class="modal {is-active : showmodal }">
			<div class="modal-background"></div>
			<div class="modal-card">
				<form action="/charge" method="post" id="payment-form">
					<header class="modal-card-head">
						<p class="modal-card-title">Upgrade to Paid plan</p>
						<button class="delete" onclick={closemodal}></button>
					</header>
					<section class="modal-card-body">
						<!-- Content ... -->

						<p class="control">
						<label for="email-element" class="label">
							E-mail
						</label>
						<input class="input" placeholder="E-mail" type="text" value={profile.Email} onkeyup={onemail}/>

						<!-- Used to display Element errors -->
						<div id="email-errors">{emailerr}</div>
						</p>

						<p class="control">
						<label class="label">
							Name on debit or credit card
						</label>
						<input class="input" name="cardholder-name" placeholder="Name on debit or credit card" type="text" />
						</p>

						<p class="control">
						<label for="card-element" class="label">
							Credit or debit card
						</label>
						<div id="card-element">
							<!-- a Stripe Element will be inserted here. -->
						</div>

						<!-- Used to display Element errors -->
						<div id="card-errors">{carderr}</div>
						</p>


					</section>
					<footer class="modal-card-foot">
						<a class="button is-success {is-loading: upgrading}" onclick={performUpgrade}>Pay 100 dkk</a>
						<a class="button" onclick={closemodal}>Cancel</a>
					</footer>

				</form>
			</div>
		</div>

	</div>
	<script>
		var self = this;
self.showmodal = false;
self.stripe = null;
self.card = null;
self.carderr = null;
self.shortnameValid = true;

self.profile = {Picture: {Links: {}}};

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

	// Create an instance of the card Element
	self.card = elements.create('card', {style: style, classes: classes});
	self.card.addEventListener('change', function(event) {
		if(event.error) {
			self.carderr = event.error.message;
		} else {
			self.carderr = null;
		}
		self.update();
	});

	// Add an instance of the card Element into the `card-element` <div>
	self.card.mount('#card-element');

	// Fetch profile!
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
			// TODO Show error if profile save failed!
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
	// TODO disable submit button
	self.stripe.createToken(self.card).then(function(result) {
		self.upgrading = false;
		if (result.error) {
			// Inform the user if there was an error
			self.carderr = result.error.message;
			self.update();
		} else {
			// Send the token to your server
			// stripeTokenHandler(result.token);
			var args = {Profile: self.profile, Token: result.token.id, Plan: 2};
			_aj.post("/api/profile/signup", args, function(data, err) {
				if( err != null ) {
					self.carderr = err;
					self.update();
					// TODO enable submit button
					return;
				}
				console.log(data);
			});
		}
	});
}

	</script>
</page-profile>
