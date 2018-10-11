// RECAPTCHA PLUGIN (v0.0.1)

var RECAPTCHASITEKEY = '6Lds9z0UAAAAAFfF0zUxizO5RB4W3GIExWCUcKW2';

function onPreRender(args) {
    if (args[0] === "/recaptcha-test") {
        return { 
            route: args[0], 
            header: args[1].replace("</head>", "<script src=\"https://www.google.com/recaptcha/api.js?render=" + RECAPTCHASITEKEY + "\"></script></head>"),
            body: args[2]
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

function main(uuid) {
    InfoLog(uuid);
}