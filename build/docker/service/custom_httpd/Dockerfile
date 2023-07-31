# Use the official httpd image as the base image
FROM httpd

# Copy the "index.html" file with "Hello, world!" content to the document root
COPY index.html /usr/local/apache2/htdocs/

# Copy the custom SSL certificate and key files
COPY server.crt /usr/local/apache2/conf/server.crt
COPY server.key /usr/local/apache2/conf/server.key

# Enable SSL module and configure HTTPS
RUN sed -i '/#LoadModule ssl_module/s/^#//g' /usr/local/apache2/conf/httpd.conf \
    && echo "Listen 443" >> /usr/local/apache2/conf/httpd.conf \
    && echo "<VirtualHost *:443>" >> /usr/local/apache2/conf/httpd.conf \
    && echo "  SSLEngine on" >> /usr/local/apache2/conf/httpd.conf \
    && echo "  SSLCertificateFile /usr/local/apache2/conf/server.crt" >> /usr/local/apache2/conf/httpd.conf \
    && echo "  SSLCertificateKeyFile /usr/local/apache2/conf/server.key" >> /usr/local/apache2/conf/httpd.conf \
    && echo "  DocumentRoot /usr/local/apache2/htdocs" >> /usr/local/apache2/conf/httpd.conf \
    && echo "  ErrorLog /dev/stderr" >> /usr/local/apache2/conf/httpd.conf \
    && echo "  CustomLog /dev/stdout common" >> /usr/local/apache2/conf/httpd.conf \
    && echo "</VirtualHost>" >> /usr/local/apache2/conf/httpd.conf
