var app = new Vue({
    el: '#app',
    data: {
      results: []
    },
    methods: {
        search(q) {
        this.$http.get('/api/search/' + q).then((response) => {
          this.results = response.data;
        });
      }
    }
  });

  app.search('complex');