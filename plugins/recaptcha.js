// RECAPTCHA PLUGIN (v0.0.1)

var RECAPTCHASITEKEY = '6Lds9z0UAAAAAFfF0zUxizO5RB4W3GIExWCUcKW2';

// args is a list, it only currently contains the URI of the requested page 
function onGetRender(args) {
    if (args[0] === "/redirect-test") {
        return {
            route: "/",
            code: 302
        };
    }

    if (args[0] === "/recaptcha-test") {
        document.Find("head").AppendHtml("<script src=\"https://www.google.com/recaptcha/api.js\" async defer></script>")
        document.Find("body").AppendHtml("<form action= \"" + args[0] + "\" method=\"post\"><input name=\"sometext\" type=\"text\"><button type=\"submit\">Send</button></form>")
        document.Find("body").AppendHtml("<div class=\"g-recaptcha\" data-sitekey=\"" + RECAPTCHASITEKEY + "\"></div>")
    }

    return null;
}

function onPostRecieve(args) {
    if (args[0] === "/recaptcha-test") {
        InfoLog("Recieved post request");
        console.log(args[1]["sometext"]);
        return {
            route: "/recaptcha-test"
        }
    }
}

function main() {
    InfoLog("Loaded plugin");
    InfoLog("Adding \"Disallow: /some-test-uri\" to robots.txt")
    AddToRobots("Disallow: /some-test-uri")
    InfoLog("Deleting \"Disallow: /some-test-uri\" from robots.txt")
    DelFromRobots("Disallow: /some-test-uri")
}
