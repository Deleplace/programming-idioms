// This depends on jQuery

$('.single-impl-tab h3, .single-impl-tab .fold-triangle').on('click',function() {
    $(this).parent().toggleClass("folded");
});