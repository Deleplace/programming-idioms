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
          this.idiom.LeadParagraphEmphasized = emphasize(this.idiom.LeadParagraph);
          for(var i=0; i<this.idiom.Implementations.length; i++){
            this.idiom.Implementations[i].AuthorCommentEmphasized = emphasize(this.idiom.Implementations[i].AuthorComment);
          }
        });
      }
    }
});

var idiomId = getParameterByName("id");
app.fetchIdiom(idiomId);