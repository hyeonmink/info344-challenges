"use strict"

var baseURL = "https://api.hyeonmin.me/v1/users";
var loginForm = document.getElementById("login-form")
loginForm.addEventListener("submit", (evt)=>{
    evt.preventDefault();

    var request = {};
    request.userName = $('#username').val();
    request.password = $('#password').val();
    request.passwordConf = $('#cpassword').val();
    request.firstName = $('#firstname').val();
    request.lastName = $('#lastname').val();
    request.email = $('#email').val();

    $.ajax({
    type: "POST",
    url: baseURL,
    data: JSON.stringify(request),
    success: function(msg){
        alert("You information has been successfully saved!");
        window.location.href="index.html"
    },
    error: function(XMLHttpRequest, textStatus, errorThrown) {
        alert(XMLHttpRequest.statusText);
    }
    });
})