<page-dash>
	<div class="section">
		<!-- 
		// TODO: Show current journals that you are writing.
		// Add entry to journal
		// Implement entry editor
		// Account page
		-->
		<div class="container">
			<div class="columns">
				<div class="column">
					<h3 class="title">My Journals</h3>
					<p>
					Welcome to journal...
					</p>
					<div class="box" each={journal in journals}>
						<article class="media">
							<div class="media-left">
								<figure class="image is-64x64">
									<img src="images/128x128.png" alt="Image">
								</figure>
							</div>
							<div class="media-content">
								<div class="content" onclick={} style="cursor:pointer;">
									<p>
									<strong>{journal.Title}</strong> <small>@jzs</small> <small>31m</small>
									<br>
									Description here...
									</p>
								</div>
								<nav class="level">
									<div class="level-left">
										<a class="level-item">
											<span class="icon is-small"><i class="fa fa-plus"></i></span>
										</a>
									</div>
								</nav>
							</div>
						</article>
					</div>
					<button class="button">New Journal</button>
				</div>
			</div>
		</div>
	</div>
	<script>
var self = this;
self.journals = [{Title: "A new beginning"},{Title: "A journey abroad"}];
	</script>
</page-dash>
