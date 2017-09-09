<edittext>
<div if={opts.isdiv} class="editfield" onkeydown={keydown} onkeyup={keyup} contenteditable style="white-space:pre">{opts.riotValue}</div>
<span if={!opts.isdiv} class="editfield" onkeydown={keydown} onkeyup={keyup} contenteditable style="white-space:pre">{opts.riotValue}</span>
<span class="is-pulled-right">
	<a class="button" onclick={startedit} if={!changes} >
		<span class="icon">
			<i class="fa fa-pencil" aria-hidden="true"></i>
		</span>
		<span>
			Edit
		</span>
	</a>
	<a class="button" onclick="{save}" if={changes}>
		<span class="icon">
			<i class="fa fa-floppy-o" aria-hidden="true"></i>
		</span>
		<span>Save</span>
	</a>
</span>
<script>
var self = this;

self.keydown = function(e) {
	if(opts.isdiv) {
		return;
	}
	switch(e.keyCode) {
		case 13:
			e.preventDefault();
			return;
	}
}

self.keyup = function(e) {
	self.edittext = e.target.innerText;
	self.changes = self.edittext != self.val;
}

self.startedit = function(e) {
	var field = self.root.getElementsByClassName("editfield")[0];
	if(opts.riotValue.length == 0) {
		field.focus();
		return;
	}

	var rng = document.createRange();
	var sel = window.getSelection();
	rng.setStart(field.childNodes[0], opts.riotValue.length);
	rng.collapse(true);
	sel.removeAllRanges();
	sel.addRange(rng);
	field.focus();
}

self.save = function(e) {
	e.preventDefault();
	if(typeof(opts.savefunc) !== 'undefined') {
		// Seems like a hack but it works. Except it breaks riot updates to value.
		var field = self.root.getElementsByClassName("editfield")[0];
		field.innerText = self.edittext;
		opts.savefunc(opts.eid, self.edittext);
	}
	self.changes = false;
}
</script>
</edittext>
