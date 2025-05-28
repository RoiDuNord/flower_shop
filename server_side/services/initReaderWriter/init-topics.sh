#!/bin/sh
set -e

BROKER="kafka:9092"
KAFKA_BIN="/opt/kafka/bin/kafka-topics.sh"

echo "Creating Kafka topics..."

$KAFKA_BIN --create --topic orders --bootstrap-server $BROKER --partitions 1 --replication-factor 1 || echo "Topic orders exists"
$KAFKA_BIN --create --topic payments --bootstrap-server $BROKER --partitions 1 --replication-factor 1 || echo "Topic payments exists"

echo "Topics created successfully."
