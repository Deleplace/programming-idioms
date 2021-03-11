$(function() {

    $("a.history-copy-imports-to-clipboard").click(function(){
        var that = $(this);
        var importsGroup = that.closest(".imports");
        var snippet = importsGroup.find("pre").text();
        if(!snippet) {
            alert("Sorry, failed to retrieve the imports code :(");
            return;
        }
        //using("TODO");
        navigator.clipboard.writeText(snippet).then(function() {
            console.log('Copying imports to clipboard was successful!');
            that.html('<i class="icon-check" title="The imports code has been copied to clipboard"></i>');

            $(".just-copied-to-clipboard").removeClass("just-copied-to-clipboard");
                importsGroup.addClass("just-copied-to-clipboard");
            }, function(err) {
                alert('Async: Could not copy imports text: ' + err);
            });
        return false;
    });

    $("a.history-copy-code-to-clipboard").click(function(){
        var that = $(this);
        var codeGroup = that.closest(".impl-code");
        var snippet = codeGroup.find("pre").text();
        if(!snippet) {
            alert("Sorry, failed to retrieve the imports code :(");
            return;
        }
        //using("TODO");
        navigator.clipboard.writeText(snippet).then(function() {
            console.log('Copying code snippet to clipboard was successful!');
            that.html('<i class="icon-check" title="The imports code has been copied to clipboard"></i>');

            $(".just-copied-to-clipboard").removeClass("just-copied-to-clipboard");
                codeGroup.addClass("just-copied-to-clipboard");
            }, function(err) {
                alert('Async: Could not copy imports text: ' + err);
            });
        return false;
    });


    $("a.history-copy-comments-to-clipboard").click(function(){
        var that = $(this);
        var commentsGroup = that.closest(".comments");
        var comment = commentsGroup.find(".diff-code-comments").text();
        if(!comment) {
            alert("Sorry, failed to retrieve the imports code :(");
            return;
        }
        comment = comment.trim();
        //using("TODO");
        navigator.clipboard.writeText(comment).then(function() {
            console.log('Copying comments to clipboard was successful!');
            that.html('<i class="icon-check" title="The imports code has been copied to clipboard"></i>');

            $(".just-copied-to-clipboard").removeClass("just-copied-to-clipboard");
                commentsGroup.addClass("just-copied-to-clipboard");
            }, function(err) {
                alert('Async: Could not copy imports text: ' + err);
            });
        return false;
    });
    
});