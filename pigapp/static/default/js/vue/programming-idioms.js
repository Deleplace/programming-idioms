var app = new Vue({
    el: '#app',
    data: {
      idiom: {
        Title: 'Declare enumeration',
        LeadParagraph: 'Create an enumerated type _Suit with 4 possible values _SPADES, _HEARTS, _DIAMONDS, _CLUBS.',
        Implementations: [
            {
                LanguageName: 'Haskell',
                CodeBlock: 'data Suit = SPADES | HEARTS | DIAMONDS | CLUBS deriving (Enum)'
            },
            {
                LanguageName: 'Go',
                CodeBlock: `const (
                    SPADES = iota
                    HEARTS
                    DIAMONDS
                    CLUBS
                )`
            }
        ]
      }
    },
    methods: {
      fetch122(resource) {
        this.$http.get('/api/idiom/122').then((response) => {
          this.idiom = response.data;
        });
      }
    }
  });

  app.fetch122();