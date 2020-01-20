
let wActive = true;
let canNotify = notify.isSupported ;
let clientId;
let ws;

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
    let titleClazz = "speakTitleContent" + suffix;
    let who = msg.UserName;
    if ( suffix == "-me" ) {
        who = "You"
    }

    titleBlockTag.setAttribute("class",'mdl-card__title-text ' + titleClazz);
	titleBlockTag.textContent = who + ' Said.';

    var iconTag = document.createElement('img');
    iconTag.setAttribute("class",'speakIcon' + suffix + ' avatar');
	iconTag.setAttribute("src","/images/icon/"+msg.UserId);
    iconTag.onerror = function() {
	  iconTag.setAttribute("src","/images/nobody.png");
    }
    titleTag.append(titleBlockTag);
    titleTag.append(iconTag);

    var contentTag = document.createElement('div');
	contentTag.setAttribute("class",'mdl-card__supporting-text speakContent');

	contentTag.innerHTML = marked(msg.Content);

    var footerTag = document.createElement('div');

    let footerClazz = "rightFooter" + suffix;
    if ( suffix == "-me" ) {
        var delBtn = document.createElement('button');
        delBtn.setAttribute('data-id',msg.Id);
        delBtn.setAttribute('class','mdl-button mdl-js-button red');
        delBtn.addEventListener('click', function(e) {
            if ( msg.Id == null ) {
              alertSpeaks("Please update once and delete.");
            } else {
              confirmSpeaks("Delete?",function(obj) {
                deleteMessage(obj);
              },msg.Id);
            }
        });

        var icon = document.createElement('i');
        icon.setAttribute("class","material-icons");
        icon.textContent = "delete_forever";
        delBtn.appendChild(icon);
	    footerTag.appendChild(delBtn);
    }

	footerTag.setAttribute("class",footerClazz + ' mdl-card__actions mdl-card--border');

    var timeTag = document.createElement('span');
	timeTag.textContent = msg.Created;
    footerTag.appendChild(timeTag);

	itemTag.appendChild(titleTag);
	itemTag.appendChild(contentTag);
	itemTag.appendChild(footerTag);
	itemTag.setAttribute("id","Message-" + msg.Id);

    return itemTag;
}

function createChangeJson() {
	var obj = new Object();
    obj.Type     = "Change";

    var userId = document.querySelector('#userId');
    var category = document.querySelector('#category');
    var project = document.querySelector('#projectKey');

    obj.UserId   = Number(userId.value);
    obj.Category = category.value;
    obj.ClientId = clientId;
    obj.Project  = project.value;
	var json = JSON.stringify(obj);
    return json;
}

function createDeleteJson(msgId) {
	var obj = new Object();
    obj.Type      = "Delete";
    obj.MessageId = Number(msgId);

    var userId = document.querySelector('#userId');
    var category = document.querySelector('#category');
    var project = document.querySelector('#projectKey');

    obj.UserId    = Number(userId.value);
    obj.Category  = category.value;
    obj.ClientId  = clientId;
    obj.Project  = project.value;

	var json = JSON.stringify(obj);
    return json;
}

function createMessageJson(msg) {
    var userId = document.querySelector('#userId');
    var category = document.querySelector('#category');
    var project = document.querySelector('#projectKey');

	var obj = new Object();

    obj.Content  = msg;
    obj.Type     = "Message";
    obj.Project  = project.value;

    obj.UserId   = Number(userId.getAttribute("value"));
    obj.Category = category.getAttribute("value");

    obj.ClientId = clientId;
	var json = JSON.stringify(obj);
    return json;
}

function createCategoryList() {

    var formData = new FormData();
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/category/list");
    xhr.addEventListener('load', function(e) {

        var resp = JSON.parse(e.target.responseText);
        if ( resp.Error != null ) {
            alertSpeaks("Error:" + resp.Error);
            return;
        }

        var ul = document.querySelector('#CategoryUL');
        while (ul.firstChild) ul.removeChild(ul.firstChild);
 
        for ( var idx = 0; idx < resp.CategoryList.length; ++idx ) {

            var category = resp.CategoryList[idx];
 
            var aTag = document.createElement('a');
            var iTag = document.createElement('i');

            aTag.setAttribute('class','mdl-navigation__link childLink');
            aTag.setAttribute('href','#');

            iTag.setAttribute('class','mdl-color-text--blue-grey-400 material-icons mdl-badge mdl-badge--overlap');
            iTag.textContent = "check_box";
            iTag.setAttribute('id','category-' + category.Key + '-icon');

            aTag.appendChild(iTag);

            var spanTag = document.createElement('span');
            spanTag.setAttribute('id',category.Key);
            spanTag.textContent = category.Name;

            //TODO なんとかならんかな。。。
            aTag.categoryKey = category.Key;
            iTag.categoryKey = category.Key;
            spanTag.categoryKey = category.Key;

            aTag.iTag = iTag;
            iTag.iTag = iTag;
            spanTag.iTag = iTag;

	        aTag.addEventListener("click",function(e) {
                e.target.iTag.removeAttribute("data-badge");
                changeCategory(e.target.categoryKey)
            });

            aTag.appendChild(spanTag);

            ul.appendChild(aTag);
       };

    }, false);

    xhr.send(formData);
}

function createMessageList(cat,lastedId) {

    var formData = new FormData();
    formData.append("lastedId",lastedId);

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/message/" + cat);
    xhr.addEventListener('load', function(e) {
        var resp = JSON.parse(e.target.responseText);
        if ( resp.Error != null ) {
            alertSpeaks("Error:" + resp.Error);
            return;
        }

        let msgs = resp.MessageList
        if ( msgs.length > 0 ) {
          var last = document.querySelector('#lastedId');
          last.setAttribute("value",msgs[msgs.length-1].Id);
          for ( var idx = 0; idx < msgs.length; ++idx ) {
            updateMessage(msgs[idx],"");
          }
        } else {
            toast("No longer exist.");
        }
    });
    xhr.send(formData);
}

function changeProject() {
    ws.send(createChangeJson());
}

function changeCategory(catKey) {

    var formData = new FormData();

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/category/view/" + catKey);
    xhr.addEventListener('load', function(e) {

      var resp = JSON.parse(e.target.responseText);
      if ( resp.Error != null ) {
        alertSpeaks("Error:" + resp.Error);
        return;
      }

      let cate = resp.Category;

      var speaks = document.querySelectorAll('.category-speak');
      for (var i=0; i< speaks.length; i++) {
          speaks[i].parentNode.removeChild(speaks[i]);
      }

      var speakTitle = document.querySelector('#speakTitle');
      speakTitle.textContent = cate.Name;

      var description = document.querySelector('#Description');
      description.textContent = cate.Description;

      var category = document.querySelector('#category');
      category.value = cate.Key;

      ws.send(createChangeJson());

      var boxes = document.querySelectorAll('.speakBox');
      for ( let idx = 0;idx < boxes.length ; ++idx )  {
          let box = boxes[idx];
          box.parentNode.removeChild(box);
      }

      createMessageList(cate.Key,"9999999999");
    });

    xhr.send(formData);


    return false;
}

function deleteMessage(msgId) {

    var formData = new FormData();
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/message/delete/" + msgId);
    xhr.addEventListener('load', function(e) {
        var resp = JSON.parse(e.target.responseText);
        if ( resp.Error != null ) {
          alertSpeaks("Error:" + resp.Error);
          return;
        }

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

//検索
function searchMessage(page) {

  var searchWord = document.querySelector('#searchWord');
  var val = searchWord.value;
  var formData = new FormData();

  var category = document.querySelector('#category');
  var project = document.querySelector('#projectKey');

  formData.append("project",project.value);
  formData.append("category",category.value);
  formData.append("search",val);
  formData.append("page",page);

  var xhr = new XMLHttpRequest();
  xhr.open("POST", "/api/search");

  xhr.addEventListener('load', function(e) {
    var resp = JSON.parse(e.target.responseText);
    if ( resp.Error != null ) {
        alertSpeaks("Error:" + resp.Error);
        return;
    }

    //描画
    let msgs = resp.MessageList
    if ( msgs.length > 0 ) {
       var btn = document.querySelector('#searchBtn');
       for ( var idx = 0; idx < msgs.length; ++idx ) {
         var messageTag = createMessageTag(msgs[idx],"");
         btn.parentNode.insertBefore(messageTag,btn);
       }
    } else {
        toast("No longer exist.");
    }

  });

  xhr.send(formData);
}

function switchSearchMode() {

   var update = document.querySelector('#updateBtn');
   var search = document.querySelector('#searchBtn');

   if ( update != null ) {
     update.remove();
   }

   var page = 0;

   var searchBox = document.querySelector('#search-field');
   var searchWord = document.querySelector('#searchWord');
   var searchResult = document.querySelector('#searchResult');

   searchWord.value = searchBox.value;

   search.style.display = "inline";
   searchResult.style.display = "inline";
   searchResult.innerHTML = "検索結果：" + searchWord.value + "<br><a href="">戻る</a>";

   search.addEventListener("click",function(e) {
       page++;
       searchMessage(page);
   });
   search.click();
}

(function() {

  var renderer = new marked.Renderer()
  renderer.code = function(code, language) {
    return '<pre><code class="hljs">' + hljs.highlightAuto(code).value + '</code></pre>';
  };
  
  marked.setOptions({
    renderer: renderer,
  });


    let userId = document.querySelector('#userId');
    ws = new WebSocket("ws://" + location.host + "/ws/" + userId.value);

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
	speakBtn.addEventListener("keydown",function(e) {
        if ( e.keyCode === 13 ) {
            e.preventDefault();
            speakBtn.click();
        }
        return false;
	});

	speakBtn.addEventListener("click",function() {
      var txtTag = document.querySelector('#speakTxt');
	  var txt = txtTag.value;
      if ( txt != "" ) {
	    ws.send(createMessageJson(txt));
        txtTag.value = '';
      }
      txtTag.focus();
      return false;
	});

    var updateBtn = document.querySelector('#updateBtn');
	updateBtn.addEventListener("click",function() {
        var lastedIdTag = document.querySelector('#lastedId');
        var categoryTag = document.querySelector('#category');
	    var lastedId = lastedIdTag.value;
	    var category = categoryTag.value;
        createMessageList(category,lastedId);
	});

    var updateFile = document.querySelector('#uploadFile');
	updateFile.addEventListener("change",function() {
        var fd = new FormData();
        var files = this.files;
        fd.append('uploadFile', files[0]);
        let localName = files[0].name;

        var xhr = new XMLHttpRequest();
        xhr.open("POST", "/api/upload");
        xhr.addEventListener('load', function(e) {

          var resp = JSON.parse(e.target.responseText);
          if ( resp.Error != null ) {
              alertSpeaks("Error:" + resp.Error);
              return;
          }

          var msg = "[" + localName +  "]\n" +
                    "http://" + location.host + "/" + resp.Result.FileName;
          insertText("#speakTxt",msg);
        });
        xhr.send(fd);
    });

    var speakFile = document.querySelector('#speakFile');
	speakFile.addEventListener("keydown",function(e) {
        if ( e.keyCode === 13 ) {
            e.preventDefault();
            speakFile.click();
        }
        return false;
	});

	speakFile.addEventListener("click",function() {
        uploadFile.click();
    });

    
    var cat = document.querySelector('#category');
    createCategoryList()

    changeCategory(cat.value);

    var speakTxt = document.querySelector('#speakTxt');
    speakTxt.focus();

    ws.onopen = function(e) {
      //初回時に必ず呼び出す
      changeProject();
    }

    ws.onmessage = function(e) {

      var msg = JSON.parse(e.data);

      var cId = msg.ClientId;
      if ( msg.Type == "Open" ) {
          clientId = cId;
          return;
      } else if ( msg.Type == "AddUser" ) {
          return;
      } else if ( msg.Type == "DeleteUser" ) {
          return;
      } else if ( msg.Type == "Delete" ) {
          var message = document.querySelector('#Message-' + msg.MessageId);
          message.parentNode.removeChild(message);
          return;
      } else if ( msg.Type == "Notify" ) {
          var projectKey =  msg.Project;
          var catKey =  msg.Category;
  
          let tag = null;
          var category = document.querySelector('#category');
          var project = document.querySelector('#projectKey');
  
          if ( projectKey == project.value ) {
            if ( catKey != category.value ) {
                if ( catKey == "Dashboard" ) {
                   //プロジェクトのところ
                   tag = document.querySelector('#project-' + projectKey + "-icon");
                } else {
                   //カテゴリのところ
                   tag = document.querySelector('#category-' + catKey + "-icon");
                }
            }
          } else {
             tag = document.querySelector('#project-' + projectKey + "-icon");
          }
  
          if ( tag != null ) {
            tag.setAttribute("data-badge","+");
          }
          return;
      }
      addMessage(msg,cId);
  }

  var speakTxt = document.querySelector('#speakTxt');
  var preview = document.querySelector('#previewSpeak');
  var previewBtn = document.querySelector('#switchPreview');

  //プレビューをOnにする
  previewBtn.addEventListener("click",function() {
      var disp = preview.style.display;
      var val = "none";
      if ( disp == "none" ) {
          val = "inline";
      }
      preview.style.display = val;
  });

  var searchBox = document.querySelector('#search-field');
  searchBox.addEventListener("keydown",function(e) {

      if ( e.keyCode != 13 ) {
          return;
      }

      var speakCard = document.querySelector('#speakCard');
      if ( speakCard != null ) {
          speakCard.remove();
      }

      var boxes = document.querySelectorAll('.speakBox');
      for ( let idx = 0;idx < boxes.length ; ++idx )  {
          let box = boxes[idx];
          box.parentNode.removeChild(box);
      }

      //searchモードへ変更
      switchSearchMode();
  });


  //テキストエリアの更新を設定
  speakTxt.addEventListener("keydown",function() {
      /*
      var disp = preview.style.display;
      if ( disp != "none" ) {
      }
      */
      var val = speakTxt.value;
      var html = marked(val);
      preview.innerHTML = html;
  });


}).call(this);


