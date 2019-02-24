var wActive = true;
var canNotify = notify.isSupported ;
var ws = new WebSocket("ws://" + location.host + "/ws/");
var clientId;

ws.onopen = function(e) {
}

function addMessage(msg,cId) {
    var messageTag = createMessageTag(msg,cId);
    $('#speakArea').after(messageTag);
}

function updateMessage(msg,cId) {
    var messageTag = createMessageTag(msg,cId);
    $('#updateBtn').before(messageTag);
}

function createMessageTag(msg,cId) {

	var suffix = "";
    if ( clientId == cId ) {
	    suffix = "-me";
    }

    var userId = Number($('#userId').val());
    if ( userId == msg.UserId ) {
	    suffix = "-me";
    }

	var itemTag = $('<div/>');
	itemTag.addClass('list-group-item');
	itemTag.addClass('category-speak');

	var iconBlockTag = $('<div/>');
	iconBlockTag.addClass('icon-block' + suffix);
	var iconTag = $('<img/>');
	iconTag.addClass('speak-icon' + suffix);
	iconTag.addClass('userIcon');
	iconTag.attr("src","/static/images/icon/"+msg.UserId);
	iconBlockTag.append(iconTag);
    iconTag.error(function() {
        $(this).attr({
            src: '/static/images/icon/nobody.png',
            alt: 'no image'
        });
    });

	var speakBlockTag = $('<div/>');
	speakBlockTag.addClass('speak-block' + suffix);
	speakBlockTag.text(msg.UserName + ' Says.');

	var speakTag = $('<pre/>');
	speakTag.addClass('speak' + suffix);

    var linkTxt = msg.Content.replace(/(http:\/\/[\x21-\x7e]+)/gi, "<a href='$1' target='_blank'>$1</a>"); 

	speakTag.html(linkTxt);
	speakBlockTag.append(speakTag);

	var footerTag = $('<footer/>');
	footerTag.addClass('text-right');
	footerTag.text(msg.Created + " ");
    if ( suffix == "-me" ) {
        var delBtn = $('<button/>');
        delBtn.addClass('btn');
        delBtn.addClass('btn-danger');
        delBtn.addClass('btn-xs');
        delBtn.css('height','26px');
        delBtn.attr('data-id',msg.Id);
        delBtn.on("click",null,deleteMessage);
        delBtn.popConfirm({
                title:"Delete Message",
                content:"Delete?",
                placement:"left"
        });
        var delSpn = $('<span/>');
        delSpn.addClass('glyphicon');
        delSpn.addClass('glyphicon-remove-sign');
        delBtn.append(delSpn);
	    footerTag.append(delBtn);
    }

	itemTag.append(iconBlockTag);
	itemTag.append(speakBlockTag);
	itemTag.append(footerTag);
	itemTag.attr("id","Message-" + msg.Id);

    return itemTag;
}

ws.onmessage = function(e) {
	var msg = $.parseJSON(e.data);
    var cId = msg.ClientId;
    if ( msg.Type == "Open" ) {
        clientId = cId;
        return;
    } else if ( msg.Type == "Delete" ) {
        $("#Message-" + msg.MessageId).remove();
        return;
    } else if ( msg.Type == "Notify" ) {
        var catKey =  msg.Category;
        var num = $("#" + catKey).text();
        if (num == "") {
            num = "0";
        }
        $("#" + catKey).text(Number(num)+1);
        return;
    }

    addMessage(msg,cId);
    createNotify("Notify","You got speak.","/static/images/notify.png");
}

function createChangeJson() {
	var obj = new Object();
    obj.Type     = "Change";
    obj.UserId   = Number($("#userId").val());
    obj.Category = $("#category").val();
    obj.ClientId = clientId;
	var json = JSON.stringify(obj);
    return json;
}

function createDeleteJson(msgId) {
	var obj = new Object();
    obj.Type      = "Delete";
    obj.MessageId = Number(msgId);
    obj.UserId    = Number($("#userId").val());
    obj.Category  = $("#category").val();
    obj.ClientId  = clientId;
	var json = JSON.stringify(obj);
    return json;
}

function createMessageJson(msg) {
	var obj = new Object();
    obj.Content  = msg;
    obj.Type     = "Message";
    obj.UserId   = Number($("#userId").val());
    obj.Category = $("#category").val();
    obj.ClientId = clientId;
	var json = JSON.stringify(obj);
    return json;
}

function getCategoryList() {
    $.ajax({
       url: "category/list",
       type: 'POST',
       data: { },
       dataType: 'json'
    }).success(function( data ) {
       var ul = $('#CategoryUL');
       if (data.length > 0 ) {
           ul.empty();
       }

       $.each(data, function(i, category){
           var li = $('<li/>');
           var aTag = $('<a/>');
           aTag.attr('href','#');
           aTag.text(category.Name);

           //<span class="badge">4</span>
           var spanTag = $('<span/>');
           spanTag.attr('id',category.Key);
           spanTag.addClass('badge');

           aTag.on("click",{key:category.Key},changeCategory);

           li.append(aTag);
           li.append(spanTag);
           ul.append(li);
       });
    }).error(function() {
        alert("Error!");
    });
}

function getMessageList(cat,lastedId) {
    $.ajax({
       url: "message/" + cat,
       type: 'POST',
       data: {
           "lastedId" : lastedId
       },
       dataType: 'json'
    }).success(function( data ) {
       if (data.length > 0 ) {
           $("#lastedId").val(data[data.length-1].Id);
       }
       $.each(data, function(i, msg){
           updateMessage(msg,"");
       });
    }).error(function() {
        alert("Error!");
    });
}

function changeCategory(evt) {

    var catKey = evt.data.key
    if ( catKey == "Dashboard") {
        $("#deleteBtn").hide();
        $("#memoBtn").hide();
    } else {
        $("#deleteBtn").show();
        $("#memoBtn").show();
    }

    $.ajax({
       url: "category/view/" + catKey,
       type: 'POST',
       data: { },
       dataType: 'json'
    }).success(function( data ) {
        $("#" + catKey).text("");
        // tag empty
        $(".category-speak").each(function() {
            $(this).remove();
        });
        // change title 
        $("#speakTitle").text(data.Name);
        $("#Description").text(data.Description);
        // change hide value 
        $("#category").val(data.Key);
        ws.send(createChangeJson());

        getMessageList(data.Key,"9999999999");
    }).error(function() {
        alert("Error!");
    });

    return false;
}

function deleteMessage(evt) {
    var msgId = $(this).attr('data-id');
    $.ajax({
       url: "message/delete/" + msgId,
       type: 'POST',
       data: { },
       dataType: 'json'
    }).success(function( data ) {
	   ws.send(createDeleteJson(msgId));
    }).error(function() {
        alert("Error!");
    });

    return false;
}

function createNotify(title,body,icon) {
    if ( !wActive && canNotify && 
            (notify.permissionLevel() == notify.PERMISSION_GRANTED) ) {
        notify.createNotification(title,{body:body,icon:icon});
    }
}

$(document).ready(function() {
    if ( canNotify ) {
        var permission = notify.permissionLevel();
        if ( permission == notify.PERMISSION_DEFAULT ) {
            notify.requestPermission();
        //} else if ( permission == notify.PERMISSION_GRANTED ) {
            //createNotify("Notify","Message!","alert.ico");
        }
        //notify.PERMISSION_DENIED
    }

    $(window).bind("focus",function(){  //フォーカスした
        wActive = true;
    }).bind("blur",function(){  //フォーカスが外れた
        wActive = false;
    }); 

	$('#updateBtn').click(function() {
	    var lastedId = $('#lastedId').val();
	    var category = $('#category').val();
        getMessageList(category,lastedId);
    });

	$('#speakBtn').click(function() {
	    var txt = $('#speakTxt').val()
        if ( txt != "" ) {
		    ws.send(createMessageJson(txt));
	        $('#speakTxt').val('')
        }
        $("#speakTxt").focus();
	});

    $('#uploadFile').change(function() {
        var fd = new FormData();
        var files = this.files;
        $.each(files, function(i, file){
            fd.append('uploadFile', file);
        });
        $.ajax({
           url: "upload",
           type: 'POST',
           data: fd,
           processData:false,
           contentType:false,
           dataType: 'json'
        }).success(function( data ) {
           var msg = "http://" + location.host + "/" + data.FileName;
		   ws.send(createMessageJson(msg));
        }).error(function() {
            alert("Error!");
        });
        $("#uploadModal").modal("hide");
        $("#speakTxt").focus();
    });

	$('#memoBtn').click(function() {
       var key = $("#category").val();
       var url = "http://" + location.host + "/memo/view/" + key;
       window.open(url, '_blank');
    });

	$('#memoBtn').hide();
	$('#deleteBtn').hide();

	$('#deleteBtn').click(function() {
        var cat = $("#category").val();
        var url = "/category/delete/" + cat;
        $.ajax({
           url: url,
           type: 'POST',
        }).success(function( data ) {
            location.href="/";
        }).error(function() {
            alert("Error!");
        });
    });
    $('#deleteBtn').popConfirm({
        title:"Delete Category",
        content:"Delete Category and All Message!It is recommended that you create a memo before you turn off.",
        placement:"bottom"
    });

    getMessageList("Dashboard","9999999999");
    getCategoryList()
    $("#speakTxt").focus();
});

