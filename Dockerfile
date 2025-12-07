FROM --platform=$BUILDPLATFORM alpine:3.20
ARG TARGETARCH
WORKDIR /app


# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the binary
COPY ./requestsink-${TARGETARCH} /app/request-sink

# Make sure the binary is executable
RUN chmod +x /app/request-sink

# Change ownership of the binary to the non-root user
RUN chown appuser:appgroup /app/request-sink

# Switch to the non-root user
USER appuser

HEALTHCHECK --interval=10s --timeout=5s --retries=3 CMD sh -c '\
  PORT=${API_PORT:-8080} && \
  wget -q --header="X-Auth-Key: $API_KEY" -O /dev/null "http://localhost:$PORT/__/startup" >/dev/null'


ENTRYPOINT ["/app/request-sink"]