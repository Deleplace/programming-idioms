
function extractSearchQueryFromCurrentLocation() {
  // e.g. https://programming-idioms.org/page/search-results.html?q=tree
  var q = getParameterByName("q");
  if( q )
    return q;

  // e.g. https://programming-idioms.org/search/tree
  // Warning this could be confused by other endpoint /api/search/tree
  var searchParams = window.location.href.match(/\/search\/(.+)$/);
  if ( searchParams ){
    q = searchParams[1];
    return q;
  }

  return null;
}

const app = new Vue({
    data: {
      q: extractSearchQueryFromCurrentLocation(),
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
      idiomURL(idiom) {
        //return "/default/html/idiom-detail.html?id=" + idiomId;
        //return "/idiom/" + idiomId;
        return "/idiom/" + idiom.Id + "/" + uriNormalize(idiom.Title);
      }
    }
}).$mount('#app');

app.search();