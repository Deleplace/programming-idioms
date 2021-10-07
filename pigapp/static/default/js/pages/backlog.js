$(".backlog .impl-actions button.view").click(function(){
    let actions = $(this).closest(".impl-actions");
    let idiomID = actions.attr("data-idiom-id");
    let implID = actions.attr("data-impl-id");
    let pageURL = `/idiom/${idiomID}/impl/${implID}`;
    window.open(pageURL);
}); 

$(".backlog .impl-actions button.edit").click(function(){
    let actions = $(this).closest(".impl-actions");
    let idiomID = actions.attr("data-idiom-id");
    let implID = actions.attr("data-impl-id");
    let pageURL = `/impl-edit/${idiomID}/${implID}`;
    window.open(pageURL);
}); 

$(".backlog .impl-actions button.flag").click( () => {
    console.log("TODO");
}); 

$(".backlog .impl-actions button.mark-good").click( () => {
    console.log("TODO");
}); 

$(".backlog .idiom-actions button.create-impl").click(function(){
    let actions = $(this).closest(".idiom-actions");
    let idiomID = actions.attr("data-idiom-id");
    let lang = actions.attr("data-missing-lang");
    let pageURL = `/impl-create/${idiomID}/${lang}`;
    window.open(pageURL);
}); 

$(".backlog .idiom-actions button.view").click(function(){
    let actions = $(this).closest(".idiom-actions");
    let idiomID = actions.attr("data-idiom-id");
    let pageURL = `/idiom/${idiomID}`;
    window.open(pageURL);
}); 