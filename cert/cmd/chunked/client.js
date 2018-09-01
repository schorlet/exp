
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
			var data = decoder.decode(chunk.value);
			console.log(data);
			return pump();
		})
		.catch(err => console.log(err));
	};
	return pump();
});
