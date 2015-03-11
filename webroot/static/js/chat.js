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

	var itemTag = $('<div/>');
	itemTag.addClass('list-group-item');

	var iconBlockTag = $('<div/>');
	iconBlockTag.addClass('icon-block' + suffix);
	var iconTag = $('<img/>');
	iconTag.addClass('speak-icon' + suffix);
	iconTag.attr("src","/static/images/nobody.png");
	iconBlockTag.append(iconTag);

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

$(document).ready(function() {

	$('#updateBtn').click(function() {
	    var lastedId = $('#lastedId').val()
        getMessage("Public",lastedId);
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

    function getMessage(cat,lastedId) {
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
    function getCategory() {
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
	          aTag.attr('href','Test');
	          aTag.text(category.Name);

	          li.append(aTag);
	          ul.append(li);
           });
        }).error(function() {
            alert("Error!");
        });
    }

    getMessage("Public","9999999999");
    getCategory()
    $("#speakTxt").focus();
});

