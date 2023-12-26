# go-json-db
This is a light weight, easily managable and beginner friendly implementation of a Database that anyone can implement directly in their Local Systems. It provides database API's to perform pretty much every operation that one can perform in a traditional database. The aim is to be a worthy contender in the choice of Database for small and medium use cases. 

```
                   _                           _ _     
  __ _  ___       (_)___  ___  _ __         __| | |__  
 / _` |/ _ \ _____| / __|/ _ \| '_ \ _____ / _` | '_ \ 
| (_| | (_) |_____| \__ \ (_) | | | |_____| (_| | |_) |
 \__, |\___/     _/ |___/\___/|_| |_|      \__,_|_.__/ 
 |___/          |__/                                   
```
# About
It goes down a No-Sql path ,quite similar to what MongoDB (Atlas) does. This project uses Directories to respresent any new collection and subsequent JSON files in their respective collections, each representing a record. 

The main idea of this project is to bring the complex database processes like caching and mutexes into the most understandable language of Developers i.e. JSON : ).

# Installation 

Currently, to install, there are two dependencies that are also to be install with the project.

This package can be installed by -

```
go get github.com/shubhexists/go-json-db
```

Cache is implemented by - 

```
go get github.com/patrickmn/go-cache
```

Logging is implemented by - 

```
github.com/jcelliott/lumber
```
# Usage 
1) Import the package by - 

```
github.com/shubhexists/go-json-db/models
```

2) Create a new db instance in your desired directory -

```
 db, cache, err := models.New("./database")
```

3) Perform Operations as desired! For reference you can take the help of `examples/main.go` of this repository

# Status 
This project was built by me as a part of learning Golang. This project is heavily insipred by the architecture that MongoDB takes to provide a database. 

This project still has lot of operations left which I would complete as soon as I get time , However as it's current stage, it is ready for usage in small projects atleast. : )

Enjoy!
