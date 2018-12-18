# Go-based gRPC Server

This repository is created for as a result of a technical interview task.

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

# Finished Tasks

- [x] Write and compile `.proto` file to define request and response methods.
- [x] Parse input URL and get title.
  - [x] In case a page does not contain `<title>` tag, the code will check every possible title candidate by checking `<h1>`, `<h2>` and `<h3>` tags.
  - [x] If code cannot find any title candidate, it will send a string about it (instead of error, but I will probably change it).
- [x] Parse input URL and get the first, *related* image's URL.
  - [x] For 3 possible HTML structure, code check the image urls and its `alt` attribute. The longest `alt` attribute value has a high probability of being related to the input page.
  - [x] If images in a page does not have any `alt` attribute or do not fit the defined structure, it returns the first image found in the page.
- [x] For dumb-downed, personal tests, a client code has been written.

# TO-DO Tasks

- [ ] Server code will take "-port" as input argument/parameter (it is a constant value in the current version)
- [ ] Parser function to parse "content" of the page.
- [ ] Learn mock and write unit-tests.
- [ ] Write an extended README for installing and using this project. 


