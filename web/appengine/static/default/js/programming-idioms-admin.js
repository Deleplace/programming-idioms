/*
 * Admin only. Do not include pages when the user is not admin.
 */

$(function() {

	$('#import-form input.upload').on("click", function(){
		// See http://stackoverflow.com/questions/166221/how-can-i-upload-files-asynchronously-with-jquery#answer-8758614
	    var formData = new FormData($('#import-form')[0]);
	    if( window.location.href.indexOf("localhost") === -1 ){
	    	var expectedSafeWord = "prod";
	    	var confirm = prompt("This is not your localhost. Please enter safeword. The safeword is \"" + expectedSafeWord + "\".", "");
	    	if( confirm != expectedSafeWord ){
	    		alert("Aborting!");
	    		return false;
	    	}
	    }
	    $.ajax({
	        url: '/admin-data-import-ajax',
	        type: 'POST',
	        xhr: function() {
	            var myXhr = $.ajaxSettings.xhr();
	            return myXhr;
	        },
	        success: function(response){
				//var count = response.imported;  ???
	        	var count = response.imported
	        	$.fn.pisuccess( count + " idioms imported.");
	        },
	        error: function(xhr, status, e){
	        	$.fn.pierror( "Import failed : " + xhr.responseText);
	        },
	        data: formData,
	        cache: false,
	        contentType: false,
	        processData: false
	    });
	});
	
	$('#refresh-toggles').on("click", function(){
	    $.ajax({
	        url: '/admin-refresh-toggles-ajax',
	        type: 'POST',
	        xhr: function() {
	            var myXhr = $.ajaxSettings.xhr();
	            return myXhr;
	        },
	        success: function(response){
	        	$.fn.pisuccess( "Toggles refreshed" );
	        },
	        error: function(xhr, status, e){
	        	$.fn.pierror( "Refresh toggles failed : " + xhr.responseText);
	        },
	    });
	});
		
	$('.toggles-list input[type=checkbox]').on("click", function(){
		var toggleName = $(this).attr('data-toggle-name');
		var	newValue = $(this).is(':checked');
		console.log(`Setting toggle ${toggleName} to ${newValue}`)

	    $.ajax({
	        url: '/admin-set-toggle-ajax',
	        type: 'POST',
	        xhr: function() {
	            var myXhr = $.ajaxSettings.xhr();
	            return myXhr;
	        },
	        success: function(response){
	        	$.fn.pisuccess( `Set toggle ${toggleName} to ${newValue}` );
	        },
	        error: function(xhr, status, e){
	        	$.fn.pierror( `Setting toggle ${toggleName} to ${newValue} failed : ${xhr.responseText}`);
	        },
	        data: {
	        	toggle: toggleName,
	        	value: newValue
	        },
	    });
	});

	$('#relation-form input.create-relation').on("click", function(){
		var idA = $("#relation-form input.idiomA").val();
		var idB = $("#relation-form input.idiomB").val();
	    $.ajax({
	        url: '/admin-create-relation-ajax',
	        type: 'POST',
	        xhr: function() {
	            var myXhr = $.ajaxSettings.xhr();
	            return myXhr;
	        },
	        success: function(response){
	        	$.fn.pisuccess( "Created relation between idioms [" + idA + ", " + idB + "]" );
	        },
	        error: function(xhr, status, e){
	        	$.fn.pierror( "Relation between idioms [" + idA + ", " + idB + "] failed : " + xhr.responseText );
	        },
	        data: {
	        	idiomAId: idA,
	        	idiomBId: idB
	        }
	    });
	});
	
	$('#reindex-form input.submit').on("click", function(){
	    $.ajax({
	        url: '/admin-reindex-ajax',
	        type: 'POST',
	        xhr: function() {
	            var myXhr = $.ajaxSettings.xhr();
	            return myXhr;
	        },
	        success: function(response){
	        	$.fn.pisuccess( response.message );
	        },
	        error: function(xhr, status, e){
	        	$.fn.pierror( "Reindex failed : " + xhr.responseText);
	        },
	        cache: false,
	        contentType: false,
	        processData: false
	    });
	});

	$('#repair-history-form input.submit').on("click", function(){
		var id = $("#repair-history-form input.idiom").val();
	    $.ajax({
	        url: '/admin-repair-history-versions',
	        type: 'POST',
	        xhr: function() {
	            var myXhr = $.ajaxSettings.xhr();
	            return myXhr;
	        },
	        success: function(response){
	        	$.fn.pisuccess( response.message );
	        },
	        error: function(xhr, status, e){
	        	$.fn.pierror( "History repair failed : " + xhr.responseText);
	        },
	        data: {
	        	idiomId: id,
	        },
	        cache: false
	    });
	});

	$('#message-for-user-form .btn.send-message-for-user').on("click", function(){
	    $.ajax({
	        url: '/admin-send-message-for-user',
	        type: 'POST',
	        xhr: function() {
	            var myXhr = $.ajaxSettings.xhr();
	            return myXhr;
	        },
	        success: function(response){
	        	$.fn.pisuccess( "Message has been sent." );
	        	$('#message-for-user-form textarea[name=message]').val('');
	        },
	        error: function(xhr, status, e){
	        	$.fn.pierror( "Message sending failed : " + xhr.responseText );
	        },
	        data: $("#message-for-user-form").serialize(),
	    });
	});

	// Flagged Contents page

	$('button.flag-mark-resolved').on("click", function(){
		let btn = $(this);
		let flagKey = btn.attr('flagkey');
		if(!flagKey) {
			console.error("no flagkey??");
			return;
		}
	    $.ajax({
	        url: '/admin-flag-resolve',
	        type: 'POST',
	        xhr: function() {
	            var myXhr = $.ajaxSettings.xhr();
	            return myXhr;
	        },
	        success: function(response){
	        	// $.fn.pisuccess( "Flag content marked resolved." );
				let tr1 = btn.closest("tr");
				tr1.addClass("resolved");
				let tr2 = tr1.next();
				if(tr2) {
					tr2.addClass("resolved");
				}
				btn.closest("td").text("✓");
	        },
	        error: function(xhr, status, e){
	        	$.fn.pierror( "Flag resolve failed : " + xhr.responseText );
	        },
	        data: {flagkey: flagKey}
	    });
	});


	$('#memcache-flush-form input.submit').on("click", function(){
	    $.ajax({
	        url: '/admin-memcache-flush',
	        type: 'POST',
	        xhr: function() {
	            var myXhr = $.ajaxSettings.xhr();
	            return myXhr;
	        },
	        success: function(response){
				$.fn.pisuccess( response.message );
	        },
	        error: function(xhr, status, e){
				$.fn.pierror( "Memcache flush failed : " + xhr.responseText);
	        },
	        cache: false
	    });
	});

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
			case 'a':
				window.open("/admin", "_blank");
				break;
		}
	});
});
