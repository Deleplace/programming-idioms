function elem(tag, clazz, html) {
    // TODO: exists something standard for this?
    var element = document.createElement(tag);
    if(clazz)
        element.className = clazz;
    if(html)
        element.innerHTML = html;
    return element;
}

function elemText(tag, clazz, text) {
    // TODO: exists something standard for this?
    var element = document.createElement(tag);
    if(clazz)
        element.className = clazz;
    if(text)
        element.textContent = text;
    return element;
}

function emphasize(raw){
    // Emphasize the "underscored" identifier
    //
    // _x -> <span class="variable">x</span>
    //
    var refined = raw.replace( /\b_([\w$]*)/gm, "<span class=\"variable\">$1</span>");
    refined = refined.replace(/\n/g,"<br/>");
    return refined;
}

function renderImpl(impl) {
    var implNode = elem("div", "implementation");
    implNode.id = "impl-" + impl.Id;
    var lg = elem("h2", "lang", impl.LanguageName);
    implNode.appendChild(lg);

    var importsAndCode = elem("div", "imports-and-code");
    if(impl.ImportsBlock){
        var imports = elem("div", "imports");
        var pre = elemText("pre", "", impl.ImportsBlock);
        imports.appendChild(pre);
        importsAndCode.appendChild(imports);
    }
    var code = elem("div", "code");
    var pre = elemText("pre", "", impl.CodeBlock);
    code.appendChild(pre);
    importsAndCode.appendChild(code);
    implNode.appendChild(importsAndCode);

    var comment = elem("div", "comment", emphasize(impl.AuthorComment));
    implNode.appendChild(comment);

    var links = elem("div", "external-links");
    var ul = elem("ul");
    if(impl.DemoURL) {
        var li = elem("li");
        var a = elem("a", "", "Demo 🗗");
        a.href = impl.DemoURL;
        a.target="_blank";
        a.rel="nofollow";
        li.appendChild(a);
        ul.appendChild(li);
    }
    if(impl.DocumentationURL) {
        var li = elem("li");
        var a = elem("a", "", "Doc 🗗");
        a.href = impl.DocumentationURL;
        a.target="_blank";
        a.rel="nofollow";
        li.appendChild(a);
        ul.appendChild(li);
    }
    if(impl.OriginalAttributionURL) {
        var li = elem("li");
        var a = elem("a", "", "Origin 🗗");
        a.href = impl.OriginalAttributionURL;
        a.target="_blank";
        a.rel="nofollow";
        li.appendChild(a);
        ul.appendChild(li);
    }
    links.appendChild(ul);
    implNode.appendChild(links);

    var impls = document.querySelector(".implementations");
    impls.appendChild(implNode);
    // console.log( implNode.id + " added!" );
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
    // TODO /default_20180923_/... ?
    hadd('<a href="/"><img src="/default/img/wheel_48x48.png" width="48" height="48" class="header_picto" /></a>');
    hadd('<h1><a href="/">Programming-Idioms</a></h1>');
    hadd('<a href="/random-idiom"><img src="/default/img/dice_32x32.png" width="32" height="32" class="picto die" title="Go to a random idiom" /></a>');
    hadd('<form class="form-search" action="/search"> \
            <input type="text" class="search-query" placeholder="Keywords..." name="q" value="" required="required"> \
            <button type="submit">Search</button> \
          </form>');
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
    // 1) Remove all "..." placeholders
    var placeholders = document.querySelectorAll(".implementation.placeholder");
    for (var ph of placeholders) {
        ph.remove();
    }

    // 2) Add each impl (if it's not there yet)
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

function highlightDie(){
    var die = document.querySelector(".die");
    var src = die.src;
    var hsrc = src.replace("dice_32x32.png", "dice_32x32_highlight.png");
    die.onmouseover=function(){this.src=hsrc;};
    die.onmouseout=function(){this.src=src;};
}

function setVisitCookie() {
    // This is not personal data, just a hint to decide
    // if resources (JS, CSS, img) should be server-pushed.
    document.cookie = "v=1;path=/; ";
}

//
// Execution!
//

renderHeader();

if(idiomPromise){
    idiomPromise
        .then(function(response) {
            // console.log("Got response");
            return response.json();
        })
        .then(function(idiom) {
            console.log("Got JSON of idiom " + idiom.Id);
            populateOtherImpls(idiom);
            decorateImpls(idiom);
            decorateSummary(idiom);
        });
}

renderFooter();
highlightDie();
setVisitCookie();