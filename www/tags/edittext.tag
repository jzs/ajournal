<edittext>
<span class="editfield" onkeydown={keydown} onkeyup={keyup} contenteditable style="white-space:pre" onblur={change}>{val}</span>
<a class="button" onclick={startedit} if={!changes} >
	<span class="icon">
		<i class="fa fa-pencil" aria-hidden="true"></i>
	</span>
</a>
<a class="button" onclick="{save}" if={changes}>
	<span class="icon">
		<i class="fa fa-floppy-o" aria-hidden="true"></i>
	</span>
	<span>Save</span></a>
<script>
var self = this;
self.val = opts.riotValue;
self.edittext = self.val;

self.on("mount", function() {
	self.val = opts.riotValue;
	self.update();
});

self.on("update", function() {
	self.val = opts.riotValue;
});

self.keydown = function(e) {
	switch(e.keyCode) {
		case 13:
			e.preventDefault();
			return;
	}
}

self.keyup = function(e) {
	self.edittext =  e.target.textContent;
	self.changes = self.edittext != self.val;
	var field = self.root.getElementsByClassName("editfield")[0];
}

self.startedit = function(e) {
	var field = self.root.getElementsByClassName("editfield")[0];
	if(self.val.length == 0) {
		field.focus();
		return;
	}
	var rng = document.createRange();
	var sel = window.getSelection();
	rng.setStart(field.childNodes[0], self.val.length);
	rng.collapse(true);
	sel.removeAllRanges();
	sel.addRange(rng);
	field.focus();
}

self.save = function(e) {
	e.preventDefault();
	if(typeof(opts.savefunc) !== 'undefined') {
		opts.savefunc(self.edittext);
		self.changes = false;
	}
}
</script>
</edittext>
