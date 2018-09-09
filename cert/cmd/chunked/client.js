
// https://developer.mozilla.org/en-US/docs/Web/API/Streams_API/Using_readable_streams

fetch('/')
.then(response => response.body.getReader())
.then(reader => {
	var decoder = new TextDecoder();
	var pump = () => {
		return reader.read().then(chunk => {
			if (chunk.done) {
				return;
			}
			var data = decoder.decode(chunk.value, {stream: true});
			console.log(data);
			return pump();
		})
		.catch(err => console.error(err));
	};
	return pump();
});

// -------------------

var race = (promise, ms) => {
	var timer = null;
	return Promise.race([
		new Promise((resolve, reject) => {
			ms = ms || 900 + Math.random()*200;
			timer = setTimeout(reject, ms, 'timeout:'+ms);
		}),
		promise.then((value) => {
			clearTimeout(timer);
			return value;
		})
	]);
};

race(fetch('/'), 2000)
.then(response => response.body.getReader())
.then(reader => {
	var decoder = new TextDecoder();
	var pump = () => {
		return race(reader.read()).then(chunk => {
			if (chunk.done) {
				return;
			}
			const data = decoder.decode(chunk.value, {stream: true});
			console.log(data);
			return pump();
		})
		.catch(err => {
			reader.cancel();
			console.error(err)
		});
	};
	return pump();
})
.catch(err => console.error(err));
