<page-entryeditor>
	<div class="section">
		<div class="container">
			<div class="columns">
				<div class="column">
					<pre style="min-height:200px" contenteditable="true" onkeyup={contentchange}>{entry.Content}</pre>
					<p>
						<button class="button">Publish</button>
						<button class="button">Save</button>
					</p>
				</div>
				<div class="column">
					<raw class="markdown" content={preview}></raw>
				</div>
			</div>
		</div>
	</div>
	<script>
var self = this;
self.preview = "";
var converter = new showdown.Converter();

self.editContent = "";
self.entry = {
	Content: "#Title\n\nThis is a paragraph"
}

self.on('mount', function() {
	self.preview = converter.makeHtml(self.entry.Content);
	self.update();
});

self.contentchange = function(e) {
	self.editContent = e.target.innerText;
	self.preview = converter.makeHtml(self.editContent);
	self.update();
}
	</script>
</page-entryeditor>
