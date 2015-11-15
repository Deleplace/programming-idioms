/*
 * Admin only. Do not include in front pages.
 */

$(function() {

	$('#import-form input.upload').on("click", function(){
		// See http://stackoverflow.com/questions/166221/how-can-i-upload-files-asynchronously-with-jquery#answer-8758614
	    var formData = new FormData($('#import-form')[0]);
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
	
	$('.toggles-list button').on("click", function(){
		var toggleName = $(this).html();
		// ?? not up-to-date yet
		var oldValue = $(this).hasClass("active");
		var	newValue = !oldValue;
		
	    $.ajax({
	        url: '/admin-set-toggle-ajax',
	        type: 'POST',
	        xhr: function() {
	            var myXhr = $.ajaxSettings.xhr();
	            return myXhr;
	        },
	        success: function(response){
	        	$.fn.pisuccess( "Set toggle " + toggleName + " to " + newValue );
	        },
	        error: function(xhr, status, e){
	        	$.fn.pierror( "Set toggle " + toggleName + " to " + newValue + " failed : " + xhr.responseText);
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
	        },
	    });
	});
	
	$('#reindex-form input.upload').on("click", function(){
		// See http://stackoverflow.com/questions/166221/how-can-i-upload-files-asynchronously-with-jquery#answer-8758614
	    $.ajax({
	        url: '/admin-reindex-ajax',
	        type: 'POST',
	        xhr: function() {
	            var myXhr = $.ajaxSettings.xhr();
	            return myXhr;
	        },
	        success: function(response){
				//var count = response.imported;  ???
	        	var count = response.indexed
	        	$.fn.pisuccess( count + " idioms reindexed.");
	        },
	        error: function(xhr, status, e){
	        	$.fn.pierror( "Reindex failed : " + xhr.responseText);
	        },
	        cache: false,
	        contentType: false,
	        processData: false
	    });
	});

	$('#message-for-user-form').on("click", function(){
	    $.ajax({
	        url: '/admin-send-message-for-user',
	        type: 'POST',
	        xhr: function() {
	            var myXhr = $.ajaxSettings.xhr();
	            return myXhr;
	        },
	        success: function(response){
	        	$.fn.pisuccess( "Message has been sent." );
	        },
	        error: function(xhr, status, e){
	        	$.fn.pierror( "Message sending failed : " + xhr.responseText );
	        },
	        data: $("#message-for-user-form").serialize(),
	    });
	});
});
