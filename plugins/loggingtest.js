//BERRYCMS PLUGIN VER(0.0.1)

function onGet(uri) {
	InfoLog("Page at " + uri + " requested");
}

function onPost(uri) {
	console.log("Page at " + uri + " posted to");
}

// main function gets called on plugin load
function main(uuid) {

}