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
    hadd('<a href="/"><img src="/default_20200205_/img/wheel.svg" width="48" height="48" class="header_picto" alt="Logo" /></a>');
    hadd('<h1><a href="/">Programming-Idioms</a></h1>');

    let nick = getCookie("Nickname");
    if(nick)
        hadd(`<span class="greeting">${nick}</span>`);
    // TODO user can click/view/remove profile cookie data
}

function renderFooter() {
    var footerz = document.getElementsByTagName("footer");
    var footer = footerz[0];
    footer.insertAdjacentHTML('beforeend', `
    <div>
		All content <a href="https://en.wikipedia.org/wiki/Wikipedia:Text_of_Creative_Commons_Attribution-ShareAlike_3.0_Unported_License" rel="license noopener">CC-BY-SA</a>
    </div>
    <div>
		<a href="/about#about-block-language-coverage" title="Coverage grid"><img src="/default_20200205_/img/coverage_icon_indexed.png" class="coverage square-loading" alt="Coverage grid" /></a>
	</div>
    <div>
		<a href="/about" class="about-link">?</a>
	</div>`);
}

function getCookie(name) {
    // From https://gomakethings.com/working-with-cookies-in-vanilla-js/
    var value = "; " + document.cookie;
    var parts = value.split("; " + name + "=");
    if (parts.length == 2) return parts.pop().split(";").shift();
}

//
// Execution!
//

renderHeader();
renderFooter();