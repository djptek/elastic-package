package install

const snapshotYml = `version: '2.3'
services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.0.0-SNAPSHOT
    healthcheck:
      test: ["CMD", "curl", "-f", "-u", "elastic:changeme", "http://127.0.0.1:9200/"]
      retries: 300
      interval: 1s
    environment:
    - "ES_JAVA_OPTS=-Xms1g -Xmx1g"
    - "network.host="
    - "transport.host=127.0.0.1"
    - "http.host=0.0.0.0"
    - "indices.id_field_data.enabled=true"
    - "xpack.license.self_generated.type=trial"
    - "xpack.security.enabled=true"
    - "xpack.security.authc.api_key.enabled=true"
    - "ELASTIC_PASSWORD=changeme"

  kibana:
    image: docker.elastic.co/kibana/kibana:8.0.0-SNAPSHOT
    depends_on:
      elasticsearch:
        condition: service_healthy
      package-registry:
        condition: service_healthy
    healthcheck:
      test: "curl -f http://localhost:5601/login | grep kbn-injected-metadata 2>&1 >/dev/null"
      retries: 600
      interval: 1s
    volumes:
      - ./kibana.config.yml:/usr/share/kibana/config/kibana.yml

  package-registry:
    image: docker.elastic.co/package-registry/package-registry:master
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080"]
      retries: 300
      interval: 1s

  elastic-agent:
    image: docker.elastic.co/beats/elastic-agent:8.0.0-SNAPSHOT
    depends_on:
      elasticsearch:
        condition: service_healthy
      kibana:
        condition: service_healthy
    environment:
    - "FLEET_ENROLL=1"
    - "FLEET_SETUP=1"
    - "KIBANA_HOST=http://kibana:5601"
`
