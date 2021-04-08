$(function() {

	//
	// Left menu
	//
	var loadAboutCentral = function(url, tempo){
		let fetched = false;
		let centralZone = $(".about-central-zone");
		centralZone.fadeOut(tempo,function(){
			if(fetched)
				return;
			centralZone.html(`<img src="/default/img/wheel.svg" class="throbber spinning-jolty2" />`);
			centralZone.show();
		});
		$(".about-left-menu li").removeClass("active");
		$(".about-left-menu a").each( function(){
			var dbu = $(this).attr("data-block-url");
			if( url == dbu ){
				var li = $(this).parent(); 
				li.addClass("active");
			}
		});
        $.get(url,function(data){
			fetched = true;
            centralZone.fadeOut(10,function(){
				centralZone.html(data);
				centralZone.fadeIn(tempo, function(){
					initLanguageTypeahead();
				});
			});
        });
	}
	
	$(".about-left-menu a").on("click", function(){
		var url = $(this).attr("data-block-url");
		loadAboutCentral(url, 250);
	});
	
	//
	// Language coverage
	//
	
	// Highlight hovered row
	$(document).on({
		mouseenter: function(){
			$("tr.highlight").removeClass("highlight");
			$(this).addClass("highlight");
		}
	}, ".language-coverage tr.highlightable");
	
	// Highlight hovered column
	$(document).on({
		mouseenter: function(){
			$("colgroup").removeClass("highlight");
			$("colgroup").eq($(this).index()).addClass("highlight");
		}
	}, ".language-coverage td");
	
	function showCoverageCellBubble(link){
		var td = link.parent();
		var tr = td.parent();
		var tbody = tr.parent();
		var table = tbody.parent();
		var thead = table.find("thead").first();
		var href = link.attr("href");
		
		var thLine = tr.children("th").first();
		var idiomId = $(thLine).attr("data-idiom-id");
		
		var index = td.index();
		var firstLine = $(thead).children("tr").first();
		var langTh = $(firstLine).children("th")[index];
		var lang = $(langTh).html();
		
		var raw = "Click to see idiom " + idiomId + " in " + lang;
		if( href.indexOf("/impl-create") >= 0 )
			raw = "Click to create implementation in " + lang;
		var content = "<div class='coverage-cell-bubble'>" + raw + "</div>";
		link.addClass("viewIdiomPop");
		
		var popo = link.popover({
			html : true,
			content : content,
			trigger: "click"
		});
		popo.popover('show');
	}
	function hideCoverageCellBubbles(){
		$(".viewIdiomPop").popover("hide");
		$(".viewIdiomPop").removeClass("viewIdiomPop");
	}

	$(document).on({
		mouseenter: function(){
			var link = $(this);
			showCoverageCellBubble(link);
		}, mouseout: function(){
			hideCoverageCellBubbles();
		}
	}, ".language-coverage td a");
	
	$(document).on({
		click: function(){
			$(".language-coverage .hidden").removeClass("hidden");
			$(".fold-unfold").hide();
		}
	}, ".fold-unfold")
	
	// Bookmarkable anchors
	if( window.location.hash.indexOf("#about-block-") != -1 ){
		var ajaxUrl = window.location.hash;
		ajaxUrl = "/" + ajaxUrl.substr(1);  // Replace first # with /
		loadAboutCentral(ajaxUrl, 10);
	}
	
	// This is REDUNDANT CODE copy-pasted from programming-idioms.js
	// TODO Remove id possible
	function initLanguageTypeahead() {
		$('.language-single-select .typeahead').typeahead({
			source : function(query, process){
				return $.get(
						'/typeahead-languages', 
						{ userInput: query }, 
						function (data) {
							let processedData = {};
							if(data && data.options)
								processedData = process(data.options);
							return processedData;
						});
			},    
			matcher: function (item) {
				// Override default behavior.
				// Show all options returned by server.
				// For example, do not remove "C#" when user has typed "csharp"
				return true;
			}
		});
	}
});