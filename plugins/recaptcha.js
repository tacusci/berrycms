// RECAPTCHA PLUGIN (v0.0.1)

var RECAPTCHASITEKEY = '6Lds9z0UAAAAAFfF0zUxizO5RB4W3GIExWCUcKW2';

// args is a list, containing: 0 -> is the page URI, 1 -> is the page header, 2 -> is the page body
function onGetRender(args) {
    if (args[0] === "/redirect-test") {
        return {
            route: "/",
            code: 302
        };
    }

    if (args[0] === "/insert-div-test") {
        //TODO: Add document modification test code
        document.Find("head").AppendHtml("<div>Div append test</div>")
        return null;
    }
    return null;
}

function onPostRecieve(args) {}

function main() {
    InfoLog("Loaded plugin " + UUID);
}
