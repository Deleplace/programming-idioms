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

self.addEventListener('fetch', (event) => {
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
    }).catch(function(error) {
        return handleOffline(event.request);
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

// Some feature are specifically designed to have a cromulent
// offline behavior.
function handleOffline(request) {

  if (request.url.indexOf("/api/random-id") != -1) {
    // Instead of asking the server to pick in the real DB,
    // pick in the local cached DB.
    console.log("Handling offline: " + request.url);
    return caches.open('v1').then(function(cache) {
      return cache.match("/api/idioms/all").then(function(resp) {
        //debugger;
        if (resp) {
          // resp is a promise
          return resp.json().then(body => {
            var fullDB = body;
            //console.log("  Offline db: " + JSON.stringify(full));
            console.log("  Offline db: " + fullDB.length + " idioms.");
            var chosen = pick(fullDB);
            console.log("  Chosen: " + JSON.stringify(chosen));
            return chosen.Id;
          });
        }else{
          console.log( "  No offline db :(");
          return;
        }
      }).catch(function(error){
        console.log("Cache error for /api/idioms/all : " + error);
        return caches.match( '/default_' + ThemeDate + '/img/dice_48x48.png' );
      });
    });
  }

  if (request.url.indexOf("/api/idiom/") != -1) {
    // Pick in the local cached DB.
    console.log("Not implemented yet: local retrieval of " + request.url);
  }

  console.log("Not found " + event.request.url + " :(");
  return caches.match( '/default_' + ThemeDate + '/img/dice_48x48.png' );
}

// Pick a random element in a list.
function pick(x) {
  return x[Math.floor(Math.random() * x.length)];
}