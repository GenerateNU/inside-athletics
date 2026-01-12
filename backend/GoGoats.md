# GO BASICS
Hello Go Goats. I understand that most of our team is not familiar with Go. While I am also not an expert, this is what I've been able to learn about Go so far that I think will be useful when you try to cook on this team. It isn't a perfect guide but a general overview/topics to look out for while you're coding. If you have any questions always feel free to ask! 

# Variables
Starting off with the basics! There are 2 ways to define a variable in go, inferred type and explicit type. Both are lit! And pretty easy 

## Inferred types
An inferred type in go means that you don't need to explicitly declare what the variable type is. In order to do this you declare it like...
``` go
    package example

    age := 7 // age is an integer
    name := "Dorchester" // string 
    bbl := true // bool
```
This same concept applies to any variable whether it's a struct, primitive, or function. Anything can be stored like dis. **BUT** if you aren't assigning the value right when you declare the variable you cannot use this method of declaration. I.E

``` go
goat // will break your code and won't count as a variable
```
You can only use inferred types within functions or other scopes, these variables cannot be used to declare gloabl variables.
However there is a way to do this...


## Explicit types
If you want to be a little fancy, or well document your stuff, you can explicitly declase what type your variable will be using the `var` keyword
``` go
    package example 

    var age int = 7 // the numba seven
    var ageAgain = 7 // you can also declare it like this

    // you can declare a variable to be assigned later 
    var ageOneMoreTime int // variable will be assigned it's "default" value (0 for int, "" for string, etc)
    ageOneMoreTime = 67
```
    Notice how you don't include the : in front of the equals! You can also declare variables as local or global variables using this method! 

## Arrays: 
Classic, arrays in go are fixed sizes. Here's how to declare em

``` go
package arrayLearning

var a [5]int // array of length 5. Each index initialized to 0
b := [5]int{1,2,3,4,5} // initialize the values of the array
c := [...]int{1,2,3,4,5} // compiler assumes length of the array 
d := [4]int{1 : 6, 3 : 9} // assigns specific indices. Output: [0, 6, 0, 9]

// referencing the value at an index
d[1] // 6

// get the length of an array 
len(a) // 5
```

## Lists = Slices
Lists in go are called slices. It's a little confusing because the syntax is almost the same as making an array. However, slices can be dynamically sized so you can add shit to them

``` go
package pizza // ba dum chh

slice := []int {1,2,3} // slize with values [1,2,3]
slice.append(4)

// you can "slice" them to get a sub-slice, not inclusive
sub := slice[1:3] // 2, 3 

// get the length of a slice
len(sub) // 2
```

## Maps (HashMaps/Dictionaries)

``` go
colors := map[string]string{
    "red": "#FF0000", 
    "blue": "#0000FF"
}
```

# Private v Public
In go you can't declare something private or public... Instead something is made private or public based on if it starts with a capital letter or not. Some examples to follow...

# Iteration
ITERATION EXAMPLES!

```go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}

i := 0
for i < 5 { // WHILE LOOP
    i++
}

// iterate through a collection (slices, maps)
goats := []string{"Suli", "LeBron", "Ur mom"}
for i, value := range goats {
    // i, the index
    // value, value at that index
}

// iterate over a map 
goatMap := map[string]string {"Goat1": "Suli", "Goat2": "LeBron"}
for key, value := range goatMap {
    // whatevah
}
```

# Structs
In go, there are no classes. There are only structs. Structs are pretty much classes but they don't have constructors but they serve the same purpose as classes in other languages allowing you to store data within them and create functions to operate on the struct

## How to make a struct
Honestly pretty simple:

``` go
package goat // we will get to this later

type Goat struct {
    Name        string      `json:"name"`
    secrets     []string    `json:"secrets"` // a slice
    secretArray [...]string `json:"secretArray"`
    // etc...
}

// creating an instance of a struct

goat := Goat{} // you can initialize the struct to store nothing 
var another Goat = Goat{Name: "Suli", secrets: []string{"loves milk", "hates water"}}
```
Here we can see that the struct name "Goat" is capital, meaning that the struct is available as a public class outside of it's package. When we take a look at the Name and secrets variables Name is public while secrets is private. 

Now something you might have noticed is this little json thingy... It lowk makes go lit asf. These are called `tags` and their super useful. They provide libraries with metadata needed to process the field. For example, if I were using a JSON library in go, I wouldn't have to map every field in the struct in a dictionary (like in python) as the library can use the tag to handle the case itself.

We will be using these tags across our codebase to simplify API calls as well as create API documentation!

## How to create a Generic Struct
A Generic Struct is a struct that is capable of "filling in" an abstracted type for an object or function. You would use this if you were making a struct/function that functions regardless of the datatypes it stores. The most basic example of this is an ArrayList in Java, since the ArrayList itself doesn't need to care about the type of data that it's storing. All of the functionality of an ArrayList can be done with the stored type in the list as a generic type. You simply pass in the type of data you want to store! 

We use this a couple of times throughout the codebase so take note of the syntax! 

``` go
package mysteriousGoat

type GenericGoat[T any] struct {
    MysteryVariable T
}

// initialize a generic
var goat GenericGoat[int] = GenericGoat{MysteryVariable: 67}

var anotherOne GenericGoat[string] = GenericGoat{MysteryVariable: "Wow generics are epic"}
```

Here I create multiple goat structs. By passing in the type of data I want the goat to store I can initialize the MysteryVariable's type to anything! The goat variable is initialized with an int MysteryVariable, while the anotherOne is initialized to a string.

Lets take the Goat struct from earlier 

# Functions
Another classic and must have! I'm sure they need no introduction, let's get into it. 

## Signatures
Declaring a function in go is simple
``` go
func GoFunc(param1 string, param2: int) string {
    // NOTE! fmt.Springf() allows you to format a string in a neat way! the %_ tags signify different types of data to put into the string. EX, %s is a string, %d is a number, etc. 
    goat := fmt.Sprintf("Stinky poop fart %s %d", param1, param2) 
    return "" 
}

func VoidFunc() { // return type of void!
    // does nothing
}
```

The structure of a function in go is simple, always start with the `func` keyword, then the name of the function, the parameters the function takes in, & the return type.

## Free floating vs attatched to struct
You can also attatch functions onto structs so the functions can directly use the data within the struct! You do this by adding one more part to the original func signature

In this example I'll create a function for a redefined Goat struct (note you can't make 2 of the same struct type I'm just remaking him for the example)
``` go
type Goat struct {
    Name string
    Silly bool
}

func (g *GenericGoat) HasBadTakes() bool { // ignore the pointer for now! We will get into that later
    if Name == "Tony" || Name == "Logan" {
        return true
    } else {
        return false
    }
}

func (g *GenericGoat) privateFunc() { // a private function on the goat struct! 
    // something 
}

// Example call:
realGoat := Goat{Name: "Suli", Silly: false}
fmt.Printf(realGoat.HasBadTakes()) // FALSE!!!! TAKE NoTE
```

## Functions as Types! 
Go makes it super easy to pass around functions as types. You can create types that represent certain function signatures! 

Here is the syntax:
``` go
type GoatFN func(goat *Goat, action string) string 
```

Now any function that follows this signature is considered a GoatFN type!

Example usage:

``` go
func exampleGoatFN(goat *Goat, action string) string {
    // logic..
    return "whatevah"
}

var func GoatFN = exampleGoatFN
func(Goat{Name: Suli, Silly: false}, "backflip")
```

# Interfaces
Go implements interfaces in a cool, but also kind of strange, way. Rather than directly telling a struct to implement a certain interface a function implements an interface if it has all of the methods in the interface defined for itself! If that sounds confusing (it is) peep this example!

``` go
type IGoat interface {
    Baaa(baaAt string)  string
    HasBadTakes()       bool
}
```

Now lets get the Goat struct we defined earlier to implment this interface!

``` go
type Goat struct {
    Name string
    Silly bool
}

func (g *GenericGoat) HasBadTakes() bool { // ignore the pointer for now! We will get into that later
    if Name == "Tony" || Name == "Logan" {
        return true
    } else {
        return false
    }
}
```

We see that the Goat struct already implements the HasBadTakes() function. In order to implement the IGoat interface we need to add the Baaa(baaA string) method!

``` go
func (g *Goat) Baaa(baaAt string) string {
    return fmt.SprintF("Baaaaaa %s", baaAt)
}
```

Bam! Let's use this as interface!

``` go
var goat IGoat = Goat{Name: "Hannah", Silly: true}
fmt.PrintF(goat.Baaa("Zainab"))
```

# Pointers

Honestly one of the things that through me off as a beginner... I will try to break it down to make it easy to understand

## What is a pointer?

Data in code is stored in memory in 2 seperate ways,

- As a value
- As an address to a value

Meaning at a given memory address it's either a value, or an address that points (that's where they get the name from ahaha) to another value

Go leverages pointers to make the code fast as fuck!

## The Basics, Syntax and Such
A point is symbolized by a * prefixing the type.

The & symbol is an operator in go that gets the memory address (pointer) to the object that follows it

For example,
``` go
var goat *Goat = &Goat{Name: "Zahra", Silly: true}
```

In this example goat is not the value of the struct created by Goat{} but rather an address that points to where the data is stored in memory

The &Goat{} part returns the address of the Goat struct we just created rather than the values of the struct

Now, in order to access the values of goat we would expand the pointer.

``` go
fmt.Printf(goat.Name)
```

Now you might be thinking... this is exactly what you would do normally for a struct. Now the best part about go issss that it will automatically expand the pointer for you! If you aren't looking to deep into the nity gritty, pointers in Go kind of don't matter too much for everyday use cases. But there are some cases you should look out for! 

## Value vs Copy
In all languages that use pointers (not just go) you either pass a value through code as a copy or a pointer

Now what does that mean...

Typically when a value (anything that isn't a pointer) is passed into a function or a struct on creation a copy of the value is made rather than passing along the original one (NOTE THIS IS HOW C WORKS, KINDA NEAT)

But this would work HORRIBLY for mutation. Passing in an object into a function or struct would be like completly creating the object again. And you would never be able to change your original dude. 

Next we will go over how to use pointers to solve this case!

### Pointers in Function Signature
If you want to change the values/have access to the original object in a function you would use pointers. Some example use cases!


**1. Passing original value into function** 
``` go
func mutateOG(goat *Goat) { // RIGHT
    goat.Name = "Mutated hehe"
}

func mutateOGCopy(goat Goat) { // WRONG!
    goat.Name = "Not mutated..."
}

goat := Goat{Name: "Avni", Silly: true}
mutateOG(&goat) // TAKE NOTE: passing in the pointer, this will mutate the original goat that we pass in
fmt.Printf(goat.Name) // OUTPUT: Mutated hehe
mutateOGCopy(goat) // NOT THE POINTER
fmt.Printf(goat.Name) // OUTPUT: Mutated hehe
```

Here when we pass in the pointer to the function the function is able to mutate the original variable rather than a copy. When we pass in the goat **value** a copy is made behind the scenes when the function is called, causing your mutation to do nothing.

See an example of this in the goat_test.go file located at backend/internal/models

**2. Struct Functions**

If you remember the Goat function example, when I attacthed the HasBadTakes function to the Goat struct I did it like (g *Goat) rather than (g Goat). While both work and will attatch the function to the struct, they both have different meanings.

When you attach a function to the pointer, ex (g *Goat) the function operates on the original struct that it's called on. When you attacth a function without the pointer, ex (g Goat) the function is called on a copy of the original

See an example of this in the goat_test.go file located at backend/internal/models

Here I will add some new functions to the Goat struct
``` go
func (g *Goat) SetName(name string) { // will set the name of the Goat it's called on 
	g.Name = name
}

func (g Goat) SetNameCopy(name string) { // will set the name of a copy of the Goat function was called on (useless)
	g.Name = name
}

func (g Goat) MakeCopy() Goat { // Creates an identical copy of the Goat the function was called on
	return g
}
```
Here we can see that setting the function to operate on the pointer or value of the struct can significantly change how you're function actually works! 

To see this example in action check out the method in goat_test.go!

### Pointers in Structs
The last example (that I can think of) was pointers in structs

When you set the values of your variables when declaring a struct the same pointer concepts apply!

``` go
type Erm struct {
    Goat *Goat
    Anotha  Goat
}
```

In this example the Erm struct had 2 variables, one is a pointer to a Goat and the other is the value of a Goat. This is important for when creating the struct! Can you guess what will happen (On some Dora shit)

``` go
goat := Goat{Name: "Varun", Silly: false}
erm := Erm{Goat: &goat, Anotha: goat}
```

The Anotha variable will be initialized to a **copy** of the goat variable. While the Goat variable will be intialized to the original instance of the goat.


# Error handling
Error handling in go follows a different pattern than in other languages. Errors in go are treated as values rather than things that "break" your code. Typically if your code can error out, lets say in a function, it's customary to return an error value in your function to indicate that an error has occured.

``` go

func errorProne(sillyName string) (string, error) {
    if sillyName == "Joe Mama" {
        return "", fmt.Errorf("The name %s is weird", sillyName)
    }

    return "AYYY", nil
}
```

**NOTEEE** I didn't go over this earlier but `tuples` exist in go. A tuple is a set of objects (any type of objects) of set size joined together. They allow you to return multiple objects in functions, store things easier, etc... Tuples are commonly used for returning the value you want along with any possible errors that can occur in your functions. You'll see this patten of (T any, error) used throughout our codebase as well. 

Typically when you error out you return the default value of your return type ("" for string, 0 for num, blank struct, etc.)

When your code is a success you return `nil` for the error value.

## Error Handling
When a function returns an error, **ALWAYS** handle it. Here's how to do it

``` go
str, err := errorProne("Joe mama...")

if err != nil {
    // handle error 
}
```

# Go Module
A Go module defines the "location" of your code. If you open up the backend/go.mod file you will see module inside-athletics. This is the root of all of your code in go!

While this isn't our use case, if you were creating a public library for go you could host your code on GitHub and people can access it just from that! You would import someone elses library (which you will see A LOT throughout our codebase) like this.

```go
import (
    "github.com/danielgtaylor/huma/v2"
)
```

Here we import the huma library straight from it's repo in GitHub!

This is a sick feature of go that we all leverage 

# Packages in Go 
A Package in go is a grouping of code. Files in the same directory (at the same level in the directory) are considered to be in the same package. Package naming convention follows a simple rule, the name of the package is the name of the direcotry the go file is in.

``` go
// If goat.go is in the Goat directory
package goat 

// rest of the code
```

The compiler doesn't enforce this, so your code won't break if you don't do this. But for convention sake, please do this! 

Packages in go let you import packages insanely easily!

## How to access code in different packages 
Simple! Import the file location of the code, whether its from GitHub or locally within your go module. 

To import code from a seperate package within your own codebase simply start with the module name/path_to_package_directory

For example if I wanted to import the models from the backend/internal/models folder

``` go
package example 

import (
    models "inside-athletics/internal/models"
)

// referencing a model 

models.Goat
```

# Dependencies
Dependency management is also super easy when using go!

If you're trying to use a library that we don't have imported within the project, start by importing the library. Then run go mod tidy in your terminal and it will update the dependency list for you as well as retrieve the code you need!

# Conclusion 
This is a very high level overview of stuff in Go. Make sure to look stuff up and ask questions if anything confuses you!


