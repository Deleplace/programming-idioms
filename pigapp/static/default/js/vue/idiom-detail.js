var app = new Vue({
    el: '#app',
    data: {
      idiom: null,
    },
    methods: {
      fetchIdiom(idiomId) {
        this.$http.get('/api/idiom/' + idiomId).then((response) => {
          this.idiom = response.data;
          this.idiom.LeadParagraphEmphasized = emphasize(this.idiom.LeadParagraph);
          for(var i=0; i<this.idiom.Implementations.length; i++){
            this.idiom.Implementations[i].AuthorCommentEmphasized = emphasize(this.idiom.Implementations[i].AuthorComment);
            this.idiom.Implementations[i].PrettyClass = Object();
            this.idiom.Implementations[i].PrettyClass["lang-" + this.idiom.Implementations[i].LanguageName.toLowerCase()] = true;
          }
          // Favlangs will be first!
          this.idiom.Implementations.sort( (a,b) => {
             if ( isFavLang(a.LanguageName) )
               return -1;
             if ( isFavLang(b.LanguageName) )
               return 1;
              
            return a.LanguageName.localeCompare(b.LanguageName);
         });
         this.$nextTick(prettyPrint);
        });
      },
      fetchIdiomImpl(idiomId, implId) {
        this.$http.get('/api/idiom/' + idiomId).then((response) => {
          this.idiom = response.data;
          this.idiom.LeadParagraphEmphasized = emphasize(this.idiom.LeadParagraph);
          var implLang = "";
          for(var i=0; i<this.idiom.Implementations.length; i++){
            this.idiom.Implementations[i].AuthorCommentEmphasized = emphasize(this.idiom.Implementations[i].AuthorComment);
            if ( this.idiom.Implementations[i].Id == implId )
            implLang = this.idiom.Implementations[i].LanguageName;
          }
          // The selected impl will be first!
          // Other impls for same lang will be second!
          // Favlangs will be third!
          this.idiom.Implementations.sort( (a,b) => {
             if ( a.Id == implId )
              return -1;
             if ( b.Id == implId )
              return 1;
              if ( a.LanguageName == implLang )
               return -1;
              if ( b.LanguageName == implLang )
               return 1;
              if ( isFavLang(a.LanguageName) )
                return -1;
              if ( isFavLang(b.LanguageName) )
                return 1;
               
             return a.LanguageName.localeCompare(b.LanguageName);
          });
          this.$nextTick(prettyPrint);
        });
      }
    }
});

function extractIdiomIdFromCurrentLocation() {
  // e.g. https://programming-idioms.org/page/idiom-detail.html?id=123
  var idiomId = getParameterByName("id");
  if( idiomId )
    return idiomId;

  // e.g. https://programming-idioms.org/idiom/127/source-code-inclusion/1616/pascal
  var idiomDetailParams = window.location.href.match(/\/idiom\/(\d+)/);
  if ( idiomDetailParams ){
    idiomId = idiomDetailParams[1];
    return idiomId;
  }

  return null;
}

function extractImplIdFromCurrentLocation() {
  // e.g. https://programming-idioms.org/idiom/127/source-code-inclusion/1616/pascal
  var idiomDetailParams = window.location.href.match(/\/idiom\/(\d+)\/[^/]+\/(\d+)/);
  if ( idiomDetailParams ) {
    var implId = idiomDetailParams[2];
    return implId;
  }

  return null;
}

function extractImplLangFromCurrentLocation() {
  // The alleged lang is extracted from URL, no hard proof to be consistent with implId.
  // e.g. https://programming-idioms.org/idiom/127/source-code-inclusion/1616/pascal
  var idiomDetailParams = window.location.href.match(/\/idiom\/(\d+)\/[^/]+\/(\d+)\/[^/]+$/);
  if ( idiomDetailParams ) {
    var implLang = idiomDetailParams[3];
    return implLang;
  }

  return null;
}

// console.log("idiom-detail.js rendering " + window.location.href + " ...");
var idiomId = extractIdiomIdFromCurrentLocation();
// console.log("  Idiom id is " + idiomId);
var implId = extractImplIdFromCurrentLocation();
// console.log("  Impl id is " + implId);
if(implId)
  app.fetchIdiomImpl( idiomId, implId );
else
  app.fetchIdiom( idiomId );
