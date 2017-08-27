<page-dash>
	<div class="section">
		<div class="container">
			<div class="columns">
				<div class="column">
					<h3 class="title">My Journals</h3>
					<section class="section" if={journals.length < 1 && !loading}>
						<p>
						<span>Welcome to a-Journal. It looks like you haven't created any journals yet. Let me help you get started.</span>
						</p>
					</section>

					<span class="help is-danger" if={err}>{err}</span>
					<div class="box" each={journal in journals}>
						<article class="media">
							<div class="media-left">
								<figure class="image is-64x64">
									<img src="images/128x128.png" alt="Image">
								</figure>
							</div>
							<div class="media-content">
								<div class="content" onclick={tojournal} style="cursor:pointer;">
									<p>
									<strong>{journal.Title}</strong> <small>@jzs</small> <small>31m</small>
									<br>
									{journal.Description}
									</p>
								</div>
								<nav class="level">
									<div class="level-left">
										<a class="level-item" onclick={createentry}>
											<span class="icon is-small"><i class="fa fa-plus"></i></span>
										</a>
									</div>
									<div class="level-right">
										<span class="level-item" if={!journal.Public}>
											Private
										</span>
									</div>

								</nav>
							</div>
						</article>
					</div>
					<button class="button" onclick={newjournal}>New Journal</button>
				</div>
			</div>
		</div>
	</div>
	<script>
var self = this;
self.journals = [];
self.loading = true;
self.err = null;

self.on('mount', function() {
	// Fetch journals...
	self.loading = true;
	_aj.get("/api/journals", function(data, err) {
		if(err != null) {
			self.err = err;
			self.loading = false;
		} else {
			self.journals = data;
			self.loading = false;
		}
		self.update();
	});
});

self.newjournal = function() {
	route("/journals/create");
};

self.tojournal = function(e) {
	route("/journals/" + e.item.journal.ID);
};
self.createentry = function(e) {
	e.preventDefault();
	route("/journals/" + e.item.journal.ID + "/entries/create");
}
	</script>
</page-dash>
