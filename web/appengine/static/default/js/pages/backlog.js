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

    $(document).on("click", ".backlog .impl-actions button.edit-doc-link", function(){
        let actions = $(this).closest(".impl-actions");
        let idiomID = actions.attr("data-idiom-id");
        let implID = actions.attr("data-impl-id");
        let pageURL = `/impl-edit/${idiomID}/${implID}#doc-url`;
        window.open(pageURL);
    }); 

    $(document).on("click", ".backlog .impl-actions button.edit-demo-link", function(){
        let actions = $(this).closest(".impl-actions");
        let idiomID = actions.attr("data-idiom-id");
        let implID = actions.attr("data-impl-id");
        let pageURL = `/impl-edit/${idiomID}/${implID}#demo-url`;
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

        let parts = endpoint.split("/");
        let lang = parts[2];
        let blockName = parts[4];
        using(`backlog/refresh/${blockName}/${lang}`);
    }); 



	//
	// On page load:
	//
    
	// #187 Auto add favorite language
	var path = window.location.pathname;
	lang = path.substring("/backlog/".length);
	console.debug("Adding favlang", lang);
	addFavlangsInCookie([lang]);
});



//
// ALL CODE BELOW IS DUPLICATED FROM programming-idioms.js
// Duplication is not great but otherwise I get
// "ReferenceError: X is not defined"
// TODO find a more effective code minif+split strategy!
//

function using(what) {
	fetch("/using/"+what, {
		method: "POST",
		body: JSON.stringify({
			page: window.location.pathname+window.location.search
		})
	});
}

function normLang(lang){
	switch(lang.toLowerCase()){
	case "c++":
		return "Cpp";
	case "c#":
		return "Csharp";
	case "cs":
		return "Csharp";
	case "golang":
		return "go";
	case "py":
		return "Python";
	case "rs":
		return "Rust";
	}
	return lang;
}

function capitalizeFirstLetter(string) {
	// https://stackoverflow.com/a/1026087
	return string.charAt(0).toUpperCase() + string.slice(1);
}

function hasFavlangInCookie(lg) {
	lg = normLang(lg);
	let cookielangs = $.cookie("my-languages");
	if(!cookielangs)
		return false;
	let favlangsConcat = cookielangs.toLowerCase();
	let favlangs = favlangsConcat.split("_");
	return favlangs.indexOf(lg.toLowerCase()) !== -1;
}

function addFavlangsInCookie(langs) {
	var langsConcat = $.cookie("my-languages") || "_";
	var newLangsConcat = langsConcat;
	for (var i = 0; i<langs.length; i++) {
		var lang = langs[i];
		lang = capitalizeFirstLetter(lang);
		lang = normLang(lang);
		if(!hasFavlangInCookie(lang)) {
			newLangsConcat += lang + "_";
		}
	}
	if(newLangsConcat != langsConcat){
		$.cookie("my-languages", newLangsConcat,{ expires : 100, path: '/' });
	}
}

//
// END OF CODE DUPLICATED FROM programming-idioms.js
//