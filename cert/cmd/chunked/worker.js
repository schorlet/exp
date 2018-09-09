
// https://developers.google.com/web/fundamentals/primers/async-functions
// https://developer.mozilla.org/en-US/docs/Web/API/Web_Workers_API/Using_web_workers

async function get(url) {
	const response = await fetch(url);
	const reader = response.body.getReader();

	const decoder = new TextDecoder();
	var chunk = await reader.read();

	while (!chunk.done) {
		let data = decoder.decode(chunk.value, {stream: true});
		postMessage(data);
		chunk = await reader.read();
	}
}

onmessage = function(e) {
	if (e.data == "start") {
		get('/').catch(err => console.error(err));
	} else {
		close();
	}
}
