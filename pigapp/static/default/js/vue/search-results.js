var app = new Vue({
    el: '#app',
    data: {
      q: getParameterByName("q"),
      results: []
    },
    methods: {
      search() {
        this.$http.get('/api/search/' + this.q).then((response) => {
          var hits = response.data;
          for (var i=0;i<hits.length;i++)
            hits[i].LeadParagraphEmphasized = emphasize(hits[i].LeadParagraph);
          this.results = hits;
        });
      },
      idiomURL(idiomId) {
        return "/default/html/idiom-detail.html?id=" + idiomId;
      }
    }
});

app.search();