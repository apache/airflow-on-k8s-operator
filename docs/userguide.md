<!--
 Licensed to the Apache Software Foundation (ASF) under one or more
 contributor license agreements.  See the NOTICE file distributed with
 this work for additional information regarding copyright ownership.
 The ASF licenses this file to You under the Apache License, Version 2.0
 (the "License"); you may not use this file except in compliance with
 the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
 -->

# User Guide

TODO

# FAQs

1. How do we refresh DAGs ?
Canonical way airflow supports refreshing DAGs is via `dag_dir_list_interval` config.
https://cwiki.apache.org/confluence/display/AIRFLOW/Scheduler+Basics#Configuration
You can set that config using `cluster.spec.config.airflow`
Set the env `AIRFLOW__SCHEDULER__ DAG_DIR_LIST_INTERVAL`
By default dags are refreshed every 5 minutes.
To enable continuous sync, use git or gcs dag source with once disabled.

```yaml
apiVersion: airflow.apache.org/v1alpha1
kind: AirflowCluster
...
spec:
  ...
  config:
    airflow:
      AIRFLOW__SCHEDULER__DAG_DIR_LIST_INTERVAL: 100 # default is 300s
  ...
  dags:
    subdir: ""
    gcs:
      bucket: "mydags"
  # OR
  dags:
    subdir: "airflow/example_dags/"
    git:
      repo: "https://github.com/apache/incubator-airflow/"
      once: false
```
