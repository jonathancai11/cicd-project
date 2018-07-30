$(document).ready(function () {
    // window.alert("Window loaded")

    var url = "http://"+window.location.host
    var APIurl = url + "/api/"

    $.ajax({ // Upon loading, get saved token
        method: "GET",
        url: APIurl + "token",
        dataType: "text",
        success: function (result, status, jgXHR) {
            $("#gittoken").text(result)
        },
    })

    $("#gitsave").click(function () { // save config
        var token = $("#gittoken").val()
        if (token == "") {
            window.alert("Enter non-empty token")
            return
        }
        var json = JSON.stringify({
            Token: token
        })
        $.ajax({
            method: "POST",
            url: APIurl + "token",
            data: json,
            dataType: "json"
        })
    });

    $("#gitrepos").click(function () { // list repos
        // var token = $("#gittoken").val()
        $.ajax({
            method: "GET",
            url: APIurl + "repos",
            success: function (result, status, jgXHR) {
                $("#gitrepositories").text(result)
            }
        })
    });

    $("#gitcreate").click(function () { // create hook
        var token = $("#gittoken").val()
        var username = $("#gitusername").val()
        var repo = $("#gitrepo").val()
        var hookname = $("#githookname").val()
        var pushconfig = document.getElementById("gitpush").checked
        var pullconfig = document.getElementById("gitpullrequest").checked
        if (!pushconfig &&  !pullconfig) {
            window.alert("Must check atleast one event type when creating a hook")
            return
        }
        if (token === "" || username === "" || repo === "" || hookname === "") {
            window.alert("Token, username, repository name, and hook name fields need to be non-empty")
            return
        }
        var json = JSON.stringify({
            Token: token,
            Username: username,
            Repo: repo,
            PushConfig: pushconfig,
            PullConfig: pullconfig,
            HookName: hookname,
        })
        $.ajax({
            method: "POST",
            url: APIurl + "create",
            dataType: "json",
            data: json,
        })
    });
    
    $("#gitlist").click(function () { // list hooks
        var token = $("#gittoken").val()
        var username = $("#gitusername").val()
        var repo = $("#gitrepo").val()
        if (token === "" || username === "" || repo === "") {
            window.alert("Token, username, and repository name fields need to be non-empty")
        }
        var json = JSON.stringify({
            Token: token,
            Username: username,
            Repo: repo,
        })
        $.ajax({
            method: "POST",
            url: APIurl + "list",
            dataType: "json",
            data: json,
            success: function (result, status, jgXHR) {
                // window.alert("Hello")
                $("#githooks").text("result")
            }
        })
    });

    // function ListHook() {
    //     return $.ajax({
    //         method: "POST",
    //         url: APIurl + "list",
    //         dataType: "json",
    //         data: json,
    //     })

    // }

    $("#gitedit").click(function () { // edit hook
        window.alert("Nothing configured yet")
        return
        // var token = $("#gittoken").val()
        // $.ajax({
        //     method: "GET",
        //     url: APIurl + "repos",
        //     success: function (result, status, jgXHR) {
        //         $("#gitrepositories").text(result)
        //     }
        // })
    });

});