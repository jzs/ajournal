<page-viewjournalentry>
	<div class="section">
		<div class="container">
			<div class="columns">
				<div class="column">
					<h1 class="title">{entry.Title}</h1>
					<h2 class="subtitle">{entry.Date}</h2>
					<hr />
					<div class="content">
						<raw class="markdown" content={entry.HtmlContent}></raw>
					</div>
				</div>
			</div>
		</div>
	</div>
	<script>
var self = this;
self.entry = {};

self.on('mount', function() {
	self.entry = {
		JournalID: parseInt(opts.journalid),
		Date: "",
		Title: "",
		Content: "",
		Tags: []
	};
	// Fetch entry
	_aj.get("/api/journals/"+self.entry.JournalID+"/entries/"+opts.entryid, function(data, err) {
		// Treat response...
		if(err != null) {
			//TODO Do something with error
			return;
		}
		self.entry = data;
		self.update();
	});

});

	</script>
</page-viewjournalentry>
