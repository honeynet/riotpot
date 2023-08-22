# Configuring FTP User Names and Passwords

**Create a Secret File**: Inside the directory where  `Dockerfile` for the FTP service is in, create a file named `ftp_users.txt` . In this file, add your FTP user credentials in the format `username|password`.

Example `ftp_users.txt` (space and | separated list):
```
ftpuser|password anotheruser|secretpass
```