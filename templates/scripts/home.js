function setCookie(cname, cvalue, exmin) {

}

function getCookie(cname) {
    return "";
}

$(document).read(function () {
    /* 定义常量 */
    // cookie持续时间为 30min
    DEFAULT_COOKIE_EXPIRE_TIME = 30;

    /* 定义变量 */
    uname = ``;
    session = ``;
    uid = 0;
    currentVideo = null;
    listedVideos = null;

    uname = getCookie(`username`)
    session = getCookie(`session`)

    /* home page 事件处理 */
    // 注册
    $("#regbtn").on('click', function (e) {

    });
    // 登录
    $("#signbtn").on('click', function (e) {

    });

    $("#signinhref").on('click', function () {
        $("#regsubmit").hide();
        $("#signinsubmit").show();
    })

    $("#registerhref").on('click', function () {
        $("#signinsubmit").hide();
        $("#regsubmit").show();
    })

    /* user home page 事件处理 */
})