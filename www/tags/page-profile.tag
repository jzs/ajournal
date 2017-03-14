<page-profile>
	<div>
		Hello world!!
		wefewf
	</div>
	<script>
var self = this;

self.on('mount', function() {
	// Fetch profile!
	_aj.get("/api/profile", null, function(data, err) {
		if ( err != null ) {
			return;
		}
		self.profile = data;
	});
});


	</script>
</page-profile>
