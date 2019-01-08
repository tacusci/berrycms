// SECOND PLUGIN (v0.0.1)

function onGetRender(args) {
    if (args[0] === "/") {
        InfoLog("Visited index page!")
    }
}

function main() {
    InfoLog("Loaded plugin")
}
