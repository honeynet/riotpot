# Configuring FTP User Names and Passwords

**Create a Secret File**: Inside the directory where  `docker-compose.yaml` is in, create a file named `.env` . In this file, add your FTP user credentials in the format `username|password`.

Example `ftp.env` (space and | separated list):
```
USERS="user|pass user2|pass2"
```