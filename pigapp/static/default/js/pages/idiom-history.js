$(function() {

    $("a.history-copy-imports-to-clipboard").on("click", function(){
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
            that.html('<i class="fas fa-clipboard-check" title="The imports code has been copied to clipboard"></i>');

            $(".just-copied-to-clipboard").removeClass("just-copied-to-clipboard");
                importsGroup.addClass("just-copied-to-clipboard");
            }, function(err) {
                alert('Async: Could not copy imports text: ' + err);
            });
        return false;
    });

    $("a.history-copy-code-to-clipboard").on("click", function(){
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
            that.html('<i class="fas fa-clipboard-check" title="The imports code has been copied to clipboard"></i>');

            $(".just-copied-to-clipboard").removeClass("just-copied-to-clipboard");
                codeGroup.addClass("just-copied-to-clipboard");
            }, function(err) {
                alert('Async: Could not copy imports text: ' + err);
            });
        return false;
    });


    $("a.history-copy-comments-to-clipboard").on("click", function(){
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
            that.html('<i class="fas fa-clipboard-check" title="The imports code has been copied to clipboard"></i>');

            $(".just-copied-to-clipboard").removeClass("just-copied-to-clipboard");
                commentsGroup.addClass("just-copied-to-clipboard");
            }, function(err) {
                alert('Async: Could not copy imports text: ' + err);
            });
        return false;
    });
   
    // Client-side diffing (pink, green) with htmldiff.js
    [
        // Idiom data
        ".idiom-left .idiom-summary-large h1 .touched",
        ".idiom-left .idiom-lead-paragraph.touched",
        ".variables .idiom-left .touched span",
        ".related-url .idiom-left .touched span",
        ".keywords .idiom-left .touched span",

        // Impl data
        ".impl-left.impl-code.touched pre",
        ".impl-left.imports.touched pre",
        ".impl-left.touched .diff-code-comments",
        ".doc-url .impl-left .field-value",
        ".origin-url .impl-left .field-value",
        ".demo-url .impl-left .field-value"

    ].forEach( selectorLeft => {
        let selectorRight = selectorLeft
            .replaceAll("idiom-left", "idiom-right")
            .replaceAll("impl-left", "impl-right");
        let leftElem = $(selectorLeft);
        let rightElem = $(selectorRight);
        if(leftElem.length>0 && rightElem.length>0) {
            // E.g.
            //     left ==  "aa bb"
            //     right == "bb cc"
            // =>  delta == "<del>aa bb</del><ins>bb cc</ins>"
            let delta = htmldiff(leftElem.text(), rightElem.text());
            leftElem.html(delta);
            rightElem.html(delta);
            // Note: the diff display is lossy for the Idiom Lead Paragraph,
            // which uses markup2CSS.
        }
    })

});