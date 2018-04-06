var app = new Vue({
    el: '#app',
    data: {
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
          return "/page/idiom-detail.html?id=" + idiomId;
        }
    }
});


// Put the whole database in cache, for offline navigation
app.$http.get('/api/idioms/all');