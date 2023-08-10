FROM delfer/alpine-ftp-server

# Add your custom user credentials here (replace username|password with your desired values)
ENV USERS="ftpuser|password"

# Create directory for FTP user
RUN mkdir -p ./ftp_data

# Uncomment the following lines to set the passive mode ports range (optional)
ENV MIN_PORT=21000
ENV MAX_PORT=21010
