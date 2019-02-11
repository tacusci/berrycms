// RECAPTCHA PLUGIN (v0.0.1)

var RECAPTCHASITEKEY = '6Lds9z0UAAAAAFfF0zUxizO5RB4W3GIExWCUcKW2';

function onGetRender(uri, vars) {
    if (uri === "/redirect-test") {
        return {
            route: "/",
            code: 302
        }
    }

    if (uri === "/recaptcha-test") {
        document.Find("head").AppendHtml("<script src=\"https://www.google.com/recaptcha/api.js\" async defer></script>")
        document.Find("body").AppendHtml("<form action= \"" + uri + "\" method=\"post\"><input name=\"sometext\" type=\"text\"><button type=\"submit\">Send</button></form>")
        document.Find("form").AppendHtml("<div class=\"g-recaptcha\" data-sitekey=\"" + RECAPTCHASITEKEY + "\"></div>")
        document.Find("body").AppendHtml("<img src='http://localhost:8080/images/logo.png'/>");
    }

    if (uri === "/plugin-page-test") {
        var data = files.Read("/main.go");
        if (data !== undefined) {
            if (typeof data === 'string') {
                document.SetHtml("<h2>" + data + "</h2>")
            }
        }
    }

    return null;
}

function onPostRecieve(uri, data) {
    if (uri === "/recaptcha-test") {
        logging.Info("Recieved post request");
        console.log(data["g-recaptcha-response"]);
        return {
            route: uri 
        }
    }
}

//this list of routes gets mapped on plugin load, before main() is called
var routesToRegister = ["/recaptcha-test", "/plugin-page-test", "/redirect-test"];

function main() {
    logging.Info("Loaded plugin")

    for (var i = 0; i < 20; i++) {
        robots.Add("Disallow: /cheesecake-test")
    }

    for (var j = 0; j < 20; j++) {
        robots.Del("Disallow: /cheesecake-test")
    }
}
