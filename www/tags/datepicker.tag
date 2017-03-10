<datepicker>
<input class="input" type="text" placeholder="yyyy/mm/dd" onblur={onblur} onfocus={onfocus} onkeydown={onkeydown} onmouseup={onmouseup}/>
<script>
var self = this;

var datestr = "yyyy/mm/dd";

self.onfocus = function(e) {
	// input field blurred.
	e.target.value = datestr;
	self.update();
};

var year = "yyyy";
var month = "mm";
var day = "dd";
var KEY_DELETE = 8;
var KEY_0 = 48;
var KEY_1 = 49;
var KEY_3 = 51;
var KEY_9 = 57;
var ARROW_LEFT = 37;
var ARROW_RIGHT = 39;

self.onkeydown = function(e) {
	e.preventDefault();
	var start = e.target.selectionStart;
	var end = e.target.selectionEnd;
	if( start <= 4 ) {
		if(e.keyCode == KEY_DELETE) {
			if(year.length > 0) {
				year = year.substring(0, year.length -1);
				datestr = year + "/" + month + "/" + day;
				e.target.value = datestr;
				e.target.setSelectionRange(start-1,start-1);
			}
		}
		if(e.keyCode >= KEY_0 && e.keyCode <= KEY_9) {
			if(year.length == 4) {
				year = "";
			}
			year = year + e.key;
			datestr = year + "/" + month + "/" + day;
			e.target.value = datestr;
			e.target.setSelectionRange(start+1,start+1);

			if(year.length >= 4) {
				e.target.setSelectionRange(5,7);
			}
		}
	} else if( start <= 7 ) {
		if(e.keyCode == KEY_DELETE) {
			if(month.length > 0) {
				month = month.substring(0, month.length -1);
				datestr = year + "/" + month + "/" + day;
				e.target.value = datestr;
				e.target.setSelectionRange(start-1,start-1);
			}
		}
		if(e.keyCode >= KEY_0 && e.keyCode <= KEY_9) {
			// We can't have a month higher than 12...
			if(month.length == 0 && e.keyCode > KEY_1) {
				return;
			}
			if(month.length == 2) {
				month = "";
			}
			month = month + e.key;
			datestr = year + "/" + month + "/" + day;
			e.target.value = datestr;
			e.target.setSelectionRange(start+1,start+1);

			if(month.length >= 2) {
				e.target.setSelectionRange(8,10);
			}

		}
		if(month.length >= 2) {
			e.target.setSelectionRange(8,10);
		}
	} else if( start <= 10 ) {
		if(e.keyCode == KEY_DELETE) {
			if(day.length > 0) {
				day = day.substring(0, day.length -1);
				datestr = year + "/" + month + "/" + day;
				e.target.value = datestr;
				e.target.setSelectionRange(start-1,start-1);
			}
		}
		if(e.keyCode >= KEY_0 && e.keyCode <= KEY_9) {
			// We can't have a month higher than 12...
			if(day.length == 0 && e.keyCode > KEY_3) {
				return;
			}
			if(day.length >= 2 && start != end) {
				day = "";
			} else if(day.length == 2) {
				return;
			}
			day = day + e.key;
			datestr = year + "/" + month + "/" + day;
			e.target.value = datestr;
			e.target.setSelectionRange(start+1,start+1);
		}
		// Else if for time mayne? ...
	} else if( start <= 13 ) { 
	} else if( start <= 16 ) {
	}
};

self.onblur = function(e) {
};

self.onmouseup = function(e) {
	e.preventDefault();
	var start = e.target.selectionStart;
	var end = e.target.selectionEnd;
	if( start <= 3 ) {
		e.target.setSelectionRange(0,4);
	} else if( start <= 6 ) {
		e.target.setSelectionRange(5,7);
	} else if( start <= 9 ) {
		e.target.setSelectionRange(8,10);
	} else if( start <= 12 ) {
		e.target.setSelectionRange(11,13);
	} else if( start <= 15 ) {
		e.target.setSelectionRange(14,16);
	}
};

self.date = function() {
	return new Date(Date.UTC(year,month,day));
};
</script>
</datepicker>
