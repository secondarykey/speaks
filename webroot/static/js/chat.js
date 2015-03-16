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
	speakBlockTag.text('say');

	var speakTag = $('<pre/>');
	speakTag.addClass('speak' + suffix);

    var linkTxt = msg.Content.replace(/(http:\/\/[\x21-\x7e]+)/gi, "<a href='$1' target='_blank'>$1</a>"); 

	speakTag.html(linkTxt);
	speakBlockTag.append(speakTag);

	var footerTag = $('<footer/>');
	footerTag.addClass('text-right');
	footerTag.text(msg.Created);

	itemTag.append(iconBlockTag);
	itemTag.append(speakBlockTag);
	itemTag.append(footerTag);

    return itemTag;
}

ws.onmessage = function(e) {
	var msg = $.parseJSON(e.data);
    var cId = msg.ClientId;
    if ( msg.Type == "Open" ) {
        clientId = cId;
        return;
    }
    addMessage(msg,cId);
}

function createMessageJson(msg) {
	var obj = new Object();
    obj.Content  = msg;
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

           aTag.on("click",{key:category.Key},changeCategory);

           li.append(aTag);
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
    $.ajax({
       url: "category/view/" + evt.data.key,
       type: 'POST',
       data: { },
       dataType: 'json'
    }).success(function( data ) {
        // tag empty
        $(".category-speak").each(function() {
            $(this).remove();
        });
        // change title 
        $("#speakTitle").text(data.Name);
        // change hide value 
        $("#category").val(data.Key);
        getMessageList(data.Key,"9999999999");
    }).error(function() {
        alert("Error!");
    });
    return false;
}

$(document).ready(function() {


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

    getMessageList("Dashboard","9999999999");
    getCategoryList()
    $("#speakTxt").focus();
});

