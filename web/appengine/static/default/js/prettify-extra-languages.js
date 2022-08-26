// Contributed by peter dot kofler at code minus cop dot org

/**
 * @fileoverview
 * Registers a language handler for Basic.
 *
 * To use, include prettify.js and this file in your HTML page.
 * Then put your code in an HTML tag like
 *      <pre class="prettyprint lang-basic">(my BASIC code)</pre>
 *
 * @author peter dot kofler at code minus cop dot org
 */

PR.registerLangHandler(
    PR.createSimpleLexer(
        [ // shortcutStylePatterns
          // "single-line-string"
          [PR.PR_STRING,        /^(?:"(?:[^\\"\r\n]|\\.)*(?:"|$))/, null, '"'],
          // Whitespace
          [PR.PR_PLAIN,         /^\s+/, null, ' \r\n\t\xA0']
        ],
        [ // fallthroughStylePatterns
          // A line comment that starts with REM
          [PR.PR_COMMENT,       /^REM[^\r\n]*/, null],
          [PR.PR_KEYWORD,       /^\b(?:AND|CLOSE|CLR|CMD|CONT|DATA|DEF ?FN|DIM|END|FOR|GET|GOSUB|GOTO|IF|INPUT|LET|LIST|LOAD|NEW|NEXT|NOT|ON|OPEN|OR|POKE|PRINT|READ|RESTORE|RETURN|RUN|SAVE|STEP|STOP|SYS|THEN|TO|VERIFY|WAIT)\b/, null],
          [PR.PR_PLAIN,         /^[A-Z][A-Z0-9]?(?:\$|%)?/i, null],
          // Literals .0, 0, 0.0 0E13
          [PR.PR_LITERAL,       /^(?:\d+(?:\.\d*)?|\.\d+)(?:e[+\-]?\d+)?/i,  null, '0123456789'],
          [PR.PR_PUNCTUATION,   /^.[^\s\w\.$%"]*/, null]
          // [PR.PR_PUNCTUATION,   /^[-,:;!<>=\+^\/\*]+/]
        ]),
    ['basic','cbm']);
// Copyright (C) 2013 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.



/**
 * @fileoverview
 * Registers a language handler Dart.
 * Loosely structured based on the DartLexer in Pygments: http://pygments.org/.
 *
 * To use, include prettify.js and this file in your HTML page.
 * Then put your code in an HTML tag like
 *      <pre class="prettyprint lang-dart">(Dart code)</pre>
 *
 * @author armstrong.timothy@gmail.com
 */

PR['registerLangHandler'](
  PR['createSimpleLexer'](
    [
      // Whitespace.
      [PR['PR_PLAIN'], /^[\t\n\r \xA0]+/, null, '\t\n\r \xA0']
    ],
    [
      // Script tag.
      [PR['PR_COMMENT'], /^#!(?:.*)/],

      // `import`, `library`, `part of`, `part`, `as`, `show`, and `hide`
      // keywords.
      [PR['PR_KEYWORD'], /^\b(?:import|library|part of|part|as|show|hide)\b/i],

      // Single-line comments.
      [PR['PR_COMMENT'], /^\/\/(?:.*)/],

      // Multiline comments.
      [PR['PR_COMMENT'], /^\/\*[^*]*\*+(?:[^\/*][^*]*\*+)*\//], // */

      // `class` and `interface` keywords.
      [PR['PR_KEYWORD'], /^\b(?:class|interface)\b/i],

      // General keywords.
      [PR['PR_KEYWORD'], /^\b(?:assert|break|case|catch|continue|default|do|else|finally|for|if|in|is|new|return|super|switch|this|throw|try|while)\b/i],

      // Declaration keywords.
      [PR['PR_KEYWORD'], /^\b(?:abstract|const|extends|factory|final|get|implements|native|operator|set|static|typedef|var)\b/i],

      // Keywords for types.
      [PR['PR_TYPE'], /^\b(?:bool|double|Dynamic|int|num|Object|String|void)\b/i],

      // Keywords for constants.
      [PR['PR_KEYWORD'], /^\b(?:false|null|true)\b/i],

      // Multiline strings, single- and double-quoted.
      [PR['PR_STRING'], /^r?[\']{3}[\s|\S]*?[^\\][\']{3}/],
      [PR['PR_STRING'], /^r?[\"]{3}[\s|\S]*?[^\\][\"]{3}/],

      // Normal and raw strings, single- and double-quoted.
      [PR['PR_STRING'], /^r?\'(\'|(?:[^\n\r\f])*?[^\\]\')/],
      [PR['PR_STRING'], /^r?\"(\"|(?:[^\n\r\f])*?[^\\]\")/],

      // Identifiers.
      [PR['PR_PLAIN'], /^[a-z_$][a-z0-9_]*/i],
      
      // Operators.
      [PR['PR_PUNCTUATION'], /^[~!%^&*+=|?:<>/-]/],

      // Hex numbers.
      [PR['PR_LITERAL'], /^\b0x[0-9a-f]+/i],

      // Decimal numbers.
      [PR['PR_LITERAL'], /^\b\d+(?:\.\d*)?(?:e[+-]?\d+)?/i],
      [PR['PR_LITERAL'], /^\b\.\d+(?:e[+-]?\d+)?/i],

      // Punctuation.
      [PR['PR_PUNCTUATION'], /^[(){}\[\],.;]/]
    ]),
  ['dart']);
// Copyright (C) 2013 Andrew Allen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


/**
 * @fileoverview
 * Registers a language handler for Erlang.
 *
 * Derived from https://raw.github.com/erlang/otp/dev/lib/compiler/src/core_parse.yrl
 * Modified from Mike Samuel's Haskell plugin for google-code-prettify
 *
 * @author achew22@gmail.com
 */

PR['registerLangHandler'](
    PR['createSimpleLexer'](
        [
         // Whitespace
         // whitechar    ->    newline | vertab | space | tab | uniWhite
         // newline      ->    return linefeed | return | linefeed | formfeed
         [PR['PR_PLAIN'],       /^[\t\n\x0B\x0C\r ]+/, null, '\t\n\x0B\x0C\r '],
         // Single line double-quoted strings.
         [PR['PR_STRING'],      /^\"(?:[^\"\\\n\x0C\r]|\\[\s\S])*(?:\"|$)/,
          null, '"'],
         
         // Handle atoms
         [PR['PR_LITERAL'],      /^[a-z][a-zA-Z0-9_]*/],
         // Handle single quoted atoms
         [PR['PR_LITERAL'],      /^\'(?:[^\'\\\n\x0C\r]|\\[^&])+\'?/,
          null, "'"],
         
         // Handle macros. Just to be extra clear on this one, it detects the ?
         // then uses the regexp to end it so be very careful about matching
         // all the terminal elements
         [PR['PR_LITERAL'],      /^\?[^ \t\n({]+/, null, "?"],

          
         
         // decimal      ->    digit{digit}
         // octal        ->    octit{octit}
         // hexadecimal  ->    hexit{hexit}
         // integer      ->    decimal
         //               |    0o octal | 0O octal
         //               |    0x hexadecimal | 0X hexadecimal
         // float        ->    decimal . decimal [exponent]
         //               |    decimal exponent
         // exponent     ->    (e | E) [+ | -] decimal
         [PR['PR_LITERAL'],
          /^(?:0o[0-7]+|0x[\da-f]+|\d+(?:\.\d+)?(?:e[+\-]?\d+)?)/i,
          null, '0123456789']
        ],
        [
         // TODO: catch @declarations inside comments

         // Comments in erlang are started with % and go till a newline
         [PR['PR_COMMENT'], /^%[^\n]*/],

         // Catch macros
         //[PR['PR_TAG'], /?[^( \n)]+/],

         /**
          * %% Keywords (atoms are assumed to always be single-quoted).
          * 'module' 'attributes' 'do' 'let' 'in' 'letrec'
          * 'apply' 'call' 'primop'
          * 'case' 'of' 'end' 'when' 'fun' 'try' 'catch' 'receive' 'after'
          */
         [PR['PR_KEYWORD'], /^(?:module|attributes|do|let|in|letrec|apply|call|primop|case|of|end|when|fun|try|catch|receive|after|char|integer|float,atom,string,var)\b/],
         
         /**
          * Catch definitions (usually defined at the top of the file)
          * Anything that starts -something
          */
         [PR['PR_KEYWORD'], /^-[a-z_]+/],

         // Catch variables
         [PR['PR_TYPE'], /^[A-Z_][a-zA-Z0-9_]*/],

         // matches the symbol production
         [PR['PR_PUNCTUATION'], /^[.,;]/]
        ]),
    ['erlang', 'erl']);
// Copyright (C) 2010 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.



/**
 * @fileoverview
 * Registers a language handler for the Go language..
 * <p>
 * Based on the lexical grammar at 
 * http://golang.org/doc/go_spec.html#Lexical_elements
 * <p>
 * Go uses a minimal style for highlighting so the below does not distinguish
 * strings, keywords, literals, etc. by design.
 * From a discussion with the Go designers:
 * <pre>
 * On Thursday, July 22, 2010, Mike Samuel <...> wrote:
 * > On Thu, Jul 22, 2010, Rob 'Commander' Pike <...> wrote:
 * >> Personally, I would vote for the subdued style godoc presents at http://golang.org
 * >>
 * >> Not as fancy as some like, but a case can be made it's the official style.
 * >> If people want more colors, I wouldn't fight too hard, in the interest of
 * >> encouragement through familiarity, but even then I would ask to shy away
 * >> from technicolor starbursts.
 * >
 * > Like http://golang.org/pkg/go/scanner/ where comments are blue and all
 * > other content is black?  I can do that.
 * </pre>
 *
 * @author mikesamuel@gmail.com
 */

PR['registerLangHandler'](
    PR['createSimpleLexer'](
        [
         // Whitespace is made up of spaces, tabs and newline characters.
         [PR['PR_PLAIN'],       /^[\t\n\r \xA0]+/, null, '\t\n\r \xA0'],
         // Not escaped as a string.  See note on minimalism above.
         [PR['PR_PLAIN'],       /^(?:\"(?:[^\"\\]|\\[\s\S])*(?:\"|$)|\'(?:[^\'\\]|\\[\s\S])+(?:\'|$)|`[^`]*(?:`|$))/, null, '"\'']
        ],
        [
         // Block comments are delimited by /* and */.
         // Single-line comments begin with // and extend to the end of a line.
         [PR['PR_COMMENT'],     /^(?:\/\/[^\r\n]*|\/\*[\s\S]*?\*\/)/],
         [PR['PR_PLAIN'],       /^(?:[^\/\"\'`]|\/(?![\/\*]))+/i]
        ]),
    ['go']);
// Copyright (C) 2008 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.



/**
 * @fileoverview
 * Registers a language handler for Common Lisp and related languages.
 *
 *
 * To use, include prettify.js and this file in your HTML page.
 * Then put your code in an HTML tag like
 *      <pre class="prettyprint lang-lisp">(my lisp code)</pre>
 * The lang-cl class identifies the language as common lisp.
 * This file supports the following language extensions:
 *     lang-cl - Common Lisp
 *     lang-el - Emacs Lisp
 *     lang-lisp - Lisp
 *     lang-scm - Scheme
 *     lang-lsp - FAT 8.3 filename version of lang-lisp.
 *
 *
 * I used http://www.devincook.com/goldparser/doc/meta-language/grammar-LISP.htm
 * as the basis, but added line comments that start with ; and changed the atom
 * production to disallow unquoted semicolons.
 *
 * "Name"    = 'LISP'
 * "Author"  = 'John McCarthy'
 * "Version" = 'Minimal'
 * "About"   = 'LISP is an abstract language that organizes ALL'
 *           | 'data around "lists".'
 *
 * "Start Symbol" = [s-Expression]
 *
 * {Atom Char}   = {Printable} - {Whitespace} - [()"\'']
 *
 * Atom = ( {Atom Char} | '\'{Printable} )+
 *
 * [s-Expression] ::= [Quote] Atom
 *                  | [Quote] '(' [Series] ')'
 *                  | [Quote] '(' [s-Expression] '.' [s-Expression] ')'
 *
 * [Series] ::= [s-Expression] [Series]
 *            |
 *
 * [Quote]  ::= ''      !Quote = do not evaluate
 *            |
 *
 *
 * I used <a href="http://gigamonkeys.com/book/">Practical Common Lisp</a> as
 * the basis for the reserved word list.
 *
 *
 * @author mikesamuel@gmail.com
 */

PR['registerLangHandler'](
    PR['createSimpleLexer'](
        [
         ['opn',             /^\(+/, null, '('],
         ['clo',             /^\)+/, null, ')'],
         // A line comment that starts with ;
         [PR['PR_COMMENT'],     /^;[^\r\n]*/, null, ';'],
         // Whitespace
         [PR['PR_PLAIN'],       /^[\t\n\r \xA0]+/, null, '\t\n\r \xA0'],
         // A double quoted, possibly multi-line, string.
         [PR['PR_STRING'],      /^\"(?:[^\"\\]|\\[\s\S])*(?:\"|$)/, null, '"']
        ],
        [
         [PR['PR_KEYWORD'],     /^(?:block|c[ad]+r|catch|con[ds]|def(?:ine|un)|do|eq|eql|equal|equalp|eval-when|flet|format|go|if|labels|lambda|let|load-time-value|locally|macrolet|multiple-value-call|nil|progn|progv|quote|require|return-from|setq|symbol-macrolet|t|tagbody|the|throw|unwind)\b/, null],
         [PR['PR_LITERAL'],
          /^[+\-]?(?:[0#]x[0-9a-f]+|\d+\/\d+|(?:\.\d+|\d+(?:\.\d*)?)(?:[ed][+\-]?\d+)?)/i],
         // A single quote possibly followed by a word that optionally ends with
         // = ! or ?.
         [PR['PR_LITERAL'],
          /^\'(?:-*(?:\w|\\[\x21-\x7e])(?:[\w-]*|\\[\x21-\x7e])[=!?]?)?/],
         // A word that optionally ends with = ! or ?.
         [PR['PR_PLAIN'],
          /^-*(?:[a-z_]|\\[\x21-\x7e])(?:[\w-]*|\\[\x21-\x7e])[=!?]?/i],
         // A printable non-space non-special character
         [PR['PR_PUNCTUATION'], /^[^\w\t\n\r \xA0()\"\\\';]+/]
        ]),
    ['cl', 'el', 'lisp', 'lsp', 'scm', 'ss', 'rkt']);
// Copyright (C) 2013 Nikhil Dabas
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


/**
 * @fileoverview
 * Registers a language handler for LLVM.
 * From https://gist.github.com/ndabas/2850418
 *
 *
 * To use, include prettify.js and this file in your HTML page.
 * Then put your code in an HTML tag like
 *      <pre class="prettyprint lang-llvm">(my LLVM code)</pre>
 *
 *
 * The regular expressions were adapted from:
 * https://github.com/hansstimer/llvm.tmbundle/blob/76fedd8f50fd6108b1780c51d79fbe3223de5f34/Syntaxes/LLVM.tmLanguage
 * 
 * http://llvm.org/docs/LangRef.html#constants describes the language grammar.
 * 
 * @author Nikhil Dabas
 */
PR['registerLangHandler'](
    PR['createSimpleLexer'](
        [
         // Whitespace
         [PR['PR_PLAIN'],       /^[\t\n\r \xA0]+/, null, '\t\n\r \xA0'],
         // A double quoted, possibly multi-line, string.
         [PR['PR_STRING'],      /^!?\"(?:[^\"\\]|\\[\s\S])*(?:\"|$)/, null, '"'],
         // comment.llvm
         [PR['PR_COMMENT'],     /^;[^\r\n]*/, null, ';']
        ],
        [
         // variable.llvm
         [PR['PR_PLAIN'],       /^[%@!](?:[-a-zA-Z$._][-a-zA-Z$._0-9]*|\d+)/],

         // According to http://llvm.org/docs/LangRef.html#well-formedness
         // These reserved words cannot conflict with variable names, because none of them start with a prefix character ('%' or '@').
         [PR['PR_KEYWORD'],     /^[A-Za-z_][0-9A-Za-z_]*/, null],

         // constant.numeric.float.llvm
         [PR['PR_LITERAL'],     /^\d+\.\d+/],
         
         // constant.numeric.integer.llvm
         [PR['PR_LITERAL'],     /^(?:\d+|0[xX][a-fA-F0-9]+)/],

         // punctuation
         [PR['PR_PUNCTUATION'], /^[()\[\]{},=*<>:]|\.\.\.$/]
        ]),
    ['llvm', 'll']);
// Copyright (C) 2008 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.



/**
 * @fileoverview
 * Registers a language handler for Lua.
 *
 *
 * To use, include prettify.js and this file in your HTML page.
 * Then put your code in an HTML tag like
 *      <pre class="prettyprint lang-lua">(my Lua code)</pre>
 *
 *
 * I used http://www.lua.org/manual/5.1/manual.html#2.1
 * Because of the long-bracket concept used in strings and comments, Lua does
 * not have a regular lexical grammar, but luckily it fits within the space
 * of irregular grammars supported by javascript regular expressions.
 *
 * @author mikesamuel@gmail.com
 */

PR['registerLangHandler'](
    PR['createSimpleLexer'](
        [
         // Whitespace
         [PR['PR_PLAIN'],       /^[\t\n\r \xA0]+/, null, '\t\n\r \xA0'],
         // A double or single quoted, possibly multi-line, string.
         [PR['PR_STRING'],      /^(?:\"(?:[^\"\\]|\\[\s\S])*(?:\"|$)|\'(?:[^\'\\]|\\[\s\S])*(?:\'|$))/, null, '"\'']
        ],
        [
         // A comment is either a line comment that starts with two dashes, or
         // two dashes preceding a long bracketed block.
         [PR['PR_COMMENT'], /^--(?:\[(=*)\[[\s\S]*?(?:\]\1\]|$)|[^\r\n]*)/],
         // A long bracketed block not preceded by -- is a string.
         [PR['PR_STRING'],  /^\[(=*)\[[\s\S]*?(?:\]\1\]|$)/],
         [PR['PR_KEYWORD'], /^(?:and|break|do|else|elseif|end|false|for|function|if|in|local|nil|not|or|repeat|return|then|true|until|while)\b/, null],
         // A number is a hex integer literal, a decimal real literal, or in
         // scientific notation.
         [PR['PR_LITERAL'],
          /^[+-]?(?:0x[\da-f]+|(?:(?:\.\d+|\d+(?:\.\d*)?)(?:e[+\-]?\d+)?))/i],
         // An identifier
         [PR['PR_PLAIN'], /^[a-z_]\w*/i],
         // A run of punctuation
         [PR['PR_PUNCTUATION'], /^[^\w\t\n\r \xA0][^\w\t\n\r \xA0\"\'\-\+=]*/]
        ]),
    ['lua']);
// Copyright (C) 2008 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.



/**
 * @fileoverview
 * Registers a language handler for OCaml, SML, F# and similar languages.
 *
 * Based on the lexical grammar at
 * http://research.microsoft.com/en-us/um/cambridge/projects/fsharp/manual/spec.html#_Toc270597388
 *
 * @author mikesamuel@gmail.com
 */

PR['registerLangHandler'](
    PR['createSimpleLexer'](
        [
         // Whitespace is made up of spaces, tabs and newline characters.
         [PR['PR_PLAIN'],       /^[\t\n\r \xA0]+/, null, '\t\n\r \xA0'],
         // #if ident/#else/#endif directives delimit conditional compilation
         // sections
         [PR['PR_COMMENT'],
          /^#(?:if[\t\n\r \xA0]+(?:[a-z_$][\w\']*|``[^\r\n\t`]*(?:``|$))|else|endif|light)/i,
          null, '#'],
         // A double or single quoted, possibly multi-line, string.
         // F# allows escaped newlines in strings.
         [PR['PR_STRING'],      /^(?:\"(?:[^\"\\]|\\[\s\S])*(?:\"|$)|\'(?:[^\'\\]|\\[\s\S])(?:\'|$))/, null, '"\'']
        ],
        [
         // Block comments are delimited by (* and *) and may be
         // nested. Single-line comments begin with // and extend to
         // the end of a line.
         // TODO: (*...*) comments can be nested.  This does not handle that.
         [PR['PR_COMMENT'],     /^(?:\/\/[^\r\n]*|\(\*[\s\S]*?\*\))/],
         [PR['PR_KEYWORD'],     /^(?:abstract|and|as|assert|begin|class|default|delegate|do|done|downcast|downto|elif|else|end|exception|extern|false|finally|for|fun|function|if|in|inherit|inline|interface|internal|lazy|let|match|member|module|mutable|namespace|new|null|of|open|or|override|private|public|rec|return|static|struct|then|to|true|try|type|upcast|use|val|void|when|while|with|yield|asr|land|lor|lsl|lsr|lxor|mod|sig|atomic|break|checked|component|const|constraint|constructor|continue|eager|event|external|fixed|functor|global|include|method|mixin|object|parallel|process|protected|pure|sealed|trait|virtual|volatile)\b/],
         // A number is a hex integer literal, a decimal real literal, or in
         // scientific notation.
         [PR['PR_LITERAL'],
          /^[+\-]?(?:0x[\da-f]+|(?:(?:\.\d+|\d+(?:\.\d*)?)(?:e[+\-]?\d+)?))/i],
         [PR['PR_PLAIN'],       /^(?:[a-z_][\w']*[!?#]?|``[^\r\n\t`]*(?:``|$))/i],
         // A printable non-space non-special character
         [PR['PR_PUNCTUATION'], /^[^\t\n\r \xA0\"\'\w]+/]
        ]),
    ['fs', 'ml']);
// Copyright (C) 2011 Zimin A.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


/**
 * @fileoverview
 * Registers a language handler for the Nemerle language.
 * http://nemerle.org
 * @author Zimin A.V.
 */
(function () {
  // http://nemerle.org/wiki/index.php?title=Base_keywords
  var keywords = 'abstract|and|as|base|catch|class|def|delegate|enum|event|extern|false|finally|'
         + 'fun|implements|interface|internal|is|macro|match|matches|module|mutable|namespace|new|'
         + 'null|out|override|params|partial|private|protected|public|ref|sealed|static|struct|'
         + 'syntax|this|throw|true|try|type|typeof|using|variant|virtual|volatile|when|where|with|'
         + 'assert|assert2|async|break|checked|continue|do|else|ensures|for|foreach|if|late|lock|new|nolate|'
         + 'otherwise|regexp|repeat|requires|return|surroundwith|unchecked|unless|using|while|yield';

  PR['registerLangHandler'](PR['createSimpleLexer'](
      // shortcutStylePatterns
      [
        [PR['PR_STRING'], /^(?:\'(?:[^\\\'\r\n]|\\.)*\'|\"(?:[^\\\"\r\n]|\\.)*(?:\"|$))/, null, '"'],
        [PR['PR_COMMENT'], /^#(?:(?:define|elif|else|endif|error|ifdef|include|ifndef|line|pragma|undef|warning)\b|[^\r\n]*)/, null, '#'],
        [PR['PR_PLAIN'], /^\s+/, null, ' \r\n\t\xA0']
      ],
      // fallthroughStylePatterns
      [
        [PR['PR_STRING'], /^@\"(?:[^\"]|\"\")*(?:\"|$)/, null],
        [PR['PR_STRING'], /^<#(?:[^#>])*(?:#>|$)/, null],
        [PR['PR_STRING'], /^<(?:(?:(?:\.\.\/)*|\/?)(?:[\w-]+(?:\/[\w-]+)+)?[\w-]+\.h|[a-z]\w*)>/, null],
        [PR['PR_COMMENT'], /^\/\/[^\r\n]*/, null],
        [PR['PR_COMMENT'], /^\/\*[\s\S]*?(?:\*\/|$)/, null],
        [PR['PR_KEYWORD'], new RegExp('^(?:' + keywords + ')\\b'), null],
        [PR['PR_TYPE'], /^(?:array|bool|byte|char|decimal|double|float|int|list|long|object|sbyte|short|string|ulong|uint|ufloat|ulong|ushort|void)\b/, null],
        [PR['PR_LITERAL'], /^@[a-z_$][a-z_$@0-9]*/i, null],
        [PR['PR_TYPE'], /^@[A-Z]+[a-z][A-Za-z_$@0-9]*/, null],
        [PR['PR_PLAIN'], /^'?[A-Za-z_$][a-z_$@0-9]*/i, null],
        [PR['PR_LITERAL'], new RegExp(
             '^(?:'
  // A hex number
             + '0x[a-f0-9]+'
  // or an octal or decimal number,
             + '|(?:\\d(?:_\\d+)*\\d*(?:\\.\\d*)?|\\.\\d\\+)'
  // possibly in scientific notation
             + '(?:e[+\\-]?\\d+)?'
             + ')'
  // with an optional modifier like UL for unsigned long
             + '[a-z]*', 'i'), null, '0123456789'],

        [PR['PR_PUNCTUATION'], /^.[^\s\w\.$@\'\"\`\/\#]*/, null]
      ]),
      ['n', 'nemerle']);
})();
// Contributed by peter dot kofler at code minus cop dot org

/**
 * @fileoverview
 * Registers a language handler for (Turbo) Pascal.
 *
 * To use, include prettify.js and this file in your HTML page.
 * Then put your code in an HTML tag like
 *      <pre class="prettyprint lang-pascal">(my Pascal code)</pre>
 *
 * @author peter dot kofler at code minus cop dot org
 */

PR.registerLangHandler(
    PR.createSimpleLexer(
        [ // shortcutStylePatterns
          // 'single-line-string'
          [PR.PR_STRING,        /^(?:\'(?:[^\\\'\r\n]|\\.)*(?:\'|$))/, null, '\''],
          // Whitespace
          [PR.PR_PLAIN,         /^\s+/, null, ' \r\n\t\xA0']
        ],
        [ // fallthroughStylePatterns
          // A cStyleComments comment (* *) or {}
          [PR.PR_COMMENT,       /^\(\*[\s\S]*?(?:\*\)|$)|^\{[\s\S]*?(?:\}|$)/, null],
          [PR.PR_KEYWORD,       /^(?:ABSOLUTE|AND|ARRAY|ASM|ASSEMBLER|BEGIN|CASE|CONST|CONSTRUCTOR|DESTRUCTOR|DIV|DO|DOWNTO|ELSE|END|EXTERNAL|FOR|FORWARD|FUNCTION|GOTO|IF|IMPLEMENTATION|IN|INLINE|INTERFACE|INTERRUPT|LABEL|MOD|NOT|OBJECT|OF|OR|PACKED|PROCEDURE|PROGRAM|RECORD|REPEAT|SET|SHL|SHR|THEN|TO|TYPE|UNIT|UNTIL|USES|VAR|VIRTUAL|WHILE|WITH|XOR)\b/i, null],
          [PR.PR_LITERAL,       /^(?:true|false|self|nil)/i, null],
          [PR.PR_PLAIN,         /^[a-z][a-z0-9]*/i, null],
          // Literals .0, 0, 0.0 0E13
          [PR.PR_LITERAL,       /^(?:\$[a-f0-9]+|(?:\d+(?:\.\d*)?|\.\d+)(?:e[+\-]?\d+)?)/i,  null, '0123456789'],
          [PR.PR_PUNCTUATION,   /^.[^\s\w\.$@\'\/]*/, null]
        ]),
    ['pascal']);
// Copyright (C) 2012 Jeffrey B. Arnold
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


/**
 * @fileoverview
 * Registers a language handler for S, S-plus, and R source code.
 *
 *
 * To use, include prettify.js and this file in your HTML page.
 * Then put your code in an HTML tag like
 *      <pre class="prettyprint lang-r"> code </pre>
 *
 * Language definition from
 * http://cran.r-project.org/doc/manuals/R-lang.html.
 * Many of the regexes are shared  with the pygments SLexer,
 * http://pygments.org/.
 *
 * Original: https://raw.github.com/jrnold/prettify-lang-r-bugs/master/lang-r.js
 *
 * @author jeffrey.arnold@gmail.com
 */
PR['registerLangHandler'](
    PR['createSimpleLexer'](
        [
            [PR['PR_PLAIN'],       /^[\t\n\r \xA0]+/, null, '\t\n\r \xA0'],
	    [PR['PR_STRING'],      /^\"(?:[^\"\\]|\\[\s\S])*(?:\"|$)/, null, '"'],
	    [PR['PR_STRING'],      /^\'(?:[^\'\\]|\\[\s\S])*(?:\'|$)/, null, "'"]
        ],
        [
            [PR['PR_COMMENT'],     /^#.*/],
	    [PR['PR_KEYWORD'],     /^(?:if|else|for|while|repeat|in|next|break|return|switch|function)(?![A-Za-z0-9_.])/],
	    // hex numbes
	    [PR['PR_LITERAL'], /^0[xX][a-fA-F0-9]+([pP][0-9]+)?[Li]?/],
	    // Decimal numbers
            [PR['PR_LITERAL'], /^[+-]?([0-9]+(\.[0-9]+)?|\.[0-9]+)([eE][+-]?[0-9]+)?[Li]?/],
	    // builtin symbols
	    [PR['PR_LITERAL'], /^(?:NULL|NA(?:_(?:integer|real|complex|character)_)?|Inf|TRUE|FALSE|NaN|\.\.(?:\.|[0-9]+))(?![A-Za-z0-9_.])/],
	    // assignment, operators, and parens, etc.
	    [PR['PR_PUNCTUATION'], /^(?:<<?-|->>?|-|==|<=|>=|<|>|&&?|!=|\|\|?|\*|\+|\^|\/|!|%.*?%|=|~|\$|@|:{1,3}|[\[\](){};,?])/],
	    // valid variable names
	    [PR['PR_PLAIN'], /^(?:[A-Za-z]+[A-Za-z0-9_.]*|\.[a-zA-Z_][0-9a-zA-Z\._]*)(?![A-Za-z0-9_.])/],
	    // string backtick
	    [PR['PR_STRING'], /^`.+`/]
        ]),
    ['r', 's', 'R', 'S', 'Splus']);
// Copyright (C) 2010 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


/**
 * @fileoverview
 * Registers a language handler for Scala.
 *
 * Derived from http://lampsvn.epfl.ch/svn-repos/scala/scala-documentation/trunk/src/reference/SyntaxSummary.tex
 *
 * @author mikesamuel@gmail.com
 */

PR['registerLangHandler'](
    PR['createSimpleLexer'](
        [
         // Whitespace
         [PR['PR_PLAIN'],       /^[\t\n\r \xA0]+/, null, '\t\n\r \xA0'],
         // A double or single quoted string 
          // or a triple double-quoted multi-line string.
         [PR['PR_STRING'],
          /^(?:"(?:(?:""(?:""?(?!")|[^\\"]|\\.)*"{0,3})|(?:[^"\r\n\\]|\\.)*"?))/,
          null, '"'],
         [PR['PR_LITERAL'],     /^`(?:[^\r\n\\`]|\\.)*`?/, null, '`'],
         [PR['PR_PUNCTUATION'], /^[!#%&()*+,\-:;<=>?@\[\\\]^{|}~]+/, null,
          '!#%&()*+,-:;<=>?@[\\]^{|}~']
        ],
        [
         // A symbol literal is a single quote followed by an identifier with no
         // single quote following
         // A character literal has single quotes on either side
         [PR['PR_STRING'],      /^'(?:[^\r\n\\']|\\(?:'|[^\r\n']+))'/],
         [PR['PR_LITERAL'],     /^'[a-zA-Z_$][\w$]*(?!['$\w])/],
         [PR['PR_KEYWORD'],     /^(?:abstract|case|catch|class|def|do|else|extends|final|finally|for|forSome|if|implicit|import|lazy|match|new|object|override|package|private|protected|requires|return|sealed|super|throw|trait|try|type|val|var|while|with|yield)\b/],
         [PR['PR_LITERAL'],     /^(?:true|false|null|this)\b/],
         [PR['PR_LITERAL'],     /^(?:(?:0(?:[0-7]+|X[0-9A-F]+))L?|(?:(?:0|[1-9][0-9]*)(?:(?:\.[0-9]+)?(?:E[+\-]?[0-9]+)?F?|L?))|\\.[0-9]+(?:E[+\-]?[0-9]+)?F?)/i],
         // Treat upper camel case identifiers as types.
         [PR['PR_TYPE'],        /^[$_]*[A-Z][_$A-Z0-9]*[a-z][\w$]*/],
         [PR['PR_PLAIN'],       /^[$a-zA-Z_][\w$]*/],
         [PR['PR_COMMENT'],     /^\/(?:\/.*|\*(?:\/|\**[^*/])*(?:\*+\/?)?)/],
         [PR['PR_PUNCTUATION'], /^(?:\.+|\/)/]
        ]),
    ['scala']);
// Copyright (C) 2009 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.



/**
 * @fileoverview
 * Registers a language handler for various flavors of basic.
 *
 *
 * To use, include prettify.js and this file in your HTML page.
 * Then put your code in an HTML tag like
 *      <pre class="prettyprint lang-vb"></pre>
 *
 *
 * http://msdn.microsoft.com/en-us/library/aa711638(VS.71).aspx defines the
 * visual basic grammar lexical grammar.
 *
 * @author mikesamuel@gmail.com
 */

PR['registerLangHandler'](
    PR['createSimpleLexer'](
        [
         // Whitespace
         [PR['PR_PLAIN'],       /^[\t\n\r \xA0\u2028\u2029]+/, null, '\t\n\r \xA0\u2028\u2029'],
         // A double quoted string with quotes escaped by doubling them.
         // A single character can be suffixed with C.
         [PR['PR_STRING'],      /^(?:[\"\u201C\u201D](?:[^\"\u201C\u201D]|[\"\u201C\u201D]{2})(?:[\"\u201C\u201D]c|$)|[\"\u201C\u201D](?:[^\"\u201C\u201D]|[\"\u201C\u201D]{2})*(?:[\"\u201C\u201D]|$))/i, null,
          '"\u201C\u201D'],
         // A comment starts with a single quote and runs until the end of the
         // line.
         // VB6 apparently allows _ as an escape sequence for newlines though
         // this is not a documented feature of VB.net.
         // http://meta.stackoverflow.com/q/121497/137403
         [PR['PR_COMMENT'],     /^[\'\u2018\u2019](?:_(?:\r\n?|[^\r]?)|[^\r\n_\u2028\u2029])*/, null, '\'\u2018\u2019']
        ],
        [
         [PR['PR_KEYWORD'], /^(?:AddHandler|AddressOf|Alias|And|AndAlso|Ansi|As|Assembly|Auto|Boolean|ByRef|Byte|ByVal|Call|Case|Catch|CBool|CByte|CChar|CDate|CDbl|CDec|Char|CInt|Class|CLng|CObj|Const|CShort|CSng|CStr|CType|Date|Decimal|Declare|Default|Delegate|Dim|DirectCast|Do|Double|Each|Else|ElseIf|End|EndIf|Enum|Erase|Error|Event|Exit|Finally|For|Friend|Function|Get|GetType|GoSub|GoTo|Handles|If|Implements|Imports|In|Inherits|Integer|Interface|Is|Let|Lib|Like|Long|Loop|Me|Mod|Module|MustInherit|MustOverride|MyBase|MyClass|Namespace|New|Next|Not|NotInheritable|NotOverridable|Object|On|Option|Optional|Or|OrElse|Overloads|Overridable|Overrides|ParamArray|Preserve|Private|Property|Protected|Public|RaiseEvent|ReadOnly|ReDim|RemoveHandler|Resume|Return|Select|Set|Shadows|Shared|Short|Single|Static|Step|Stop|String|Structure|Sub|SyncLock|Then|Throw|To|Try|TypeOf|Unicode|Until|Variant|Wend|When|While|With|WithEvents|WriteOnly|Xor|EndIf|GoSub|Let|Variant|Wend)\b/i, null],
         // A second comment form
         [PR['PR_COMMENT'], /^REM\b[^\r\n\u2028\u2029]*/i],
         // A boolean, numeric, or date literal.
         [PR['PR_LITERAL'],
          /^(?:True\b|False\b|Nothing\b|\d+(?:E[+\-]?\d+[FRD]?|[FRDSIL])?|(?:&H[0-9A-F]+|&O[0-7]+)[SIL]?|\d*\.\d+(?:E[+\-]?\d+)?[FRD]?|#\s+(?:\d+[\-\/]\d+[\-\/]\d+(?:\s+\d+:\d+(?::\d+)?(\s*(?:AM|PM))?)?|\d+:\d+(?::\d+)?(\s*(?:AM|PM))?)\s+#)/i],
         // An identifier.  Keywords can be turned into identifers
         // with square brackets, and there may be optional type
         // characters after a normal identifier in square brackets.
         [PR['PR_PLAIN'], /^(?:(?:[a-z]|_\w)\w*(?:\[[%&@!#]+\])?|\[(?:[a-z]|_\w)\w*\])/i],
         // A run of punctuation
         [PR['PR_PUNCTUATION'],
          /^[^\w\t\n\r \"\'\[\]\xA0\u2018\u2019\u201C\u201D\u2028\u2029]+/],
         // Square brackets
         [PR['PR_PUNCTUATION'], /^(?:\[|\])/]
        ]),
    ['vb', 'vbs']);
