# Use the official MongoDB image
FROM mongo:latest

# Set the environment variable for MongoDB to accept connections from all IP addresses
ENV MONGO_INITDB_ROOT_USERNAME=mongo_user
ENV MONGO_INITDB_ROOT_PASSWORD=mongo_password

# Expose the default MongoDB port
EXPOSE 27017

# The entrypoint is already set in the MongoDB image, so no need to define CMD
# Run with the following command
# docker run -d -p 27017:27017 -v /my/local/dir:/data/db my-mongodb
