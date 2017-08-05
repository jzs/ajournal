<page-viewjournal>
	<section class="section">
		<div class="container">
			<section class="section">
				<h3 class="title">Journal: {journal.Title}</h3>
				<p>
				{journal.Description}
				</p>
			</section>
			<section class="section">

				<div class="box" each={entry in entries.Entries} onclick={onentry} style="cursor:pointer;">
					<article class="media">
						<div class="media-content">
							<div class="content">
								<p>
								<h3 class="title">{entry.Title}</h3>
								<small>@{user.Username}</small> <small>{entry.Date}</small>
								<raw class="markdown" content={entry.preview} />
								<br>
								<p if="{entry.Content.length > 2000}"><span >Read more</span></p>
								<span each={tag in parent.Tags}>{tag}</span>
								</p>
							</div>
							<!--
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
							-->
						</div>
					</article>
				</div>
				<button class="button is-primary" if={entries.HasNext} onclick={loadMore}>Load more</button>
			</section>
		</div>
	</section>
	<script>
var self = this;

self.user = {};
self.entries = {Entries: []};
self.journal = {
	Tags: [],
	Entries: []
};
//simpleLineBreaks
var converter = new showdown.Converter();
converter.setFlavor('github');

self.on('mount', function() {
	_aj.get("/api/journals/" + opts.journalid, function(data, err) {
		if(err != null) {
			self.err = err;
			self.update();
			return;
		}
		self.journal = data;

		self.update();

		// Initial fetch of entries...
		getEntries(null);

		_aj.get("/api/users/" + opts.username, function(data, err) {
			if(err != null) {
				self.err = err;
				self.update();
				return;
			}
			self.user = data;
			self.update();
		});
	});
});

self.onentry = function(e) {
	route("/view/"+opts.journalid+"/entries/"+e.item.entry.ID);
};

self.loadMore = function(e) {
	getEntries(self.entries.Next);
}

var getEntries = function(from) {
	var fromstr = "?limit=10";
	if(from != null) {
		fromstr += "&from=" +from;
	}
	_aj.get("/api/journals/" + opts.journalid + "/entries" + fromstr, function(data, err) {
		if(err != null) {
			self.err = err;
			self.update();
			return;
		}

		for(var i = 0; i < data.Entries.length; i++) {
			var j = data.Entries[i];
			j.preview = converter.makeHtml(j.Content.substring(0, 2000));
			//self.journal.Entries[i] = j
		}


		Array.prototype.push.apply(self.entries.Entries, data.Entries);
		self.entries.HasNext = data.HasNext;
		self.entries.Next = data.Next;
		self.update();
	});
};

	</script>

</page-viewjournal>
