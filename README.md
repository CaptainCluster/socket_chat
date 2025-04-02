How to install go
---

Since this application uses Golang, rather than Python that seems to be the popular choice for this course, you have 
to ensure you have Go installed as a dependency on your system.

[Click here for more information](https://go.dev/doc/install)

Running the application
---

The server must be run in order to have an X amount of clients interacting with it.

1. Starting the server involves going to the server directory
> cd server
>
> go run .

2. Same goes for the client script
> cd client
>
> go run .

Ensure that _port 8000_ is free! If not, then you must alter const variables in both
of the file (client.go & server.go).


Sources
---

Since this is my first time using _Golang_ (completely my own choice),
I used various sources to help me with the creation of this application.

Golang documentation - [click here (go.dev)](https://go.dev/doc/)

Golang book - [click here (golang-book.com)](https://www.golang-book.com/public/pdf/gobook.pdf)

How to make sockets with Go - [click here (dev.to)](https://dev.to/alicewilliamstech/getting-started-with-sockets-in-golang-2j66)

Structs in Go - [click here (gobyexample.com)](https://gobyexample.com/structs)

TCP client - [click here (golinuxcloud.com)](https://www.golinuxcloud.com/golang-tcp-server-client/)

Socket client - [click here (golangr.com)](https://golangr.com/socket-client)

Getting string to uppercase - [click here (geeksforgeeks.com)](https://www.geeksforgeeks.org/how-to-convert-a-string-in-uppercase-in-golang/)

Making an empty array in Golang - [click here (stackoverflow.com)](https://stackoverflow.com/questions/45317074/best-practices-constructing-an-empty-array)

Goroutine, aka Go's thread - [click here (go.dev)](https://go.dev/tour/concurrency/1)

Converting int to string without value loss - [click here (golang.cafe)](https://golang.cafe/blog/golang-int-to-string-conversion-example.html)

Receiving inputs, even with whitespace - [click here (stackoverflow.com)](https://stackoverflow.com/questions/27414598/golang-accepting-input-with-spaces)

Random number library - [click here (golangdocs.com)](https://golangdocs.com/generate-random-numbers-in-golang)

Slicing to delete array elements - [click here (geeksforgeeks.com)](https://www.geeksforgeeks.org/delete-elements-in-a-slice-in-golang/)

Mutex in Golang - [click here (geeksforgeeks.com)](https://www.geeksforgeeks.org/mutex-in-golang-with-examples/)
