# Testing!!!
I know this isn't everyone's favorite part of coding... And probably most people's least favorite... BUTTTT it's super important for maintaining a long term codebase

Creating tests (unit & integration) is super important as you change things within your codebase. By having a foundation of tests, you can make sure that any new changes that you implement doesn't break any old code, thus maintaining the functionality of our application.

Lucky for you guys, our tests will be relatively simple! And if you do them as you create your endpoints, it really won't be hard to keep up with the testing while also ensuring our app doesn't break somewhere down the line.

We will be going over some examples on how to write integration and unit tests!

## General Testing Rules
In go testing is built into the language. Which means that it's super easy to use!

All you have to do to make a test is end the filename with _test.go and start all of your test functions in the given class with Test<Test_Name>

The function also needs to take in a `*testing.T` object, which will let you report back errors to go when a test fails.

See references in the test files for usage!

## Integration Tests
This will probably be the test that most of yall will be using as it allows you to test your endpoints locally without relying on our production database. 

This process is called "mocking" and most of you should be familiar with this from OOD.

We are essentially creating dummy databases and API's that still function the exact same way our production code would, but instead of being connected to the real data (aka stuff you can break and mess up) it's connected to a local database that harms NOTHING if anything goes wrong. Through this you can create endpoints and fully test them without worrying about anything else breaking.

Note, there are many things that can break... From migrations to updating internal representations of data, there is much room for error. So it's SUPERRR important to test your shit to make sure that it doesn't go wrong. 

### Where they at? `tests/route_tests`
This is where all of the integration tests are stored. For each route grouping there is a `<name_of_route_group>_test.go` file. This is where your integration tests are going to go.

And honestly... The setup from there is pretty easy (aka we have infastructure to setup the test env for you) SO PLEASE USE IT IT'S SO EASY

I'll use the `health_test.go` file as an example to show yall

``` go
package routeTests

import (
	"inside-athletics/internal/handlers/health"
	"strings"
	"testing"
)

func TestGetGreeting(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t) 
	api := testDB.API // get access to the API from test DB

	resp := api.Get("/api/v1/health/") // MAKE CALL TO ENDPOINT!!!

	var health health.HealthResponse // create struct to parse data into

	DecodeTo(&health, resp) // SIMPLE DECODE FUNCTION TO CREATE STRUCT!

    // CALL AND MAKE SURE YOUR DATA IS WHAT YOU EXPECT
	if !strings.Contains(health.Message, "Welcome to Inside Athletics API Version 1.0.0") {
		t.Fatalf("Unexpected response: %s", resp.Body.String()) // report error back if we get something unexpected
	}
}
```

I'll go over the key components as these will be the same in literally EVERY test. Wow what amazing and easy test infrastructure that we should all totally use......

- `SetupTestDB`: This call is going to setup your local Postgres DB!!! This is an actual fully functioning DB that is created in docker and destroyed when the test stops running! All you have to do to make it is call this method and it'll handle everything for you
- `defer testDB.Teardown(t)`: The defer keyword in go is a call that is made when the go function stops running. It serves as "cleanup" essentially this tears down the DB connections and Docker instance so it doesn't exist forever. MAKE SURE TO CALL THIS!
- `Decode()`: This function takes the bytes returned from the call and parses it into the actual structs that we expect them to be. If you take a look at the return type of api.Get() its a http.ResponseRecorder and not the struct that you actually return from the DB. Buttt by simply calling this func and passing in a struct to write to you can easily extract the information you need.

Now I know that it seems a little silly and that this example is small... BUT, once things get more complicated it's important that we have measures in place to track the funcionality of all our things. 

And the best part, LOOK HOW EZ THIS WAS. If you take a look at the `user_test.go` file, even for a more complicated endpoint (that I literally made another endpoint inside of) ITS STILL SO SHORT AND EZ TO WRITE. We did this on purpose. Please please please, test your stuff it's so easy.

## Unit Tests `tests/unit_tests`
This one is a lot simpler and functions almost identical to unit tests in other languages! 

It honestly follows the same thing as the integration tests but instead of calling it on a route (no need to setup any test DBS for these) you're just testing functions that you custom define for models within our project.

I think the best way to learn for this one is just to read through a simple unit test! Peep `tests/unit_tests/goat_test.go`