

if ('serviceWorker' in navigator) {
  // Register a service worker hosted at the root of the
  // site using the default scope.
  navigator.serviceWorker.register('/service-worker.js').then(function(registration) {
    console.log('Service worker registration succeeded:', registration);
  }).catch(function(error) {
    console.log('Service worker registration failed:', error);
  });
} else {
  console.log('Service workers are not supported.');
}

function getParameterByName(name, url) {
  if (!url) url = window.location.href;
  name = name.replace(/[\[\]]/g, "\\$&");
  var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
      results = regex.exec(url);
  if (!results) return null;
  if (!results[2]) return '';
  return decodeURIComponent(results[2].replace(/\+/g, " "));
}

// Client-side formatting.
function emphasize(raw){
  // Emphasize the "underscored" identifier
  //
  // _x -> <span class="variable">x</span>
  //
  if(!raw)
    return "";
  var refined = raw.replace( /\b_([\w$]*)/gm, "<span class=\"variable\">$1</span>");
  refined = refined.replace(/\n/g,"<br/>");
  return refined;
}

function uriNormalize(s) {
  if(!s)
    return "";
	s = s.trim();
	s = s.replace(/\[/g, "-");
	s = s.replace(/ /g, "-");
	s = s.replace(/\]/g, "-");
	s = s.replace(/,/g, "-");
	s = s.replace(/;/g, "-");
	s = s.replace(/--/g, "-");
  s = s.replace(/--/g, "-"); // Again
  s = s.replace(/[-\/ ]+$/, "");
  s = s.toLowerCase();
	return s
}