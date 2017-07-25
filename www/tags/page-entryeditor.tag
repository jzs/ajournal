<page-entryeditor>
	<div class="section">
		<div class="container">
			<div class="columns">
				<div class="column">
					<label class="label">Title</label>
					<p class="control">
					<input class="input" type="text" placeholder="Title" onkeyup={onTitle} value={entry.Title}>
					</p>
					<label class="label">Date</label>
					<p>
					<datepicker date={entry.Date}></datepicker>
					</p>
					<label class="label">Content</label>
					<p class="control">
					<textarea style="min-height: 200px;" class="textarea" placeholder="Textarea" onkeyup={contentchange} ondrop={drop}>{entry.Content}</textarea>
					<span if={err} class="help is-danger">{err}</span>
					</p>
					<p>
					<p class="control">
					<ul>
						<li each={blobs}>
						<img class="thumb160" src="{Links.Orig}" dragable="true" ondragstart={drag}/>
						</li>
					</ul>
					</p>
					<p class="control">
					<input id="blob" name="blob" type="file" accept="image/png, image/jpeg" onsubmit={uploadfile} onchange={selectfile}/>
					<span if={uploadingfile}>Uploading</span>
					</p>
					<br/>
					<a class="button {is-link : showpreview}" onclick={togglepreview}>Preview</a>
					<button class="button is-pulled-right {is-loading : saving}" onclick={saveEntry}>Save</button>
					</p>
				</div>
				<div class="column" if={showpreview}>
					<label class="label">Preview</label>
					<raw class="markdown" content={preview}></raw>
				</div>
			</div>
		</div>
	</div>
	<script>
var self = this;
self.showpreview = false;
self.preview = "";
self.saving = false;
var converter = new showdown.Converter();

self.entry = {};

self.on('mount', function() {

	self.entry = {
		JournalID: parseInt(opts.journalid),
		Date: "",
		Title: "",
		Content: "",
		Tags: []
	};

	if( opts.entryid != 'create' ) {
		// Fetch entry
		_aj.get("/api/journals/"+self.entry.JournalID+"/entries/"+opts.entryid, function(data, err) {
			// Treat response...
			if(err != null) {
				//TODO Do something with error
				return;
			}
			self.entry = data;
			self.editContent = self.entry.Content;
			self.update();
		});

		_aj.get("/api/journals/"+self.entry.JournalID+"/blobs", function(data, err) {
			if(err != null) {
				return;
			}
			self.blobs = data;
			self.update();
		});
	}


	self.preview = converter.makeHtml(self.entry.Content);
	self.update();
});

self.editContent = "";

self.contentchange = function(e) {
	self.editContent = e.target.value;
	if(self.showpreview) {
		self.preview = converter.makeHtml(self.editContent);
	}
	self.update();
};

self.onTitle = function(e) {
	self.entry.Title = e.target.value;
};

self.saveEntry = function(e) {
	self.saving = true;
	self.entry.Date = self.tags.datepicker.date().toISOString();
	self.entry.Content = self.editContent;
	if(typeof(self.entry.ID) != 'undefined') {
		// Update
		_aj.post("/api/journals/"+self.entry.JournalID+"/entries/"+self.entry.ID, self.entry, function(data, err) {
			self.saving = false;
			// Handle data, err
			if( err != null ) {
				// Present error!
				self.err = err;
				self.update();
				return
			}
			// TODO Handle response
			self.update();
		});

	} else {
		// create new!
		_aj.post("/api/journals/"+self.entry.JournalID+"/entries", self.entry, function(data, err) {
			self.saving = false;
			// Handle data, err
			if( err != null ) {
				// Present error!
				self.err = err;
				self.update();
				return
			}
			self.entry.ID = data.ID;

			// TODO Handle response
			self.update();
		});
	}
	self.entry.Content = self.editContent;
};

self.togglepreview = function(e) {
	self.showpreview = !self.showpreview;
	if( self.showpreview ) {
		self.preview = converter.makeHtml(self.editContent);
	}
	self.update();
}

self.selectfile = function(e) {
	console.log("hello");
	//Submit form
	self.uploadfile(e);
}

self.uploadfile = function(e) {
	e.preventDefault();

	var files = e.target.files;

	var formData = new FormData();
	// Loop through each of the selected files.
	for (var i = 0; i < files.length; i++) {
		var file = files[i];
		// Check the file type.
		if (!file.type.match('image.*')) {
			continue;
		}
		// Add the file to the request.
		formData.append('blobs', file, file.name);
	}
	var xhr = new XMLHttpRequest();
	xhr.open('POST', '/api/journals/'+ self.entry.JournalID +'/blobs', true);
	xhr.onload = function () {
		if (xhr.status === 200) {
			// File(s) uploaded.
			self.uploadingfile = false;
			var data = JSON.parse(xhr.response);
			self.blobs.push(data);
			self.update();
		} else {
			alert('An error occurred!');
		}
	};

	self.uploadingfile = true;
	self.update();

	xhr.send(formData);
};


self.dragItem = null;
self.drag = function(e) {
	self.dragItem = e.item;
};
self.drop = function(e) {
	console.log(e);
	e.preventDefault();
	console.log(self.dragItem);
	console.log(e.target.selectionStart);
	e.toElement.value += "![](" + self.dragItem.Links.Orig + ")"
};

	</script>
</page-entryeditor>
