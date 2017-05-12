## README

## Primary parts

* randomUrl Generator.
	* randomize numbers and use the logic to return charaters.
	* initial version could have only numbers extend if time permits.

* URLFactory
	* This is the type which handles the Url and all its assosciated actions.
	* Generate random url, this is the

* Handlers
        * GET /{shortURL}
        	* Check in couchdb if the value exists.
		* Permantent redirect if it exists.
	* POST /add/{longURL}
		* Uses the randomURL generator to create a tiny url.
		* Adds the longURL to the the URLFactory.
		* Returns with a JSON of

## RUN

```

Extract the source code here as `go get` does not work.

$ go get github.com/georgethomas111/urlshortner
$ cd $GOPATH/src/github.com/georgethomas111/urlshortner
$ cd main
$ go build
$ ./main
```

The backend database used is couchdb. So change address in main.go


## TESTS

# with curl

```
$ curl -X GET 127.0.0.1:8080/abcd
$ curl -X POST 127.0.0.1:8080/add/google.com
```

# Go tests

```
$ cd $GOPATH/src/github.com/georgethomas111/urlshortner/random
$ go test -v
# Tests for random

$ cd $GOPATH/src/github.com/georgethomas111/urlshortner
$ go test -v
# Tests to test the factory i.e interaction with the database.
```
		
## FIXME
handler.go:// FIXME urlLength should be taken from the config
handler.go:// FIXME Establish the right method to log errors.
handler.go:// FIXME The arguments to the function Start should be taken in the config.

## ISSUES

 * The library used `github.com/gospackler/caddyshack-couchdb` has errors. I used it as I have personally worked on it. I currently use a chromebook so I was not able to get couchdb runnning locally. So the code is not tested with a real DB. The compilation works and the code may execute in one go or with a small tweaks with the library.
 
 * The ORM which is based on reflection uses the the `by` tag and creates a view with it as the key. This makes sure views are created without explicitly specifying the needs for creating it.
 

## Requirements

```
A Short Url:

1. Has one long url

>>
# As the ADD method always creates a new url, once a randomURL is generated, it is checked in the factory to see if it exists. If it already exists, another random url is generated. The complexity of the conflict is low based on the sample space chosen.
# Chances of conflict is < 0.01 % as the random function produes 10^9 possible urls and this can be adjusted with the length passed to the random generator.

2. Permanent; Once created

>>
# This is true as it gets added to the database with the key and if it already exists another url is chosen.

3. Is Unique; If a long url is added twice it should result in two different short urls.

>> This is true. As mapping is one to one.

Random is based on `crypto/rand` or `/dev/urandom`

4. Not easily discoverable; incrementing an already existing short url should have a low

>>
#  Random is based on `crypto/rand` or `/dev/urandom` in Unix. So this secure and incrementing a value is not going to give another value. The sample space for the url is huge such that the chances of conflict is low even for million users. 10^6/10^9


Your solution must support:

1. Generating a short url from a long url
>> POST endpoint takes care of this.

2. Redirecting a short url to a long url within 10 ms.
>> http Redirect takes care of this with a call to the view of a database which should return the data quickly.

3. Listing the number of times a short url has been accessed in the last 24 hours, past

>> Not implemented but the infrastructure exists and can be implemented through logs in the handler or writing the data about access to the database in parallel.

4. Persistence (data must survive computer restarts)

>> This is true as the servers are stateless and the state exists in the database. Even if an application crashes, another application pointing to the same data should do the trick.
