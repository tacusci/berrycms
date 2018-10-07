//BERRYCMS PLUGIN VER(0.0.1)

var pluginGUUID = '';

// onLoad gets called during a page load, but before render
function onLoad() {
	InfoLog("This is logging some information");
	DebugLog("This is logging some debug information");
	ErrorLog("This is logging an error");
}

// main function gets called on plugin load
function main() {
	console.log("Plugin of UUID -> ");
}