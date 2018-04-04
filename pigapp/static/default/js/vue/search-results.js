var app = new Vue({
    el: '#app',
    data: {
      q: getParameterByName("q"),
      results: []
    },
    methods: {
      search() {
        this.$http.get('/api/search/' + this.q).then((response) => {
          this.results = response.data;
        });
      },
      idiomURL(idiomId) {
        return "/default/html/idiom-detail.html?id=" + idiomId;
      }
    }
});

app.search();