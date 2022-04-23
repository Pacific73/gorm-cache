# gorm-cache
## Work in progress, don't use.

`gorm-cache` aims to provide a look-aside, easy-to-use cache solution to gorm v2 users.

When we use gorm as our orm framework in projects, there're lots of interactions with database, where `gorm-cache` comes in and improves the performance of our program.

We provide 2 types of cache storage here:

1. Memory, where all cached data stores in memory of a single server
2. Redis, where cached data stores in redis (if you have multiple servers running the same procedure, they don't share the same cached data)

