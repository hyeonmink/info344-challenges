"use strict"

const baseURL = "https://api.hyeonmin.me/v1/summary?url="

var queryResults = document.querySelector(".query-results");
var searchForm = document.querySelector(".search-form");
var searchInput = searchForm.querySelector("input");
var searchButton = searchForm.querySelector("button");
var spinner = searchForm.querySelector("header .mdl-spinner");

//renders img, title, and description in result by using the data passed from main function.
function renderWebsite(site){
    var result = document.createElement("div");
    result.className = "result"
    var img = document.createElement("img");
    img.src = (site.image) ? site.image : "./img/no_img.jpg"
    img.alt = site.title;
    img.title = img.alt;

    var text = document.createElement("div")
    text.className = "text";
    
    var title = document.createElement("h4");
    title.innerHTML = (site.title) ? site.title : "no title found";

    var descr = document.createElement("p");
    descr.innerHTML = (site.description) ? site.description : "no description found";

    text.appendChild(title)
        .appendChild(descr);
    result.appendChild(img)
    result.appendChild(text);
    queryResults.appendChild(result);
}

//make sure there is no error before render it.
//if there is error, go to renderError Function.
function render(data){
    if(data.error){
        renderError(data.error);
    } else {
        renderWebsite(data)
    }
}

//handles any error and print it to result.
function renderError(err){
    var message = document.createElement("p");
    message.classList.add("error-message");//add message style

    message.textContent = err.message;
    queryResults.appendChild(message);
}

//when user press the submit, fetch the data.
searchForm.addEventListener("submit", function(evt){
    evt.preventDefault();
    queryResults.innerHTML = ""; //reset result before running the function
    var query = searchInput.value.trim();
    if(query.length <= 0){
        return false;
    }

    fetch(baseURL + query)
        .then((response)=>response.json())
        .then(render)
        .catch(renderError);
    return false;
});
