# binMap --> Requirements
 
### RAM
However much space you want for your db plus the baseline for your OS and other services

### CPU
Pretty much anything will do fine if you don't have a lot of users

### Disk space
However much space you want for your db

### Networking
However fast would you need to create a entry with the following parameters:
  - A key of `EMERGENCY HOME DIR BACKUP`
  - The command `tar -cf - ~` piped into stdin as the value
