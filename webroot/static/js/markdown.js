function change() {
        var src = $("#editor").val();
        var html = marked(src);
        $('#result').html(html);
        $('pre code').each(function(i, block) {
            hljs.highlightBlock(block);
        });
}

$(document).ready(function() {
    marked.setOptions({
        langPrefix: ''
    });

    $('#editor').keyup(function() {
        change();
    });
    change();
});
