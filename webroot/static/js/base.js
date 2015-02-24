var ws = new WebSocket("ws://" + location.host + "/ws/");
var clientId;

ws.onopen = function(e) {
}

ws.onmessage = function(e) {
	var msg = $.parseJSON(e.data);
    var cId = msg.ClientId;

    if ( msg.Type == "Open" ) {
        clientId = cId;
        return;
    }

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

	speakTag.text(msg.Content);
	speakBlockTag.append(speakTag);

	var footerTag = $('<footer/>');
	footerTag.addClass('text-right');
	footerTag.text('xxxx-xx-xx xx:xx:xx');

	itemTag.append(iconBlockTag);
	itemTag.append(speakBlockTag);
	itemTag.append(footerTag);
	$('#speakArea').after(itemTag);
}

function createMessage(msg) {
	var obj = new Object();
    obj.Content  = msg;
    obj.UserId   = Number($("#userId").val());
    obj.Category = $("#category").val();
    obj.ClientId = clientId;
	var json = JSON.stringify(obj);
    return json;
}

$(document).ready(function() {
	$('#speakBtn').click(function() {
	    var txt = $('#speakTxt').val()
		ws.send(createMessage(txt));
	    $('#speakTxt').val('')
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
		   ws.send(createMessage(msg));
        }).error(function() {
            alert("Error!");
        });
        $("#uploadModal").modal("hide");
    });
});

