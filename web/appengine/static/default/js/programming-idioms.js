$(function() {

	var YYYYMMDD = function(dateStr) {
		var date = new Date(dateStr);
	    var d = date.getDate();
	    var m = date.getMonth() + 1;
	    var y = date.getFullYear();
	    return '' + y + '-' + (m<=9 ? '0' + m : m) + '-' + (d <= 9 ? '0' + d : d);
	}

	var HHmm = function(dateStr) {
		var date = new Date(dateStr);
	    var h = date.getHours();
	    var m = date.getMinutes();
	    return '' + (h <= 9 ? '0' + h : h) + ':' + (m<=9 ? '0' + m : m);
	}

	var YYYYMMDDHHmm = function(dateStr) {
		return YYYYMMDD(dateStr) + ' ' + HHmm(dateStr);
	}

	//
	// jQuery stuff activation
	//
	
	$('button').button();
	$(".sortable-lang").sortable({
        // cursor: 'move',
        update: function( event, ui ) {
        	updateFavlangCookie();
        }
	});
	$(".implementations-tabs").tabs({
		activate: function( event, ui ) {
			$('pre[data-content]').popover("show"); // Fix (0,0) popovers of hidden tabs
		}
	});
	
	$('.togglabe').on('click',function() {
		$(this).toggleClass("active");
	});
	$('div[data-content]').popover({
		html : true
	}).popover('show');
	$('textarea[data-content]').popover({
		html : true,
		trigger: 'manual'
	}).popover('show');
	$('pre[data-content]').popover({
		html : true
	}).popover('show');
	$(document).on("click", ".popover-content", function (){
		// Attached to <pre>: this is the detail view.
		// We want to hide the bubble on bubble click.
		let pre = $(this).closest(".picode").children("pre");
		if(pre.length == 1) {
			let bubble = $(this).closest(".popover");
			bubble.hide( "slide", {direction: "left"}, 200, function(){
				pre.popover("toggle");
			});
		}
		// Attached to <textarea>: this is the edit view.
		// We don't want to hide the bubble.
	});
	// $('a').popover('show');
	$('input[data-content]').popover({
		trigger: 'manual'
	}).popover('show');
	$('button.show-popover').popover('show');
	$('.popover-on-hover').popover({
		trigger : 'hover focus'
	});


	
	$(".idiom-picture img").on("load", function() {
		// Repaint some bubbles when idiom picture has finished disploying
		$('pre[data-content]').popover("show");
	});
	
	$(window).on("resize", function () {
		// Repaint some bubbles on window resize
		$('pre[data-content]').popover("show");
	});
		
	$('.input-suggest-language').typeahead({
		source : function(query, process){
	        return $.get(
	        		'/typeahead-languages', 
	        		{ userInput: query }, 
	        		function (data) {
						data.options = data.options || [];
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
			addFavlang(item);
		}
	});
	$('.language-single-select .typeahead').typeahead({
		source : function(query, process){
	        return $.get(
	        		'/typeahead-languages', 
	        		{ userInput: query }, 
	        		function (data) {
						let processedData = {};
						// console.log( "Before process:"+ JSON.stringify(data) );
						if(data && data.options)
	        				processedData = process(data.options);
	        			// console.log( "After process:"+ JSON.stringify(processedData) );
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

		var data = {};
		if( $(this).hasClass("reason-needed") ) {
			var reason = window.prompt("Why?");
		 	if( reason===null )
				return; // Clicked Cancel
			data = {reason: reason};
		}
		
		var url = $(this).attr("data-url");
	    $.ajax({
	        url: url,
			type: 'POST',
			data: data,
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
	// Identification
	// (weak, no proper authentication)
	//
	var logged = function(){
		var nick = $.cookie("Nickname");
		if( nick )
			return true;
		else
			return false;
	}

	var username = function(){
		return $.cookie("Nickname");
	}
	
	$("#modal-nickname .form-nickname").on("submit", function(){
		var nick = $(this).find("input.nickname").val();
		if( nick.length>30 )
			nick = nick.substring(0,30);
		$.cookie("Nickname", nick, { expires : 100, path: '/' });
		$(".greetings").html('<i class="fas fa-user"> '+ nick +'</i> <a href="#" class="remove-nickname fas fa-times" alt="Clear username" title="Clear username"></a>').show();
		$("#modal-nickname").modal("hide");
	});

	// New-school "Live" binding
	$(document).on("click", ".remove-nickname", function(){
		using("remove-nickname");
		$.removeCookie("Nickname", { path: '/' });
		$(".greetings").hide();
	});

	$(".user-info-link a").on("click", function() {
		var headerCode =  '<div class="modal-header">'
						+ '	<button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>'
						+ '	<h3>Cookie contents</h3>'
						+ '</div>';
		var header = $(headerCode);
		var body = $("<div>").addClass("modal-body");
		var dlNickname = $("<dl>");
		if(logged()) {
			var dt1 = $("<dt><tt>Nickname</tt></dt>");
			var removeBtn = $("<button>").text("Delete this cookie");
			dt1.append(removeBtn);
			var dd1 = $("<dd>").text(username());
			removeBtn.on("click", function(){
				using("cookie-box/delete-cookie-nickname");
				$.removeCookie("Nickname", { path: '/' });
				dlNickname.hide("slow", function(){ dlNickname.remove(); });
				takeaway.hide("slow", function(){ takeaway.remove(); });
				updateProfileUrl();
				$("p.greetings").hide("slow", function(){ $("p.greetings").remove(); });
			});
			dlNickname.append(dt1).append(dd1);
			body.append(dlNickname);
		}

		var dlLangs = $("<dl>");
		var langsConcat = $.cookie("my-languages");
		if( langsConcat ){
			var dt = $("<dt><tt>my-languages</tt></dt>");
			var removeBtn = $("<button>").text("Delete this cookie");
			removeBtn.on("click", function(){
				using("cookie-box/delete-cookie-favlangs");
				$.removeCookie("my-languages", { path: '/' });
				dlLangs.hide("slow", function(){ dlLangs.remove(); });
				updateProfileUrl();
				$("ul.favorite-languages li").hide("slow", function(){ $("ul.favorite-languages li").remove(); });
			});
			dt.append(removeBtn);
			dlLangs.append(dt);
			var langs = langsConcat.split(/_/);
			langs.forEach(function(lang) {
				var dd = $("<dd>").text(lang);
				dlLangs.append(dd);
			});
			body.append(dlLangs);
		}
		var takeaway = $("<div>").addClass("profile-take-away");
		takeaway.append("<h4>To take your profile with you</h4>")
		takeaway.append("<p>Copy this URL. Profiles are not stored on server, only in cookies or in this URL.</p>")
		urlbox = $("<input>").attr("type", "text").addClass("profile-url");
		takeaway.append(urlbox);
		let copyProfileLink = $('<a href="#" class="copy-url-clipboard"><i class="far fa-copy" title="Copy to the clipboard"></i></a>');
		copyProfileLink.on("click", function(){ 
			using("cookie-box/copy-profile-url");
			navigator.clipboard.writeText(urlbox.val()).then(function() {
				console.log('Copying to clipboard was successful!');
				copyProfileLink.html('<i class="fas fa-clipboard-check" title="The profile link has been copied to clipboard"></i>');
			  }, function(err) {
				alert('Async: Could not copy text: ' + err);
			  });
		});
		takeaway.append(copyProfileLink);
		body.append(takeaway);
		var fullhost = location.protocol+'//'+location.hostname+(location.port ? ':'+location.port: '');
		var updateProfileUrl = function(){
			var nick = username();
			var lgs = $.cookie("my-languages");
			if(!lgs)
				lgs = "";
			var profileUrl = fullhost
				+ "/my/" + encodeURIComponent(nick)
				+ "/_" + lgs;
			urlbox.val(profileUrl);
		}
		updateProfileUrl();
		$("<div>").addClass("modal")
			.addClass("profile-box")
			.append(header)
			.append(body)
			.modal("show");
		urlbox.select();
		return false;
	});
	
	// 
	// Widgets click events
	//
	
	$('.idiom_cover .count').on("click",
			function() {
				$(this).children('i').toggleClass(
						'icon-chevron-right icon-chevron-down');
				$(this).parent().children('.full').toggle();
			});

	
	$('.voting-idiom').on('click', function() {
		if( !logged() ){
			$('#modal-nickname').modal({
				keyboard: true
			});
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

	function addFavlang(lg) {
		let lgDisplay = niceLang(lg);
		lg = normLang(lg);

		// Favlang bar
		var li = $(`
			<li class="active" data-language="${lg}">
				<span class="badge badge-success">
					${lgDisplay}
					<a href="#" class="favorite-language-remove fas fa-times"></a>
				</span>
			</li>`);
		li.appendTo($(".favorite-languages"));

		// Right column cheatsheet links
		li = $(`<li>
					<i class="fas fa-list"></i> <a href="/cheatsheet/${lg}">${lgDisplay}</a>
				</li>`);
		li.appendTo($(".right_column .cheatsheet-links ul"));

		// Cookie
		updateFavlangCookie();
	}

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
		
	$(document).on('click', '.show-languages-pool', function(){
		$('.addible-languages-pool').show(200);
	});
	
	// New-school "Live" binding
	$(document).on("click", ".addible-languages-pool button", function(){
		var li = $(this).parent();
		var lg = li.attr('data-language');
		var lgDisplay = $(this).html();
		var li = $('<li class="active" data-language="'+lg+'"><span class="badge badge-success">'+lgDisplay+' <a href="#" class="favorite-language-remove fas fa-times" title="Remove fav"></a></span></li>');
		li.hide().appendTo($(".favorite-languages")).show('normal');
		using("fav-lang-bar-add/" + lg);
    	updateFavlangCookie();
		$(this).hide('normal');

		/* $(".btn-favorite-language-remove").show().removeClass("hidden"); */
		$(".btn-see-non-favorite").show().removeClass("disabled");
	});
	
		
	$(document).on('click', ".favorite-language-remove", function(){
		var a = $(this);
		var span = a.parent();
		a.remove();
		var lgDisplay = span.html();
		var li = span.parent();
		var lg = li.attr('data-language');
		using("fav-lang-bar-remove/" + lg);
		li.removeAttr('data-language');
		li.slideUp(500, function(){ li.remove(); } );
		updateFavlangCookie();

		var liStock = $('<li data-language="'+lg+'"><button class="btn btn-primary btn-mini active togglabe">'+lgDisplay+'</button></li>');
		liStock.hide().prependTo($(".addible-languages-pool ul")).show('normal');			
	});

	$('.btn-see-non-favorite').on('click', function(){
		oldValue = $(this).hasClass('active');
		if( oldValue ) {
			$.cookie("see-non-favorite", "0", { expires : 100, path: '/' });
			using("see-non-favorite/set/0");
		} else {
			$.cookie("see-non-favorite", "1", { expires : 100, path: '/' });
			using("see-non-favorite/set/1");
		}
		location.reload();
	});
	
	//
	// Idiom detail
	//
	
	$(document).on('click', '.selector-language', function(){
		var lg = $(this).closest("li").attr("data-language");
		var nicelg = niceLang(lg);
		$(this).closest(".language-single-select").find("input[type=text]").val(nicelg).change();
		return false;
	});

	function isIdiomPage() {
		// Diff pages are not concerned
		if ( /\/diff\//.test(window.location.pathname) )
			return false;

		return /^\/idiom\//.test(window.location.pathname);
	}

	function isIdiomDetailWithLang() {
		// Diff pages are not concerned
		if ( /\/diff\//.test(window.location.pathname) )
			return false;

		// E.g. "/idiom/52/check-if-map-contains-value/2870/csharp"
		return /\/idiom\/[0-9]+\/[^/]+\/[0-9]+\/[^/]+/.test(window.location.pathname);
	}

	function extractLangFromURL() {
		let parts = window.location.pathname.split(/\//);
		let lang = parts[parts.length-1];
		lang = capitalizeFirstLetter(lang);
		return lang;
	}

	if( isIdiomDetailWithLang() ) {
		// #112 Auto add favorite languages
		let lang = extractLangFromURL();
		if(!hasFavlangInCookie(lang)) {
			addFavlang(lang);
		}
	}
		
	// Lame client-side trick.
	// We should be able to set first tab as "active" server-side.
	// And why do we have to manage click event ourselves?
	$(".implementations-tabs li:first-child").addClass("active");
	$(".implementations-tabs li").on("click", function(){ 
		$(this).parent().children("li").removeClass("active"); 
		$(this).addClass("active"); 
	});

	// Impl grid view (exposé-like) for current idiom.
	function showImplGrid(){
		var idiomID = $(this).closest(".idiom-summary-large").attr("data-idiom-id");
		using("grid/" + idiomID);
		$(".modal-impl-grid").modal();
	};
	$('.show-impl-grid').on('click', function(){
		showImplGrid.call(this);
	});

	$(".copy-imports-to-clipboard a").on("click", function(){
		var that = $(this);
		var impl = that.closest(".implementation");
		var piimports = impl.find(".piimports");
		var snippet = piimports.find("pre").text();
		if(!snippet) {
			alert("Sorry, failed to retrieve the imports code :(");
			return false;
		}
		using("impl/copy-imports-to-clipboard/" + impl.attr("data-idiom-id") + "/" + impl.attr("data-impl-id") + "/" + impl.attr("data-impl-lang"));
		navigator.clipboard.writeText(snippet).then(function() {
			console.log('Copying imports to clipboard was successful!');
			that.html('<i class="fas fa-clipboard-check" title="The imports code has been copied to clipboard"></i>');

			$(".just-copied-to-clipboard").removeClass("just-copied-to-clipboard");
			piimports.addClass("just-copied-to-clipboard");
		  }, function(err) {
			alert('Async: Could not copy imports text: ' + err);
		  });
		return false;
	});

	$("a.copy-code-to-clipboard").on("click", function(){
		var button = $(this);
		var impl = button.closest(".implementation");
		clipboardCopyCodeOfImpl(button, impl);
		using("impl/copy-to-clipboard/" + impl.attr("data-idiom-id") + "/" + impl.attr("data-impl-id") + "/" + impl.attr("data-impl-lang"));
		return false;
	});

	function clipboardCopyCodeOfImpl(copyButtonElement, implElement) {
		var picode = implElement.find(".picode");
		var snippet = picode.find("pre").text();
		if(!snippet) {
			alert("Sorry, failed to retrieve the snippet code :(");
			return false;
		}
		navigator.clipboard.writeText(snippet).then(function() {
			console.log('Copying to clipboard was successful!');
			copyButtonElement.html('<i class="fas fa-clipboard-check" title="The snippet code has been copied to clipboard"></i>');

			$(".just-copied-to-clipboard").removeClass("just-copied-to-clipboard");
			implElement.addClass("just-copied-to-clipboard");
		  }, function(err) {
			alert('Async: Could not copy text: ' + err);
		  });
	}

	$(".impl-external-links .demo a").on("click", function(event) {
		var impl = $(this).closest(".implementation");
		using("impl/demo/"+ impl.attr("data-idiom-id") + "/" + impl.attr("data-impl-id") + "/" + impl.attr("data-impl-lang"));
	});

	$(".impl-external-links .doc a").on("click", function(event) {
		var impl = $(this).closest(".implementation");
		using("impl/doc/"+ impl.attr("data-idiom-id") + "/" + impl.attr("data-impl-id") + "/" + impl.attr("data-impl-lang"));
	});

	$(".impl-external-links .origin a").on("click", function(event) {
		var impl = $(this).closest(".implementation");
		using("impl/origin/"+ impl.attr("data-idiom-id") + "/" + impl.attr("data-impl-id") + "/" + impl.attr("data-impl-lang"));
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

	// $("input[name=impl_language]").on("autocompletechange", function(event,ui) {
	// 	alert("autocompletechange");
	//  });

	$("input[name=impl_language]").on("change", function() {
		let inputField = $(this);
		let userinput = inputField.val();
		let userinputlower = userinput.toLowerCase();
		let group = inputField.closest(".control-group");
		if(group.length == 0)
			return;
		let message = group.find(".help-inline");
		if(userinput === "") {
			// Input is empty, not worth an error message right now.
			// Field is still required though, and submitting the form would be explicitly denied.
			group.removeClass("error");
			message.text("");
			return;
		}
		$.get('/supported-languages', 
				{}, 
				function(response) {
					window.setTimeout( () => {
						userinput = inputField.val(); // in some cases, this value has changed here
						userinputlower = userinput.toLowerCase();
						if( response.languages ) {
							if( response.languages.find( function(lang) { return lang.toLowerCase() === userinputlower; }) ){
								console.debug(`${userinput} is a legit lang, clearing any error`);
								group.removeClass("error");
								message.text("");
							}else{
								console.warn( userinput + " is currently not a supported language");
								group.addClass("error");
								message.text(userinput + " is currently not a supported language. Supported languages are " + response.languages.join(", "));
								group.find(".language-single-select").popover("hide");
							}
						}
					}, 200); // The change event is fired twice with possibly a different input, unfortunately. A small delay helps.
				});
	});

	$("textarea[name=impl_code").on("change", function() {
		let warnZone = $(".warning-code-cromulence");
		warnZone.hide().empty();
		function warn(line) {
			warnZone.append( $("<div>").html(line) ).show();
		}
		let code = $(this).val();

		let expectedVarsComma = $(this).attr("data-variables");
		if(expectedVarsComma) {
			expectedVarsComma = expectedVarsComma.replace(/[_ ]/g,"");
			let vars = expectedVarsComma.split(",");
			let missing = [];
			for(let i=0;i<vars.length;i++) {
				if ( !new RegExp('\\b' + vars[i].toLowerCase() + '\\b').test(code.toLowerCase()) ) {
					missing.push(vars[i]);
				}
			}
			if(missing.length >= 1) {
				let plural = (missing.length) >= 2 ? "s" : "";
				let missingBold = missing.map(function(v){return "<em>" + v + "</em>" });
				let warning = "The code <i>should</i> contain the identifier" + plural + " " + missingBold.join(", ") + ".";
				warn(warning);
			}
		}
		
		if( /\bmain\b/.test(code) ) {
			warn("Are you sure about <em>main</em>? We usually don't want a whole program.");
		}
	});
	
	// Being able to insert <tab> characters in code
	// See https://stackoverflow.com/questions/6140632/how-to-handle-tab-in-textarea#answer-6140696
	$("textarea").on("keydown", function(e) {
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

	// Impl flag (to the admin)
	$(document).on("click", ".btn-flag-impl", function(e){
		let btn = $(e.target).closest(".btn-flag-impl");
		let rationale;
		do {
			rationale = window.prompt("I'd like to report this implementation because:");
			if( rationale===null )
				   return; // Clicked Cancel
		} while(rationale === "");

		let idiomId = btn.attr('data-idiom-id');
		let implId = btn.attr('data-impl-id');
		let idiomVersion = btn.attr('data-idiom-version');

		$.post(
			'/ajax-impl-flag/' + idiomId + '/' + implId,
			{
				"idiomVersion": idiomVersion,
				"rationale": rationale
			}
		).done(function() {
			alert( "Thanks :)" );
		}).fail(function() {
			alert( "Unfortunately we could not save this report :(" );
		});
	});

	// Impl create, impl edit : show other implementations below,
	// read-only, in a defered ajax block
	//
	// 2015-12-23  ajax fetch deactivated because
	// doesn't play well with escaping of bubbles text.
	/*
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
							$('pre[data-content]').popover("show"); // Fix (0,0) popovers of hidden tabs
						}
					});
					otherImplDiv.find("li:first-child").addClass("active");
					otherImplDiv.find("li").on("click", function(){ 
						$(this).parent().children("li").removeClass("active"); 
						$(this).addClass("active"); 
					});
					//displayCodeCommentBubble();
					$('pre[data-content]').popover("show");
	        	});
		// }, 3000 );
	});
    */
	
	//
	// Impl create, impl edit : [Preview] button injects values
	// in modal window.
	//

	// This client-side formatting should be rarely used : only in Previews.
	function emphasize(raw){
		// Emphasize the "underscored" identifier
		//
		// _x -> <em>x</em>
		//
		var refined = raw.replace( /\b_([\w$]*)/gm, "<em>$1</em>");
		refined = refined.replace(/\n/g,"<br/>");
		return refined;
	}

	function showImplCreatePreview(){
			$('pre[data-content]').popover("hide"); // Hide (0,0) popovers of hidden tabs
			var m = $('.modal-impl-preview');
			var lang = $(".form-impl-creation input[name=impl_language]").val();
			m.find(".lang-tab span.label").html(lang);
			var imports = $(".form-impl-creation textarea.imports").val();
			if( imports )
				m.find(".piimports pre > code").text( imports ).show();
			else
				m.find(".piimports pre > code").text( imports ).hide();
			m.find(".picode pre > code").text( $(".form-impl-creation textarea.impl-code").val() );
			var comment = $(".form-impl-creation textarea[name=impl_comment]").val();
			var escapedComment = $("<div>").text(comment).html();
			var refinedComment = emphasize(escapedComment);
			m.find(".picode pre").attr("data-content", refinedComment);
			var extDocURL = $(".form-impl-creation input[name=impl_doc_url]").val();
			if( extDocURL )
				m.find("a.impl-doc").attr("href", extDocURL).show();
			else
				m.find("a.impl-doc").attr("href", "#").hide();
			var extDemoURL = $(".form-impl-creation input[name=impl_demo_url]").val();
			if( extDemoURL )
				m.find("a.impl-demo").attr("href", extDemoURL).show();
			else
				m.find("a.impl-demo").attr("href", "#").hide();
			var extAttributionURL = $(".form-impl-creation input[name=impl_attribution_url]").val();
			if( extAttributionURL )
				m.find("a.impl-attribution").attr("href", extAttributionURL).show();
			else
				m.find("a.impl-attribution").attr("href", "#").hide();
			m.modal();
			window.setTimeout(function(){
				$('pre[data-content]').popover("show"); // Fix and show (0,0) popovers of hidden tabs
			}, 800);
	}

	$(".btn-impl-create-preview").on("click", function(){
		showImplCreatePreview();
		return false;
	})

	function showImplEditPreview(){
			$('pre[data-content]').popover("hide"); // Hide (0,0) popovers of hidden tabs
			var m = $('.modal-impl-preview');
			var lang = $(".form-impl .badge").html();
			m.find(".lang-tab span.label").html(lang);
			var imports = $(".form-impl textarea.imports").val();
			if( imports )
				m.find(".piimports pre > code").text( imports ).show();
			else
				m.find(".piimports pre > code").text( imports ).hide();
			m.find(".picode pre > code").text( $(".form-impl textarea.impl-code").val() );
			var comment = $(".form-impl textarea[name=impl_comment]").val();
			var escapedComment = $("<div>").text(comment).html();
			var refinedComment = emphasize(escapedComment);
			m.find(".picode pre").attr("data-content", refinedComment);
			var extDocURL = $(".form-impl input[name=impl_doc_url]").val();
			if( extDocURL )
				m.find("a.impl-doc").attr("href", extDocURL).show();
			else
				m.find("a.impl-doc").attr("href", "#").hide();
			var extDemoURL = $(".form-impl input[name=impl_demo_url]").val();
			if( extDemoURL )
				m.find("a.impl-demo").attr("href", extDemoURL).show();
			else
				m.find("a.impl-demo").attr("href", "#").hide();
			var extAttributionURL = $(".form-impl input[name=impl_attribution_url]").val();
			if( extAttributionURL )
				m.find("a.impl-attribution").attr("href", extAttributionURL).show();
			else
				m.find("a.impl-attribution").attr("href", "#").hide();
			m.modal();
			window.setTimeout(function(){
				$('pre[data-content]').popover("show"); // Fix and show (0,0) popovers of hidden tabs
			}, 800);
	}

	$(".btn-impl-edit-preview").on("click", function(){
		showImplEditPreview();
		return false;
	});

	//
	// Idiom create : [Preview] button injects values
	// in modal window.
	//

	function showIdiomCreatePreview(){
			$('pre[data-content]').popover("hide"); // Hide (0,0) popovers of hidden tabs
			var m = $('.modal-idiom-preview');

			var title = $(".form-idiom-creation input[name=idiom_title]").val();
			m.find(".idiom-title").html(title);
			var lead = $(".form-idiom-creation textarea[name=idiom_lead]").val();
			var escapedLead = $("<div>").text(lead).html();
			var refinedLead = emphasize(escapedLead);
			m.find(".idiom-lead-paragraph").html(refinedLead);

			var lang = $(".form-idiom-creation input[name=impl_language]").val();
			m.find(".lang-tab span.label").html(lang);
			var imports = $(".form-idiom-creation textarea.imports").val();
			if( imports )
				m.find(".piimports pre > code").text( imports ).show();
			else
				m.find(".piimports pre > code").text( imports ).hide();
			m.find(".picode pre > code").text( $(".form-idiom-creation textarea.impl-code").val() );
			var comment = $(".form-idiom-creation textarea[name=impl_comment]").val();
			var escapedComment = $("<div>").text(comment).html();
			var refinedComment = emphasize(escapedComment);
			m.find(".picode pre").attr("data-content", refinedComment);
			var extDocURL = $(".form-idiom-creation input[name=impl_doc_url]").val();
			if( extDocURL )
				m.find("a.impl-doc").attr("href", extDocURL).show();
			else
				m.find("a.impl-doc").attr("href", "#").hide();
			var extDemoURL = $(".form-idiom-creation input[name=impl_demo_url]").val();
			if( extDemoURL )
				m.find("a.impl-demo").attr("href", extDemoURL).show();
			else
				m.find("a.impl-demo").attr("href", "#").hide();
			var extAttributionURL = $(".form-idiom-creation input[name=impl_attribution_url]").val();
			if( extAttributionURL )
				m.find("a.impl-attribution").attr("href", extAttributionURL).show();
			else
				m.find("a.impl-attribution").attr("href", "#").hide();
			m.modal();
			window.setTimeout(function(){
				$('pre[data-content]').popover("show"); // Fix and show (0,0) popovers of hidden tabs
			}, 800);
	}

	$(".btn-idiom-create-preview").on("click", function(){
		showIdiomCreatePreview();
		return false;
	})

	//
	// Messages sent from admin to user
	//
	setTimeout(function(){
		if( logged() ){
			$.get('/ajax-user-message-box', 
				function(response) {
					var messages = response.messages;
					if(messages.length > 0){
						var zone = $(".user-messages");
						messages.forEach(function(message) {
							var item = $("<div>").addClass("user-message alert");
							item.append( $("<div>").addClass("dismissal").html( $("<button>")
								.attr("type", "button")
								.addClass("close")
								.html("&times; dismiss")
								.attr("key", message.key)
								.on("click", function(event){
									console.log("Dismissing " + $(this).attr("key"));
									var hideBtn = $(this);
								    $.ajax({
								        url: "/ajax-dismiss-user-message",
								        type: 'POST',
								        data: {key:hideBtn.attr("key")},
								        xhr: function() {
								            var myXhr = $.ajaxSettings.xhr();
								            return myXhr;
								        },
								        success: function(response){
								        	console.log("User message dismissed.");
								        	hideBtn.closest(".user-message").hide("fast");
								        },
								        error: function(xhr, status, e){
								        	$.fn.pierror( xhr.responseText );
								        },
	   								});
								})
							) );
							item.append( $("<h4>").html("Message for " + username()) );
							item.append( $("<div>").addClass("date").append($("<small>").html(YYYYMMDDHHmm(message.creationDate))) );
							var content = message.message.replace(/\n/g, "<br/>");
							item.append( $("<div>").addClass("content").html(content) );
							zone.append(item);
						});
						zone.show("fast");
					}
			});
		}
	},300);

	function preload(arrayOfImages) {
    	$(arrayOfImages).each(function(){
        	$('<img/>')[0].src = this;
    	});
	}

	// Prefetch the highlight version of the
	// dice icon "Go To Random Idiom".
	$("img.dice").each(function(){
		var src = $(this).attr('src');
		if(src.indexOf('_highlight') !== -1)
			return;
		var srcHighlight = src.replace('.png', '_highlight.png').replace('.svg', '_highlight.svg');
		preload([
			srcHighlight
		]);
	});

	//
	// Cheatsheet language select screen
	//
	$(document).on('submit', 'form.two-languages-select', function(){
		var lang1 = $(this).find("input[name=lang1]").val();
		var lang2 = $(this).find("input[name=lang2]").val();
		lang1=normLang(lang1);
		lang2=normLang(lang2);
		window.location.href = "/cheatsheet/" + lang1 + "/" + lang2;
		return false;
	});

	// Restore a previous version of the Idiom
	// (only the admin can do this)
	$("form.idiom-restore-version > input.presubmit").on("click", function(e) {
		var form = $(this).closest("form");
		var whyPlaceholder = "";
		if( form.find("input[name=why]").val() ){
			whyPlaceholder = form.find("input[name=why]").val();
		}
		var reason = window.prompt("Why?", whyPlaceholder);
		if( reason===null )
	   		return; // Clicked Cancel
		form.find("input[name=why]").val(reason);
		form.submit();
	});

	function generateQAStructuredData(idiom) {
		let lang = extractLangFromURL();
		let questionName = document.title;
		let questionText = $(".idiom-lead-paragraph").text();
		if(!questionText) {
			// In the rare case where the description (lead paragraph) is empty,
			// let's stuff "text" with the question title, because "text" is mandatory.
			questionText = questionName;
		}
		let questionCreation = idiom.CreationDate.substring(0,10);
		let questionAuthor = idiom.Author;
		// let questionText = idiom.LeadParagraph;
		let answers = [];
		idiom.Implementations.forEach(function(impl){
			if(impl.LanguageName.toLowerCase() != lang.toLowerCase())
				return;
			let parts = window.location.href.split(/\//);
			let n = parts.length;
			parts[n-2] = impl.Id;
			let implAbsoluteURL = parts.join("/");
			let answer = {
				"@type": "Answer",
				"text": impl.CodeBlock,
				"dateCreated": impl.CreationDate.substring(0,10),
				"url": implAbsoluteURL,
				"author": {
					"@type": "Person",
					"name": impl.Author
				}
			};
			 answers.push(answer);
		});

		const script = document.createElement('script');
		script.setAttribute('type', 'application/ld+json');
		let sd = {
			"@context": "https://schema.org",
			"@type": "QAPage",
			"mainEntity": {
				"@type": "Question",
				"name": questionName,
				"text": questionText,
				"answerCount": answers.length,
				"dateCreated": questionCreation,
				"author": {
                    "@type": "Person",
                    "name": questionAuthor
                },
				"license": "https://en.wikipedia.org/wiki/Wikipedia:Text_of_Creative_Commons_Attribution-ShareAlike_3.0_Unported_License",
				"suggestedAnswer": answers
			}
		};
		script.textContent = JSON.stringify(sd, null, "  ");
		document.head.appendChild(script);
		// console.debug("Just injected", script.textContent);
		console.debug("Generated Q&A structured data");
	}
	
	if( isIdiomDetailWithLang() ) {
		let idiomID = $(".idiom-summary-large").attr("data-idiom-id"); // Or extract from the URL?
		fetch('/api/idiom/' + idiomID)
		.then(function(response){ return response.json()})
		.then(function(idiom) {
			generateQAStructuredData(idiom);
		});
	}

	$(".dice").on("click", function() {
		$(this).addClass("spinning-jolty");
	});

	// #170 horizontal bar for an empty list of tabs is nor great.
	// CSS doesn't have a selector for "empty but containing whitespace" so
	// instead we hide the list after page load...
	$('.implementations-tabs ul.nav.nav-tabs.language-names').each(function(){
		let ul = $(this);
		if(!/[\S]/.test(ul.html())) {
			console.debug("Hiding empty nav-tabs");
			ul.hide();
		}
	});

	function showKeyboardShortcuts() {
		$("<div>").addClass("modal").attr("tabindex","-1")
			.addClass("keyboard-shortcuts")
			.append(`<div class="modal-header">
						<button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>
						<h3>Keyboard shortcuts <span class="keyboard-icon">⌨</span></h3>
					</div>`)
			.append(`<div class="modal-body">
				<dl class="dl-horizontal">
					<dt>r</dt> <dd>Go to a random idiom</dd>
					<dt>/</dt> <dd>Search</dd>
					<dt>p</dt> <dd>Go to "previous" idiom</dd>
					<dt>n</dt> <dd>Go to "next" idiom</dd>
					<dt>e</dt> <dd>Edit current implementation</dd>
					<dt>c</dt> <dd>Copy snippet code to the clipboard</dd>
					<dt>d</dt> <dd>Follow demo link</dd>
					<dt>i</dt> <dd>Follow documentation link</dd>
					<dt>x</dt> <dd>Create a new implementation</dd>
					<dt>g</dt> <dd>Grid view</dd>
					<dt>h</dt> <dd>History</dd>
					<dt>z</dt> <dd>Cheatsheet</dd>
					<dt>?</dt> <dd>Show keyboard shortcuts</dd>
				</dl>
			</div>`)
			.modal("show");
	}

	if( isIdiomPage() ) {
		let idiomID = $(".idiom-summary-large").attr("data-idiom-id");
		// Conventionally, the "current impl" is the top impl in the page.
		// This does not capture all desired use cases, but probably the
		// most frequent one.
		let impl = $("div.implementation").first();
		let implID = impl.attr("data-impl-id");
		let implLang = impl.attr("data-impl-lang");
		let implCopyButton = impl.find("a.copy-code-to-clipboard");

		// #192 Keyboard shortcuts
		$(document).on("keydown", function(e) {
			if ( e.target.tagName.toLowerCase() === 'input' ||
			     e.target.tagName.toLowerCase() === 'textarea' ) { 
				// Do not mess with the search text box
				return;
			}
			if ( e.ctrlKey || e.altKey || e.metaKey ) {
				// Do not mess with popular shortcuts like Ctrl+R, etc.
				return;
			}

			switch(e.key) {
				case '?':
					using(`keyboard/show-keyboard-shortcuts`);
					showKeyboardShortcuts();
					break;
				case 'r':
					using("keyboard/random-idiom");
					$(".dice").addClass("spinning-jolty");
					window.location = "/random-idiom";
					break;
				case '/':
					using("keyboard/search");
					$("form.form-search > input[name=q]").first().focus();
					e.preventDefault();
					break;
				case 'p':
					using("keyboard/previous-idiom");
					window.location = `/previous-idiom/${idiomID}`;
					break;
				case 'n':
					using("keyboard/next-idiom");
					window.location = `/next-idiom/${idiomID}`;
					break;
				case 'e':
					using(`keyboard/impl-edit/${idiomID}/${implID}`);
					window.location = `/impl-edit/${idiomID}/${implID}`;
					break;
				case 'c':
					using(`keyboard/impl/copy-to-clipboard/${idiomID}/${implID}/${implLang}`);
					clipboardCopyCodeOfImpl(implCopyButton, impl);
					break;
				case 'd':
					let demoUrl = impl.find("li.demo > a").attr("href");
					if (demoUrl) {
						using(`keyboard/impl/demo/${idiomID}/${implID}/${implLang}`);
						window.open(demoUrl, "_blank");
					}
					break;
				case 'i':
					let docUrl = impl.find("li.doc > a").attr("href");
					if (docUrl) {
						using(`keyboard/impl/doc/${idiomID}/${implID}/${implLang}`);
						window.open(docUrl, "_blank");
					}
					break;
				case 'x':
					using(`keyboard/impl-create/${idiomID}`);
					window.location = `/impl-create/${idiomID}`;
					break;
				case 'g':
					using(`keyboard/grid/${idiomID}`);
					$(".modal-impl-grid").modal();
					break;
				case 'h':
					using(`keyboard/history/${idiomID}`);
					window.location = `/history/${idiomID}`;
					break;
				case 'z':
					if(implLang) {
						using(`keyboard/cheatsheet/${implLang}`);
						window.open(`/cheatsheet/${implLang}`, "_blank");
					}
					break;
			}
		});

		// Special tooltips for the very first impl
		let copyIcon = implCopyButton.find("i");
		copyIcon.attr("title", copyIcon.attr("title") + " (c)");
		let editBtn = impl.parent().find("a.btn-edit");
		editBtn.attr("title", editBtn.attr("title") + " (e)");
		let demoLink = impl.find("li.demo > a");
		if(demoLink)
			demoLink.attr("title", "Online demo (d)");
		let docLink = impl.find("li.doc > a");
		if(docLink)
			docLink.attr("title", "Documentation page (i)");
	}

	$("a.show-keyboard-shortcuts").click(function(){
		using(`/right/show-keyboard-shortcuts`);
		showKeyboardShortcuts();
	});

	// #205: auto-summary when adding a demo URL
	$("form.form-impl.form-edit input[name=impl_doc_url]").change(function() {
		let docField = $(this);
		let summaryField = $("form.form-impl.form-edit input[name=edit_summary]");
		if( docField.val() && !summaryField.val() )
			summaryField.val("Added Doc URL");
	});
	// #206: auto-summary when adding a demo URL
	$("form.form-impl.form-edit input[name=impl_demo_url]").change(function() {
		let demoField = $(this);
		let summaryField = $("form.form-impl.form-edit input[name=edit_summary]");
		if( demoField.val() && !summaryField.val() )
			summaryField.val("Added Demo URL");
	});
});

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

function niceLang(lang){
	switch(lang){
	case "Cpp":
		return "C++";
	case "Csharp":
		return "C#";
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
