
window.addEventListener('load', function() {
	new Vue({
		el: '#todo',
		components: {
			TodoApp: httpVueLoader('./components/todo-app.vue')
		}
	});
});

/*
if ('serviceWorker' in navigator) {
	window.addEventListener('load', function() {
		navigator.serviceWorker.register('/sw.js')
		.then(function(registration) {
			console.log('sw registration successful:', registration);

		}, function(err) {
			console.log('sw registration failed:', err);
		});
	});
} else {
	console.log('sw registration skipped: not supported');
}
*/
