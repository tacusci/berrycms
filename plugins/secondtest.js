//BERRYCMS PLUGIN VER(0.0.1)

function onGet(uri) {
	InfoLog("Page at " + uri + " requested");
}

function onPreRender(page) {
	InfoLog("Page content looks like: " + page);
}

function onPost(uri) {
	InfoLog("Page at " + uri + " posted to");
}

// main function gets called on plugin load
function main(uuid) {

}