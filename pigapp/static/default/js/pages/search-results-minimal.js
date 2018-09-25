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
    hadd('<a href="/"><img src="/default/img/wheel_48x48.png" width="48" height="48" class="header_picto" /></a>');
    hadd('<h1><a href="/">Programming-Idioms</a></h1>');
    hadd('<a href="/random-idiom"><img src="/default/img/dice_32x32.png" width="32" height="32" class="picto die" title="Go to a random idiom" /></a>');
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
    var footerz = document.getElementsByTagName("footer");
    var footer = footerz[0];
    footer.insertAdjacentHTML('beforeend', '<div> \
		All content <a href="http://en.wikipedia.org/wiki/Wikipedia:Text_of_Creative_Commons_Attribution-ShareAlike_3.0_Unported_License" rel="license">CC-BY-SA</a> \
    </div>');
    footer.insertAdjacentHTML('beforeend', '<div> \
		<a href="/about" class="about-link">?</a> \
	</div>');
}

function highlightDie(){
    var die = document.querySelector(".die");
    var src = die.src;
    var hsrc = src.replace("dice_32x32.png", "dice_32x32_highlight.png");
    die.onmouseover=function(){this.src=hsrc;};
    die.onmouseout=function(){this.src=src;};
}

//
// Execution!
//

renderHeader();
populateSearchQuery();
renderFooter();
highlightDie();