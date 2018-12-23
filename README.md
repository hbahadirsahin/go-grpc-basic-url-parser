# Go-based gRPC Server

This repository is created for a technical interview task.

Note that I have no past Go, gRPC and webpage parsing experience before.

# Introduction

The aim of the task is to write a gRPC server (Nope, I won't give too much detail, just google it) which will serve a "ParseURL" 
method to parse a given URL's HTML page.

ParseURL method takes a "url" as an input parameter and returns the parsed "title," "thumbnail_url" and "content". 

- The "url" can be either a newspage or a blog page. 
- "title" is the `<title>` of the page. 
- "thumbnail_url" is an image URL which is parsed from the page as a thumbnail image. 
- "content" is the all text content of the page.

This repository contains:
- A "parserproto" folder, which contains the `parser.proto` file for this task as well as its compiled version for GO which is `parser.pb.go`
- A server main code `parser_server_main.go`
- A client (for tests) main code `parser_client_main.go`

# How to Run

Assuming you install gRPC and Protocol Buffer related stuffs for Go:

- Open a command window, if you are against using any kind of IDEs, and type `go run parser_server_main.go`.
- Open another command window, and type `go run parser_client_main.go`. 
- If you are using an IDE, just press the run/build/compile whatever button you have for both main.go files.

Since, the default url is hardcoded, you will get parsed results of the input url. Note that server uses "50051" port and the client connects "localhost:50051" as their default values. But all these constant default variables will become input arguments in very near future.

As a note: Since compiling `parser.proto` was a pain for me (in Windows), I pushed the compiled version, too. So, you do not need to compile it, but you need gRPC related libraries installed for Go to run this repository (For how to install gRPC: https://grpc.io/docs/quickstart/go.html)


# TO-DO Tasks

- [ ] Extend unit tests.
- [ ] More edge case checker and better error handling.
- [ ] Write an extended README for installing and using this project. 

# Finished Tasks 

## Update 23.12.2018

- [x] Learn mock and write unit-tests.
  - [x] Mock for client is created by using "mockgen".
  - [x] A basic test is implemented for testing the first unit test I've written in Go. 

## Update 21.12.2018

- [x] Client code will take "-address" and "-url" as input argument/parameter

## Update 20.12.2018

- [x] Parser function to parse "content" of the page.
  - [x] 2 more specific parsers are added. By specific, I mean the added parsing functionality is probably not gonna work for many different webpages.
  - [x] A general purpose content parser functionality added. It gets the values of `<p>` tags first, then `<ol>` tags and finally `<ul>` tags. This method does not guarantee the order of the text if the content contains list(s). 

## Update 19.12.2018

- [x] Server code will take "-port" as input argument/parameter 
- [x] A basic parser function is implemented to parse text "content" of "Medium Blog" pages. Again, since all pages are different, and I can't find a way to parse them all with a unified method, this function will be updated with 3-4 different page cases at most (and will return empty content for the rest).

## Update 18.12.2018

- [x] Write and compile `.proto` file to define request and response methods.
- [x] Parse input URL and get title.
  - [x] In case a page does not contain `<title>` tag, the code will check every possible title candidate by checking `<h1>`, `<h2>` and `<h3>` tags.
  - [x] If code cannot find any title candidate, it will send a string about it (instead of error, but I will probably change it).
- [x] Parse input URL and get the first, *related* image's URL.
  - [x] For 3 possible HTML structure, code check the image urls and its `alt` attribute. The longest `alt` attribute value has a high probability of being related to the input page.
  - [x] If images in a page does not have any `alt` attribute or do not fit the defined structure, it returns the first image found in the page.
- [x] For dumb-downed, personal tests, a client code has been written.



