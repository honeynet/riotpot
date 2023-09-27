FROM alpine:latest

# Install OpenSSH
RUN apk add --no-cache --update openssh

# Add non-root users and set passwords
COPY ssh_users.txt /tmp/ssh_users.txt
RUN cat /tmp/ssh_users.txt | while IFS=: read -r user password; do \
      if [ "$user" = "root" ]; then \
        echo "root:$password" | chpasswd; \
      else \
        adduser --gecos "" "$user"; \
        echo "$user:$password" | chpasswd; \
        sed -i "s#^$user:x:#$user:x:/root:/bin/bash#" /etc/passwd; \
      fi; \
    done && rm /tmp/ssh_users.txt

# Enable root login via password
RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

# Expose SSH port
EXPOSE 22

# Set SSH as the entry point
CMD ["/usr/sbin/sshd", "-D"]
