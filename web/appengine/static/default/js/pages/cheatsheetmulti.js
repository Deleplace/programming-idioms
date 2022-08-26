$(function() {

	function applyCheatsheetFilters() {
		var word = $("#filter").val();
		var restrictorLangs = new Set();
		$("input.restrict-having").each(function(){
			var lang = $(this).attr("data-lang")
			if(!lang)
				return;
			var checked = $(this).is(':checked');
			if(!checked)
				return;
			restrictorLangs.add(lang);
		});

		$("tr.cheatsheet-line").each(function(){
			var line = $(this);
			var show = true;
		
			// Full-text (raw, no tokenization)
			if(word){
				var lowerHtml = line.html().toLowerCase();
				var lowerWord = word.toLowerCase();
				if( lowerHtml.indexOf(lowerWord) === -1 ){
					show = false;
				}
			}

			// Restrict to existing impls
			restrictorLangs.forEach(function(lang) {
				var cell = line.find("td.lang-" + lang);
				if(cell.length==0) {
					console.log("Table cell for " + lang + " not found!");
					return;
				}
				if(cell.text().trim() === "") {
					// No impl in this language, let's hide the whole line
					show = false;
				}
			});

			if(show) {
				line.show('normal');
			} else {
				line.hide('normal');
			}
		});
	}

	$("button.page-print").on("click", function(){
		using("print");
		window.print();
	});

	$(".cheatsheet-lines button.close").on("click", function(){
		var line = $(this).closest("tr");
		var idiomID = line.find("th.idiom-id").text();
		line.remove();	
		using("cheatsheet/remove-line/" + idiomID);
	});

	$(".page-cheatsheet #showIdiomId").on("change", function(){
		if( $(this).is(':checked') ){
			$("th.idiom-id").show();
			using("cheatsheet/options/idiom-id/show");
		}else{
			$("th.idiom-id").hide();
			using("cheatsheet/options/idiom-id/hide");
		}
	});

	$(".page-cheatsheet #showImports").on("change", function(){
		if( $(this).is(':checked') ){
			$(".piimports").show();
			using("cheatsheet/options/imports/show");
		}else{
			$(".piimports").hide();
			using("cheatsheet/options/imports/hide");
		}
	});

	$(".page-cheatsheet #showComments").on("change", function(){
		if( $(this).is(':checked') ){
			$(".impl-comment").show();
			using("cheatsheet/options/comments/show");
		}else{
			$(".impl-comment").hide();
			using("cheatsheet/options/comments/hide");
		}
	});

	$(".page-cheatsheet #showExternalLinks").on("change", function(){
		if( $(this).is(':checked') ){
			$(".impl-external-links").show();
			using("cheatsheet/options/impl-external-links/show");
		}else{
			$(".impl-external-links").hide();
			using("cheatsheet/options/impl-external-links/hide");
		}
	});

	$(".page-cheatsheet #filter").on("change", function(){
		applyCheatsheetFilters();
		var word = $("#filter").val();
		using("cheatsheet/options/filter/" + word); // spaces will appear as %20
	});

	$("input.restrict-having").on("change", function(){
		applyCheatsheetFilters();
		var lang = $(this).attr("data-lang");
		if( $(this).val() )
			using("cheatsheet/restrict/" + lang);
		else
			using("cheatsheet/unrestrict/" + lang);
	});

	//
	// On page load:
	//
    
	// #134 Auto add favorite languages
	var path = window.location.pathname;
	path = path.substring("/cheatsheet/".length);
	let langs = path.split(/\//);
	console.debug("Adding favlangs", langs);
	addFavlangsInCookie(langs);

    // #186 In Firefox, checkboxes may come already checked from last visit
	if( $(".page-cheatsheet #showImports").is(':checked') ){
		$(".piimports").show();
	};
	if( $(".page-cheatsheet #showComments").is(':checked') ){
		$(".impl-comment").show();
    }
	if( $(".page-cheatsheet #showExternalLinks").is(':checked') ){
		$(".impl-external-links").show();
    }

	// Filter out lines "assumed uninteresting"
	applyCheatsheetFilters();
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