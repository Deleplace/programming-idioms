$(function() {

    $(document).on("click", ".backlog .impl-actions button.view", function(){
        let actions = $(this).closest(".impl-actions");
        let idiomID = actions.attr("data-idiom-id");
        let implID = actions.attr("data-impl-id");
        let pageURL = `/idiom/${idiomID}/impl/${implID}`;
        window.open(pageURL);
    }); 

    $(document).on("click", ".backlog .impl-actions button.edit", function(){
        let actions = $(this).closest(".impl-actions");
        let idiomID = actions.attr("data-idiom-id");
        let implID = actions.attr("data-impl-id");
        let pageURL = `/impl-edit/${idiomID}/${implID}`;
        window.open(pageURL);
    }); 

    $(document).on("click", ".backlog .impl-actions button.mark-good", function(){
        // TODO: mark good only if user has a Nickname. Include the Nickname in the "vote" log.
        let actions = $(this).closest(".impl-actions");
        let idiomID = actions.attr("data-idiom-id");
        let implID = actions.attr("data-impl-id");
        using(`backlog/mark-as-good/${idiomID}/impl/${implID}`);
        alert( "Thank you for this positivity :)" );
    }); 

    $(document).on("click", ".backlog .idiom-actions button.create-impl", function(){
        let actions = $(this).closest(".idiom-actions");
        let idiomID = actions.attr("data-idiom-id");
        let lang = actions.attr("data-missing-lang");
        let pageURL = `/impl-create/${idiomID}/${lang}`;
        window.open(pageURL);
    }); 

    $(document).on("click", ".backlog .idiom-actions button.view", function(){
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

    $(".btn.block-data-refresh").click(function(){
        let btn = $(this)
        let endpoint = btn.attr('data-block-endpoint');
        if(!endpoint) {
            console.error(`No endpoint, no block refresh!`);
            return;
        }
        let target = $(this).siblings('.block-data-contents');
        if(!target) {
            console.error(`Couldn't find the block-data-contents`);
            return;
        }
        target.addClass("refreshing");
        btn.addClass("refreshing");


		$.get(endpoint, 
            {}, 
            function(response) {
                target.html(response);
                target.removeClass("refreshing");
                btn.removeClass("refreshing");
                $('pre[data-content]').popover({
                    html : true
                }).popover('show');
            });
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