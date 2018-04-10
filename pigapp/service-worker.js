var ThemeDate = "20171211_";

self.addEventListener('install', function(event) {
    event.waitUntil(
      caches.open('v1').then(function(cache) {
        return cache.addAll([
          '/',
          '/page/home.html',
          '/page/idiom-detail.html',
          '/page/search-results.html',
          '/default_' + ThemeDate + '/css/bootstrap-combined.no-icons.min.css',
          '/default_' + ThemeDate + '/css/font-awesome/css/font-awesome.css',
          '/default_' + ThemeDate + '/css/prettify.css',
          '/default_' + ThemeDate + '/css/programming-idioms.css',
          '/default_' + ThemeDate + '/img/dice_48x48.png',
          '/default_' + ThemeDate + '/img/disconnected.png',
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
        console.log("Catched " + error);
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
    return withLocalIdiomsDB( (db) => {
      var chosen = pick(db);
      console.log("  Chosen: #" + chosen.Id + " " + chosen.Title);
      return new Response(chosen.Id);
    });
  }

  if (request.url.indexOf("/api/idiom/") != -1) {
    // Pick in the local cached DB.
    var chunks = request.url.split("/");
    var idiomIdStr = chunks[chunks.length-1];
    console.log("Local retrieval of #" + idiomIdStr);
    return withLocalIdiomsDB( (db) => {
      for(var i = 0; i < db.length; i++) {
        var idiom = db[i];
        if (idiom.Id == idiomIdStr){
          console.log("Idiom #" + idiomIdStr + " found in Offline DB :)");
          return new Response( JSON.stringify(idiom) );
        }
      }
      console.log("Idiom #" + idiomIdStr + " not found in Offline DB of " + db.length + " idioms.");
      return;
    });
  }

  if (request.url.indexOf("/api/search/") != -1) {
    // Search in the local cached DB.
    // Warning: result matches and order may differ from the equivalent server search.
    var chunks = request.url.split("/");
    var q = chunks[chunks.length-1];  // TODO what if the original q contains a "/" ?
    q = decodeURIComponent(q);
    console.log("Local search for " + q);
    return withLocalIdiomsDB( (db) => {
      var hits = [];
      for(var i = 0; i < db.length; i++) {
        var idiom = db[i];
        if (matches(idiom, q)){
          hits.push(idiom);
          continue;
        }
      }
      return new Response( JSON.stringify(hits) );
    });
  }

  if (request.url.indexOf("/page/idiom-detail.html") != -1) {
    // Regarless the ?id=123 param, the HTML is in cache.
    return caches.match("/page/idiom-detail.html");
  }

  if (request.url.indexOf("/page/search-results.html") != -1) {
    // Regarless the ?q=foo param, the HTML is in cache.
    return caches.match("/page/search-results.html");
  }

  console.log("Not found " + request.url + " :(");
  return disconnectedResponse();
}

// Fetch list of all idioms in local cache, and call f.
// f must take the db as argument, and return a Response object.
function withLocalIdiomsDB(f) {
  return caches.match("/api/idioms/all").then(function(resp) {
    //debugger;
    if (resp) {
      // resp is a promise
      return resp.json().then(body => {
        var fullDB = body;
        return f(fullDB);
      });
    }else{
      console.log( "  No offline db :(");
      return;
    }
  }).catch(function(error){
    console.log("Cache error for /api/idioms/all : " + error);
    return disconnectedResponse();
  });
}

function matches(idiom, q) {
  q = q.toLowerCase();

  // If more than 1 word, the rule is "all words must match, independently"
  var words = q.split(/ +/);
  for(var i=0;i<words.length;i++) {
    var word = words[i];
    if (!word)
      continue;
    if (!matchesWord(idiom, word))
      return false;
  }
  return true;
}

function matchesWord(idiom, word) {
  var fieldMatches = (str) => str.toLowerCase().indexOf(word) !== -1;

  if (fieldMatches(idiom.Title)) {
    return true;
  }
  if (fieldMatches(idiom.LeadParagraph)) {
    return true;
  }
  for(var i=0;i<idiom.Implementations.length;i++){
    var impl = idiom.Implementations[i];
    if (fieldMatches(impl.LanguageName)) {
      return true;
    }
    if (fieldMatches(impl.CodeBlock)) {
      return true;
    }
    if (fieldMatches(impl.AuthorComment)) {
      return true;
    }
  }
  return false;
  // TODO return a score, for ranking
  // TODO match with language names synonyms
}

function disconnectedResponse() {
  return caches.match( '/default_' + ThemeDate + '/img/disconnected.png' );
}

// Pick a random element in a list.
function pick(x) {
  return x[Math.floor(Math.random() * x.length)];
}