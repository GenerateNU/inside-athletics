# Backend Overview
**WELCOMEEEE TO ISNIDE ATHLETICS** Here is a general overview of our repo, some of our design decisions, and how to use it. 

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
cd backend/cmd
go run .
```

The server should then start and you will be able to ping some endpoints! Try `/api/v1/health/`

## Implementing an Endpoint
In order to get the whole repo to come together 

### Creating a route

### Creating a model

### Implementing a service

- 

### Database transaction

### 

