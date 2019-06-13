# Anomaly detection data generator

I gouldn't get polyfractal's Rust demo working in Docker, so decided to rewrite it in Go.

Original Rust version:

https://github.com/polyfractal/playground/tree/master/hotcloud

This version takes just under an hour to generate the data with default settings, bulk inserting into a single Elasticsearch running locally in Docker.

Also this version includes keyword fields for node, query, and metric. You can use these keyword fields to split data, and as key fields (influencers) when setting up machine learning jobs in Kibana/Elasticsearch.

## Options

| Environment variable | Description                                             | Default               |
|----------------------|---------------------------------------------------------|-----------------------|
| CONFIG_FILE          | The path+name of the TOML config file to load           | ./config.toml         |
| ELASTICSEARCH_HOST   | Elasticsearch host including scheme, hostname, and port | http://localhost:9200 |
| LOG_LEVEL            | The level of logs to show (uses logrus log levels)      | debug                 |

## Known issues:

- Hangs after final bulk request is sent, has to be exited with CTRL+C/SIGINT. Probably an easy fix
- Does not generate the hotcloud search template etc. Didn't need that part for what I was doing (using Elasticsearch's built-in anomaly detection features) so didn't finish implementing it

