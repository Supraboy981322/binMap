# binMap --> Requirements
 
### RAM
However much space you want for your database plus the baseline for your OS and other services.

At start, the full database is read from the gomn-as-a-binary file, then held in memory, meaning if you have a 16GB database, the server is going to fill 16GBs of RAM.

There are plans to add an option to disable the in-memory database (with a warning of degraded performance), and for clearing it when it reaches a certain size. Since these have not yet been implemented, for now, you can auto-clear the database after an amount of time like so:

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
