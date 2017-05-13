<page-viewuser>
	<section class="section">
		<div class="container">
			<section class="section">
				<p>
				{user.Username}
				{profile.Description}
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
	</script>
</page-viewuser>
