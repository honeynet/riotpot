# Use the Ubuntu base image
FROM ubuntu

# Install Telnet
RUN apt-get update && apt-get install -y telnet

# Expose port 23 for Telnet
EXPOSE 23

# Set Telnet as the entry point
CMD ["telnet", "localhost"]
