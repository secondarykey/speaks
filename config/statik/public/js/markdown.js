function change(result) {

  if ( result !== null ) {
    let editor = document.querySelector('#editor');
    let src = editor.value;
    let html = marked(src);
    result.innerHTML = html;
  }

  let codes = document.querySelectorAll('pre code');
  for ( var idx = 0; idx < codes.length; ++idx ) {
    let code = codes[idx];
    hljs.highlightBlock(code);
  }
}

let viewer = null;
(function() {

  let viewerBtn = document.querySelector('#viewer');
  if ( viewerBtn !== null ) {
    viewerBtn.addEventListener("click",function(e) {
      let id = viewerBtn.getAttribute("data-id");
      viewer = window.open("/memo/" + id,"viewer");
    });
  }

  marked.setOptions({
    langPrefix: ''
  });

  let editor = document.querySelector('#editor');
  editor.addEventListener("keyup",function() {
    if ( viewer != null ) {
      let result = viewer.document.querySelector("#result");
      change(result);
    }
  });

  change(document.querySelector("#result"));

}).call(this);
