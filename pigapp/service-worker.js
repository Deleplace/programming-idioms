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

// This may be abusing response.clone() ... not sure about that.
self.addEventListener('fetch', function(event) {
    event.respondWith(
      caches.match(event.request).then(function(resp) {
        if (resp) {
            console.log("Found in cache: " + event.request.url);
            return resp;
        }
        return fetch(event.request).then(function(response) {
          console.log("Fetched " + event.request.url);
          // Cache 200-299, but not 302, 404, 500...
          if(response.ok) {
            caches.open('v1').then(function(cache) {
                console.log("Caching response for " + event.request.url);
                cache.put(event.request, response.clone());
            });
          }
          return response.clone();
      });
    }).catch(function() {
        console.log("Not found " + event.request.url + " :(");
        return caches.match( '/default_' + ThemeDate + '/img/dice_48x48.png' );
      })
    );
  });