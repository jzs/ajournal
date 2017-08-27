<latestjournals>
<article class="media" each={j in journals}>
	<div class="media-content">
		<div class="content">
			<p>
			<strong>{j.Entry.Title}</strong>
			<br/>
			<small>{j.Title}</small> <small>{moment(j.Entry.Date).format('YYYY/MM/DD')}</small>
			<br/>
			<small>
				<a href="/app#view/{j.ID}/entries/{j.Entry.ID}">Read more</a>
			</small>
			</p>
		</div>
	</div>
</article>
<script>
var self = this;
self.entries = {Entries: []};
self.journals = [];
_aj.get("/api/journals/latest?limit=3", function(data, err) {
	if(err != null) {
		self.err = err;
		self.update();
		return;
	}
	self.journals = data.Journals;
	self.update();
});

</script>
</latestjournals>
