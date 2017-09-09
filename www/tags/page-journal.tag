<page-journal>
	<section class="section">
		<div class="container">
			<section class="section">
				<h3 class="title">Journal: <edittext eid="Title" savefunc={savejournal} value={journal.Title} /></h3>
				<p>
				<edittext eid="Description" isdiv="true" savefunc={savejournal} value={journal.Description}/>
				</p>
				<br/>
				<div class="is-clearfix">
				<a class="button" href="#/users/{opts.username}/journals/{opts.journalid}">View Journal</a>
				</div>
			</section>
			<section class="section">
				<button class="button" onclick={newentry}>New Entry</button>
			</section>
			<section class="section">
				<div class="box" each={entry in entries.Entries} onclick={onentry} style="cursor:pointer;">
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
				<button class="button is-primary" if={entries.HasNext} onclick={loadMore}>Load more</button>
			</section>
		</div>
	</section>
	<script>
var self = this;
self.entries = {Entries: []};
self.journal = {};

self.on('mount', function() {
	_aj.get("/api/journals/" + opts.journalid, function(data, err) {
		if(err != null) {
			self.err = err;
			self.update();
			return;
		}
		self.journal = data;
		self.update();
		self.trigger('data', data);
		getEntries(null);
	});
});

self.newentry = function(e) {
	route("/journals/"+opts.journalid+"/entries/create");
}

self.onentry = function(e) {
	route("/journals/"+opts.journalid+"/entries/"+e.item.entry.ID);
};

self.loadMore = function(e) {
	getEntries(self.entries.Next);
}

self.savejournal = function(id, e) {
	self.journal[id] = e;
	self.update();
	_aj.post("/api/journals/" + self.journal.ID, self.journal, function(data, err) {
		console.log(data);
		console.log(err);
	})

	// Consider what to do... reset contenteditable?

	//
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
		Array.prototype.push.apply(self.entries.Entries, data.Entries);
		self.entries.HasNext = data.HasNext;
		self.entries.Next = data.Next;
		self.update();
	});
};

	</script>
</page-journal>
