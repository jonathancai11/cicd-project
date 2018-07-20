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

    $("#gitsave").click(function () {
        // window.alert("Clicked list repos")
        var token = $("#gittoken").val()
        if (token === "") { // If empty token, alert
            window.alert("Enter non-empty token")
            return
        }
        // window.alert(url + token)
        var json = JSON.stringify({
            Token: token,
        })
        // Send POST request to api
        $.ajax({
            method: "POST",
            url: APIurl,
            data: json,
            dataType: "json",
        })
    });

    $("#gitrepos").click(function () {
        // window.alert("Clicked list repos")
        $.ajax({
            method: "GET",
            url: APIurl,
            success: function (result, status, jgXHR) {
                $("#gitcommit").text(result)
            }
        })
    });

});