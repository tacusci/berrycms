
function onGetRender(uri) {
    if (uri === "/main.go") {
        var pageAndResult = session.Get("login_page");
        if (pageAndResult[1]) {
            document.SetHtml(pageAndResult[0]);
        }
    }

    //if the URI starts with images, ignoring the image placeholder section for now
    if (uri.lastIndexOf("/images/", 0) === 0) {
        var imageAndResult = session.Get("logo_image");
        if (imageAndResult[1]) {
            return {
                data: imageAndResult[0]
            }
        }
    }
}

function onPostRecieve(uri, data) {}

//this list of routes gets mapped on plugin load, before main() is called
var routesToRegister = ["/main.go", "/images/{imgfile}"];

function main() {
    var mainPage = files.Read("./main.go");
    if (mainPage !== undefined) {
        if (typeof mainPage === 'string') {
            session.Set("main_page", mainPage);
        }
    }

    var logoImage = files.ReadBytes("./plugins/logo.png");
    if (logoImage !== undefined) {
        if (Array.isArray(logoImage)) {
            session.Set("logo_image", logoImage);
        }
    }
}
