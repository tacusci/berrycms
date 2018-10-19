// RECAPTCHA PLUGIN (v0.0.1)

var RECAPTCHASITEKEY = '6Lds9z0UAAAAAFfF0zUxizO5RB4W3GIExWCUcKW2';

// args is a list, containing: 0 -> is the page URI, 1 -> is the page header, 2 -> is the page body
function onPreRender(args) {
    if (args[0] === "/recaptcha-test") {
        return { 
            route: args[0], 
            header: args[1].replace("</head>", "<script src=\"https://www.google.com/recaptcha/api.js\" async defer></script></head>"),
            body: args[2] + "<div class=\"g-recaptcha\" data-sitekey=\"" + RECAPTCHASITEKEY + "\"></div>",
            code: 200
        };
    }
    return null;
}

function onPostRender(args) {
    if (args[0] === "/recaptcha-test") {
        page = args[1]
        DebugLog(page + RECAPTCHASITEKEY);
    }
}

function main() {
    InfoLog("Loaded plugin " + UUID);
}