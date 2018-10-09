//BERRYCMS PLUGIN VER(0.0.1)

function onGet(uri, page) {
	InfoLog("Page at " + uri + " requested");
}

function onPreRender(uri, page) {
	InfoLog("Page content looks like: " + page);
}

function onPost(uri) {
	InfoLog("Page at " + uri + " posted to");
}

// main function gets called on plugin load
function main(uuid) {

}