FROM debian:sid

# Copy the binary server
ADD server /server
RUN cp /server /usr/local/bin/wappalyzer-server

# Copy the JS files
ADD extraction/js/ extraction/js/

EXPOSE 3001

ENTRYPOINT ["wappalyzer-server"]
