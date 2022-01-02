$(function() {

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
    
});