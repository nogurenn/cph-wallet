# This Dockerfile does not build the app itself. Please use `make build` to build the app binary first.
FROM golang:1.17


# Use non-root user for better security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN mkdir -p /app/bin
RUN chown -R appuser:appgroup /app
USER appuser

COPY --chown=appuser:appgroup ./bin/wallet /app/bin
WORKDIR /app/bin

ENTRYPOINT ["/app/bin/wallet"]
