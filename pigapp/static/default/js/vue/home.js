var app = new Vue({
    el: '#app',
    data: {
        q: ""
    },
    methods: {
        gotoRandomIdiom(){
            this.$http.get('/api/random-id').then((response) => {
                var idiomId = response.data;
                window.location = this.idiomURL(idiomId);
            });

            // TODO: if we have the whole list locally, then pick locally
            // instead of asking the server.
        },
        idiomURL(idiomId) {
          //return "/page/idiom-detail.html?id=" + idiomId;
          return "/idiom/" + idiomId;
          // TODO return "/idiom/" + idiomId + "/" + normalize title;
        },
        search(event) {
            if(!this.q) {
                console.warn("No search term :/");
                return;
            }
            console.log("Let's search for [" + this.q + "]");
            // TODO normalize q
            window.location = "/search/" + this.q;
            event.preventDefault();
        }
    }
});


// Put the whole database in cache, for offline navigation
console.log("Fetching the full DB");
app.$http.get('/api/idioms/all');