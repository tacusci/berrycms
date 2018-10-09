// RECAPTCHA PLUGIN (v0.0.1)

var RECAPTCHASITEKEY = '6Lds9z0UAAAAAFfF0zUxizO5RB4W3GIExWCUcKW2';

function onPostRender(args) {
    if (args[0] === "/recaptcha-test") {
        page = args[1]
        DebugLog(page + RECAPTCHASITEKEY);
    }
}

function main(uuid) {

}