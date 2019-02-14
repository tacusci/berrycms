
function onGetRender(uri, vars) {
    if (uri === "/main.go") {
        var pageAndResult = session.Get("main_page");
        if (pageAndResult[1]) {
            document.SetHtml(pageAndResult[0]);
        }
    }

    if (uri.lastIndexOf("/images/", 0) === 0) {
        //retrieve image to serve from in memory cache
        if (vars["imgfilename"] === "logo.png") {
            var imageAndResult = session.Get("logo_image");
            if (imageAndResult[1]) {
                return {
                    data: imageAndResult[0]
                }
            }
        }

        //load image to serve from disk
        if (vars["imgfilename"] === "logofromdisk.png") {
            var logoImage = files.ReadBytes("./plugins/logo.png");
            if (logoImage !== undefined) {
                //is the object a byte array basically
                return {
                    data: logoImage
                }
            }
        }

        return {
            code: 404
        }
    }
}

function onPostRecieve(uri, data) {}

//this list of routes gets mapped on plugin load, before main() is called
var routesToRegister = ["/main.go", "/images/{imgfilename}"];

function main() {
    var mainPage = files.Read("./main.go");
    if (mainPage !== undefined) {
        if (typeof mainPage === 'string') {
            session.Set("main_page", mainPage);
        }
    }

    var logoImage = files.ReadBytes("./plugins/logo.png");
    if (logoImage !== undefined) {
        //is the object a byte array basically
        if (Array.isArray(logoImage)) {
            session.Set("logo_image", logoImage);
        }
    }
}
