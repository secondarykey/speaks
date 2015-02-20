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

$(document).ready(function() {
	$('#speakBtn').click(function() {
	    var txt = $('#speakTxt').val()
		var obj = new Object();

	    obj.Content  = txt;
	    obj.UserId   = Number($("#userId").val());
	    obj.Category = $("#category").val();
	    obj.ClientId = clientId;

        console.log(obj);

		var json = JSON.stringify(obj);
        console.log(json);
		ws.send(json);
	    $('#speakTxt').val('')
	});

});

