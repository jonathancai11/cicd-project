$(document).ready(function () {
    // window.alert("Window loaded")

    var url = "http://"+window.location.host
    var APIurl = url + "/api/"

    // GITHUB
    // $.ajax({ // Upon loading, get saved urlconfig
    //     method: "GET",
    //     url: APIurl,
    //     dataType: "text",
    //     success: function (result, status, jgXHR) {
    //         $("#url").text(result)
    //     },
    // })

    $("#gitrepos").click(function () {
        // window.alert("Clicked list repos")

        var url = $("#giturl").val()
        var token = $("#gittoken").val()
        if (url === "") { // If empty url, alert
            window.alert("Enter non-empty url")
            return
        }
        if (token === "") { // If empty token, alert
            window.alert("Enter non-empty url")
            return
        }
        // window.alert(url + token)
        var json = JSON.stringify({
            Url: url,
            Token: token,
        })
        // Send POST request to api
        $.ajax({
            method: "POST",
            url: APIurl,
            data: json,
            dataType: "json",
            // success: success
        })
    });

});