let wActive = true;
let canNotify = notify.isSupported ;
let ws = new WebSocket("ws://" + location.host + "/ws/");
let clientId;

ws.onopen = function(e) {
}

function addMessage(msg,cId) {
    var messageTag = createMessageTag(msg,cId);
    var area = document.querySelector('#speakArea');
    area.parentNode.insertBefore(messageTag,area.nextSibling);
}

function updateMessage(msg,cId) {
    var messageTag = createMessageTag(msg,cId);

    var btn = document.querySelector('#updateBtn');
    btn.parentNode.insertBefore(messageTag,btn);
}

function createMessageTag(msg,cId) {

	var suffix = "";
    if ( clientId == cId ) {
	    suffix = "-me";
    }

    var userId = document.querySelector('#userId');
    var userId = Number(userId.getAttribute("value"));
    if ( userId == msg.UserId ) {
	    suffix = "-me";
    }

    var itemTag = document.createElement('div');
    itemTag.setAttribute("id",'Message-' + msg.Id);
    itemTag.setAttribute("class",'mdl-card mdl-shadow--2dp speakBox');

    var titleTag = document.createElement('div');
    titleTag.setAttribute("class",'mdl-card__title mdl-card--expand speakTitle');

    var titleBlockTag = document.createElement('h6');
	titleBlockTag.setAttribute("class",'mdl-card__title-text speakTitleContent');
	titleBlockTag.textContent = msg.UserName + ' Says.';

    var iconTag = document.createElement('img');
    iconTag.setAttribute("class",'speakIcon' + suffix + ' avatar');
	iconTag.setAttribute("src","/static/images/icon/"+msg.UserId);
    iconTag.onerror = function() {
	  iconTag.setAttribute("src","/static/images/icon/nobody.png");
    }
    titleTag.append(titleBlockTag);
    titleTag.append(iconTag);

    var contentTag = document.createElement('div');
	contentTag.setAttribute("class",'mdl-card__supporting-text speakContent');

    var linkTxt = msg.Content.replace(/(http:\/\/[\x21-\x7e]+)/gi, "<a href='$1' target='_blank'>$1</a>"); 
	contentTag.textContent = linkTxt;

    var footerTag = document.createElement('div');

	footerTag.setAttribute("class",'text-right mdl-card__actions mdl-card--border');
	footerTag.textContent = msg.Created + " ";
    if ( suffix == "-me" ) {
        var delBtn = document.createElement('button');
        delBtn.setAttribute('data-id',msg.Id);
        delBtn.addEventListener('click', function(e) {
            deleteMessage(msg.Id);
        });
	    footerTag.appendChild(delBtn);
    }

	itemTag.appendChild(titleTag);
	itemTag.appendChild(contentTag);
	itemTag.appendChild(footerTag);
	itemTag.setAttribute("id","Message-" + msg.Id);

    return itemTag;
}

ws.onmessage = function(e) {

    var msg = JSON.parse(e.data);

    var cId = msg.ClientId;
    if ( msg.Type == "Open" ) {
        clientId = cId;
        return;
    } else if ( msg.Type == "Delete" ) {
        var message = document.querySelector('#Message-' + msg.MessageId);
        message.parentNode.removeChild(message);
        return;
    } else if ( msg.Type == "Notify" ) {
        var catKey =  msg.Category;

        var numTag = document.querySelector('#' + catKey);
        var num = numTag.textContent;
        if (num == "") {
            num = "0";
        }
        numTag.textContent = Number(num) + 1;
        return;
    }

    addMessage(msg,cId);
    createNotify("Notify","You got speak.","/static/images/notify.png");
}

function createChangeJson() {
	var obj = new Object();
    obj.Type     = "Change";

    var userId = document.querySelector('#userId');
    var category = document.querySelector('#category');

    obj.UserId   = Number(userId.value);
    obj.Category = category.value;
    obj.ClientId = clientId;
	var json = JSON.stringify(obj);
    return json;
}

function createDeleteJson(msgId) {
	var obj = new Object();
    obj.Type      = "Delete";
    obj.MessageId = Number(msgId);

    var userId = document.querySelector('#userId');
    var category = document.querySelector('#category');

    obj.UserId    = Number(userId.value);
    obj.Category  = category.value;
    obj.ClientId  = clientId;
	var json = JSON.stringify(obj);
    return json;
}

function createMessageJson(msg) {
    var userId = document.querySelector('#userId');
    var category = document.querySelector('#category');

	var obj = new Object();
    obj.Content  = msg;
    obj.Type     = "Message";

    obj.UserId   = Number(userId.getAttribute("value"));
    obj.Category = category.getAttribute("value");

    obj.ClientId = clientId;
	var json = JSON.stringify(obj);
    return json;
}

function getCategoryList() {

    var formData = new FormData();
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/category/list");
    xhr.addEventListener('load', function(e) {

       var resp = JSON.parse(e.target.responseText);
       var ul = document.querySelector('#CategoryUL');
       while (ul.firstChild) ul.removeChild(ul.firstChild);

       for ( var idx = 0; idx < resp.length; ++idx ) {

           var category = resp[idx];
           var li = document.createElement('li');
           var aTag = document.createElement('a');

           aTag.setAttribute('href','#');
           aTag.textContent = category.Name;

           var spanTag = document.createElement('span');
           spanTag.setAttribute('id',category.Key);
           spanTag.setAttribute('class','badge');

           aTag.data = category.Key;
	       aTag.addEventListener("click",function(e) {
               changeCategory(e.target.data)
           });

           li.append(aTag);
           li.append(spanTag);

           ul.appendChild(li);
       };

    }, false);
    xhr.send(formData);
}

function getMessageList(cat,lastedId) {

    var formData = new FormData();
    formData.append("lastedId",lastedId);

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/message/" + cat);
    xhr.addEventListener('load', function(e) {
       var resp = JSON.parse(e.target.responseText);
       if ( resp.length > 0 ) {
         var last = document.querySelector('#lastedId');
         last.setAttribute("value",resp[resp.length-1].Id);
         for ( var idx = 0; idx < resp.length; ++idx ) {
           updateMessage(resp[idx],"");
         }
       }
    });
    xhr.send(formData);
}

function changeCategory(catKey) {

    var formData = new FormData();

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/category/view/" + catKey);
    xhr.addEventListener('load', function(e) {
      var resp = JSON.parse(e.target.responseText);
      var catTag = document.querySelector('#' + catKey);
      catTag.textContent = "";
      
      var speaks = document.querySelectorAll('.category-speak');
      for (var i=0; i< speaks.length; i++) {
          speaks[i].parentNode.removeChild(speaks[i]);
      }

      var speakTitle = document.querySelector('#speakTitle');
      var description = document.querySelector('#Description');
      speakTitle.textContent = resp.Name;
      description.textContent = resp.Description;

      var category = document.querySelector('#category');
      category.value = resp.Key;

      ws.send(createChangeJson());

      var boxes = document.querySelectorAll('.speakBox');
      for ( let idx = 0;idx < boxes.length ; ++idx )  {
          let box = boxes[idx];
          box.parentNode.removeChild(box);
      }

      getMessageList(resp.Key,"9999999999");
    });
    xhr.send(formData);

    return false;
}

function deleteMessage(msgId) {

    var formData = new FormData();
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/message/delete/" + msgId);
    xhr.addEventListener('load', function(e) {
	   ws.send(createDeleteJson(msgId));
    });
    xhr.send(formData);

    return false;
}

function createNotify(title,body,icon) {
    if ( !wActive && canNotify && 
            (notify.permissionLevel() == notify.PERMISSION_GRANTED) ) {
        notify.createNotification(title,{body:body,icon:icon});
    }
}

function insertText(name,text) {
    let areaTag = document.querySelector(name);
    let v = areaTag.value;
    let cur = areaTag.selectionStart;
    let v1 = v.substr(0,cur);
    let v2 = v.substr(cur);
    areaTag.value = v1 + "\n" + text + "\n" + v2;
    return;
}

(function() {

    if ( canNotify ) {
        var permission = notify.permissionLevel();
        if ( permission == notify.PERMISSION_DEFAULT ) {
            notify.requestPermission();
        } else if ( permission == notify.PERMISSION_GRANTED ) {
            createNotify("Notify","Message!","alert.ico");
        }
        notify.PERMISSION_DENIED
    }

    window.addEventListener("focus",function(){
        wActive = true;
    });

    window.addEventListener("blur",function(){
        wActive = false;
    }); 

    var speakBtn = document.querySelector('#speakBtn');
	speakBtn.addEventListener("click",function() {
      var txtTag = document.querySelector('#speakTxt');
	  var txt = txtTag.value;
      if ( txt != "" ) {
	    ws.send(createMessageJson(txt));
        txtTag.value = '';
      }
      txtTag.focus();
	});

    var updateBtn = document.querySelector('#updateBtn');
	updateBtn.addEventListener("click",function() {
        var lastedIdTag = document.querySelector('#lastedId');
        var categoryTag = document.querySelector('#category');
	    var lastedId = lastedIdTag.value;
	    var category = categoryTag.value;
        getMessageList(category,lastedId);
	});

    var updateFile = document.querySelector('#uploadFile');
	updateFile.addEventListener("change",function() {
        var fd = new FormData();
        var files = this.files;
        fd.append('uploadFile', files[0]);

        var xhr = new XMLHttpRequest();
        xhr.open("POST", "/upload");
        xhr.addEventListener('load', function(e) {
          var resp = JSON.parse(e.target.responseText);
          var msg = "http://" + location.host + "/" + resp.FileName;

          insertText("#speakTxt",msg);
        });
        xhr.send(fd);
    });


/*
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
*/

    getMessageList("Dashboard","9999999999");
    getCategoryList()

    var speakTxt = document.querySelector('#speakTxt');
    speakTxt.focus();

}).call(this);

