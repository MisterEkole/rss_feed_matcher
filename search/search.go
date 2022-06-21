package search

import (
	"log"
	"sync"
)

//creating a map for registered matcher for searching purpose
var matchers = make(map[string]Matcher)

//fun Run to perform search logic

func Run(searchTerm string) {
	//getting the list of registered feeds
	//log.Println("Searching for", searchTem)
	feeds, err := RetrieveFeeds()
	if err != nil {
		log.Fatal(err)
	}

	//unbuferred channel to receive matcher
	results := make(chan *Result)

	//Wait group
	var waiter sync.WaitGroup

	//set number of goroutines to wait for
	waiter.Add(len(feeds))

	for _, feed := range feeds {

		matcher, exist := matchers[feed.Type]

		if !exist {
			matcher = matchers["default"]

		}
		//for each feed start a goroutine to search
		go func(matcher Matcher, feed *Feed) {
			//perform the search
			match(matcher, feed, searchTerm, results)
			//decrement the wait group
			waiter.Done()
		}(matcher, feed)
	}
	go func() {
		waiter.Wait()
		close(results)
	}()
	Display(results)
}

func Register(feedType string, matcher Matcher) {
	if _, exist := matchers[feedType]; exist {
		log.Fatal("matcher already registered")
	}
	log.Println("Registering matcher", matcher, "for type", feedType)
	matchers[feedType] = matcher
}
