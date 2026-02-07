# Create a shared Docker network (once)
```
docker network create observability
```

# Run Tempo
```shell
docker run -d --name tempo \
  --network observability \
  -p 3200:3200 \
  -p 4317:4317 \
  -v $(pwd)/tempo.yml:/etc/tempo.yml \
  grafana/tempo:2.4.1 \
  -config.file=/etc/tempo.yml
```

# Run Prometheus
```
docker run -d --name prometheus \
  --network observability \
  -p 9090:9090 \
  -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus
```

# Run Grafana
```
docker run -d --name grafana \
  --network observability \
  -p 3000:3000 \
  grafana/grafana
```
