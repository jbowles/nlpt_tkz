# nlpt\_tkz
Natural language tokenizer; supports a simple white space tokenizer (preferable for simple tasks), a unicode pattern-matcher, and state function lexer.

There are 2 toplevel functions: 

* `TokenizeStr(string, string) []string *Digest`
* `TokenizeBytes([]byte, string) *Digest`. 

Both functions require you to select which tokenizer type to be used (the second argument). The `TokenizeBytes` function only supports `lex` and 'unicode` options and will only return a `Digest` struct with a `Bytes` field; it is for use in instances where dealing with strings in not preferred or possible. The former function is for dealing direclty with strings and supports `lex`, `unicode`, `whitespace`; it returns a slice of tokens as well as digest, and depending on the tokenizer used the digest will contain different data fields.

## What Kind of text can I tokenize
Any kind. During development I often used publicly available corpora with particularly nasty text. For example, the [20 Newsgroups data set](http://qwone.com/~jason/20Newsgroups/) has archived news and emails from early internet groups.... Specifically, the `comp.windows.x` data set is wealth of archived emails or early internet text containing technical manuals, source code, and even full bitmaps. The unicode and lexer parsers can handle these easily, though in nuanced ways (keep reading for an example).

Also, any non-IndoEuropean languages must be tokenized via unicode or with apporpriate lexers as whitespace is not a significant delimiter for atomic units. I have not done thorough experimentation on this. In this regard, using this project for an Arabic parser or Chinese tokenizer may uncover some inconsistencies.

## Get it

```sh
go get github.com/jbowles/nlpt_tkz
```

## Use it

```go
package main

import (
	"fmt"
	tkz "github.com/jbowles/nlpt_tkz"
)

func main() {
	s := "From: mathew <mathew@mantis.co.uk> \nSubject: Alt.Atheism FAQ: Atheist Resources\n\nArchive-name: atheism/resources\nAlt-atheism-archive-name: resources\nLast-modified: 11 December 1992\nVersion: 1.0"
	b := []byte(s)
	digest1 := tkz.TokenizeBytes(b, "lex")
	digest2 := tkz.TokenizeBytes(b, "unicode")
	_, digest3 := tkz.TokenizeStr(s, "lex")
	_, digest4 := tkz.TokenizeStr(s, "unicode")
	_, digest5 := tkz.TokenizeStr(s, "whitespace")

	fmt.Printf("-----printed digest-----")
	fmt.Printf("\n\n")
	fmt.Printf("LexBytes \n %+v\n\n", digest1)  //Notice we are printing out the full struct via +v
	fmt.Printf("UnicodeBytes \n %+v\n\n", digest2)
	fmt.Printf("LexStr \n %+v\n\n", digest3)
	fmt.Printf("UnicodeStr \n %+v\n\n", digest4)
	fmt.Printf("WhitespaceStr \n %+v\n\n", digest5)
	fmt.Printf("---------------------")
	fmt.Printf("\n\n\n")
	fmt.Printf("-----printed bytes-----")
	fmt.Printf("\n\n")

	fmt.Printf("++++++ LexBytes Printed +++++++ \n\n %s\n", digest1.Bytes)
	fmt.Printf("\n\n")
	fmt.Printf("+++++ UnicodeBytes Printed ++++++ \n\n %s\n", digest2.Bytes)
}
```

The program above will produce this:

```sh
-----printed digest-----

LexBytes
 &{Tokens:[] TokenBytes:map[] Bytes:[70 114 111 109 32 58 32 32 109 97 116 104 101 119 32 60 109 97 116 104 101 119 64 109 97 110 116 105 115 32 46 32 99 111 32 46 32 117 107 62 32 10 83 117 98 106 101 99 116 32 58 32 32 65 108 116 32 46 32 65 116 104 101 105 115 109 32 70 65 81 32 58 32 32 65 116 104 101 105 115 116 32 82 101 115 111 117 114 99 101 115 10 10 65 114 99 104 105 118 101 45 110 97 109 101 32 58 32 32 97 116 104 101 105 115 109 47 114 101 115 111 117 114 99 101 115 10 65 108 116 45 97 116 104 101 105 115 109 45 97 114 99 104 105 118 101 45 110 97 109 101 32 58 32 32 114 101 115 111 117 114 99 101 115 10 76 97 115 116 45 109 111 100 105 102 105 101 100 32 58 32 32 49 49 32 68 101 99 101 109 98 101 114 32 49 57 57 50 10 86 101 114 115 105 111 110 32 58 32 32 49 32 46 32 48] SpaceCount:0 CharCount:0 Letter:[] Title:[] Number:[] Punct:[] Space:[] Symbol:[] TokenCount:0 PunctCount:0 LineCount:0 EmptyLine:false LastTokenType:3}

UnicodeBytes
 &{Tokens:[] TokenBytes:map[] Bytes:[70 114 111 109 32 58 32 32 109 97 116 104 101 119 32 32 60 32 109 97 116 104 101 119 32 64 32 109 97 110 116 105 115 32 46 32 99 111 32 46 32 117 107 32 62 32 32 10 83 117 98 106 101 99 116 32 58 32 32 65 108 116 32 46 32 65 116 104 101 105 115 109 32 70 65 81 32 58 32 32 65 116 104 101 105 115 116 32 82 101 115 111 117 114 99 101 115 10 10 65 114 99 104 105 118 101 32 45 32 110 97 109 101 32 58 32 32 97 116 104 101 105 115 109 32 47 32 114 101 115 111 117 114 99 101 115 10 65 108 116 32 45 32 97 116 104 101 105 115 109 32 45 32 97 114 99 104 105 118 101 32 45 32 110 97 109 101 32 58 32 32 114 101 115 111 117 114 99 101 115 10 76 97 115 116 32 45 32 109 111 100 105 102 105 101 100 32 58 32 32 49 49 32 68 101 99 101 109 98 101 114 32 49 57 57 50 10 86 101 114 115 105 111 110 32 58 32 32 49 32 46 32 48] SpaceCount:0 CharCount:0 Letter:[] Title:[] Number:[] Punct:[] Space:[] Symbol:[] TokenCount:0 PunctCount:0 LineCount:0 EmptyLine:false LastTokenType:0}

LexStr
 &{Tokens:[From mathew <mathew@mantis co uk> Subject Alt Atheism FAQ Atheist Resources Archive-name atheism/resources Alt-atheism-archive-name resources Last-modified 11 December 1992 Version 1 0] TokenBytes:map[Last-modified:[76 97 115 116 45 109 111 100 105 102 105 101 100] 1992:[49 57 57 50] From:[70 114 111 109] mathew:[109 97 116 104 101 119] .:[46] atheism/resources:[97 116 104 101 105 115 109 47 114 101 115 111 117 114 99 101 115] Alt-atheism-archive-name:[65 108 116 45 97 116 104 101 105 115 109 45 97 114 99 104 105 118 101 45 110 97 109 101] resources:[114 101 115 111 117 114 99 101 115] 11:[49 49] 1:[49] <mathew@mantis:[60 109 97 116 104 101 119 64 109 97 110 116 105 115] Alt:[65 108 116] Atheism:[65 116 104 101 105 115 109] FAQ:[70 65 81] Archive-name:[65 114 99 104 105 118 101 45 110 97 109 101] December:[68 101 99 101 109 98 101 114] Version:[86 101 114 115 105 111 110] ::[58] uk>:[117 107 62] Subject:[83 117 98 106 101 99 116] Atheist:[65 116 104 101 105 115 116] co:[99 111] Resources:[82 101 115 111 117 114 99 101 115] 0:[48]] Bytes:[70 114 111 109 58 32 109 97 116 104 101 119 32 60 109 97 116 104 101 119 64 109 97 110 116 105 115 46 99 111 46 117 107 62 32 10 83 117 98 106 101 99 116 58 32 65 108 116 46 65 116 104 101 105 115 109 32 70 65 81 58 32 65 116 104 101 105 115 116 32 82 101 115 111 117 114 99 101 115 10 10 65 114 99 104 105 118 101 45 110 97 109 101 58 32 97 116 104 101 105 115 109 47 114 101 115 111 117 114 99 101 115 10 65 108 116 45 97 116 104 101 105 115 109 45 97 114 99 104 105 118 101 45 110 97 109 101 58 32 114 101 115 111 117 114 99 101 115 10 76 97 115 116 45 109 111 100 105 102 105 101 100 58 32 49 49 32 68 101 99 101 109 98 101 114 32 49 57 57 50 10 86 101 114 115 105 111 110 58 32 49 46 48] SpaceCount:19 CharCount:193 Letter:[] Title:[] Number:[] Punct:[: . . : . : : : : : .] Space:[] Symbol:[] TokenCount:22 PunctCount:11 LineCount:7 EmptyLine:false LastTokenType:3}

UnicodeStr
 &{Tokens:[From mathew mathewmantiscouk  Subject AltAtheism FAQ Atheist Resources  Archivename atheismresources Altatheismarchivename resources Lastmodified  December  Version ] TokenBytes:map[] Bytes:[] SpaceCount:0 CharCount:0 Letter:[F r o m ,  m a t h e w ,  m a t h e w m a n t i s c o u k ,  ,  S u b j e c t ,  A l t A t h e i s m ,  F A Q ,  A t h e i s t ,  R e s o u r c e s ,  ,  A r c h i v e n a m e ,  a t h e i s m r e s o u r c e s ,  A l t a t h e i s m a r c h i v e n a m e ,  r e s o u r c e s ,  L a s t m o d i f i e d ,  ,  D e c e m b e r ,  ,  V e r s i o n , ] Title:[] Number:[1 1 1 9 9 2 1 0] Punct:[: @ . . : . : - : / - - - : - : : .] Space:[] Symbol:[< >] TokenCount:0 PunctCount:0 LineCount:0 EmptyLine:false LastTokenType:0}

WhitespaceStr
 &{Tokens:[From: mathew <mathew@mantis.co.uk>
Subject: Alt.Atheism FAQ: Atheist Resources

Archive-name: atheism/resources
Alt-atheism-archive-name: resources
Last-modified: 11 December 1992
Version: 1.0] TokenBytes:map[] Bytes:[] SpaceCount:13 CharCount:193 Letter:[] Title:[] Number:[] Punct:[] Space:[] Symbol:[] TokenCount:0 PunctCount:0 LineCount:0 EmptyLine:false LastTokenType:0}

---------------------


-----printed bytes-----

++++++ LexBytes Printed +++++++

 From :  mathew <mathew@mantis . co . uk>
Subject :  Alt . Atheism FAQ :  Atheist Resources

Archive-name :  atheism/resources
Alt-atheism-archive-name :  resources
Last-modified :  11 December 1992
Version :  1 . 0


+++++ UnicodeBytes Printed ++++++

 From :  mathew  < mathew @ mantis . co . uk >
Subject :  Alt . Atheism FAQ :  Atheist Resources

Archive - name :  atheism / resources
Alt - atheism - archive - name :  resources
Last - modified :  11 December 1992
Version :  1 . 0
```


## Note
Only the Lexer for `TokenizeStr` returns metadata such as line, token, punctation count.
There is no special reason for this, nor is there any rationale behind the different ways metadata is collected. Since nobody I know, including me, is using these tokenizers I've been reluctant to lock down any real consistency between the metatdata collection as the usefulness of metatada shoudl be determined by usage.

## Tokenizers so far...
The white space tokenizer is merely a wrapper around `strings.Split(" ")` with some digest content for counts of tokens and such.

The `lex` (technically a State Function Lexer) and `unicode` (technically using go's `unicode` package for matching unicode code points) **can and probably will** return very differnt tokens sets of tokens (a token is defined as a unique atomic unit). For example, and email subject heading with the Lexer tokenizer will return one token `\<mathew@mantis` but the Unicode tokenizer will return `mathewmantiscouk`. For most people or tasks these differences don't matter. The rule of thumb is that if you just need "words", i.e. atomic units, and have a text that is pretty clean use `unicode`; it is the fastest of the two.


## Benchmarks

```sh
go test -bench .

go test -bench=.

go test -bench Lex

## or use the project go bench ui
gobenchui .
gobenchui github.com/jbowles/nlpt_tkz
```

```sh
PASS
BenchmarkStateFnTknzGoodStr-8       	   10000	    112024 ns/op
BenchmarkStateFnTknzBadStr-8        	   20000	     63397 ns/op
BenchmarkStateFnTknzBytesGoodStr-8  	   20000	     87202 ns/op
BenchmarkStateFnTknzBytesBadStr-8   	   30000	     49868 ns/op
BenchmarkLexStrGood-8               	   10000	    111359 ns/op
BenchmarkUnicodeStrGood-8           	  100000	     17056 ns/op
BenchmarkWhitespaceStrGood-8        	 1000000	      1324 ns/op
BenchmarkLexBytesGood-8             	   20000	     88690 ns/op
BenchmarkUnicodeBytesGood-8         	  300000	      5984 ns/op
BenchmarkLexStrBad-8                	   20000	     63092 ns/op
BenchmarkUnicodeStrBad-8            	  200000	     11048 ns/op
BenchmarkWhitespaceStrBad-8         	 2000000	       822 ns/op
BenchmarkLexBytesBad-8              	   30000	     50150 ns/op
BenchmarkUnicodeBytesBad-8          	  500000	      3726 ns/op
BenchmarkUncdMatchTknzGoodStr-8     	  100000	     17184 ns/op
BenchmarkUncdMatchTnkzBadStr-8      	  200000	     10954 ns/op
BenchmarkUncdMatchTknzBytesGoodStr-8	  300000	      5682 ns/op
BenchmarkUncdMatchTnkzBytesBadStr-8 	  500000	      3401 ns/op
BenchmarkWhiteSpaceTknzGoodStr-8    	 3000000	       536 ns/op
BenchmarkWhiteSpaceTknzBadStr-8     	 2000000	       706 ns/op
ok  	github.com/jbowles/nlpt_tkz	39.137s
```

