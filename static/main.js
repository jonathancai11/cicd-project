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
            url: APIurl + "config",
            data: json,
            dataType: "json"
        })
    });

    $("#gitrepos").click(function () { // list repos
        $.ajax({
            method: "GET",
            url: APIurl + "repos",
            success: function (result, status, jgXHR) {
                $("#gitrepositories").text(result)
            }
        })
    });

    $("#gitcreate").click(function () { // create hook
        var hookurl = $("#githookurl").val()
        var pushconfig = document.getElementById("gitpush").checked
        var pullconfig = document.getElementById("gitpullrequest").checked
        if (!pushconfig &&  !pullconfig) {
            window.alert("Must check atleast one event type when creating a hook")
            return
        }
        if (hookurl === "") {
            window.alert("Hook name field need to be non-empty")
            return
        }
        var json = JSON.stringify({
            PushConfig: pushconfig,
            PullConfig: pullconfig,
            HookURL: hookurl,
        })
        $.ajax({
            method: "POST",
            url: APIurl + "create",
            dataType: "json",
            data: json,
        })
    });
    
    $("#gitlist").click(function () { // list hooks
        $.ajax({
            method: "GET",
            url: APIurl + "list",
            dataType: "json",
            success: function (result, status, jgXHR) {
                $("#githooks").empty() // Wipe out all previously loaded hooks
                var p = result
                count = 0
                for (var key in p) {
                    if (p.hasOwnProperty(key)) {
                        count++
                        var hook = "<div class='hook'><p>  Hook #  " +count.toString()+ "-  URL:  "+ key + ", Events:  " +p[key] +"</p></div>";
                        $("#githooks").append(hook)
                        console.log(key + " -> " + p[key]);
                    }                
                }
            }
        })
    });

    $("#gitedit").click(function () { // edit hook
        window.alert("Nothing configured yet")
        return
    });

});