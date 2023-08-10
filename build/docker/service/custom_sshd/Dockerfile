FROM ubuntu

# Install OpenSSH
RUN apt-get update && apt-get install -y openssh-server

# Install OpenSSH server
RUN apt-get update && apt-get install -y openssh-server

# Copy the SSH public key to the container's root directory
# COPY $SSH_KEY_PUB_PATH /root/.ssh/authorized_keys
RUN mkdir -p /root/.ssh
#RUN echo "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIPp/jbeIDrxDXW18vBuUxOj+N5msJe/elIJlgI4znfi4 riotpot@riotpot.com" >> /root/.ssh/authorized_keys
COPY keys/riotpot_ed25519.pub /root/.ssh/authorized_keys


# Set up SSH configuration
RUN mkdir /var/run/sshd
RUN echo 'root:password' | chpasswd
RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

# Expose SSH port
EXPOSE 22

# Set SSH as the entry point
CMD ["/usr/sbin/sshd", "-D"]
