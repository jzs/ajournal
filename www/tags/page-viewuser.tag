<page-viewuser>
	<section class="section">
		<div class="container">
			<section class="section">
				<div class="columns">
					<div class="column is-one-quarter">
					<img src="images/profile-placeholder.png"/>
					</div>
					<div class="column">
						<h3 class="title">{profile.Name} ({user.Username})</h3>
						{profile.Description}
					</div>
				</div>
				<p>
				</p>
				<h3 class="title">Journals</h3>
				<div class="box" each={journal in journals}>
					<article class="media">
						<div class="media-content">
							<div class="content" onclick={tojournal} style="cursor:pointer;">
								<p>
								<strong>{journal.Title}</strong> 
								<br>
								{journal.Description}
								</p>
							</div>
						</div>
					</article>
				</div>

			</section>
		</div>
	</section>
	<script>
var self = this;
self.journals = [];
self.user = {Username: opts.username};
self.profile = {};

self.on('mount', function() {
	// Fetch journals for user...
	_aj.get("/api/users/" + opts.username, function(data, err) {
		if(err != null) {
			self.err = err;
			self.update();
			return;
		}
		self.user = data;
		self.update();

		_aj.get("/api/users/" + data.ID + "/profile", function(data, err) {
			if(err != null) {
				return;
			}
			self.profile = data;
			self.update();
		});

		_aj.get("/api/users/" + data.ID + "/journals", function(data, err) {
			if(err != null) {
				self.err = err;
				self.update();
				return;
			}
			self.journals = data;
			self.update();
		});
	});
});

self.tojournal = function(e) {
	route("/users/"+opts.username+"/journals/" + e.item.journal.ID);
};
	</script>
</page-viewuser>
