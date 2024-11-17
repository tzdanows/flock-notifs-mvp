# Running an Arroyo cluster on Kubernetes

![alt text](image-1.png)

primarily following https://doc.arroyo.dev/deployment/kubernetes#example-local-configuration & https://doc.arroyo.dev/tutorial/kafka

```bash
# add the arroyo helm repo to your helm installation
helm repo add arroyo https://arroyosystems.github.io/helm-repo

> helm search repo arroyo
NAME            CHART VERSION   APP VERSION     DESCRIPTION                                       
arroyo/arroyo   0.12.0          0.12.0          Helm chart for the Arroyo stream processing engine
```

```bash
helm show values arroyo/arroyo

# relevant options to change are:
postgresql.deploy: 
    Whether to deploy a new Postgres instance. If set to false, the chart will expect a Postgres instance to be available with the connection settings determined by postgresql.externalDatabase configurations (by default: postgres://arroyo:arroyo@localhost:5432/arroyo).
artifactUrl & checkpointUrl: 
    Configures where pipeline artifacts and checkpoints are stored. See the overview for more details on how these are configured. If this is set to a local directory (when running a local k8s cluster), you will need to configure volumes and volumeMounts to make this directory available on all of the pods.
existingConfigMap:
    allows you to set environment variables on the Arroyo pods.
```
values.yaml
```yaml
# Default values for arroyo.

nameOverride: ""
fullnameOverride: ""
clusterDomain: cluster.local

image:
  repository: ghcr.io/arroyosystems/arroyo
  tag: "0.12"
imagePullPolicy: IfNotPresent
imagePullSecrets: []

controller:
  resources:
    limits: {}
    requests:
      memory: 256Mi
      cpu: 500m

  service:
    grpcPort: 5116
    compilerPort: 5117
    adminPort: 5114
    httpPort: 80

worker:
  resources:
    limits: {}
    requests:
      memory: 490Mi
      cpu: 900m
  slots: 16

postgresql:
  # set to true to deploy a postgres instance in-cluster
  deploy: true

  # otherwise, configure the connection to an external postgres instance
  externalDatabase:
    host: localhost
    port: 5432
    name: arroyo
    user: arroyo
    password: arroyo

  # postgresql configuration
  auth:
    username: arroyo
    password: arroyo
    database: arroyo

  primary:
    initdb:
      user: postgres

  serviceAccount:
    create: true
    name: arroyo-postgresql


prometheus:
  # set to true to enable prometheus annotations on Arroyo cluster resources; if false,
  # you will need to configure prometheus separately to discover Arroyo pods
  setAnnotations: true

s3:
  region: ""

artifactUrl: ""
checkpointUrl: ""

serviceAccount:
  # Specifies whether a service account should be created
  create: true
  # Annotations to add to the service account
  annotations: {}
  # The name of the service account to use.
  # If not set and create is true, a name is generated using the fullname template
  name: ""

role:
  # Specified whether the role should be created
  create: true


podAnnotations: {}

podSecurityContext: {}
  # fsGroup: 2000

securityContext: {}
  # capabilities:
  #   drop:
  #   - ALL
  # readOnlyRootFilesystem: true
  # runAsNonRoot: true
  # runAsUser: 1000



controllerResources: {}

nodeSelector: {}

tolerations: []

affinity: {}

volumes: []
volumeMounts: []

existingConfigMap: ""
```

setup for local cluster without S3
```yaml
artifactUrl: "/tmp/arroyo-test"
checkpointUrl: "/tmp/arroyo-test"

volumes:
  - name: checkpoints
    hostPath:
      path: /tmp/arroyo-test
      type: DirectoryOrCreate

volumeMounts:
  - name: checkpoints
    mountPath: /tmp/arroyo-test
```

once you have your config file values.yaml, install the helm chart
```bash
helm install arroyo arroyo/arroyo -f values.yaml
# then
kubectl get pods

NAME                                        READY   STATUS             RESTARTS      AGE
arroyo-compiler-ccd6b7bdb-752vt             1/1     Running            0             36s
arroyo-controller-75587f886b-k9drg          1/1     Running            1 (18s ago)   36s
arroyo-postgresql-0                         1/1     Running            0             26s
arroyo-api-5dccb89967-zl727                 1/1     Running            2 (17s ago)   36s
```

```bash
# access the ui
open "http://$(kubectl get service/arroyo -o jsonpath='{.spec.clusterIP}')"

# can proxy the service to local machine with:
kubectl port-forward service/arroyo 5115:80
```

TODO: verify we can deploy postgres & minio/s3
