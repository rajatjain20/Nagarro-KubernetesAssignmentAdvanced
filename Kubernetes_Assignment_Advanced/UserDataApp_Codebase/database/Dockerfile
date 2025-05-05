# Base Image
FROM mcr.microsoft.com/mssql/server:2019-latest

USER root

# Set working directory
WORKDIR /usr/config

# Copy scripts
COPY --chown=root ./scripts/ /usr/config/

# Grant permissions to our scripts to be executable
RUN chmod +x /usr/config/entrypoint.sh
RUN chmod +x /usr/config/configure-db.sh

# Set Environment Variables
ENV ACCEPT_EULA=Y \
    TZ="Asia/Kolkata"

# Expose SQL SERVER PORT    
EXPOSE 1433

# Entry Point script to be executed
ENTRYPOINT ["./entrypoint.sh"]