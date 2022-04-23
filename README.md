# gorm-cache
## Work in progress, don't use.

`gorm-cache` aims to provide a look-aside, easy-to-use cache solution to gorm v2 users.

`gorm-cache` 旨在为gormv2用户提供一个即插即用的旁路缓存解决方案。

We provide 2 types of cache storage here:

我们提供2种存储介质：

1. Memory, where all cached data stores in memory of a single server
2. Redis, where cached data stores in redis (if you have multiple servers running the same procedure, they don't share the same redis storage space)


1. 内存 (所有数据存储在单服务器的内存中)
2. Redis (所有数据存储在redis中，如果你有多个实例使用本缓存，那么他们不共享redis存储空间)

There're mainly 5 kinds of operations in gorm (gorm function names in brackets):

在gorm中主要有5种操作（括号中是gorm中对应函数名）：

1. Query (First/Take/Last/Find/FindInBatches/FirstOrInit/FirstOrCreate/Count/Pluck)
2. Create (Create/CreateInBatches/Save)
3. Delete (Delete)
4. Update (Update/Updates/UpdateColumn/UpdateColumns/Save)
5. Row (Row/Rows/Scan)


