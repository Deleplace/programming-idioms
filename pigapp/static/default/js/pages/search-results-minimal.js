function renderHeader() {
    var hh = document.getElementsByTagName("header");
    if(!hh.length) {
        console.error("Couldn't find header element");
        return;
    }
    var h = hh[0];
    while(h.firstChild) {
        h.removeChild(h.firstChild);
    }

    var hadd = function(code){
        h.insertAdjacentHTML('beforeend', code);
    }
    hadd('<a href="/"><img src="/default_20171211_/img/wheel_48x48.png" width="48" height="48" class="header_picto" /></a>');
    hadd('<h1><a href="/">Programming-Idioms</a></h1>');
    hadd('<a href="/random-idiom"><img src="/default_20171211_/img/dice_32x32.png" width="32" height="32" class="picto die" title="Go to a random idiom" /></a>');
    hadd('<form class="form-search" action="/search"> \
            <input type="text" class="search-query" placeholder="Keywords..." name="q" value="" required="required"> \
            <button type="submit">Search</button> \
          </form>');
}

function populateSearchQuery() {
    var nodes = document.getElementsByClassName("results-idioms");
    if(!nodes || !nodes.item(0))
        return;
    var q = nodes.item(0).getAttribute("data-search-query");
    nodes = document.getElementsByClassName("search-query");
    if(!nodes || !nodes.item(0))
        return;
    nodes.item(0).setAttribute("value", q);
}

function renderFooter() {
    // TODO
}


//
// Execution!
//

renderHeader();

populateSearchQuery();

renderFooter();