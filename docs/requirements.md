# binMap --> Requirements
 
### RAM
However much space you want for your database plus the baseline for your OS and other services. The full database is read from the gomn-as-a-binary file, then held in memory, meaning if you have a 16GB database, the server is going to fill 16GBs of storage. There are plans to add an option to disable storing the database in memory (with a warning of degraded performance), but this has not yet been implemented.

### CPU
Pretty much anything will do fine if you don't have a lot of users.

### Disk space
However much space you want for your database.

### Networking
However fast you would need to create a entry with the following parameters:
  - A key of `EMERGENCY HOME DIR BACKUP`
  - The command `tar -cf - ~` piped into stdin as the value
