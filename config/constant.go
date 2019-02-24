package config
var ConstantCategory = `
{{define "JavaScript"}}
{{end}}
{{define "content"}}

<form action="#" method="post" class="form-horizontal">

  <input type="hidden" name="key" value="{{.CategoryKey}}">

  <div class="form-group">
    <label for="CategoryName" class="col-sm-2 control-label">Name</label>
    <div class="col-sm-7">
       <input type="text" name="name" class="form-control" id="CategoryName" placeholder="Category Name">
    </div>
  </div>

  <div class="form-group">
    <label for="Description" class="col-sm-2 control-label">Description</label>
    <div class="col-sm-7">
    <textarea id="Description" name="description" class="form-control" rows="3" placeholder="Description"></textarea>
    </div>
  </div>

  <div class="form-group">
    <div class="col-sm-offset-2 col-sm-7">
      <button type="submit" class="btn btn-primary">Create</button>
    </div>
  </div>
</form>
{{end}}

`
var ConstantChat = `
{{define "JavaScript"}}
    <script src="/static/js/desktop-notify-min.js"></script>
    <script src="/static/js/chat.js"></script>
{{end}}

{{define "content"}}
<div class="list-group" id="speakList">
  <div id="speakTitle" class="list-group-item active">Dashboard</div>

  <div id="Description" class="list-group-item small">
  </div>

  <div class="form-group list-group-item" id="speakArea">
    <input type="hidden" id="userId" name="userId" value="{{.User.Id}}"/>
    <input type="hidden" id="category" name="category" value="Dashboard"/>
    <input type="hidden" id="lastedId" name="lastedId" value=""/>

    <textarea id="speakTxt" class="form-control" rows="3" placeholder="現在の課題をしゃべりましょう"></textarea>
    <button type="button" class="btn btn-info form-control" id="speakBtn"><span class="glyphicon glyphicon-edit" aria-hidden="true">Speak</span></button>

    <div class="btn-group">
    <button type="button" class="btn btn-lg btn-default" data-toggle="modal" data-target="#uploadModal"><span class="glyphicon glyphicon-picture" aria-hidden="true"></span></button>
    <button type="button" class="btn btn-lg btn-success" id="memoBtn"><span class="glyphicon glyphicon-list-alt"></span></button>
    <button type="button" class="btn btn-lg btn-danger" id="deleteBtn"><span class="glyphicon glyphicon-remove-sign"></span></button>
    </div>
  </div>

  <div class="list-group-item text-center btn btn-success" id="updateBtn">
      <span class="glyphicon glyphicon-refresh" aria-hidden="true"></span>Update 
  </div>

  <div class="modal fade" id="uploadModal" tabindex="-1" role="dialog" aria-labelledby="uploadModalLabel" aria-hidden="true">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
          <h4 class="modal-title" id="uploadModalLabel">選択するとそのまま発言します</h4>
        </div>

        <div class="modal-body">
          <form action="upload" method="post" enctype="multipart/form-data">
            <div class="input-group">
              <input type="file" id="uploadFile" name="uploadFile" class="form-control">
            </div>
          </form>
        </div>

        <div class="modal-footer">
          <button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
        </div>

      </div>
    </div>
  </div>
{{end}}

`
var ConstantDatabase = `
{{define "JavaScript"}}
{{end}}

{{define "content"}}
<div class="panel panel-default">
    <div class="panel-heading">Database Console</div>

    <div class="panel-body">
      <form method="post" >
      <textarea name="SQL" class="form-control" rows="3" placeholder="SQL"></textarea>
      <button type="submit" class="btn btn-success form-control"><span class="glyphicon glyphicon-search" aria-hidden="true">Search</span></button>
      </form>
    </div>

    <table class="table table-striped">
       <thead>
         <tr id="tableHeader">
{{range .Columns}}
<th>{{.}}</th>
{{end}}
         </tr>
       </thead>
       <tbody id="tableData">
{{range .Records}}
<tr>
    {{range .}}
    <td>{{.}}</td>
    {{end}}
</tr>
{{end}}
       </tbody>
    </table>
</div>
{{end}}

`
var ConstantLayout = `
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="description" content="">
    <meta name="author" content="">
    <link rel="icon" href="/static/images/favicon.ico">

    <title>SpeakAll</title>

    <link href="/static/css/bootstrap.min.css" rel="stylesheet">
    <link href="/static/css/base.css" rel="stylesheet">

    <script src="/static/js/jquery-1.11.2.min.js"></script>
    <script src="/static/js/bootstrap.min.js"></script>
    <script src="/static/js/jquery.popconfirm.js"></script>

{{template "JavaScript" .}}

   <script>
      $(document).ready(function() {
            var curentFile = window.location.pathname.split("/").pop();
            $('ul.nav > li > a[href="/' + curentFile + '"]').parent().addClass('active');
      });
   </script>

  </head>
  <body>
    <div class="container">
    {{template "menu" .}}
    </div>
  </body>
</html>

`
var ConstantLogin = `
{{define "JavaScript"}}
{{end}}
{{define "menu"}}
	  <form class="form-signin" method="post" action="login">
        <h2 class="form-signin-heading">Please sign in</h2>
        <label for="inputEmail" class="sr-only">Email address</label>
        <input type="email" id="inputEmail" class="form-control" placeholder="Email address" required="" autofocus="" name="email">
        <label for="inputPassword" class="sr-only">Password</label>
        <input type="password" id="inputPassword" name="password" class="form-control" placeholder="Password" required="">
<!--
        <div class="checkbox">
          <label>
            <input type="checkbox" value="remember-me"> Remember me
          </label>
        </div>
-->
        <button class="btn btn-lg btn-primary btn-block" type="submit">Sign in</button>
      </form>
{{end}}

`
var ConstantMenu = `
{{define "menu"}}
    <nav class="navbar navbar-inverse navbar-fixed-top">
      <div class="container-fluid">
        <div class="navbar-header">

          <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
            <span class="sr-only">Toggle navigation</span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          </button>

          <a class="navbar-brand" href="/">SpeakAll</a>

        </div>

        <div id="navbar" class="navbar-collapse collapse">

          <ul class="nav navbar-nav navbar-left">
            <li><a href="/memo">Archive</a></li>
          </ul>

          <ul class="nav navbar-nav navbar-right">
            <li><a href="/me">{{.User.Name}}</a></li>
            <li><a href="/logout">Logout</a></li>
          </ul>
        </div>
      </div>
    </nav>

    <div class="container-fluid">
      <div class="row">
        <div class="col-sm-3 col-md-3 sidebar">
          <ul class="nav nav-sidebar">
{{if .User.IsSpeaker}}
            <li><a href="/">Speak<span id="Dashboard" class="badge"></span></a>
            <ul id="CategoryUL">
            </ul>
            </li>
{{end}}

{{if .User.IsChairman}}
            <li><a href="/category">Catefgory<span class="sr-only">(current)</span></a></li>
{{end}}

<!--
            <li><a href="/memo">Memo<span class="sr-only">(current)</span></a></li>
-->

{{if .User.IsAdmin}}
            <li><a>Management</a>
              <ul id="ManagementUL" class="nav nav-sidebar">
                 <li><a href="/user">User</a></li>
                 <li><a href="/database">Database</a></li>
              </ul>
            </li>
{{end}}

          </ul>
        </div>

        <div class="col-sm-9 col-sm-offset-3 col-md-9 col-md-offset-3 main">
{{template "content" .}}
        </div>

      </div>
    </div>
  </div>

{{end}}

`
var ConstantUser = `
{{define "JavaScript"}}
{{end}}
{{define "content"}}

<form action="{{.URL}}" method="post" class="form-horizontal">

  <div class="form-group">
    <label for="name" class="col-sm-2 control-label">Name</label>
    <div class="col-sm-10">
       <input name="dispName" type="text" class="form-control" id="name" placeholder="Name" value="{{.EditUser.Name}}">
    </div>
  </div>

  <div class="form-group">
    <label for="email" class="col-sm-2 control-label">Email</label>
    <div class="col-sm-10">
       <input name="email" type="email" class="form-control" id="email" placeholder="Email" value="{{.EditUser.Email}}">
    </div>
  </div>

{{if ne .URL "/user" }}
  <div class="form-group">
    <label for="inputPassword" class="col-sm-2 control-label">Password</label>
    <div class="col-sm-10">
      <input type="password" name="password" class="form-control" id="inputPassword" placeholder="Password">
    </div>
  </div>

  <div class="form-group">
    <label for="verifiedPassword3" class="col-sm-2 control-label">Password(verified)</label>
    <div class="col-sm-10">
      <input type="password" name="verifiedPassword" class="form-control" id="verifiedPassword" placeholder="Password">
    </div>
  </div>
{{end}}

  <div class="form-group">
    <div class="col-sm-offset-2 col-sm-10">
      <button type="submit" class="btn btn-primary">Save</button>
    </div>
  </div>

{{if eq .URL "/user" }}
<div class="panel panel-default">
  <!-- Default panel contents -->
  <div class="panel-heading">User List</div>
  <table class="table">
  <tr> 
    <th>Name</th>
    <th>EMail</th>
    <th>URL</th>
  </tr>
  {{range .UserList}}
  <tr> 
    <td>{{.Name}}</td>
    <td>{{.Email}}</td>
    <td>{{.Password}}</td>
  </tr>
  {{end}}
  </table>
</div>
{{end}}
</form>

{{if eq .URL "/me" }}
<form action="/me/upload" method="post" enctype="multipart/form-data">
    <div class="input-group">
      <input type="file" id="uploadFile" name="uploadFile" class="form-control">
      <span class="input-group-btn">
        <button type="submit" class="btn btn-primary">Upload</button>
      </span>
    </div>
</form>
{{end}}

{{end}}

`
var ConstantEdit = `
{{define "JavaScript"}}
<link href="/static/css/github.css" rel="stylesheet">
<script src="/static/js/marked.min.js"></script>
<script src="/static/js/highlight.pack.js"></script>
<script src="/static/js/markdown.js"></script>
{{end}}

{{define "menu"}}
  <div class="row" style="height:100%;">
    <div class="col-xs-6">
    <form action="#" method="post" id="memoForm">
        <input type="text" class="form-control" name="Name" value="{{.Memo.Name}}" />
        <textarea name="Content" id="editor" class="form-control">{{.Memo.Content}}</textarea>
        <button type="submit" class="btn btn-success form-control">Save</button>
    </form>
    </div>

    <div class="col-xs-6">
        <button type="button" id="deleteBtn" class="btn btn-danger form-control">Delete</button>
        <div id="result"></div>
    </div>
  </div>

{{end}}

`
var ConstantList = `
{{define "JavaScript"}}
<script>
$(document).ready(function() {
    $(".memoBtn").click(function() {
        var key = $(this).attr('data-id');
        var url = "http://" + location.host + "/memo/view/" + key;
        window.open(url, '_blank');
    });
});
</script>
{{end}}
{{define "content"}}

<div class="panel panel-default">
  <!-- Default panel contents -->
  <div class="panel-heading">ArhciveMemo List</div>
  <table class="table">
  {{range .MemoList}}
  <tr> 
    <td><button style="width:100%" type="button" class="btn btn-success memoBtn" data-id="{{.Key}}">{{.Name}}</td>
  </tr>
  {{end}}
  </table>
</div>
{{end}}

`
var ConstantView = `
{{define "JavaScript"}}
<link href="/static/css/github.css" rel="stylesheet">
<script src="/static/js/marked.min.js"></script>
<script src="/static/js/highlight.pack.js"></script>
<script src="/static/js/markdown.js"></script>
{{end}}

{{define "menu"}}
<h1>{{.Memo.Name}}</h1>

  <input type="hidden" id="editor" value="{{.Memo.Content}}"/>
  <div class="row" style="background-color:white;">
    <div id="result" style="width:100%;min-height:200px;"></div>
  </div>

<a class="btn btn-success form-control" href="/memo/edit/{{.Memo.Key}}" role="button">Edit</a>
{{end}}

`
