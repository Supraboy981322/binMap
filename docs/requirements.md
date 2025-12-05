# binMap --> Requirements
 
### RAM
However much space you want for your database plus the baseline for your OS and other services.

By defauly, for performance, when the server starts the full database is read from the gomn-as-a-binary file, then held in memory, meaning if you have a 16GB database, the server is going to fill 16GBs of RAM.

To save on RAM, you can do one of the following:
- Disable the in-memory database
  - Find `["use in-memory db"]` in your `conf.gomn`
  - Set it to any value greater-than `0` (in seconds). Eg:
    ```gomn
    ["use in-memory db"] := 86400
    ```
    To clear the database once-a-day, every-day

- Clear the database when it reaches a certain size 
  - Find `["clear db if size is n MB"]` in your `conf.gomn`
  - Set it to any value greater-than `0` (in seconds). Eg:
    ```gomn
    ["clear db if size is n MB"] := 86400
    ```
    To clear the database once-a-day, every-day

- Clear the database after an amount of time
  - Find `["clear db every n seconds"]` in your `conf.gomn`
  - Set it to any value greater-than `0` (in seconds). Eg:
    ```gomn
    ["clear db every n seconds"] := 86400
    ```
    To clear the database once-a-day, every-day


### CPU
Pretty much anything will do fine if you don't have a lot of users.


### Disk space
However much space you want for your database.


### Networking
However fast you would need to create a entry with the following parameters:
  - A key of `EMERGENCY HOME DIR BACKUP`
  - The command `tar -cf - ~` piped into stdin as the value
