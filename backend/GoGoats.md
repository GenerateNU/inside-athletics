# SOME BASICS
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

## Lists
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

# Private v Public
In go you can't declare something private or public... Instead something is made private or public based on if it starts with a capital letter or not. Some examples to follow...

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

goat := Goat{} // you can initialize the struct to nothing 
var another Goat = Goat{Name: "Suli", secrets: []string{"loves milk", "hates water"}}
```
Here we can see that the struct name "Goat" is capital, meaning that the struct is available as a public class outside of it's package. When we take a look at the Name and secrets variables Name is public while secrets is private. 

Now something you might have noticed is this little json thingy... It lowk makes go lit asf. These are called `tags` and their super useful. They provide libraries with metadata needed to process the field. For example, if I were using a JSON library in go, I wouldn't have to map every field in the struct in a dictionary (like in python) as the library can use the tag to handle the case itself.

We will be using these tags across our codebase to simplify API calls as well as create API documentation!

## How to create a Generic
Something you will have to do less often but if you were curious.

``` go
package mysteriousGoat

type GenericGoat[T any] struct {
    MysteryVariable T
}

// initialize a generic
var goat GenericGoat[int] = GenericGoat{MysterVariable: 67}
```

Here I create a generic struct

# Functions

## Signatures

## Free floating vs attatched to struct

# Pointers

Honestly one of the things that through me off as a beginner... I will try to break it down to make it easy to understand

## What is a pointer?

Data in code is stored in memory in 2 seperate ways,

- As a value
- As an address to a value

Meaning at a given memory address it's either a value, or an address that points (that's where they get the name from ahaha) to another value

## Value vs Copy

## Primitives

## Apply to function vs apply to value

## Passing in pointers vs value


# Interfaces

# Iteration

# Error handling

# Packages in Go 

# Testing 

## Directories are king

# Dependencies 

