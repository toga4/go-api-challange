# Create GKE Cluster

- Enable GKE
- Enable Cloud DNS

```
$ gcloud beta container clusters create "sandbox-toga4-gke-cluster01-tokyo" \
 --project "sandbox-toga4-gke" \
 --zone "us-central1-c" \
 --no-enable-basic-auth \
 --cluster-version "1.21.4-gke.1801" \
 --release-channel "rapid" \
 --machine-type "e2-micro" \
 --image-type "COS_CONTAINERD" \
 --disk-type "pd-standard" \
 --disk-size "10" \
 --metadata disable-legacy-endpoints=true \
 --scopes "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring","https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" \
 --max-pods-per-node "110" \
 --spot \
 --num-nodes "3" \
 --logging=SYSTEM,WORKLOAD \
 --monitoring=SYSTEM \
 --enable-private-nodes \
 --master-ipv4-cidr "172.16.0.0/28" \
 --enable-ip-alias \
 --network "projects/sandbox-toga4-gke/global/networks/default" \
 --subnetwork "projects/sandbox-toga4-gke/regions/us-central1/subnetworks/default" \
 --enable-intra-node-visibility \
 --default-max-pods-per-node "110" \
 --enable-dataplane-v2 \
 --enable-master-authorized-networks \
 --master-authorized-networks 0.0.0.0/0 \
 --addons HorizontalPodAutoscaling,HttpLoadBalancing,GcePersistentDiskCsiDriver \
 --enable-autoupgrade \
 --enable-autorepair \
 --max-surge-upgrade 1 \
 --max-unavailable-upgrade 0 \
 --workload-pool "sandbox-toga4-gke.svc.id.goog" \
 --enable-shielded-nodes \
 --node-locations "us-central1-c" \
 --cluster-dns clouddns \
 --cluster-dns-scope cluster
```

# Config Connector

## Enable config connector
https://cloud.google.com/config-connector/docs/how-to/install-upgrade-uninstall?hl=ja#enabling_the
```
$ gcloud container clusters update sandbox-toga4-gke-cluster01-tokyo --update-addons ConfigConnector=ENABLED
```

## Setup service account
https://cloud.google.com/config-connector/docs/how-to/install-upgrade-uninstall?hl=ja#identity
```
$ gcloud iam service-accounts create config-connector \
    --project=sandbox-toga4-gke

$ gcloud projects add-iam-policy-binding sandbox-toga4-gke \
    --project=sandbox-toga4-gke \
    --member="serviceAccount:config-connector@sandbox-toga4-gke.iam.gserviceaccount.com" \
    --role="roles/editor"

$ gcloud iam service-accounts add-iam-policy-binding \
    config-connector@sandbox-toga4-gke.iam.gserviceaccount.com \
    --project=sandbox-toga4-gke \
    --member="serviceAccount:sandbox-toga4-gke.svc.id.goog[cnrm-system/cnrm-controller-manager]" \
    --role="roles/iam.workloadIdentityUser"
```

## Install config connector

```
$ kubectl apply -f config-connector.yaml
$ kubectl annotate namespace default cnrm.cloud.google.com/project-id=sandbox-toga4-gke
```

**Config Connector can't use due to insufficient memory.**

# Setup Workload Identity pool for Github Actions

```
$ PROJECT_ID=sandbox-toga4-gke
$ gcloud iam service-accounts create "github-actions" \
  --project ${PROJECT_ID}

$ gcloud services enable iamcredentials.googleapis.com \
  --project ${PROJECT_ID}

$ gcloud iam workload-identity-pools create "github-actions-pool" \
  --project ${PROJECT_ID} \
  --location="global" \
  --display-name="github actions pool"

$ gcloud iam workload-identity-pools providers create-oidc "github-actions-provider" \
  --project ${PROJECT_ID} \
  --location="global" \
  --workload-identity-pool="github-actions-pool" \
  --display-name="github actions provider" \
  --display-name="Demo provider" \
  --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.aud=assertion.aud,attribute.repository=assertion.repository" \
  --issuer-uri="https://token.actions.githubusercontent.com"

$ WIP_ID=$(gcloud iam workload-identity-pools describe "github-actions-pool" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --format="value(name)")

$ gcloud iam service-accounts add-iam-policy-binding "github-actions@${PROJECT_ID}.iam.gserviceaccount.com" \
  --project="${PROJECT_ID}" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/${WIP_ID}/attribute.repository/toga4/go-api-challange"
```
