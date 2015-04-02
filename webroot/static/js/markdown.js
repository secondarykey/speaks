function change() {
        var src = $("#editor").val();
        var html = marked(src);
        $('#result').html(html);
        $('pre code').each(function(i, block) {
            hljs.highlightBlock(block);
        });
}

$(document).ready(function() {

     var box = $('#memoForm');
     if ( box.length != 0 ) {
         var boxTop = box.offset().top;
         $(window).scroll(function() {
             if( $(window).scrollTop() >= boxTop - 30 ) {
                 box.addClass('memoFix');
             } else {
                 box.removeClass('memoFix');
             }
         });
     }


    marked.setOptions({
        langPrefix: ''
    });

    $('#editor').keyup(function() {
        change();
    });
    change();
});
