# Airflow On K8S Operator
[![Go Report Card](https://goreportcard.com/badge/github.com/apache/airflow-on-k8s-operator)](https://goreportcard.com/report/github.com/apache/airflow-on-k8s-operator)

## Community

* Join [Airflow Slack](https://apache-airflow-slack.herokuapp.com) and the dedicated #sig-kubernetes channel.

## Project Status

*Alpha*

The Airflow Operator is still under active development and has not been extensively tested in production environment. Backward compatibility of the APIs is not guaranteed for alpha releases.

## Prerequisites
* Version >= 1.9 of Kubernetes.
* Uses 1.9 of Airflow (1.10.1+ for k8s executor)
* Uses 4.0.x of Redis (for celery operator)
* Uses 5.7 of MySQL

## Get Started

[One Click Deployment](https://console.cloud.google.com/marketplace/details/google/airflow-operator) from Google Cloud Marketplace to your [GKE cluster](https://cloud.google.com/kubernetes-engine/)

Get started quickly with the Airflow Operator using the [Quick Start Guide](docs/quickstart.md)

For more information check the [Design](docs/design.md) and detailed [User Guide](docs/userguide.md)

## Airflow Operator Overview
Airflow Operator is a custom [Kubernetes operator](https://coreos.com/blog/introducing-operators.html) that makes it easy to deploy and manage [Apache Airflow](https://airflow.apache.org/) on Kubernetes. Apache Airflow is a platform to programmatically author, schedule and monitor workflows. Using the Airflow Operator, an Airflow cluster is split into 2 parts represented by the `AirflowBase` and `AirflowCluster` custom resources.
The Airflow Operator performs these jobs:
* Creates and manages the necessary Kubernetes resources for an Airflow deployment.
* Updates the corresponding Kubernetes resources when the `AirflowBase` or `AirflowCluster` specification changes.
* Restores managed Kubernetes resources that are deleted.
* Supports creation of Airflow schedulers with different Executors
* Supports sharing of the `AirflowBase` across mulitple `AirflowClusters`

Checkout out the [Design](docs/design.md)

![Airflow Cluster](docs/airflow-cluster.png)


## Development

Refer to the [Design](docs/design.md) and [Development Guide](docs/development.md).

## History
This repo has been donated to Apache foundation.
It was originally developed here at [GoogleCloud repo](https://github.com/GoogleCloudPlatform/airflow-operator)
