# Spelling Corrector

You can run the script by using the command `go run script.go`.

## Credits

This algorithm's originally written in Python by [Peter Norvig](https://norvig.com/spell-correct.html), transcripted to Go Lang by [Yi Wang](https://cxwangyi.wordpress.com/2012/02/15/peter-norvigs-spelling-corrector-in-go/) and adapted by me.

The list of common misspellings is [powered by Wikipedia](https://en.wikipedia.org/wiki/Wikipedia:Lists_of_common_misspellings/For_machines) and also adapted by me.

## Results

I got a total of **1303 errors** and **2816 successful** cases (68.36611% success) in 175 seconds of processing.

## Future Work

~~In order to make the script run faster, it would be interesting to use something like BK-Trees or Trie, instead of iterating through a plain text file.~~

In order to improve the successful cases ratio, we could fill our dictionary files with new words. We could take context into consideration, as well. Problems like `Expect yearm->year but got years` can only be solved with context.