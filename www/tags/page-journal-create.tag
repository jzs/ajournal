<page-journal-create>
	<section class="section">
		<div class="container">
			<h3 class="title">New Journal</h3>
			<label class="label">Title</label>
			<p class="control">
			<input class="input" type="text" placeholder="Title" onkeyup={onTitle} value="">
			</p>
			<label class="label">Description</label>
			<p class="control">
			<textarea class="textarea" placeholder="Description" onkeyup={onDescription}></textarea>
			<span if={errmsg} class="help is-danger">{errmsg}</span>
			</p>
			<p class="control">
			<label class="checkbox">
				<input type="checkbox" checked={journal.Public} onchange={onPublic}>
				Public
			</label>
			</p>

			<p class="control">
			<label class="label">Tags</label>
			<p class="control has-addons">
			<input class="input" type="text" placeholder="Tagname" onkeyup={onjournaltag}>
			<a class="button is-info" onclick={addJournalTag}>
				Add
			</a>
			</p>
			<span class="tag is-large" each={t in journal.Tags}>
				{t}
				<button class="delete" onclick={deleteTag}></button>
			</span>
			</p>

			<button class="button is-success is-pulled-right" onclick={create}>Create</button>
		</div>
	</section>
	<script>
var self = this;
self.errmsg = null;

self.journal = {
	Tags: [],
	Title: "",
	Description: "",
	Public: false
};

self.onTitle = function(e) {
	self.journal.Title = e.target.value;
};
self.onDescription = function(e) {
	self.journal.Description = e.target.value;
};
self.onPublic = function(e) {
	self.journal.Public = e.target.checked;
};

self.journaltag = "";
self.onjournaltag = function(e) {
	e.preventDefault();
	self.journaltag = e.target.value;
};
self.addJournalTag = function() {
	self.journal.Tags.push(self.journaltag);
	self.journaltag = "";
	self.update();
};
self.deleteTag = function(e) {
	var index = self.journal.Tags.indexOf(e.item.t);
	self.journal.Tags.splice(index, 1);
}

self.create = function() {
	_aj.post("/api/journals", self.journal, function(data,err) {
		if( err != null ) {
			self.errmsg = err;
			self.update();
			return;
		}
		route("/journals/" + data.ID);
	});
};
	</script>
</page-journal-create>
