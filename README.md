This is my solution to http://golang-challenge.org/go-challenge1/

Initial commit from
https://github.com/golangchallenge/golang-challenge/blob/gh-pages/data/ch1/golang-challenge-1-drum_machine.zip

Impressions on my First go code

* I think it is a bit tricky the first challenge, as the part of figuring out the file format, needs fluent progmamming skills. I personally used Python to reverse engineer the format, after trying clumsily few hours with Go
* It feels painful after coming from Python (no types, list comprehensions) to need 3 variables and 6 lines to read a string for a byte stream of size x (to get the track name). I hope it can be rewritten in a simpler way
* Error handling imho if this was a library it would be good to raise library spesific errors, to do that in Go seems we need to add ifs in all the calls that have errors??!? Or am I missing sth?

