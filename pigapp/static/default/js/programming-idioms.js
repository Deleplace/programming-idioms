$(function() {

	//
	// jQuery stuff activation
	//
	
	$('button').button();
	$(".sortable-y").sortable({
		handle : ".handle",
        cursor: 'move'
	});
	$(".sortable-lang").sortable({
        cursor: 'move',
        update: function( event, ui ) {
        	updateFavlangCookie();
        }
	});
	$(".implementations-tabs").tabs({
		activate: function( event, ui ) {
			$('pre').popover("show"); // Fix (0,0) popovers of hidden tabs
			emphasize();
		}
	});

	$("pre[data-toggle=popover]").each(function(){
		// Newlines are allowed in Author Comments
		var $this = $(this);
		var content = $this.attr("data-content");
		if(content)
			$this.attr("data-content", "<div class='code-bubble'>" + content.replace(/</g,"&lt;").replace(/\n/g,"<br/>") + "</div>");
	});
	
	$('.togglabe').on('click',function() {
		$(this).toggleClass("active");
	});
	$('div').popover({
		html : true
	}).popover('show');
	$('textarea').popover({
		html : true,
		trigger: 'manual'
	}).popover('show');
	$('pre').popover({
		html : true
	}).popover('show');
	$('a').popover('show');
	$('input').popover({
		trigger: 'manual'
	}).popover('show');
	$('button.show-popover').popover('show');
	$('.popover-on-hover').popover({
		trigger : 'hover focus'
	});
	

	
	$(".idiom-picture img").load(function() {
		// Repaint some bubbles when idiom picture has finished disploying
		$('pre').popover("show");
	});
	
	$(window).resize(function () {
		// Repaint some bubbles on window resize
		$('pre').popover("show");
	});
	
	// Popover: hide on click
	$(document).on("click", ".code-bubble", function(){
		var codeBubble = $(this);
		var popoverContent = codeBubble.parent();
		var popover = popoverContent.parent();
		var pre = popover.prev();
		pre.popover("hide");
	});
	
	$('.input-suggest-language').typeahead({
		source : function(query, process){
	        return $.get(
	        		'/typeahead-languages', 
	        		{ userInput: query }, 
	        		function (data) {
	        			return process(data.options);
	        		});
		},
		matcher: function (item) {
			// Override default behavior.
			// Show all options returned by server.
			// For example, do not remove "C#" when user has typed "csharp"
		    return true;
		},
		updater : function(item){
			var lgDisplay = item;
			var lg = normLang(item);
			var li = $('<li class="active" data-language="'+lg+'"><span class="badge badge-success">'+lgDisplay+' <a href="#" class="favorite-language-remove icon-remove"></a></span></li>');
			li.appendTo($(".favorite-languages"));
	    	updateFavlangCookie();
		}
	});
	$('.language-single-select .typeahead').typeahead({
		source : function(query, process){
	        return $.get(
	        		'/typeahead-languages', 
	        		{ userInput: query }, 
	        		function (data) {
	        			//console.log( "Before process:"+ JSON.stringify(data) );
	        			var processedData = process(data.options)
	        			//console.log( "After process:"+ JSON.stringify(processedData) );
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
	
	//
	// Messages
	//
	
	$.fn.clearMessages = function(){
		$(".message-zone .pimessage").html("");
	}

	$.fn.pisuccess = function(msg){
		$.fn.clearMessages();
		$(".message-zone .alert-success").html(msg);
	}

	$.fn.pierror = function(msg){
		$.fn.clearMessages();
		$(".message-zone .alert-error").html(msg);
	}

	$.fn.piinfo = function(msg){
		$.fn.clearMessages();
		$(".message-zone .alert-info").html(msg);
	}
	
	 $(".ajax-generic-action").on("click", function(){
		 if( $(this).hasClass("confirm-needed") )
			 if( ! window.confirm("Are you sure?") )
				 return;
		
		var url = $(this).attr("data-url");
	    $.ajax({
	        url: url,
	        type: 'POST',
	        xhr: function() {
	            var myXhr = $.ajaxSettings.xhr();
	            return myXhr;
	        },
	        success: function(response){
	        	$.fn.pisuccess( "OK!! " + JSON.stringify(response) );
	        },
	        error: function(xhr, status, e){
	        	$.fn.pierror( xhr.responseText );
	        },
	    });
	 });
	
	// 
	// Authentication (weak)
	//
	var logged = function(){
		var nick = $.cookie("Nickname");
		if( nick )
			return true;
		else
			return false;
	}
	
	$("#modal-nickname .form-nickname").on("submit", function(){
		var nick = $(this).find("input.nickname").val();
		$.cookie("Nickname", nick, { expires : 100, path: '/' });
		$(".greetings").html('<i class="icon-user"> '+ nick +'</i> <a href="#" class="remove-nickname"><i class="icon-remove"></i></a>').show();
		$("#modal-nickname").modal("hide");
	});

	// New-school "Live" binding
	$(document).on("click", ".remove-nickname", function(){
		$.removeCookie("Nickname", { path: '/' });
		$(".greetings").hide();
	});
	
	// 
	// Widgets click events
	//
	
	$('.idiom_cover .count').click(
			function() {
				$(this).children('i').toggleClass(
						'icon-chevron-right icon-chevron-down');
				$(this).parent().children('.full').toggle();
			});

	
	$('.voting-idiom').on('click', function() {
		if( !logged() ){
			$('#modal-nickname').modal();
			return;
		}
		
		$.ajaxSetup({
			  error: function(xhr, status, error) {
				  $.fn.pierror( "Error: " + error);
			  }
		});
		
		var clickedButton = this;
		var clickedButtonWrapper = $(clickedButton);
		clickedButtonWrapper.button('loading');
		var span_voting_score = $(this).parent().next();
		span_voting_score.removeClass("hidden");
		var star = span_voting_score.children("i");
		var idiomId = $(this).attr('data-idiom-id');
		var choice = $(this).attr('data-vote-choice');
		$.get('/ajax-idiom-vote', 
				{idiomId : idiomId,	choice : choice}, 
				function(response) {
					var newScore = response.rating;
					var myVote = response.myVote;
					star.html(" " + newScore);
					clickedButtonWrapper.button('reset');
					$.fn.updateVoteButtonsActiveState(clickedButtonWrapper.parent(), myVote);
				});
	});

	$('.voting-impl').on('click', function() {
		if( !logged() ){
			$('#modal-nickname').modal();
			return;
		}
		
		$.ajaxSetup({
			  error: function(xhr, status, error) {
				  $.fn.pierror( "Error: " + error);
			  }
		});
		
		var clickedButtonWrapper = $(this);
		clickedButtonWrapper.button('loading');
		var span_voting_score = $(this).parent().next();
		span_voting_score.removeClass("hidden");
		var star = span_voting_score.children("i");
		var implId = $(this).attr('data-impl-id');
		var choice = $(this).attr('data-vote-choice');
		$.get('/ajax-impl-vote', 
				{implId : implId,	choice : choice}, 
				function(response) {
					var newScore = response.rating;
					var myVote = response.myVote;
					star.html(" " + newScore);
					clickedButtonWrapper.button('reset');
					$.fn.updateVoteButtonsActiveState(clickedButtonWrapper.parent(), myVote);
				});
	});

	$.fn.updateVoteButtonsActiveState = function(buttonsDiv, voteValue){
		if( voteValue == 1 )
			buttonsDiv.children("[data-vote-choice='up']").addClass("active");
		else
			buttonsDiv.children("[data-vote-choice='up']").removeClass("active");
		if( voteValue == -1 )
			buttonsDiv.children("[data-vote-choice='down']").addClass("active");
		else
			buttonsDiv.children("[data-vote-choice='down']").removeClass("active");
	}

	//
	// Favorite languages
	//

	var updateFavlangCookie = function(){
		var container = $(".favorite-languages")[0];
		var langs = "";
		$(container).children().each( function(i,e){
			var lg = $(e).attr('data-language');
			if(lg)
				langs += lg + "_";
		});
		$.cookie("my-languages", langs,{ expires : 100, path: '/' });
		
		if(langs==""){
			// No favorite langs? Then you really need to see the other langs
			$.cookie("see-non-favorite", "1", { expires : 100, path: '/' });
		}
	};
	
	var normLang = function(lang){
		switch(lang){
		case "C++":
			return "Cpp";
		case "C#":
			return "Csharp";
		}
		return lang;
	}

	var niceLang = function(lang){
		switch(lang){
		case "Cpp":
			return "C++";
		case "Csharp":
			return "C#";
		}
		return lang;
	}
	
	$('.show-languages-pool').on('click', function(){
		$('.addible-languages-pool').show(200);
	});
	
	// New-school "Live" binding
	$(document).on("click", ".addible-languages-pool button", function(){
		var li = $(this).parent();
		var lg = li.attr('data-language');
		var lgDisplay = $(this).html();
		var li = $('<li class="active" data-language="'+lg+'"><span class="badge badge-success">'+lgDisplay+' <a href="#" class="favorite-language-remove icon-remove"></a></span></li>');
		li.hide().appendTo($(".favorite-languages")).show('normal');
    	updateFavlangCookie();
		$(this).hide('normal');

		/* $(".btn-favorite-language-remove").show().removeClass("hidden"); */
		$(".btn-see-non-favorite").show().removeClass("hidden");
	});
	
		
	$(document).on('click', ".favorite-language-remove", function(){
		var a = $(this);
		var span = a.parent();
		a.remove();
		var lgDisplay = span.html();
		var li = span.parent();
		var lg = li.attr('data-language');
		li.removeAttr('data-language');
		li.slideUp(500, function(){ li.remove(); } );
		updateFavlangCookie();

		var liStock = $('<li data-language="'+lg+'"><button class="btn btn-primary btn-mini active togglabe">'+lgDisplay+'</button></li>');
		liStock.hide().prependTo($(".addible-languages-pool ul")).show('normal');			
	});

	$('.btn-see-non-favorite').on('click', function(){
		oldValue = $(this).hasClass('active');
		if( oldValue )
			$.cookie("see-non-favorite", "0", { expires : 100, path: '/' });
		else
			$.cookie("see-non-favorite", "1", { expires : 100, path: '/' });
		location.reload();
	});
	
	//
	// Idiom detail
	//
	function emphasize(){
		$(".identifier-emphasize,.picode div.popover-content").each(function(){
			// Emphasize the "underscored" identifier
			//
			// _x -> <span class="variable">x</span>
			//
			var lead = $(this).html();
			lead = lead.replace( /\b_(\w*)/gm, "<span class=\"variable\">$1</span>");
			$(this).html(lead);
		});
	}
	emphasize();
	
	$('.selector-language').on('click', function(){
		var lg = $(this).closest("li").attr("data-language");
		var nicelg = niceLang(lg);
		$(this).closest(".language-single-select").find("input[type=text]").val(nicelg);
	});
		
	// Lame client-side trick.
	// We should be able to set first tab as "active" server-side.
	// And why do we have to manage click event ourselves?
	$(".implementations-tabs li:first-child").addClass("active");
	$(".implementations-tabs li").on("click", function(){ 
		$(this).parent().children("li").removeClass("active"); 
		$(this).addClass("active"); 
	});
	
	//
	// Forms : idiom creation, impl creation
	//
	$(".form-idiom-creation .language-choices a, .form-impl-creation .language-choices a").on("click", function(){
		var form = $(this).closest(".form-idiom-creation, .form-impl-creation");
		newLang =  form.find("input[name=impl_language]").attr("value");
		$.get('/ajax-demo-site-suggest', 
				{lang : newLang}, 
				function(response) {
					if( response.suggestion )
						form.find("input[name=impl_demo_url]").attr("placeholder", response.suggestion)
				});
	});
	
	// Being able to insert <tab> characters in code
	// See http://stackoverflow.com/questions/6140632/how-to-handle-tab-in-textarea#answer-6140696
	$("textarea").keydown(function(e) {
	    if(e.keyCode === 9) { // tab was pressed
	    	if(! e.ctrlKey){ // but not Ctrl+tab (do not prevent the default browser shortcut)
		        // get caret position/selection
		        var start = this.selectionStart;
		        var end = this.selectionEnd;
	
		        var $this = $(this);
		        var value = $this.val();
	
		        // set textarea value to: text before caret + tab + text after caret
		        $this.val(value.substring(0, start)
		                    + "\t"
		                    + value.substring(end));
	
		        // put caret at right position again (add one for the tab)
		        this.selectionStart = this.selectionEnd = start + 1;
	
		        // prevent the focus lose
		        e.preventDefault();
	    	}
	    }
	});

	// Impl create, impl edit : show other implementations below,
	// read-only, in a defered ajax block
	$(".other-impl-placeholder").each(function(){
		var otherImplDiv = $(this);
		otherImplDiv.html("<i class='icon-spinner icon-spin'></i>");
		var idiomId = otherImplDiv.attr("data-idiom-id");
		var excludedImplId = otherImplDiv.attr("data-excluded-impl-id");
		// window.setTimeout(function(){
		$.get(
	        	'/ajax-other-implementations', 
	        	{ idiomId: idiomId,
	        	  excludedImplId: excludedImplId }, 
	        	function (data) {
					otherImplDiv.html(data);
					otherImplDiv.tabs({
						activate: function( event, ui ) {
							$('pre').popover("show"); // Fix (0,0) popovers of hidden tabs
							emphasize();
						}
					});
					otherImplDiv.find("li:first-child").addClass("active");
					otherImplDiv.find("li").on("click", function(){ 
						$(this).parent().children("li").removeClass("active"); 
						$(this).addClass("active"); 
					});
					$('pre').popover("show");
	        	});
		// }, 3000 );
	});
	
});