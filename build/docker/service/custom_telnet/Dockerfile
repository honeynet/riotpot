# Use the Ubuntu base image
FROM ubuntu

# Install Telnet
RUN apt-get update && apt-get install -y telnetd

# Expose port 23 for Telnet
EXPOSE 23

# Start Telnet server using telnetd
CMD ["/usr/sbin/in.telnetd", "-debug", "23"]
