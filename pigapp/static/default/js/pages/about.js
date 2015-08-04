$(function() {

	//
	// Left menu
	//
	var loadAboutCentral = function(url, tempo){
		$(".about-central-zone").fadeOut(tempo);
		$(".about-left-menu li").removeClass("active");
		$(".about-left-menu a").each( function(){
			var dbu = $(this).attr("data-block-url");
			if( url == dbu ){
				var li = $(this).parent(); 
				li.addClass("active");
			}
		});
        $.get(url,function(data){
            $(".about-central-zone").fadeOut(10,function(){$(".about-central-zone").html(data);});
            $(".about-central-zone").fadeIn(tempo);
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
		var langA = $(langTh).children("a").first();
		var lang = langA.html();
		
		var raw = "Click to see idiom "+idiomId+" in "+lang;
		if( href.indexOf("/impl-create") >= 0 )
			raw = "Click to create implementation in "+lang;
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
	
});