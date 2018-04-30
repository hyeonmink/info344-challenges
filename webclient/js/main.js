"use strict"
var form = document.getElementById("user_info")
var authkey = localStorage.getItem('authkey')
if (authkey) {
    const userURL = "https://api.hyeonmin.me/v1/users"
    const logOutURL = "https://api.hyeonmin.me/v1/sessions/mine"
    const channelURL = "https://api.hyeonmin.me/v1/channels"
    const messageURL = "https://api.hyeonmin.me/v1/messages"
    const chatbotURL = "https://api.hyeonmin.me/v1/bot"
    var nav = $('#info')
    var eventsDiv = $('#events')
    var chanList = $('#chanList')
    var channelSelecter = $('#channelSelecter')
    var websock = new WebSocket("wss://api.hyeonmin.me/v1/websocket")
    var login_info = {}
    var currentUser
    let currentChan = {};
    let channelCollection = {}
    let allUsers = {};

    $.ajax({
        type: "GET",
        contentType: "application/json",
        url: userURL,
        headers: {
            "Authorization": authkey
        },
        success: function(msg) {
            for (var i = 0; i < msg.length; i++) {
                allUsers[msg[i].id] = msg[i]
            }
        },
        error: function(XMLHttpRequest, textStatus, errorThrown) {
            alert(XMLHttpRequest.statusText);
        }
    });



    $.ajax({
        contentType: "application/json",
        url: userURL + "/me",
        headers: {
            "Authorization": authkey
        },
        complete: (rest) => {
            login_info = JSON.parse(rest.responseText)
            currentUser = login_info.User
            display(currentUser)
        },
        error: function(XMLHttpRequest, textStatus, errorThrown) {
            window.location.href = "index.html"
        }
    });

    function renderChan() {
        $.ajax({
            type: "GET",
            contentType: "application/json",
            url: channelURL,
            headers: {
                "Authorization": authkey
            },
            complete: (rest) => {
                displayList(JSON.parse(rest.responseText))
            },
            error: function(XMLHttpRequest, textStatus, errorThrown) {
                alert(XMLHttpRequest.statusText);
            }
        });
    }

    function displayList(channels) {
        $('#chanList').empty();
        if (!currentChan.name) {
            currentChan = channels[0]
            getMessages()
        }
        channelCollection = {}
        for (var i = 0; i < channels.length; i++) {
            var channel = channels[i]
            channelCollection[channel.id] = channel
            $('#chanList').append(`<li><button id=${channel.id} class="mdl-button mdl-js-button mdl-button--raised channelbtn">${channel.name}</button></li>`)

            $(`#${channel.id}`).click(function() {
                currentChan = channelCollection[this.id]
                getMessages()
                renderEditPanel()
            });
        }
        renderEditPanel()
    }

    function renderEditPanel() {
        channelSelecter.html(currentChan.name)
        $('#editChan').prop('disabled', currentChan.creatorID != currentUser.id)
        $('.editChanContainer').addClass("is-dirty")
        $('#chanEdit #editChanName').val(currentChan.name)
        $('#chanEdit #editChanDescr').val(currentChan.descr)
        if (currentChan.private) {
            $('#edit-public').removeClass("is-checked");
            $('#edit-private').addClass("is-checked");
        } else {
            $('#edit-private').removeClass("is-checked");
            $('#edit-public').addClass("is-checked");
        }

        $('.mdl-menu__container').removeClass('is-visible')

        currentChan = channelCollection[currentChan.id]
        if (currentChan) {
            $('#send').prop('disabled', !contains(currentChan, currentUser))
            $('#textBox').prop('disabled', !contains(currentChan, currentUser))

            if (currentChan.creatorID == currentUser.id) {
                $('#join-leave').css("background-color","#bdbdbd")
                $('#join-leave').prop('disabled', true)
                $('#join-leave').html("You can't leave")
            } else {
                $('#join-leave').prop('disabled', false)
                if ($('#send').prop('disabled')) {
                    $('#join-leave').html("JOIN")
                    $('#join-leave').css("background-color","#f06292")
                } else {
                    $('#join-leave').html("LEAVE")
                    $('#join-leave').css("background-color","#5c6bc0")

                };
            }


        }
        renderPeople()
    }

    function display(user) {
        nav.empty()
            .append(`<img src=${user.photoURL} alt="">`)
            .append(`<p>Hello ${user.firstName}!</p>`)
    }

    $("#signout").click(function() {
        event.preventDefault();
        $.ajax({
            type: "DELETE",
            contentType: "application/json",
            url: logOutURL,
            headers: {
                "Authorization": localStorage.getItem('authkey')
            },
            complete: (rest) => {
                alert(rest.responseText);
                localStorage.removeItem('authkey')
                window.location.href = "index.html"
            },
            error: function(XMLHttpRequest, textStatus, errorThrown) {
                alert(XMLHttpRequest.statusText);
            }
        });
    });

    $("#editbtn").click(function() {
        event.preventDefault();
        window.location.href = "edit.html"
    })


    $("#change_btn").click(function() {
        event.preventDefault();
        var request = {};
        request.firstName = $('#firstname').val();
        request.lastName = $('#lastname').val();
        $.ajax({
            type: "PATCH",
            contentType: "application/json",
            url: userURL + "/me",
            data: JSON.stringify(request),
            headers: {
                "Authorization": localStorage.getItem('authkey')
            },
            success: function(msg) {
                alert("You information has been successfully changed!");
                window.location.href = "main.html"
            },
            error: function(XMLHttpRequest, textStatus, errorThrown) {
                alert(XMLHttpRequest.statusText);
            }
        });
    })

    $("#send").click(function(e) {
        e.preventDefault();
        var request = {
            ChannelID: currentChan.id,
            body: $('#textBox').val()
        }
        $.ajax({
            type: "POST",
            contentType: "application/json",
            url: messageURL,
            data: JSON.stringify(request),
            headers: {
                "Authorization": authkey
            },
            success: function(msg) {},
            error: function(XMLHttpRequest, textStatus, errorThrown) {
                alert(XMLHttpRequest.statusText);
            }
        });
        $('#textBox').val('');
    })

    $("#newChanAdd").click(function(e) {
        e.stopPropagation();
    })

    $("#chanEdit").click(function(e) {
        e.stopPropagation();
    })

    $("#addPersonList").click(function(e) {
        e.stopPropagation();
    })

    $("#chanEdit").click(function(e) {
        e.stopPropagation();
    })

    $("#newChanAdd").submit(function(e) {
        e.preventDefault()
        var request = {
            name: ($("#chanName").val()) ? $("#chanName").val() : " ",
            descr: ($("#chanDescr").val()) ? $("#chanDescr").val() : " ",
            private: ($('input[name=radio]:checked', '#newChanAdd').val() == "private")
        }
        $.ajax({
            type: "POST",
            contentType: "application/json",
            url: channelURL,
            data: JSON.stringify(request),
            headers: {
                "Authorization": localStorage.getItem('authkey')
            },
            success: function(msg) {
                alert("New Channel has been successfully added!");
                currentChan = msg;
                $('.mdl-menu__container').removeClass("is-visible")
                renderChan();
                getMessages();
            },
            error: function(XMLHttpRequest, textStatus, errorThrown) {
                alert(XMLHttpRequest.statusText);
            }
        });
    })

    $('#editChanBtn').click((e) => {
        e.preventDefault();
        var editURL = `${channelURL}/${currentChan.id}`
        var chanName = $('#chanEdit #editChanName').val() ? $('#chanEdit #editChanName').val() : currentChan.name
        var chanDescr = $('#chanEdit #editChanDescr').val() ? $('#chanEdit #editChanDescr').val() : currentChan.descr
        var privBtn = ($('input[name=radio]:checked', '#chanEdit').val() == "private");
        var request = {
            name: chanName,
            descr: chanDescr,
            private: privBtn
        }

        $.ajax({
            type: "PATCH",
            contentType: "application/json",
            url: editURL,
            data: JSON.stringify(request),
            headers: {
                "Authorization": localStorage.getItem('authkey')
            },
            success: function(msg) {
                alert("your channel information has been successfully changed!");
                currentChan = msg;
                channelCollection[msg.id] = currentChan
                renderChan();
            },
            error: function(XMLHttpRequest, textStatus, errorThrown) {
                alert(XMLHttpRequest.statusText);
            }
        });
    })


    $('#deleteChanBtn').click((e) => {
        e.preventDefault();
        var deleteUrl = `${channelURL}/${currentChan.id}`
        $.ajax({
            type: "DELETE",
            contentType: "application/json",
            url: deleteUrl,
            headers: {
                "Authorization": localStorage.getItem('authkey')
            },
            success: function(msg) {
                currentChan = channelCollection[Object.keys(channelCollection)[0]];
                getMessages();
                renderChan();
            },
            error: function(XMLHttpRequest, textStatus, errorThrown) {
                alert(XMLHttpRequest.statusText);
            }
        });
    })

    $('#join-leave').click((e) => {
        //$('#join-leave').prop('joined') ? "LINK" : "UNLINK"
        e.preventDefault();
        var linkURL = `${channelURL}/${currentChan.id}`
        $.ajax({
            type: $('#send').prop('disabled') ? "LINK" : "UNLINK",
            contentType: "application/json",
            url: linkURL,
            headers: {
                "Authorization": authkey
            },
            success: function(msg) {
                renderChan();
            },
            error: function(XMLHttpRequest, textStatus, errorThrown) {
                alert(XMLHttpRequest.statusText);
            }
        });
    })

    function contains(channel, user) {
        return channel.members.indexOf(user.id) != -1
    }

    function getMessages() {
        var linkURL = `${channelURL}/${currentChan.id}`
        $.ajax({
            type: "GET",
            contentType: "application/json",
            url: linkURL,
            headers: {
                "Authorization": authkey
            },
            success: function(msg) {
                getMessagesHelper(msg)
            },
            error: function(XMLHttpRequest, textStatus, errorThrown) {
                alert(XMLHttpRequest.statusText);
            }
        });
    }

    function getMessagesHelper(msg) {
        $('#msglist').empty();
        for (var i = 0; i < msg.length; i++) {
            var message = msg[i]
            $('#msglist').append(
                `<ul id = "outMsgBox" class="demo-list-three mdl-list">
                    <li id = "innerMsgBox"class="mdl-list__item mdl-list__item--three-line">
                        <span class="mdl-list__item-primary-content">
                            <img src="${allUsers[message.creatorID].photoURL}" class="material-icons mdl-list__item-avatar">
                            <span id = "message-username">${allUsers[message.creatorID].firstName} ${allUsers[message.creatorID].lastName} 
                                <a id="edit${message.id}"><i class="material-icons msgEditBtn">edit</i></a>
                                <a id="del${message.id}"><i class="material-icons msgDeleteBtn">delete</i></a></span>
                            <span id="body${message.id}"class="mdl-list__item-text-body">
                                ${message.body}
                            </span>
                        </span>
                    </li>
                <ul>`
            )
            if (message.creatorID != currentUser.id) {
                $(`#edit${message.id}`).remove()
                $(`#del${message.id}`).remove()
            }
            $(`#edit${message.id}`).click(function() {
                var id = $(this)[0].id
                id = id.substring(4, 28)
                var val = $(`#body${id}`)[0].innerText;
                $(`#body${id}`).empty();
                $(`#body${id}`).append(`
                    <input type="text" class="mdl-textfield__input" id="textbox${id}" value="${val}"class="mdl-checkbox__input messageEdit"/></input>
                `)

                $(`#textbox${id}`).keyup(function(e) {
                    if (e.which == 13) {
                        var value = $(this)[0].value
                        var request = {
                            body: value
                        }
                        $.ajax({
                            type: "PATCH",
                            contentType: "application/json",
                            url: `${messageURL}/${id}`,
                            data: JSON.stringify(request),
                            headers: {
                                "Authorization": authkey
                            },
                            success: function(msg) {
                                getMessages()
                            },
                            error: function(XMLHttpRequest, textStatus, errorThrown) {
                                alert(XMLHttpRequest.statusText);
                            }
                        });
                    }
                })
            })



            $(`#del${message.id}`).click(function() {
                var id = $(this)[0].id
                id = id.substring(3, 27)
                $.ajax({
                    type: "DELETE",
                    contentType: "application/json",
                    url: `${messageURL}/${id}`,
                    headers: {
                        "Authorization": authkey
                    },
                    success: function(msg) {
                        getMessages()
                    },
                    error: function(XMLHttpRequest, textStatus, errorThrown) {
                        alert(XMLHttpRequest.statusText);
                    }
                });



            })


            document.getElementById('msglist').scrollTop = document.getElementById('msglist').scrollHeight;
        }


    }


    function renderPeople() {
        var creator = allUsers[currentChan.creatorID];
        currentChan = channelCollection[currentChan.id]
        $('#addPersonList').empty();
        $('#addPersonList').append(`
            <li class="mdl-list__item mdl-menu__item--full-bleed-divider">
                <span class="mdl-list__item-primary-content">
                <img src="${creator.photoURL}" class="material-icons mdl-list__item-avatar">                
                ${creator["firstName"]} ${creator["lastName"]}
                </span>
                <span class="mdl-list__item-secondary-action">
                <i class="material-icons">star</i>
                </span>
            </li>
        `)

        for (var i = 0; i < Object.keys(allUsers).length; i++) {
            var tempUser = allUsers[Object.keys(allUsers)[i]]
            if (creator.id != tempUser.id) {
                var id = `list-${tempUser.id}`
                $('#addPersonList').append(`
                    <li class="mdl-list__item">
                        <span class="mdl-list__item-primary-content">
                        <img src="${tempUser["photoURL"]}" class="material-icons mdl-list__item-avatar">                
                        ${tempUser["firstName"]} ${tempUser["lastName"]}
                        </span>
                        <span class="mdl-list__item-secondary-action">
                        <label class="mdl-checkbox mdl-js-checkbox mdl-js-ripple-effect" for="member-checkbox">
                            <input type="checkbox" id="${id}" class="mdl-checkbox__input" checked disabled/>
                        </label>
                        </span>
                    </li>
                `)
                if (!contains(currentChan, tempUser)) {
                    $(`#${id}`).removeProp("checked")
                }

                if (currentChan.creatorID == currentUser.id) {
                    $(`#${id}`).removeProp("disabled")
                }

                $(`#${id}`).click(function() {
                    var flag = $(this).prop('checked');

                    var linkURL = `${channelURL}/${currentChan.id}`

                    $.ajax({
                        type: flag ? "LINK" : "UNLINK",
                        contentType: "application/json",
                        url: linkURL,
                        headers: {
                            "Authorization": authkey,
                            "Link": (this.id).substring(5, 29)
                        },
                        success: function(msg) {
                            renderPeople();
                        },
                        error: function(XMLHttpRequest, textStatus, errorThrown) {
                            alert(XMLHttpRequest.statusText);
                        }
                    });
                });
            }
        }
    }

    $('#botSend').click((evt)=>{
        evt.preventDefault()
        var text = $('#chatbotText')
        $.ajax({
            type: "POST",
            contentType: "text/plain",
            url: chatbotURL,
            data: text.val(),
            headers: {
                "Authorization": authkey
            },
            success: function(msg) {
                text.val('')
                $('#robomsg').html(msg)
            },
            error: function(XMLHttpRequest, textStatus, errorThrown) {
                alert(XMLHttpRequest.statusText);
            }
        });

    })

    websock.addEventListener("message", function(e) {
        var data = JSON.parse(e.data)
        switch (data.type) {
            case "newChannelcreated":
                renderChan();
                break;
            case "Channelupdated":
                renderChan();
                break;
            case "Channeldeleted":
                renderChan();
                break;
            case "UserJoinedChannel":
                renderChan();
                renderPeople();
            case "UserLeftChannel":
                renderChan();
                renderPeople();
            case "NewMessagePosted":
                getMessages();
            case "MessageDeleted":
                getMevssages();
            case "MessageUpdated":
                getMessages();

        }
    })

    renderChan()

} else {
    window.location.href = "index.html"
}