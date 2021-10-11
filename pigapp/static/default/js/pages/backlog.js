$(function() {

    $(".backlog .impl-actions button.view").click(function(){
        let actions = $(this).closest(".impl-actions");
        let idiomID = actions.attr("data-idiom-id");
        let implID = actions.attr("data-impl-id");
        let pageURL = `/idiom/${idiomID}/impl/${implID}`;
        window.open(pageURL);
    }); 

    $(".backlog .impl-actions button.edit").click(function(){
        let actions = $(this).closest(".impl-actions");
        let idiomID = actions.attr("data-idiom-id");
        let implID = actions.attr("data-impl-id");
        let pageURL = `/impl-edit/${idiomID}/${implID}`;
        window.open(pageURL);
    }); 

    $(".backlog .impl-actions button.mark-good").click(function(){
        // TODO: mark good only if user has a Nickname. Include the Nickname in the "vote" log.
        let actions = $(this).closest(".impl-actions");
        let idiomID = actions.attr("data-idiom-id");
        let implID = actions.attr("data-impl-id");
        using(`backlog/mark-as-good/${idiomID}/impl/${implID}`);
        alert( "Thank you for this positivity :)" );
    }); 

    $(".backlog .idiom-actions button.create-impl").click(function(){
        let actions = $(this).closest(".idiom-actions");
        let idiomID = actions.attr("data-idiom-id");
        let lang = actions.attr("data-missing-lang");
        let pageURL = `/impl-create/${idiomID}/${lang}`;
        window.open(pageURL);
    }); 

    $(".backlog .idiom-actions button.view").click(function(){
        let actions = $(this).closest(".idiom-actions");
        let idiomID = actions.attr("data-idiom-id");
        let pageURL = `/idiom/${idiomID}`;

        // "View full idiom" may be better if it shows this impl at the top
        let implID = actions.attr("data-impl-id");
        if(implID) {
            pageURL += `/impl/${implID}`;
        }
        window.open(pageURL);
    }); 


    // THIS IS DUPLICATED FROM programming-idioms.js
    // Duplication is not great but otherwise I get
    // "ReferenceError: using is not defined"
    function using(what) {
        fetch("/using/"+what, {
            method: "POST",
            body: JSON.stringify({
                page: window.location.pathname+window.location.search
            })
        });
    }
});