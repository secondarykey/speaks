var ws = new WebSocket("ws://" + location.host + "/ws/");

ws.onmessage = function(e) {
	var msg = $.parseJSON(e.data);
console.log(e.data)
console.log(msg)
	var linkTag = $('<a/>');
	linkTag.text(msg.Content);
	linkTag.addClass('list-group-item');
	$('#speakTitle').after(linkTag);
}

$(document).ready(function() {
	$('#speakBtn').click(function() {
	    var txt = $('#speakTxt').val()
		var obj = new Object();
	    obj.Content  = txt;
	    obj.Category = "Dashboard";
		var json = JSON.stringify(obj);
		ws.send(json);
	    $('#speakTxt').val('')
	});

});

