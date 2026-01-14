# Backend Overview
**WELCOMEEEE TO ISNIDE ATHLETICS** Here is a general overview of our repo, some of our design decisions, and how to use it. We REALLLY want you guys to be able to learn as much as you can so please feel free to ask questions if you have any! About go, backend structure, ANYTHINGGGG! 

## Tech Stack

### GO
Our backend is written in Go! Go is a very fast language that has a lot of built in features and libraries that simplify the backend development environment HELLA!

Some libraries that we have choosen to implement

- **Routing**: Fiber + Huma. Using Huma we are able to automate the creation of API documentation as you're coding your endpoints! It also allows for insanely easy parameter validation and more!
- **DataBase**: For our database we are using Supabase! Supabase is a super easy to interact with database where you can view your data, edit tables, etc. It's a postgres database.
- **GORM**: To communicate with this database we are using an ORM called GORM. GORM simplifies the endpoint creation process by communicating with our DB for us, essentially, you won't have to code any giant SQL strings when you wanna talk to Supabase.
- **Atlas**: Tbh I personally don't know much but it helps us with DB migrations! Section will be more detailed soon...

## Organization
Here we will go over what each section of the backend means, as well as the purpose of each section. 

### cmd 
Here the main file lives. This is where we run the backend. This main functions sets up our DB connections, starts the server, handles shutdown, etc...

### config
For now this empty... As Zainab and I create configuration files to deploy the backend this will get populated (aka you guys don't rlly need to look into this)

### internal
This is mainly where you guys will live. Here we initialize all of the endpoints required for our app to function. In this folder we have the following directories: 

- handlers: Endpoint logic including database transactions and services
- migrations: Migration logs for our DB, as we make migrations they will be automatically stored here using Atlas!
- models: Where we store all of the internal representations of our data in the system! Any model you need to make (user, post, etc) will be defined in here
- server: Sets up the routes and server with Huma + Fiber.
- tests: Testing directory for unit tests and integration tests!
- utils: Utility functions (for abstraction, code used across multiple packages, etc..)

## Setup

### Installing Go
If you haven't already install go on your computer!

### env file
In the backend directory create a .env file! We will send the contents of the file during the meeting

## How to Run
``` bash
cd backend
go make start-prod
```

The server should then start and you will be able to ping some endpoints! Try `/api/v1/health/`

NOTE: This is a temporary thing... once the dev scripts are in you won't have to do this

## Implementing an Endpoint
Now for the part that really matters. MAKING AN ENDPOINT! I will be going over the purpose of each part of the endpoint to show why they matter as well if you were curious!

Also, if you haven't already, download the go extension in VS Code. It will make your life easier.

The structure of making an endpoint is pretty simple. Here's the breakdown:

 - **Creating the Routes**: Mapping the routes from the website path to the actual functionality
 - **Creating the service**: This is the "Business Logic" of the endpoint. Here you will be creating the logic for what the endpoint actually does. This typically looks like getting data from Supabase via DB struct, transforming it in some way (applying functions to the data to get a certain output or something), and returning it in a way that Huma can understand
 - **Database Transaction**: This is where the code interacts with Supabase. Here you will utilize the `db` variable (our connection to the prod database) in order to query, update, or delete data that we store.

### Setting up File Structure
We are going to start by laying out an endpoint file wise

Navigate to the backend/internal/handlers directory, this is where all of the logic for an endpoint is going to be stored!

Lets start by creating a folder called `goat` this will store all of the respective files for our endpoint

Next lets create the files for the endpoint. Every endpoint is going to be structured the same, so take note! Create the following files in the goat directory:

``` go
- goat_db.go // database transaction file
- goat_route.go // routing function file
- goat_service.go // service functions (business logic, ima keep saying this everywhere Zahra)
- goat_types.go // Endpoint types
```

### Creating a Route, `goat_route.go`
Honestly the easiest part of creating an endpoint! Here we establish the route paths and the functions they map to!

The structure of these route paths is simple, the base URL (for now http://127.0.0.1:8080/) with the extension to the endpoint you want to hit. For example, if I want to hit the health check endpoint it would be, http://127.0.0.1:8080/api/v1/health/

If this is completely new we are essentially creating a path to the function in our code that web browsers can use via their URLs.

Now we will define the routes for our goat endpoints in the `goat.go` file

``` go
package goat // take note that the package matches the directory name

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {

}
```

Your route definitions **MUST** follow this signature otherwise you won't be able to add the routes (the name can be GoatRoutes if you wanted but I think this is easier)

Next we will add some logic to the route using the `api` variable 

``` go
func Route(api huma.API, db *gorm.DB) {
	{
		grp := huma.NewGroup(api, "api/v1/goat") // this creates a goat route group!
	}
}
```

A Route group allows us to organize our endpoints in a structure similar to file systems (think of the route group /goat as the parent directory). The route group will store all of the other paths for the goats functionality

Next we will add our first actual endpoint!

``` go
func Route(api huma.API, db *gorm.DB) {
	{
		grp := huma.NewGroup(api, "api/v1/goat")
		huma.Get(grp, "/", nil) // pass the group into the huma get method
	}
}
```

This will (kind of) create your first endpoint! Now your huma.Get call should be red... That's because we need to pass in a `function` for huma to call when this route is hit. As of now, we haven't defined anything so we can't pass anything yet...

Lastly, navigate to the backend/internal/server.go file

In the CreateRoutes() function add your goat.Route function to the list of route functions

``` go
func CreateRoutes(db *gorm.DB, api huma.API) {
	// Create all the routing groups:
	routeGroups := [...]RouteFN{health.Route, user.Route, goat.Route} // Here!!!
	for _, fn := range routeGroups {
		fn(api, db)
	}
}
```

So the next logical step is to define some `Business Logic` in the service file!!!!

### Implementing a service, `goat_service.go`
The service file has 2 main parts. Defining the struct of the business logic functions belong to and the business logic.


Lets start by creating the `GoatService` struct

``` go
package goat

type GoatService struct { // blank for now... we will add something later, stay hyped
}
```

The GoatService struct will hold all of the goat endpoint functionality

Lets create a simple service

``` go
func (g *GoatService) Ping(ctx context.Context, input *utils.EmptyInput) (_, error) {

}
```
Now your function should be highlighted red, yet again... The reason is we have yet to define a return type for our service. When we use Huma, it's important to define structs for our input and output types. As you can see, since this route isn't taking anything in as an input we are using our utils.EmptyInput struct (which is just a blank struct)

The `input` variable to this function represents the route parameters associated with this route. In the future this will be a struct that actually has values mapping to the path params it takes in. Stay tuned for the next example!

Lets create our first response type for the service!

Navigate to the `goat_types.go` file and define this struct

``` go
package goat

type GoatResponse struct {
	Id   int8   `json:"id" example:"1" doc:"id of this goat"`
	Name string `json:"name" example:"Suli" doc:"Name of this goat"`
	Age  string `json:"age" example:"67" doc:"Age of this goat"`
}
```

The `json` tag represents a mapping of the given variable to how it would be represented in JSON. This is used by many libraries to map responses in json to our local struct definitions. The `example` and `doc` tags are for huma to create API documentation for the response (we will take a look at this later). 

Now we can go back into our `goat.service` file and define the return type!

``` go
func (g *GoatService) Ping(ctx context.Context, input *utils.EmptyInput) (*utils.ResponseBody[GoatResponse], error) {

}
```

Now this endpoint returns a GoatResponse and error!

You may be wondering why we have to return a ResponseBody of type GoatResponse (*utils.ResponseBody[GoatResponse]). The reason for this is Huma expects a specific struct format to determine what the body of the response is, hence we wrap our response in this struct

A quick note about how errors are handled in go. Errors are handled on a value basis. This allows you to return errors as values and have them handled accordingly rather than terminating your entire code as it runs. We will take a look at this flow later, in this example nothing should break...

Now lets actually define some logic! 

Add the following line to the `Ping` function:

``` go
func (g *GoatService) Ping(ctx context.Context, input *utils.EmptyInput) (*utils.ResponseBody[GoatResponse], error) {
	resp := &GoatResponse{
		Id:   1,
		Name: "Suli",
		Age:  67,
	}
	return &utils.ResponseBody[GoatResponse]{
		Body: resp,
	}, nil
}
```
First we create a hardcoded GOAT, and then we return it. Simple!

As you can see, for the error we are returning nil (null value in go) for the error because there is no error value!

Navigate to `goat_service.go` and add this function to your endpoint

``` go
func Route(api huma.API, db *gorm.DB) {
	goatService := GoatService{} // define an instance of the goat service
	{
		grp := huma.NewGroup(api, "api/v1/goat")
		huma.Get(grp, "/", goatService.Ping) // add functionality of the endpoint!
	}
}
```

Now we can actually run this!

From the backend directory call `make start-prod`. Open the local host in your browser and go to the http://localhost:8080/api/v1/goat/ url (if you have postman I highly recommend just using that instead and calling it from there)

You should see the GoatResponse we created!

### Creating a model `backend/internal/models/goat.go`
For more complex queries (i.e getting stuff back from the DB) we need to define some models to match our database schemas!

Lucky for you guys I've created a goat model with one small thing to add

Navigate to the backend/internal/models folder

This is where we will store all of our models for the project! This folder represents all of the database models we have translated into go to use within the codebase!

Peep the `goat.go` file, here we have a struct defined for a goat.

``` go
type Goat struct {
	gorm.Model // Id, created_at, deleted_at, updated_at
	Name string `json:"name" example:"Suli" doc:"The name of a goat"`
	Age  int8   `json:"age" example:"67" doc:"The age of this goat"`
}
```

After you've done that we can jump into the awesome part...

### Database transaction `goat_db.go`
This is where we will interact with the database to get live information!

First we start by creating the struct to attach all of our database functionality onto

``` go
package goat

import "gorm.io/gorm"

type GoatDB struct {
	db *gorm.DB // our connection to Supabase
}
```

Next we will define a simple `GET` transaction from our goat table!

``` go
func (g *GoatDB) GetGoat(id int8) (*models.Goat, error) {

}
```

Notice here how we are returning the `models.Goat` and not the GoatResponse we defined earlier. This is because the GoatResponse struct will only return data that we want the user of the endpoint to see. For example, the models.Goat structure has a created_at, deleted_at, updated_at, etc... These variables might not be needed to be seen by someone who is actually calling our endpoint and they only exist for our internal usage. Therefore we don't return them and only return a subset of the values through our endpoints (ofc this can be updated and def will be updated throughout the sem).

Next lets actually make the call to our database!

``` go
func (g *GoatDB) GetGoat(id uint) (*models.Goat, error) {
	var goat models.Goat
	dbResponse := g.db.Where("id = ?", id).First(&goat) // gets the entry with the given id and copies the data into goat
	return utils.HandleDBError(&goat, dbResponse.Error) // Error handling utility function!!!
}
```

CONGRATS! This is your first DB transaction!!!

Now we want our service to have access to these database functions

In the `goat_service.go` file add the following line to the struct:

``` go
type GoatService struct {
	goatDB *GoatDB // give the service a DB connection!
}
```

Now in the `goat_route.go` file add the db to your service

``` go
func Route(api huma.API, db *gorm.DB) {
	goatDB := &GoatDB{db: db} // create instance of the db!
	goatService := GoatService{goatDB: goatDB} // add the transactions to the service! 
	{
		grp := huma.NewGroup(api, "api/v1/goat")
		huma.Get(grp, "/", goatService.Ping)
	}
}
```

## ANOTHA ENDPOINT!
Now that you are familiar with the file structure of the codebase it's time to implement some real functionality!

We just defined a DB transaction, but none of our endpoints use it... 

Lets fix that!

This endpoint will implement one more step of logic!

### Create the route to the endpoint
Just like the previous endpoint go to your `goat_route.go` endpoint and create a new goat route

``` go
func Route(api huma.API, db *gorm.DB) {
	goatDB := &GoatDB{db: db}
	goatService := GoatService{goatDB: goatDB}
	{
		grp := huma.NewGroup(api, "api/v1/goat")
		huma.Get(grp, "/", goatService.Ping)
		huma.Get(grp, "/{id}", _) // NEW ENDPOINT!
	}
}
```

See how this endpoint ends in {id}, this is called a `path param` and it allows us to pass in variables, similar to a function, into our endpoint! 

Now we need to define the service functionality for our endpoint

In the `goat_service.go` file: 

``` go
func (g *GoatService) GetGoat(ctx context.Context, _) (utils.ResponseBody[GoatResponse], error) {
	
}
```

Now... how do we get access to that path parameter we defined earlier?

That would be through the input struct! All we need to do is create a type to represent the input data for the path and sub it in for the _

In the `goat_types.go` file add the following struct

``` go
type GetGoatParams struct {
	Id int8 `path:"id" example:"1" doc:"id of this goat"`
}
```

NOTE: that we are using the tag path:"id" to tell huma to map the param ID to this variable in the struct

In the service we can now add: 

``` go
func (g *GoatService) GetGoat(ctx context.Context, input *GetGoatParams) (*utils.ResponseBody[GoatResponse], error) {

}
```

Now we can access the Id param that we pass in through the path!

``` go
func (g *GoatService) GetGoat(ctx context.Context, input *GetGoatParams) (*utils.ResponseBody[GoatResponse], error) {
	resp := &utils.ResponseBody[GoatResponse]{}
	id := input.Id
	goat, err := g.goatDB.GetGoat(id)

	if err != nil { // if we have an err getting the goat return the empty response and the error
		return resp, err
	}

	resp.Body = &GoatResponse{
		Id:   int8(goat.ID),
		Name: goat.Name,
		Age:  goat.Age,
	}

	return resp, nil
}
```

Now go back into the route and add the service:

``` go
func Route(api huma.API, db *gorm.DB) {
	goatDB := &GoatDB{db: db}
	goatService := GoatService{goatDB: goatDB}
	{
		grp := huma.NewGroup(api, "api/v1/goat")
		huma.Get(grp, "/", goatService.Ping)
		huma.Get(grp, "/{id}", goatService.GetGoat) // ADD THE FUNCTIONALITY
	}
}
```

YOUVE MADE A ROUTE THAT CAN TAKE IN INPUTTTT

There is wayyy more to learn and as we start making routes you will encounter the other cases but this is a general overview!



