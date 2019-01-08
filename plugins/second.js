// SECOND PLUGIN (v0.0.1)

function onGetRender(args) {
    if (args[0] === "/") {
        //insert form to post
        document.Find("body").AppendHtml("<form action= \"" + args[0] + "\" method=\"post\"><button type=\"submit\">Send</button></form>")
    }
}

function main() {
    InfoLog("Loaded plugin")
}
