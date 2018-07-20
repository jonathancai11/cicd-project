$(document).ready(function () {
    var url = "http://"+window.location.host
    var APIurl = url + "/api/"
    $.ajax({ // Upon loading, get saved urlconfig
        method: "GET",
        url: APIurl,
        dataType: "text",
        success: function (result, status, jgXHR) {
            $("#url").text(result)
        },
    })

    $("#sve").click(function () {
        // Save the config 
        var urlconfig = $("#url").val() // Obtain url value
        urlconfig = urlconfig.replace(/\s/g, "") // Replace white space
        if (urlconfig === "") { // If empty url, alert
            window.alert("Enter non-empty url")
            return
        } else { // post request to save the data
            var json_str = JSON.stringify({
                Url: urlconfig,
            })
            $.ajax({
                method: "POST",
                url: APIurl + "save",
                data: json_str,
            })
        }
    });

    $("#del").click(function () {
        // Delete the config
        $.ajax({
            method: "DELETE",
            url: APIurl,
        })
    });

    $("#sub").click(function () {
        // Send a message
        var msg = $("#msg").val() // Obtain message from textbox
        var _url = $("#url").val() // Obtain URL from textbox
        var json_str = JSON.stringify({ // JSONify string
            Url: _url,
            Text: msg
        });
        window.alert(json_str)
        if (msg === "") { // If message is empty
            window.alert("enter nonempty message")
        } else {
            // when change dataType to jsonp and have no response header written, 
            // the request sent by jQuery becomes GET
            $.ajax({
                url: APIurl + "msg",
                method: "POST",
                data: json_str,
            })
        }
    })
});