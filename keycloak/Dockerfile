# Use a lightweight base image
FROM alpine:latest

# Install required packages
RUN apk --no-cache add curl jq

# Copy the setup script into the container
COPY setup.sh /setup.sh

# Make the script executable
RUN chmod +x /setup.sh

# Default command to run the setup script
ENTRYPOINT ["sh", "-c", "sleep 60 && /setup.sh"]
#ENTRYPOINT ["sh", "-c", "while true; do sleep 1000; done"]
