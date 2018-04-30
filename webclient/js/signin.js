"use strict"
var baseURL = "https://api.hyeonmin.me/v1/sessions";
var loginForm = document.getElementById("login-form")

loginForm.addEventListener("submit", (evt)=>{
    evt.preventDefault();
    var request = {};
    request.email = $('#email').val();
    request.password = $('#password').val();


    $.ajax({
    type: "POST",
    url: baseURL,
    data: JSON.stringify(request),
    success: function(msg){
        alert("Welcome to MinChat!");
        window.location.href="main.html"
    },
    complete: (rest)=>{
        var authkey = rest.getResponseHeader("authorization");
        localStorage.setItem('authkey',authkey);
        
    },
    error: function(XMLHttpRequest, textStatus, errorThrown) {
        alert(XMLHttpRequest.statusText);
    }
    });
});