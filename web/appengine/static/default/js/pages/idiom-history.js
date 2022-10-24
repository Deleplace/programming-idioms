$(function() {

    $("a.history-copy-imports-to-clipboard").on("click", function(){
        var that = $(this);
        var importsGroup = that.closest(".imports");

        let markup = importsGroup.find("pre").clone();

        if( importsGroup.closest(".impl-left").length > 0 ) {
            // We're in the LEFT column
            // INS parts are hidden, we don't want to copy them
            markup.find("ins").remove();
        }
        if( importsGroup.closest(".impl-right").length > 0 ) {
            // We're in the RIGHT column
            // DEL parts are hidden, we don't want to copy them
            markup.find("del").remove();
        }
        let snippet = markup.text();

        if(!snippet) {
            alert("Sorry, failed to retrieve the imports code :(");
            return false;
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

        let markup = codeGroup.find("pre").clone();

        if( codeGroup.closest(".impl-left").length > 0 ) {
            // We're in the LEFT column
            // INS parts are hidden, we don't want to copy them
            markup.find("ins").remove();
        }
        if( codeGroup.closest(".impl-right").length > 0 ) {
            // We're in the RIGHT column
            // DEL parts are hidden, we don't want to copy them
            markup.find("del").remove();
        }
        let snippet = markup.text();
        if(!snippet) {
            alert("Sorry, failed to retrieve the snippet code :(");
            return false;
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

        let markup = commentsGroup.find(".diff-code-comments").clone();
        if( commentsGroup.closest(".impl-left").length > 0 ) {
            // We're in the LEFT column
            // INS parts are hidden, we don't want to copy them
            markup.find("ins").remove();
        }
        if( commentsGroup.closest(".impl-right").length > 0 ) {
            // We're in the RIGHT column
            // DEL parts are hidden, we don't want to copy them
            markup.find("del").remove();
        }
        let comment = markup.text();

        if(!comment) {
            alert("Sorry, failed to retrieve the comments :(");
            return false;
        }
        comment = comment.trim();
        //using("TODO");
        navigator.clipboard.writeText(comment).then(function() {
            console.log('Copying comments to clipboard was successful!');
            that.html('<i class="fas fa-clipboard-check" title="The comments have been copied to the clipboard"></i>');

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
        ".impl-left.impl-code.touched pre > code",
        ".impl-left.imports.touched pre > code",
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
            const left = leftElem.text().replaceAll("<", "&lt;").replaceAll(">", "&gt;");
            const right = rightElem.text().replaceAll("<", "&lt;").replaceAll(">", "&gt;");
            let delta = htmldiff(left, right);
            leftElem.html(delta);
            rightElem.html(delta);
            // Note: the diff display is lossy for the Idiom Lead Paragraph,
            // which uses markup2CSS.
        }
    })

});