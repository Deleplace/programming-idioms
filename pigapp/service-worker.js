var ThemeDate = "20171211_";

self.addEventListener('install', function(event) {
    event.waitUntil(
      caches.open('v1').then(function(cache) {
        return cache.addAll([
          '/',
          '/default_' + ThemeDate + '/css/bootstrap-combined.no-icons.min.css',
          '/default_' + ThemeDate + '/css/font-awesome/css/font-awesome.css',
          '/default_' + ThemeDate + '/css/prettify.css',
          '/default_' + ThemeDate + '/css/programming-idioms.css',
          '/default_' + ThemeDate + '/img/dice_48x48.png',
        ]);
      })
    );
});

self.addEventListener('fetch', function(event) {
    event.respondWith(
      caches.match(event.request).then(function(resp) {
        if (resp) {
            console.log("Found in cache: " + event.request.url);
            return resp;
        }
        return fetch(event.request).then(function(response) {
          console.log("Fetched " + event.request.url);
          if (cachable(event.request, response)) {
            let responseClone = response.clone();
            caches.open('v1').then(function(cache) {
                console.log("Caching response for " + event.request.url);
                cache.put(event.request, responseClone);
            });
          }
          return response;
      });
    }).catch(function() {
        console.log("Not found " + event.request.url + " :(");
        return caches.match( '/default_' + ThemeDate + '/img/dice_48x48.png' );
      })
    );
  });

function cachable(request, response) {
  if (!response.ok) {
    // Cache 200-299, but not 302, 404, 500...
    return false;
  }

  // Random endpoints are specifically to **not** cache
  if (request.url.indexOf("/api/random-id") != -1)
    return false;
  if (request.url.indexOf("/random-idiom") != -1)
    return false;

  // JS, CSS, PNG.
  // TODO: handle updates...
  if (request.url.indexOf("/default/") != -1)
    return true;

  // Limit to SPA pages & services, for now
  if (request.url.indexOf("/page/") != -1)
    return true;
  if (request.url.indexOf("/api/") != -1)
    return true;

  return false;
}