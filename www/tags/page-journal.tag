<page-journal>
	<section class="section">
		<div class="container">
			<section class="section">
				<h3 class="title">Journal: <edittext eid="Title" savefunc={savejournal} value={journal.Title} /></h3>
				<p>
				<edittext eid="Description" isdiv="true" savefunc={savejournal} value={journal.Description}/>
				</p>
			</section>
			<section class="section">
				<a class="button" href="#/users/{opts.username}/journals/{opts.journalid}">View Journal</a>
				<button class="button" onclick={newentry}>New Entry</button>
			</section>
			<section class="section">
				<div class="box" each={entry in entries.Entries} onclick={onentry} style="cursor:pointer;">
					<article class="media">
						<div class="media-content">
							<div class="content">
								<p>
								<strong>{entry.Title}</strong> 
								<br/>
								<small>{moment(entry.Date).format('LL')}</small>
								<br>
								<div class="entry-preview-text">{entry.Content.substring(0, 400)}</div>
								<br>
								<span each={tag in parent.Tags}>{tag}</span>
								</p>
							</div>
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
