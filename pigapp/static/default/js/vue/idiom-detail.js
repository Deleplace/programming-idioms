var app = new Vue({
    el: '#app',
    data: {
      idiom: null,
    },
    methods: {
      fetch122() {
        this.fetchIdiom(122);
      },
      fetchIdiom(idiomId) {
        this.$http.get('/api/idiom/' + idiomId).then((response) => {
          this.idiom = response.data;
        });
      }
    }
});

var idiomId = getParameterByName("id");
app.fetchIdiom(idiomId);