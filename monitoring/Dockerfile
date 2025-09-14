FROM prom/prometheus:latest

# Copy our config file
COPY prometheus.yml /etc/prometheus/prometheus.yml

# Expose the default Prometheus port
EXPOSE 9090

# Use the default Prometheus entrypoint with our config
CMD ["--config.file=/etc/prometheus/prometheus.yml", \
     "--storage.tsdb.path=/prometheus", \
     "--web.console.libraries=/etc/prometheus/console_libraries", \
     "--web.console.templates=/etc/prometheus/consoles", \
     "--web.enable-lifecycle", \
     "--web.enable-admin-api"]