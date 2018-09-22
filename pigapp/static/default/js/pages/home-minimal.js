function elem(tag, clazz, html) {
    // TODO: exists something standard for this?
    var element = document.createElement(tag);
    if(clazz)
        element.className = clazz;
    if(html)
        element.innerHTML = html;
    return element;
}

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
/*
    hadd('<a href="/random-idiom"><img src="/default_20171211_/img/dice_32x32.png" width="32" height="32" class="picto die" title="Go to a random idiom" /></a>');
    hadd('<form class="form-search" action="/search"> \
            <input type="text" class="search-query" placeholder="Keywords..." name="q" value="" required="required"> \
            <button type="submit">Search</button> \
          </form>');
*/
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

// Server-side rendering already includes the HTML for only
// a few impls.
// populateOtherImpls does client-side rendering of all other
// impls.
function populateOtherImpls(idiom) {
    idiom.Implementations.forEach(function(impl) {
        var nodeId = "impl-" + impl.Id;
        if( document.getElementById(nodeId) ){
            // console.log("Skipping existing " + nodeId);
        }else{
            renderImpl(impl);
        }
    });
}

function decorateImpls(idiom) {
    idiom.Implementations.forEach(function(impl) {
        var nodeId = "impl-" + impl.Id;
        var implNode = document.getElementById(nodeId);
        if( !implNode ){
            console.error("Couldn't find " + nodeId);
            return;
        }
        implNode.insertAdjacentHTML('beforeend',
            '<a href="/impl-edit/' 
            + idiom.Id + '/' 
            + impl.Id + '" class="edit hide-on-mobile" title="Edit this implementation">Edit</a>');
    });
}

function decorateSummary(idiom) {
    var nodes = document.getElementsByClassName("summary-large");
    if (!nodes || !nodes.item(0))
        return;
    nodes.item(0).insertAdjacentHTML('beforeend', 
        '<a href="/idiom-edit/' + 
        idiom.Id + 
        '" title="Edit the idiom statement" class="edit hide-on-mobile">Edit</a>');
}

//
// Execution!
//

renderHeader();

renderFooter();