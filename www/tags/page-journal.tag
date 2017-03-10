<page-journal>
	<section class="section">
		<div class="container">
			<section class="section">
				<h3 class="title">Journal: {journal.Title}</h3>
				<p>
				{journal.Description}
				</p>
			</section>
			<section class="section">
				<button class="button" onclick={newentry}>New Entry</button>
			</section>
			<section class="section">

				<div class="box" each={entry in journal.Entries} onclick={onentry} style="cursor:pointer;">
					<article class="media">
						<div class="media-content">
							<div class="content">
								<p>
								<strong>{entry.Title}</strong> <small>@jzs</small> <small>31m</small>
								<br>
								<pre>{entry.Content.substring(0, 200)}...</pre>
								<br>
								<span each={tag in parent.Tags}>{tag}</span>
								</p>
							</div>
							<nav class="level">
								<div class="level-left">
									<a class="level-item">
										<span class="icon is-small"><i class="fa fa-reply"></i></span>
									</a>
									<a class="level-item">
										<span class="icon is-small"><i class="fa fa-retweet"></i></span>
									</a>
									<a class="level-item">
										<span class="icon is-small"><i class="fa fa-heart"></i></span>
									</a>
								</div>
							</nav>
						</div>
					</article>
				</div>

			</section>
		</div>
	</section>
	<script>
var self = this;
self.journal = {
	Title: "Journal title",
	Description: "A fine description of a fine wine",
	Tags: [],
	Entries: [{
		ID: 1,
		Date: "",
		Title: "My first entry",
		Content: "#My first content\n\nHello world",
		Tags: ["diary"],
		Created: "",
		IsPublised: false
	}, {
		ID: 2,
		Title: "My second entry",
		Content: "",
		Tags: []
	}],
	Created: ""
};

self.on('mount', function() {
	_aj.get("/api/journals/" + opts.journalid, function(data, err) {
		if(err != null) {
			self.err = err;
			self.update();
			return;
		}
		self.journal = data;
		self.update();
	});
});

self.newentry = function(e) {
	route("/journals/"+opts.journalid+"/entries/create");
}

self.onentry = function(e) {
	route("/journals/"+opts.journalid+"/entries/"+e.item.entry.ID);
};
	</script>
</page-journal>
