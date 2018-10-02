FROM alpine:3.8

RUN apk --no-cache add bash

ENV USER=appuser
ENV HOME=/home/$USER

RUN adduser -S $USER

# Copy files to container
COPY toRoot/* $HOME/

# Configure permissions
RUN chmod 755 $HOME/*.sh && \
    chown -R $USER $HOME

## Expose application port
EXPOSE 8080 8080
EXPOSE 8081 8081

USER $USER
WORKDIR /home/$USER
CMD ["./runapp.sh"]