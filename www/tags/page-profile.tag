<page-profile>
	<div class="container">
		<h3 class="title">Hi {profile.Name}</h3>
		Hello world!!
		wefewf
	</div>
	<script>
var self = this;
self.profile = {};

self.on('mount', function() {
	// Fetch profile!
	_aj.get("/api/profile", function(data, err) {
		if ( err != null ) {
			return;
		}
		self.profile = data;
	});
});


	</script>
</page-profile>
