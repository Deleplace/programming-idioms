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

// Below:
// - works on Firefox
// - not in Chromium : Uncaught (in promise) TypeError: Failed to execute 'clone' on 'Response': Response body is already used
self.addEventListener('fetch', function(event) {
    event.respondWith(
      caches.match(event.request).then(function(resp) {
        return resp || fetch(event.request).then(function(response) {
          caches.open('v1').then(function(cache) {
            cache.put(event.request, response.clone());
          });
          return response;
        });
      }).catch(function() {
        return caches.match( '/default_' + ThemeDate + '/img/dice_48x48.png' );
      })
    );
  });