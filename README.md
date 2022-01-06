# platform-services
A public PoC with functioning services using a simple Istio Mesh running on Kubernetes

Version : <repo-version>0.9.0-main-aaaagtcfmwj</repo-version>

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/karlmutch/platform-services/blob/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/karlmutch/platform-services)](https://goreportcard.com/report/karlmutch/platform-services)

This project is intended as a sand-box for experimenting with Istio and some example services using robust SaaS services for security and Observability.  It also provides a good way of packaging and testing out non-proprietary platform functions while collaborating with other parties such as vendors and customers.

Because this project is intended to mirror production examples of services deploy it requires that the user has an account and a registered internet domain name.  A service such as domains.google.com, or cloudflare are good places to start.  Your DNS registered host will be used to issue certificates on your behalf to secure the public connections that are exposed by services to the internet, and more specifically to secure the username and password based access exposed by the service mesh.

# Purpose

This proof of concept (PoC) implementation is intended as a means by which commercial teams can experiment with features of Service Mesh, PaaS, and SaaS platforms provided by third parties.  This project serves as a way of exercising services so that code can be openly shared while testing external services and technologies, and for support in relation to external open source offerings in a public support context.

In its current form the PoC is used to deploy two main services, an experiment service and a downstream service. These services are provisioned with a gRPC API and leverage an Authorizationm Athentication, and Accounting (AAA) capability and an Observability platform integration to services offered by thrid parties.

# Installation

These instructions were used with Kubernetes 1.22.x, and Istio 1.12.1

## Development and Building from source

Clone the repository using the following instructions when this is part of a larger project using multiple services:
<pre><code><b>mkdir ~/project
cd ~/project
export GOPATH=`pwd`
export PATH=$GOPATH/bin:$PATH
mkdir -p src/github.com/karlmutch
cd src/github.com/karlmutch
git clone https://github.com/karlmutch/platform-services
cd platform-services
</b></code></pre>

To boostrap development you will need a copy of Go 1.17+. Go installation instructions can be found at, https://golang.org/doc/install.

# Build environment

Builds of the PoC software are intended to be done using Rancher Desktop with a Kubernetes hosted buildkit instance.

## Rancher Desktop (rancher)

Rancher Desktop will install a fully featured Kubernetes distribution and CLI tools include, nerdctl (container image wrangling), kubectl, trivy (image scanning), helm and a docker CLI.  Rancher Desktop can be deployed for Windows Services for Linux 2, Max OSX, LInux and Arch.

Instructions for installation of rancher can be found at https://rancherdesktop.io/.

Rancher will automatically deploy Kubernetes for you and update your kubectl configuration file so that you can immediately access the cluster.

## Running the build using the container

Creating a build container to isolate the build into a versioned environment.  nerdctl will use buildkit which rancher installs when deployed.

<pre><code><b>nerdctl build -t platform-services:latest --build-arg USER=$USER --build-arg USER_ID=`id -u $USER` --build-arg USER_GROUP_ID=`id -g $USER` .
</b></code></pre>

Prior to doing the build a GitHub OAUTH token needs to be defined within your environment if you wish to release the build artifacts.  Use the github admin pages for your account to generate a token.
<pre><code><b>nerdctl run -e GITHUB_TOKEN=$GITHUB_TOKEN -v $GOPATH:/project platform-services ; echo "Done" ; docker container prune -f</b>
</code></pre>

A combined build script is provided 'platform-services/build.sh' to allow all stages of the build including producing container images using nerdctl in order that rancher desktop has access to them.


Download development workflow dependencies into your development environment.

<pre><code><b>
go install github.com/karlmutch/duat/cmd/semver@0.16.0
go install github.com/karlmutch/duat/cmd/github-release@0.16.0
go install github.com/karlmutch/duat/cmd/stencil@0.16.0
go install github.com/karlmutch/duat/cmd/license-detector@0.16.0
go install github.com/karlmutch/petname/cmd/petname@latest
</b></code></pre>

## Software Bill of Materials

The go compiler tools can be used to produce a simple list of all dependencies within this project without explicit licensing information using the following command:

<pre><code><b>
go list -json --mod=readonly ./...
</b></code></pre>

A Software Bill Of Materials can be generated using spdx project tooling as follows:

<pre><code><b>
mkdir spdx-output

# The following command will print a short summary of the types of licenses detected and the highest confidence for each type.
license-detector -r -s -o spdx-output/summary.rpt

license-detector -r -o spdx-output/full.rpt

</b></code></pre>

The license-detector tool uses fuzzy license matching to deduce the license being used.

If cyclonedx is being used as an alternative to SPDX then the cycleonedx-gomod tooling can be used to generate CycloneDX compliant files for use with dependency tracker and TeamCity like tools as follows:

<pre><code><b>
go install github.com/CycloneDX/cyclonedx-gomod@latest
cyclonedx-gomod mod -output bom.xml -assert-licenses -licenses
</b></code></pre>

The dependency tracker API and service can be used with the following commands:

</b></code></pre>
go install github.com/ozonru/dtrack-audit/cmd/dtrack-audit@latest
</b></code></pre>


If SPDX is being used then the following tool uses attestations to generate the license report but these are not common and so the information is very patchy and not as useful.

<pre><code><b>
nerdctl run -v $GOPATH:/repository -v "$(pwd)/spdx-output:/out" spdx/spdx-sbom-generator -p /repository/src/github.com/karlmutch/platform-services -o /out</b>
WARN[0000] kernel support for cgroup blkio weight missing, weight discarded
INFO[2021-12-21T22:15:43Z] Starting to generate SPDX ...
INFO[2021-12-21T22:15:46Z] Running generator for Module Manager: `go-mod` with output `/out/bom-go-mod.spdx`
INFO[2021-12-21T22:15:46Z] Current Language Version go version go1.16.5 linux/amd64
INFO[2021-12-21T22:15:50Z] Command completed successful for below package managers
INFO[2021-12-21T22:15:50Z] Plugin go-mod generated output at /out/bom-go-mod.spdx
</code></pre>

# Deployment Environment

These deployment instructions are intended for use with the Ubuntu 20.04 LTS distribution.

The following instructions make use of the stencil tool for templating configuration files.

This major section describes two basic alternatives for deployment locally hosted Rancher Desktop, and AWS EKS.  Other Kubernetes distribution and deployment models should work but are not explicitly described here.

## Rancher Desktop (rancher)

The development environment makes use of Rancher Desktop so this step will have been already been completed if you have been performing builds.

Rancher Desktop will install a fully featured Kubernetes distribution and CLI tools include, nerdctl (container image wrangling), kubectl, trivy (image scanning), helm and a docker CLI.  Rancher Desktop can be deployed for Windows Services for Linux 2, Max OSX, LInux and Arch.

Instructions for installation of rancher can be found at https://rancherdesktop.io/.

Rancher will automatically deploy Kubernetes for you and update your kubectl configuration file so that you can immediately access the cluster.

## Installing AWS Kubernetes

When using AWS tools some of the developer focused utilities are installed explicitly.

Instructions for installing the various tools can be found in the following list :

* nerdctl an image wrangling tool can be found at, https://github.com/containerd/nerdctl#readme.
* trivy a tool for image scanning can be found at, https://aquasecurity.github.io/trivy/v0.21.3/getting-started/installation/
* helm a tool for Kubernetes package management can be found at, https://helm.sh/docs/intro/install/
* kubectl a CLI tool for running commands against Kubernetes clusters can be found at, https://kubernetes.io/docs/tasks/tools/

The current preferred approach for deploying on AWS is to make use of EKS, via the eksctl tool.

To install eksctl the following should be done.

```
$ curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
$ sudo mv /tmp/eksctl /usr/local/bin
$ eksctl version
```

### Using eksctl with auto-scaling

A basic cluster with auto scaling can be initialized using eksctl and then the addition of the auto-scaler from the Kubernetes project can be used to scale out the project.  The example eks-cluster.yaml file contains the definitions of a cluster within the us-west-2 region, named platform-services.  Before deploying a long lived cluster it is worth while considering cost savings options which are described at the following URL, https://aws.amazon.com/ec2/cost-and-capacity/.

Cluster creation can be performed using the following:

<pre><code><b>export AWS_ACCOUNT=`aws sts get-caller-identity --query Account --output text`
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin $AWS_ACCOUNT.dkr.ecr.us-west-2.amazonaws.com
export AWS_ACCESS_KEY=xxx
export AWS_SECRET_ACCESS_KEY=xxx
export AWS_DEFAULT_REGION=xxx
sudo ntpdate ntp.ubuntu.com
export KUBECONFIG=~/.kube/config
export AWS_CLUSTER_NAME=test-eks
eksctl create cluster -f eks-cluster.yaml</b>
2021-04-06 13:39:37 [ℹ]  eksctl version 0.43.0
2021-04-06 13:39:37 [ℹ]  using region us-west-2
2021-04-06 13:39:37 [ℹ]  subnets for us-west-2a - public:192.168.0.0/19 private:192.168.96.0/19
2021-04-06 13:39:37 [ℹ]  subnets for us-west-2b - public:192.168.32.0/19 private:192.168.128.0/19
2021-04-06 13:39:37 [ℹ]  subnets for us-west-2d - public:192.168.64.0/19 private:192.168.160.0/19
2021-04-06 13:39:37 [ℹ]  nodegroup "overhead" will use "ami-0a93391193b512e5d" [AmazonLinux2/1.19]
2021-04-06 13:39:37 [ℹ]  using SSH public key "/home/kmutch/.ssh/id_rsa.pub" as "eksctl-test-eks-nodegroup-overhead-be:07:a0:27:44:d8:27:04:c2:ba:28:fa:8c:47:7f:09"
2021-04-06 13:39:37 [ℹ]  using Kubernetes version 1.19
2021-04-06 13:39:37 [ℹ]  creating EKS cluster "test-eks" in "us-west-2" region with un-managed nodes
2021-04-06 13:39:37 [ℹ]  1 nodegroup (overhead) was included (based on the include/exclude rules)
2021-04-06 13:39:37 [ℹ]  will create a CloudFormation stack for cluster itself and 1 nodegroup stack(s)
2021-04-06 13:39:37 [ℹ]  will create a CloudFormation stack for cluster itself and 0 managed nodegroup stack(s)
2021-04-06 13:39:37 [ℹ]  if you encounter any issues, check CloudFormation console or try 'eksctl utils describe-stacks --region=us-west-2 --cluster=test-eks'
2021-04-06 13:39:37 [ℹ]  Kubernetes API endpoint access will use default of {publicAccess=true, privateAccess=false} for cluster "test-eks" in "us-west-2"
2021-04-06 13:39:37 [ℹ]  2 sequential tasks: { create cluster control plane "test-eks", 3 sequential sub-tasks: { 3 sequential sub-tasks: { wait for control plane to
 become ready, tag cluster, update CloudWatch logging configuration }, create addons, create nodegroup "overhead" } }
2021-04-06 13:39:37 [ℹ]  building cluster stack "eksctl-test-eks-cluster"
2021-04-06 13:39:38 [ℹ]  deploying stack "eksctl-test-eks-cluster"
2021-04-06 13:40:08 [ℹ]  waiting for CloudFormation stack "eksctl-test-eks-cluster"
2021-04-06 13:40:38 [ℹ]  waiting for CloudFormation stack "eksctl-test-eks-cluster"
...
2021-04-06 13:52:39 [ℹ]  waiting for CloudFormation stack "eksctl-test-eks-cluster"
2021-04-06 13:53:39 [ℹ]  waiting for CloudFormation stack "eksctl-test-eks-cluster"
2021-04-06 13:53:39 [✔]  tagged EKS cluster (environment=test-eks)
2021-04-06 13:53:40 [ℹ]  waiting for requested "LoggingUpdate" in cluster "test-eks" to succeed
2021-04-06 13:53:57 [ℹ]  waiting for requested "LoggingUpdate" in cluster "test-eks" to succeed
2021-04-06 13:54:14 [ℹ]  waiting for requested "LoggingUpdate" in cluster "test-eks" to succeed
2021-04-06 13:54:33 [ℹ]  waiting for requested "LoggingUpdate" in cluster "test-eks" to succeed
2021-04-06 13:54:34 [✔]  configured CloudWatch logging for cluster "test-eks" in "us-west-2" (enabled types: audit, authenticator, controllerManager & disabled types
: api, scheduler)
2021-04-06 13:54:34 [ℹ]  building nodegroup stack "eksctl-test-eks-nodegroup-overhead"
2021-04-06 13:54:34 [ℹ]  deploying stack "eksctl-test-eks-nodegroup-overhead"
2021-04-06 13:54:34 [ℹ]  waiting for CloudFormation stack "eksctl-test-eks-nodegroup-overhead"
2021-04-06 13:54:50 [ℹ]  waiting for CloudFormation stack "eksctl-test-eks-nodegroup-overhead"
...
2021-04-06 13:57:48 [ℹ]  waiting for CloudFormation stack "eksctl-test-eks-nodegroup-overhead"
2021-04-06 13:58:06 [ℹ]  waiting for CloudFormation stack "eksctl-test-eks-nodegroup-overhead"
2021-04-06 13:58:06 [ℹ]  waiting for the control plane availability...
2021-04-06 13:58:06 [✔]  saved kubeconfig as "/home/kmutch/.kube/config"
2021-04-06 13:58:06 [ℹ]  no tasks
2021-04-06 13:58:06 [✔]  all EKS cluster resources for "test-eks" have been created
2021-04-06 13:58:06 [ℹ]  adding identity "arn:aws:iam::613076437200:role/eksctl-test-eks-nodegroup-overhea-NodeInstanceRole-Q1RVPO36W4VJ" to auth ConfigMap
2021-04-06 13:58:08 [ℹ]  kubectl command should work with "/home/kmutch/.kube/config", try 'kubectl get nodes'
2021-04-06 13:58:08 [✔]  EKS cluster "test-eks" in "us-west-2" region is ready
<b>kubectl get pods --namespace kube-system</b>
NAME                       READY   STATUS    RESTARTS   AGE
coredns-6548845887-h82sj   1/1     Running   0          10m
coredns-6548845887-wz7dm   1/1     Running   0          10m
</code></pre>

Now the auto scaler can be deployed.

<pre><code><b>
kubectl apply -f eks-scaler.yaml</b>
serviceaccount/cluster-autoscaler created
clusterrole.rbac.authorization.k8s.io/cluster-autoscaler created
role.rbac.authorization.k8s.io/cluster-autoscaler created
clusterrolebinding.rbac.authorization.k8s.io/cluster-autoscaler created
rolebinding.rbac.authorization.k8s.io/cluster-autoscaler created
deployment.apps/cluster-autoscaler created
</code></pre>

## Istio 1.12.x

Istio affords a control layer on top of the k8s data plane.  Instructions for deploying Istio are the vanilla instructions that can be found at, https://istio.io/docs/setup/getting-started/#install.  Istio was at one time a Helm based installation but has since moved to using its own methodology, this is the reason we dont use arkade to install it.

<pre><code><b>cd $HOME
curl -q -LO https://github.com/istio/istio/releases/download/1.12.1/istio-1.12.1-linux-amd64.tar.gz
tar xzf istio-1.12.1-linux-amd64.tar.gz
export ISTIO_DIR=$HOME/istio-1.12.1
export PATH=$ISTIO_DIR/bin:$PATH
cd -
istioctl install --set profile=demo -y
</b>
✔ Istio core installed
✔ Istiod installed
✔ Egress gateways installed
✔ Ingress gateways installed
✔ Installation complete
</code></pre>

In order to access you cluster you will need to define some environment variables that will be used later in these instructions:

<pre><code><b>
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
</b></code></pre>

## Helm 3 Kubernetes package manager

Helm is used by several packages that are deployed using Kubernetes.  Helm can be installed using instructions found at, https://helm.sh/docs/using\_helm/#installing-helm.  For snap based linux distributions the following can be used as a quick-start.

<pre><code><b>
curl https://baltocdn.com/helm/signing.asc | sudo apt-key add -
sudo apt-get install apt-transport-https --yes
echo "deb https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
sudo apt-get update
sudo apt-get install helm

kubectl create serviceaccount --namespace kube-system tiller
kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
helm repo update
</b></code></pre>

## Encryption

This section describes how to secure traffic into the service mesh ingress.  Lets Encrypt is an internet based service for provisioning certificates.  If you are using a locally deployed mesh within rancher desktop for example you will need to use the minica approach.

### minica

Minica is a simple CA intended for use in situations where the CA operator also operates each host where a certificate will be used. It automatically generates both a key and a certificate when asked to produce a certificate. It does not offer OCSP or CRL services. Minica is appropriate, for instance, for generating certificates for RPC systems or microservices.

More information about minica can be found at, https://github.com/jsha/minica.

<pre><code>
$ go install github.com/jsha/minica@latest
$ mkdir minica
$ cd minica
$ minica -domains platform-services.karlmutch.com
$ cd -
$ tree minica
minica
├── minica-key.pem
├── minica.pem
└── platform-services.karlmutch.com
    ├── cert.pem
    └── key.pem

1 directory, 4 filess
</code></pre>

### Lets Encrypt

letsencrypt is a public SSL/TLS certificate provider intended for use with the open internet that is being used to secure our service mesh for this project. The lets encrypt provisioning tool can be installed from github and accessed to produce TLS certificates for your service.  If you are deploying on rancher desktop or a local installation then lets encrypt is not supported and using minica, https://github.com/jsha/minica, is recommended instead for testing purposes.

Prior to running the lets encrypt tools you should identify the desired DNS hostname and email you wish to use for your example service cluster.  In our example we have a domain registered as, karlmutch.com.  This domain is available to us as an administrator, and we have choosen to use the host name platform-service.karlmutch.com as the services hostname.

The first step is to add a registered hosts entry for the platform-services.karlmutch.com host into the DNS account, if the host is unknown add an IP address such as 127.0.0.1.  During the generation process you will be prompted to add a DNS TXT record into the custom resource records for the domain, this requires the dummy entry to be present.

Setting up and initiating this process can be done using the following:

<pre><code><b>git clone https://github.com/letsencrypt/letsencrypt</b>
Cloning into 'letsencrypt'...
remote: Enumerating objects: 255, done.
remote: Counting objects: 100% (255/255), done.
remote: Compressing objects: 100% (188/188), done.
remote: Total 71278 (delta 135), reused 110 (delta 65), pack-reused 71023
Receiving objects: 100% (71278/71278), 23.55 MiB | 26.53 MiB/s, done.
Resolving deltas: 100% (52331/52331), done.
<b>cd letsencrypt</b>
<b>./letsencrypt-auto certonly --rsa-key-size 4096 --agree-tos --manual --preferred-challenges=dns --email=karlmutch@cognizant.com -d platform-services.karlmutch.com</b>
</code></pre>

You will be prompted with the IP address logging when starting the script, you should choose 'Y' to enabled the logging as this assists auditing of DNS changes on the internet by registras and regulatory bodies.

<pre><code>
Are you OK with your IP being logged?
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
(Y)es/(N)o: <b>Y
</b></code></pre>

After this step you will be asked to add a text record to your DNS records proving that you have control over the domain, this will appear much like the following:

<pre><code>
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
Please deploy a DNS TXT record under the name
acme-challenge.platform-services.karlmutch.com with the following value:

mbUa9_gb4RVYhTumHy3zIi3PIXFh0k_oOgCie4NvhqQ

Before continuing, verify the record is deployed.
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
Press Enter to Continue
</code></pre>

You should wait 10 to 20 minutes for the TXT record to appear in the database of your DNS provider before selecting continue otherwise the verification will fail and you will need to restart it.

<pre><code>
Waiting for verification...
Cleaning up challenges

IMPORTANT NOTES:
 - Congratulations! Your certificate and chain have been saved at:
   /etc/letsencrypt/live/platform-services.karlmutch.com/fullchain.pem
   Your key file has been saved at:
   /etc/letsencrypt/live/platform-services.karlmutch.com/privkey.pem
   Your cert will expire on 2020-02-23. To obtain a new or tweaked
   version of this certificate in the future, simply run
   letsencrypt-auto again. To non-interactively renew *all* of your
   certificates, run "letsencrypt-auto renew"
 - If you like Certbot, please consider supporting our work by:

   Donating to ISRG / Let's Encrypt:   https://letsencrypt.org/donate
   Donating to EFF:                    https://eff.org/donate-le
</code></pre>

Once the certificate generation is complete you will have the certificate saved at the location detailed in the 'IMPORTANT NOTES' section.  Keep a record of where this is you will need it later.

After the certificate has been issue feel free to delete the TXT record that served as proof of ownership as it is no longer needed.

## Configuration of secrets and cluster access

This project makes use of several secrets that are used to access resources under its control, including the Postgres Database, the Honeycomb service, and the lets encrypt issues certificate.

The experiment service Honeycomb observability solution uses a key to access Datasets defined by the Honeycomb account and store events in the same.  Configuring the service is done by creating a Kubernetes secret.  For now we can define the Honeycomb Dataset, and API key using an environment variable and when we deploy the secrets for the Postgres Database the secret for the API will be injected using the stencil tool.

<pre><code><b>export O11Y_KEY=a55676f62a47474b8491a
export O11Y_DATASET=platform-services
</b></code></pre>

The services also use a postgres Database instance to persist experiment data, this is installed later in the process.  The following shows an example of what should be defined for Postgres support prior to running the stencil command:

<pre><code><b>export PGRELEASE=$USER-poc
export PGHOST=$PGRELEASE-postgresql.default.svc.cluster.local
export PORT=5432
export PGUSER=postgres
export PGPASSWORD=p355w0rd
export PGDATABASE=platform
</b></code></pre>

<pre><code><b>
stencil < cmd/experimentsrv/secret.yaml | kubectl apply -f -
</b></code></pre>

The last set of secrets that need to be stored are related to securing the mesh for third party connections using TLS.  This secret contains the full certificate chain and private key needed to implement TLS on the gRPC connections exposed by the mesh.

If you used Lets Encrypt then the following applies:

<pre><code><b>
sudo kubectl create -n istio-system secret generic platform-services-tls-cert \
    --from-file=key=/etc/letsencrypt/live/platform-services.karlmutch.com/privkey.pem \
    --from-file=cert=/etc/letsencrypt/live/platform-services.karlmutch.com/fullchain.pem
</b></code></pre>

If minica was used then the following would be used:

<pre><code><b>
kubectl create -n istio-system secret generic platform-services-tls-cert \
    --from-file=key=./minica/platform-services.karlmutch.com/key.pem \
    --from-file=cert=./minica/platform-services.karlmutch.com/cert.pem
</b></code></pre>

If you are using AWS ACM to manage your certificates the platform-services-tls-cert secret is not required.

## Deploying the Istio Ingress configuration

Istio provides an ingress resource that can be used to secure the service using either secrets (certificates) or using cloud provider provisioned certificates.

Using your own secrets you will use the default ingress yaml that will point at the platform-services-tls-cert Kubernetes provisioned in the previous section.

<pre><code><b>
cp ingress.yaml ingress-$USER.yaml
</b></code></pre>

You MUST also modify the ingress.yaml file to contain the appropriate information for the JWT values that are specific to your Auth0 applicaiton.

<pre><code><b>
vim ingress-$USER.yaml
</b></code></pre>

<pre><code><b>
kubectl apply -f ingress-$USER.yaml
</b></code></pre>

### Deploying using AWS Route 53 and ACM Public CA

Create a DNS hosted zone for your cluster.  Add an A record for your clusters ingress that points at your cluster using the load balancer address identified via the following commands:

<pre><code><b>
kubectl get svc istio-ingressgateway -n istio-system</b>
NAME                   TYPE           CLUSTER-IP     EXTERNAL-IP                                                               PORT(S)
                                       AGE
istio-ingressgateway   LoadBalancer   10.100.90.52   acaaf2fd65f1844b09ed91ee1b409811-1454371712.us-west-2.elb.amazonaws.com   15021:31220/TCP,80:31317/TCP,443:31536/TCP,31400:31243/TCP,15443:30506/TCP   38s
<b>nslookup</b>
> <b>acaaf2fd65f1844b09ed91ee1b409811-1454371712.us-west-2.elb.amazonaws.com</b>
Server:         127.0.0.53
Address:        127.0.0.53#53

Non-authoritative answer:
Name:   acaaf2fd65f1844b09ed91ee1b409811-1454371712.us-west-2.elb.amazonaws.com
Address: 35.163.132.217
Name:   acaaf2fd65f1844b09ed91ee1b409811-1454371712.us-west-2.elb.amazonaws.com
Address: 54.70.157.123
> <b>[ENTER]</b>
<b>cat <<EOF >changes.json
{
            "Comment": "Add LB to the Route53 records so that certificates can be validated against the host name ",
            "Changes": [{
            "Action": "UPSERT",
                        "ResourceRecordSet": {
                                    "Name": "platform-services.karlmutch.com",
                                    "Type": "A",
                                    "TTL": 300,
                                 "ResourceRecords": [{ "Value": "35.163.132.217"}, {"Value": "54.70.157.123"}]
}}]
}
EOF
aws route53 list-hosted-zones-by-name --dns-name karlmutch.com</b>

{
    "HostedZones": [
        {
            "Id": "/hostedzone/Z1UZMEEDVXVHH3",
            "Name": "karlmutch.com.",
            "CallerReference": "RISWorkflow-RD:97f5f182-4b86-41f1-9a24-46218ab70d25",
            "Config": {
                "Comment": "HostedZone created by Route53 Registrar",
                "PrivateZone": false
            },
            "ResourceRecordSetCount": 16
        }
    ],
    "DNSName": "karlmutch.com",
    "IsTruncated": false,
    "MaxItems": "100"
}
<b>aws route53 change-resource-record-sets --hosted-zone-id Z1UZMEEDVXVHH3  --change-batch file://changes.json</b>
{
    "ChangeInfo": {
        "Id": "/change/C029593935XMLAA0T5OHH",
        "Status": "PENDING",
        "SubmittedAt": "2021-04-07T22:12:17.619Z",
        "Comment": "Add LB to the Route53 records so that certificates can be validated against the host name "
    }
}
<b>aws route53  get-change --id /change/C029593935XMLAA0T5OHH</b>
{
    "ChangeInfo": {
        "Id": "/change/C029593935XMLAA0T5OHH",
        "Status": "INSYNC",
        "SubmittedAt": "2021-04-07T22:12:17.619Z",
        "Comment": "Add LB to the Route53 records so that certificates can be validated against the host name "
    }
}
</code></pre>

Now we are in a position to generate the certificate:
<pre><code>
<b>aws acm request-certificate --domain-name platform-services.karlmutch.com --validation-method DNS --idempotency-token 1234</b>
{
    "CertificateArn": "arn:aws:acm:us-west-2:613076437200:certificate/dcd2ca31-f27c-45ac-85a7-539688f8e4cb"
}
{
    "Certificate": {
        "CertificateArn": "arn:aws:acm:us-west-2:613076437200:certificate/dcd2ca31-f27c-45ac-85a7-539688f8e4cb",
        "DomainName": "platform-services.karlmutch.com",
        "SubjectAlternativeNames": [
            "platform-services.karlmutch.com"
        ],
        "DomainValidationOptions": [
            {
                "DomainName": "platform-services.karlmutch.com",
                "ValidationDomain": "platform-services.karlmutch.com",
                "ValidationStatus": "PENDING_VALIDATION",
                "ResourceRecord": {
                    "Name": "_a9d4b51d79d2b08121a0796cbfbb7a68.platform-services.karlmutch.com.",
                    "Type": "CNAME",
                    "Value": "_e327e9f51160630a9f0056fd3eb56a74.bbfvkzsszw.acm-validations.aws."
                },
                "ValidationMethod": "DNS"
            }
        ],
        "Subject": "CN=platform-services.karlmutch.com",
        "Issuer": "Amazon",
        "CreatedAt": 1617837918.0,
        "Status": "PENDING_VALIDATION",
        "KeyAlgorithm": "RSA-2048",
        "SignatureAlgorithm": "SHA256WITHRSA",
        "InUseBy": [],
        "Type": "AMAZON_ISSUED",
        "KeyUsages": [],
        "ExtendedKeyUsages": [],
        "RenewalEligibility": "INELIGIBLE",
        "Options": {
            "CertificateTransparencyLoggingPreference": "ENABLED"
        }
    }
}
<b>cat <<EOF >changes.json
{
            "Comment": "Add the certificate issuance validation to the Route53 records so that certificates can be validated",
            "Changes": [{
            "Action": "UPSERT",
                        "ResourceRecordSet": {
                                 "Name": "_a9d4b51d79d2b08121a0796cbfbb7a68.platform-services.karlmutch.com.",
                                 "Type": "CNAME",
                                 "TTL": 300,
                                 "ResourceRecords": [{ "Value": "_e327e9f51160630a9f0056fd3eb56a74.bbfvkzsszw.acm-validations.aws."}]
}}]
}
EOF
aws route53 change-resource-record-sets --hosted-zone-id Z1UZMEEDVXVHH3  --change-batch file://changes.json</b>
{
    "ChangeInfo": {
        "Id": "/change/C084369336HI4IZ12CDU9",
        "Status": "PENDING",
        "SubmittedAt": "2021-04-07T23:46:04.831Z",
        "Comment": "Add the certificate issuance validation to the Route53 records so that certificates can be validated"
    }
}
<b>aws route53  get-change --id /change/C084369336HI4IZ12CDU9</b>
{
    "ChangeInfo": {
        "Id": "/change/C084369336HI4IZ12CDU9",
        "Status": "INSYNC",
        "SubmittedAt": "2021-04-07T23:46:04.831Z",
        "Comment": "Add the certificate issuance validation to the Route53 records so that certificates can be validated"
    }
}
# Now we wait for the certificate to be issued:
<b>aws acm describe-certificate --certificate-arn arn:aws:acm:us-west-2:613076437200:certificate/dcd2ca31-f27c-45ac-85a7-539688f8e4cb<b>
{
    "Certificate": {
        "CertificateArn": "arn:aws:acm:us-west-2:613076437200:certificate/dcd2ca31-f27c-45ac-85a7-539688f8e4cb",
        "DomainName": "platform-services.karlmutch.com",
        "SubjectAlternativeNames": [
            "platform-services.karlmutch.com"
        ],
        "DomainValidationOptions": [
            {
                "DomainName": "platform-services.karlmutch.com",
                "ValidationDomain": "platform-services.karlmutch.com",
                "ValidationStatus": "SUCCESS",
                "ResourceRecord": {
                    "Name": "_a9d4b51d79d2b08121a0796cbfbb7a68.platform-services.karlmutch.com.",
                    "Type": "CNAME",
                    "Value": "_e327e9f51160630a9f0056fd3eb56a74.bbfvkzsszw.acm-validations.aws."
                },
                "ValidationMethod": "DNS"
            }
        ],
        "Serial": "07:df:33:6b:78:11:e2:a3:ee:f1:54:51:3f:81:78:28",
        "Subject": "CN=platform-services.karlmutch.com",
        "Issuer": "Amazon",
        "CreatedAt": 1617837918.0,
        "IssuedAt": 1617839861.0,
        "Status": "ISSUED",
        "NotBefore": 1617753600.0,
        "NotAfter": 1651881599.0,
        "KeyAlgorithm": "RSA-2048",
        "SignatureAlgorithm": "SHA256WITHRSA",
        "InUseBy": [],
        "Type": "AMAZON_ISSUED",
        "KeyUsages": [
            {
                "Name": "DIGITAL_SIGNATURE"
            },
            {
                "Name": "KEY_ENCIPHERMENT"
            }
        ],
        "ExtendedKeyUsages": [
            {
                "Name": "TLS_WEB_SERVER_AUTHENTICATION",
                "OID": "1.3.6.1.5.5.7.3.1"
            },
            {
                "Name": "TLS_WEB_CLIENT_AUTHENTICATION",
                "OID": "1.3.6.1.5.5.7.3.2"
            }
        ],
        "RenewalEligibility": "INELIGIBLE",
        "Options": {
            "CertificateTransparencyLoggingPreference": "ENABLED"
        }
    }
}
</code></pre>

If you are using AWS to provision certificates to secure the ingress connections then use aws-ingress.yaml

<pre><code><b>
kubectl apply -f aws-ingress.yaml
</b></code></pre>

Then we patch the ingress so that it uses our newly issued ACM certificate via the AWS ARN.

<pre><code>
<b>arn="arn:aws:acm:us-west-2:613076437200:certificate/dcd2ca31-f27c-45ac-85a7-539688f8e4cb"</b>
<b>kubectl -n istio-system patch service istio-ingressgateway --patch "$(cat<<EOF
metadata:
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-ssl-cert: $arn
    service.beta.kubernetes.io/aws-load-balancer-backend-protocol: tcp
    service.beta.kubernetes.io/aws-load-balancer-ssl-ports: "https"
    service.beta.kubernetes.io/aws-load-balancer-connection-idle-timeout: "3600"
EOF
)"</b>
</code></pre>


A test of the certificate on an empty ingress will then appear as follows:

<pre><code>
<b>curl -Iv https://platform-services.karlmutch.com</b>
* Rebuilt URL to: https://platform-services.karlmutch.com/
*   Trying 35.163.132.217...
* TCP_NODELAY set
* Connected to platform-services.karlmutch.com (35.163.132.217) port 443 (#0)
* ALPN, offering h2
* ALPN, offering http/1.1
* successfully set certificate verify locations:
*   CAfile: /etc/ssl/certs/ca-certificates.crt
  CApath: /etc/ssl/certs
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
* TLSv1.3 (IN), TLS handshake, Server hello (2):
* TLSv1.2 (IN), TLS handshake, Certificate (11):
* TLSv1.2 (IN), TLS handshake, Server key exchange (12):
* TLSv1.2 (IN), TLS handshake, Server finished (14):
* TLSv1.2 (OUT), TLS handshake, Client key exchange (16):
* TLSv1.2 (OUT), TLS change cipher, Client hello (1):
* TLSv1.2 (OUT), TLS handshake, Finished (20):
* TLSv1.2 (IN), TLS handshake, Finished (20):
* SSL connection using TLSv1.2 / ECDHE-RSA-AES128-GCM-SHA256
* ALPN, server did not agree to a protocol
* Server certificate:
*  subject: CN=platform-services.karlmutch.com
*  start date: Apr  7 00:00:00 2021 GMT
*  expire date: May  6 23:59:59 2022 GMT
*  subjectAltName: host "platform-services.karlmutch.com" matched cert's "platform-services.karlmutch.com"
*  issuer: C=US; O=Amazon; OU=Server CA 1B; CN=Amazon
*  SSL certificate verify ok.
> HEAD / HTTP/1.1
> Host: platform-services.karlmutch.com
> User-Agent: curl/7.58.0
> Accept: */*
>
* TLSv1.2 (IN), TLS alert, Client hello (1):
* Empty reply from server
* Connection #0 to host platform-services.karlmutch.com left intact
curl: (52) Empty reply from server
</code></pre>

## Deploying the Observability proxy server

This proxy server is used to forward tracing and metrics from your istio mesh based deployment to the Honeycomb service.

<pre><code><b>
helm repo rm honeycomb || true
helm repo add honeycomb https://honeycombio.github.io/helm-charts
helm install opentelemetry-collector honeycomb/opentelemetry-collector --set honeycomb.apiKey=$O11Y_KEY --set honeycomb.dataset=$O11Y_DATASET
</b></code></pre>

In order to instrument the base Kubernetes deployment for use with honeycomb you should follow the instructions found at https://docs.honeycomb.io/getting-data-in/integrations/kubernetes/.

The dataset used by the istio and services deployed within this project also needs configuration to allow the Honeycomb platform to identify import fields.  Once data begins flowing into the data set you can navigate to the definitions section for the dataset and set the 'Name' item to the name field, 'Parent span ID' item to parentId, 'Service name' to serviceName, 'Span duration' to durationMs, 'Span ID' to id, and finally 'Trace ID' to traceId.

### Postgres DB

To deploy the platform experiment service a database must be present.  The PoC is intended to use an in-cluster DB designed that is dropped after the service is destroyed.

If you wish to use Aurora then you will need to use the AWS CLI Tools or the Web Console to create your Postgres Database appropriately, and then to set your environment variables PGHOST, PGPORT, PGUSER, PGPASSWORD, and PGDATABASE appropriately.  You will also be expected to run the sql setup scripts yourself.

The first step is to install the postgres 14 client on your system and then to populate the schema on the remote database:

<pre><code><b></b>
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" >> /etc/apt/sources.list.d/pgdg.list'
sudo apt-get -y update
sudo apt-get -y --upgrade install postgresql-client-14
</code></pre>

### Deploying an in-cluster Database

This section gives guidence on how to install an in-cluster database for use-cases where data persistence beyond a single deployment is not a concern.  These instructions are therefore limited to testing only scenarios.  For information concerning Kubernetes storage strategies you should consult other sources and read about stateful sets in Kubernetes.  In production using a single source of truth then cloud provider offerings such as AWS Aurora are recommended.

A secrets file containing host information, passwords and other secrets is assumed to have already been applied using the instructions several sections above.  The secrets are needed to allows access to the postgres DB, and/or other external resources.  YAML files will be needed to populate secrets into the service mesh, individual services document the secrets they require within their README.md files found on github and provide examples, for example https://github.com/karlmutch/platform-services/cmd/experimentsrv/README.md.

In order to deploy Postgres this document describes a helm based approach.  The bitnami postgresql distribution can be installed using the following:

<pre><code><b>
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install $PGRELEASE \
  --set postgresqlPassword=$PGPASSWORD,postgresqlDatabase=postgres\
  bitnami/postgresql
</b>
NAME: wondrous-sturgeon
LAST DEPLOYED: Mon Dec 28 12:08:20 2020
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
** Please be patient while the chart is being deployed **

PostgreSQL can be accessed via port 5432 on the following DNS name from within your cluster:

    wondrous-sturgeon-postgresql.default.svc.cluster.local - Read/Write connection

To get the password for "postgres" run:

    export POSTGRES_PASSWORD=$(kubectl get secret --namespace default wondrous-sturgeon-postgresql -o jsonpath="{.data.postgresql-password}" | base64 --decode)

To connect to your database run the following command:

    kubectl run wondrous-sturgeon-postgresql-client --rm --tty -i --restart='Never' --namespace default --image docker.io/bitnami/postgresql:11.10.0-debian-10-r24 --env="PGPASS
WORD=$POSTGRES_PASSWORD" --command -- psql --host wondrous-sturgeon-postgresql -U postgres -d postgres -p 5432



To connect to your database from outside the cluster execute the following commands:

    kubectl port-forward --namespace default svc/wondrous-sturgeon-postgresql 5432:5432 &
    PGPASSWORD="$POSTGRES_PASSWORD" psql --host 127.0.0.1 -U postgres -d postgres -p 5432
</code></pre>

Special note should be taken of the output from the helm command it has a lot of valuable information concerning your postgres deployment that will be needed when you load the database schema.

Setting up the proxy will be needed prior to running the SQL database provisioning scripts.  When doing this prior to running the postgres client set the PGHOST environment variable to 127.0.0.1 so that the proxy on the localhost is used.  The proxy will timeout after inactivity and shutdown so be prepared to restart it when needed.

<pre><code><b>
kubectl wait --for=condition=Ready pod/$PGRELEASE-postgresql-0
kubectl port-forward --namespace default svc/$PGRELEASE-postgresql 5432:5432 &amp;
sleep 2
PGHOST=127.0.0.1 PGDATABASE=platform psql -f sql/platform.sql -d postgres
</b></code></pre>

Further information about how to deployed the service specific database for the experiment service for example can be found in the cmd/experiment/README.md file.

## Deploying into the Istio mesh

This section describes the activities for deployment of the services provided by the Plaform PoC.  The first two sections provide a description of how to deploy the TLS secured ingress for the service mesh, the first being for cloud provisioned systems and the second for the localized rancher desktop based deployment.

### Configuring the Rancher Desktop ingress (Critical Step)

<pre><code><b>
export INGRESS_HOST=$(kubectl get po -l istio=ingressgateway -n istio-system -o jsonpath='{.items[0].status.hostIP}')
export CLUSTER_INGRESS=$INGRESS_HOST:$SECURE_INGRESS_PORT
</b></code></pre>

### Configuring the service cloud based DNS (Critical Step)

When using this mesh instance with a TLS based deployment the DNS domain name used for the LetsEncrypt certificate (CN), will need to have its address record (A) updated to point at the AWS load balancer assigned to the Kubernetes cluster.  In AWS this is done via Route53:

<pre><code><b>
INGRESS_HOST=`kubectl get svc --namespace istio-system -o go-template='{{range .items}}{{range .status.loadBalancer.ingress}}{{.hostname}}{{printf "\n"}}{{end}}{{end}}'`
dig +short $INGRESS_HOST
</b></code></pre>

Take the IP addresses from the above output and use these as the A record for the LetsEncrypt host name, platform-services.karlmutch.com, and this will enable accessing the mesh and validation of the common name (CN) in the certificate.  After several minutes, you should test this step by using the following command to verify that the certificate negotiation is completed.

### Service deployment overview

Platform services use Dockerfiles to encapsulate their build steps which are documented within their respective README.md files.  Building services are single step CLI operations and require only the installation of Docker, and any version of Go 1.7 or later.  Builds will produce containers and will upload these to your current AWS account users ECS docker registry.  Deployments are staged from this registry.  

<pre><code><b>kubectl get nodes</b>
NAME                                           STATUS    ROLES     AGE       VERSION
ip-172-20-118-127.us-west-2.compute.internal   Ready     node      17m       v1.9.3
ip-172-20-41-63.us-west-2.compute.internal     Ready     node      17m       v1.9.3
ip-172-20-55-189.us-west-2.compute.internal    Ready     master    18m       v1.9.3
</code></pre>

Once secrets are loaded individual services can be deployed from a checked out developer copy of the service repo using a command like the following :

<pre><code><b>cd ~/project/src/github.com/karlmutch/platform-services</b>
<b>cd cmd/[service] ; kubectl apply -f \<(stencil [service].yaml 2>/dev/null); cd - </b>
</code></pre>

A minimal set of servers for test is to use the downstream and experiment servers as follows:

```
kubectl apply -f <(cd cmd/downstream > /dev/null ; stencil < downstream.yaml) && \
kubectl apply -f <(cd cmd/experimentsrv > /dev/null ; stencil < experimentsrv.yaml)
```

In order to locate the image repository the stencil tool will test for the presence of AWS credentials and if found will use the account as the source of AWS ECR images.  In the case where the credentials are not present then the default localhost registry will be used for image deployment.

Once the application is deployed you can discover the gateway points within the kubernetes cluster by using the kubectl commands as documented in the cmd/experimentsrv/README.md file.

More information about deploying a real service and using the experimentsrv server can be found at, https://github.com/karlmutch/platform-services/blob/master/cmd/experimentsrv/README.md.

### Testing connectivity using Rancher Desktop

Once the initial services have been deployed the connectivity can be tested using the following:

<pre><code><b>curl -Iv --cacert minica/platform-services.karlmutch.com/../minica.pem --header "Host: platform-services.karlmutch.com"  --connect-to "platform-services.karlmutch.com:$SECURE_INGRESS_PORT:$CLUSTER_INGRESS" https://platform-services.karlmutch.com:$SECURE_INGRESS_PORT
</b>
* Rebuilt URL to: https://platform-services.karlmutch.com:30017/
* Connecting to hostname: 172.20.0.2
* Connecting to port: 30017
*   Trying 172.20.0.2...
* TCP\_NODELAY set
* Connected to 172.20.0.2 (172.20.0.2) port 30017 (#0)
* ALPN, offering h2
* ALPN, offering http/1.1
* successfully set certificate verify locations:
*   CAfile: minica/platform-services.karlmutch.com/../minica.pem
  CApath: /etc/ssl/certs
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
* TLSv1.3 (IN), TLS handshake, Server hello (2):
* TLSv1.3 (IN), TLS Unknown, Certificate Status (22):
* TLSv1.3 (IN), TLS handshake, Unknown (8):
* TLSv1.3 (IN), TLS handshake, Certificate (11):
* TLSv1.3 (IN), TLS handshake, CERT verify (15):
* TLSv1.3 (IN), TLS handshake, Finished (20):
* TLSv1.3 (OUT), TLS change cipher, Client hello (1):
* TLSv1.3 (OUT), TLS Unknown, Certificate Status (22):
* TLSv1.3 (OUT), TLS handshake, Finished (20):
* SSL connection using TLSv1.3 / TLS\_AES\_256\_GCM\_SHA384
* ALPN, server accepted to use h2
* Server certificate:
*  subject: CN=platform-services.karlmutch.com
*  start date: Dec 23 18:10:54 2020 GMT
*  expire date: Jan 22 18:10:54 2023 GMT
*  subjectAltName: host "platform-services.karlmutch.com" matched cert's "platform-services.karlmutch.com"
*  issuer: CN=minica root ca 34d475
*  SSL certificate verify ok.
* Using HTTP2, server supports multi-use
* Connection state changed (HTTP/2 confirmed)
* Copying HTTP/2 data in stream buffer to connection buffer after upgrade: len=0
* TLSv1.3 (OUT), TLS Unknown, Unknown (23):
* TLSv1.3 (OUT), TLS Unknown, Unknown (23):
* TLSv1.3 (OUT), TLS Unknown, Unknown (23):
* Using Stream ID: 1 (easy handle 0x5584d85765c0)
* TLSv1.3 (OUT), TLS Unknown, Unknown (23):
\> HEAD / HTTP/2
\> Host: platform-services.karlmutch.com
\> User-Agent: curl/7.58.0
\> Accept: */*
\>
* TLSv1.3 (IN), TLS Unknown, Certificate Status (22):
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
* TLSv1.3 (IN), TLS Unknown, Unknown (23):
* Connection state changed (MAX\_CONCURRENT\_STREAMS updated)!
* TLSv1.3 (OUT), TLS Unknown, Unknown (23):
* TLSv1.3 (IN), TLS Unknown, Unknown (23):
\< HTTP/2 404
HTTP/2 404
\< date: Mon, 04 Jan 2021 19:50:31 GMT
date: Mon, 04 Jan 2021 19:50:31 GMT
\< server: istio-envoy
server: istio-envoy

\<
* Connection #0 to host 172.20.0.2 left intact
</code></pre>

### Testing connectivity using Cloud based solutions

A very basic test of the TLS negotiation can be done using the curl command:

<pre><code><b>
curl -Iv https://helloworld.letsencrypt.org
curl -Iv https://platform-services.karlmutch.com:$SECURE_INGRESS_PORT
export INGRESS_HOST=platform-services.karlmutch.com:$SECURE_INGRESS_PORT
</b></code></pre>

### Debugging

There are several pages of debugging instructions that can be used for situations when grpc failures to occur without much context, this applies to unexplained GRPC errors that reference C++ files within envoy etc.  These pages can be found by using the search function on the Istio web site at, https://istio.io/search.html?q=debugging.

You might find the following use cases useful for avoiding using hard coded pod names etc when debugging.

The following example shows enabling debugging for http2 and rbac layers within the Ingress Envoy instance.

<pre><code><b>
kubectl exec $(kubectl get pods -l istio=ingressgateway -n istio-system -o jsonpath='{.items[0].metadata.name}') -c istio-proxy -n istio-system -- curl -X POST "localhost:15000/logging?rbac=debug" -s
kubectl exec $(kubectl get pods -l istio=ingressgateway -n istio-system -o jsonpath='{.items[0].metadata.name}') -c istio-proxy -n istio-system -- curl -X POST "localhost:15000/logging?filter=debug" -s
kubectl exec $(kubectl get pods -l istio=ingressgateway -n istio-system -o jsonpath='{.items[0].metadata.name}') -c istio-proxy -n istio-system -- curl -X POST "localhost:15000/logging?http2=debug" -s
</b></code></pre>

After making a test request the log can be retrieved using something like the following:

<pre><code><b>
kubectl logs $(kubectl get pods --namespace istio-system -l istio=ingressgateway -o jsonpath='{.items[0].metadata.name}') --namespace istio-system</b></code></pre>

When debugging the istio proxy side cars for services you can do the following to enable all of the modules within the proxy:

<pre><code><b>
kubectl exec $(kubectl get pods -l app=experiment -o jsonpath='{.items[0].metadata.name}') -c istio-proxy -- curl -X POST "localhost:15000/logging?level=debug" -s</b></code></pre>

And then the logs can be captured during the testing using the following:

<pre><code><b>
kubectl logs $(kubectl get pods -l app=experiment -o jsonpath='{.items[0].metadata.name}') -c istio-proxy</b></code></pre>

A good way to diagnose resource consumption issues is to make use of ktop.

<pre><code><b>
go install github.com/vladimirvivien/ktop@latesta</b></code></pre>

# Logging and Observability

Currently the service mesh is deployed with Observability tools.  These instruction do not go into Observability at this time.  However we do address logging.

Individual services do offering logging using the systemd facilities and these logs are routed to Kubernetes.  Logs can be obtained from pods and containers. The 'kubectl get services' command can be used to identify the running platform services and the 'kubectl get pod' command can be used to get the health of services.  Once a pod isidentified with a running service instance the logs can be extract using a combination of the pod instance and the service name together, for example:

<pre><code><b>kuebctl get services</b>
NAME          TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)     AGE
experiments   ClusterIP   100.68.93.48   <none>        30001/TCP   12m
kubernetes    ClusterIP   100.64.0.1     <none>        443/TCP     1h
<b>kubectl get pod</b>
NAME                             READY     STATUS    RESTARTS   AGE
experiments-v1-bc46b5d68-tltg9   2/2       Running   0          12m
<b>kubectl logs experiments-v1-bc46b5d68-tltg9 experiments</b>
./experimentsrv built at 2018-01-18\_15:22:47+0000, against commit id 34e761994b895ac48cd832ac3048854a671256b0
2018-01-18T16:50:18+0000 INF experimentsrv git hash version 34e761994b895ac48cd832ac3048854a671256b0 _: [host experiments-v1-bc46b5d68-tltg9]
2018-01-18T16:50:18+0000 INF experimentsrv database startup / recovery has been performed dev-platform.cluster-cff2uhtd2jzh.us-west-2.rds.amazonaws.com:5432 name platform _: [host experiments-v1-bc46b5d68-tltg9]
2018-01-18T16:50:18+0000 INF experimentsrv database has 1 connections  dev-platform.cluster-cff2uhtd2jzh.us-west-2.rds.amazonaws.com:5432 name platform dbConnectionCount 1 _: [host experiments-v1-bc46b5d68-tltg9]
2018-01-18T16:50:33+0000 INF experimentsrv database has 1 connections  dev-platform.cluster-cff2uhtd2jzh.us-west-2.rds.amazonaws.com:5432 name platform dbConnectionCount 1 _: [host experiments-v1-bc46b5d68-tltg9]
</code></pre>

The container name can also include the istio mesh and kubernetes installed system containers for indepth debugging purposes.

### Kubernetes Web UI and console

In addition to the kops information for a cluster being hosted on S3, the kubectl information for accessing the cluster is stored within the ~/.kube directory.  The web UI can be deployed using the instruction at https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/#deploying-the-dashboard-ui, the following set of instructions include the deployment as it stood at k8s 1.9.  Take the opportunity to also review the document at the above location.

Kubectl service accounts can be created at will and given access to cluster resources.  To create, authorize and then authenticate a service account the following steps can be used:

```
kubectl create -f https://raw.githubusercontent.com/kubernetes/heapster/release-1.5/deploy/kube-config/influxdb/influxdb.yaml
kubectl create -f https://raw.githubusercontent.com/kubernetes/heapster/release-1.5/deploy/kube-config/influxdb/heapster.yaml
kubectl create -f https://raw.githubusercontent.com/kubernetes/heapster/release-1.5/deploy/kube-config/influxdb/grafana.yaml
kubectl create -f https://raw.githubusercontent.com/kubernetes/heapster/release-1.5/deploy/kube-config/rbac/heapster-rbac.yaml
kubectl create -f https://raw.githubusercontent.com/kubernetes/dashboard/master/src/deploy/recommended/kubernetes-dashboard.yaml
kubectl create serviceaccount studioadmin
secret_name=`kubectl get serviceaccounts studioadmin -o json | jq '.secrets[] | [.name] | join("")' -r`
secret_kube=`kubectl get secret $secret_name -o json | jq '.data.token' -r | base64 --decode`
# The following will open up all service accounts for admin, review the k8s documentation specific to your
# install version of k8s to narrow the roles
kubectl create clusterrolebinding serviceaccounts-cluster-admin --clusterrole=cluster-admin --group=system:serviceaccounts
```

The value in secret kube can be used to login to the k8s web UI.  First start 'kube proxy' in a terminal window to create a proxy server for the cluster.  Use a browser to navigate to http://localhost:8001/ui.  Then use the value in the secret\_kube variable as your 'Token' (Service Account Bearer Token).

You will now have access to the Web UI for your cluster with full privs.


# AAA using Auth0

Platform services are secured using the Auth0.com service.  Auth0 is a service that provides support for headless machine to machine authentication.  Auth0 is being used initially to provide Bearer tokens for both headless and CLI clients to platform services proof of concept.

Auth0 supports the ability to create a hosted database for storing user account and credential information.  You should navigate to the Connections -> Database section and create a database with the name of "Username-Password-Authentication".  This will be used later when creating applications as your source of user information.

Auth0 authorizations can be done using a Demo auth0.com account.  To do this you will need to add a custom API to the Auth0 account, call it something like "Experiments API" and give it an identifier of "http://api.karlmutch.com/experimentsrv", you should also enable RBAC and the "Add Permissions in the Access Token" options.  Then use the save button to persist the new API.  Identifiers are used as the 'audience' setting when generating tokens via web calls against the AAA features of the Auth0 platform.

The next stop is to use the menu bar to select the "Permissions" tab.  This tab allows you to create a scope to be used for the permissions granted to user.  Create a scope called "all:experiments" with a description, and select the Add button.  This scope will become available for use by authenticated user roles to allow them to access the API.

Next, click the "Machine To Machine Applications" tab.  This should show that a new Test Application has been created and authorized against this API.  To the right of the Authorized switch is a drop down button that can be used to expose more detailed information related to the scopes that are permitted via this API.  You should see that the all:experiments scope is not yet selected, select it and then use the update button.

Now navigate using the left side panel to the Applications screen.  Click to select your new "Experiments (Test Application)".  The screen displayed as a result will show the Client ID, and the Client secret that will be needed later on, take a note of thes values as they will be needed during AAA operation.  Go to the bottom of the page and you will be able to expose some advanced settings".  Inside the advanced settings you will see a ribbon bar with a "Grant Types" tab that can be clicked on revealing the selections available for grant type, ensure that the "password" radio button is selected to enable passwords for authentication, and click on the Save Changes button to save the selection.

The first API added by the Auth0 platform will be the client that accesses the Auth0 service itself providing per user authentication and token generation. When you begin creating a client you will be able to select the "Auth0 Management API" as one of the APIs you wish to secure.

The next step is to create a User and assign the user a role.  The left hand panel has a "Users & Roles" menu.  Using the menu you can select the "User" option and then use the "CREATE USER" button on the right side of the screen.  This where the real power of the Auth0 platform comes into play as you can use your real email address and perform operations related to identity management and passwords resets without needing to implement these features yourself.  When creating the user the connection field should be filled in with the Database connection that you created initially in these instructions. "Username-Password-Authentication".  After creating your User you can go to the Users panel and click on the email, then click on the permissions tab.  Add the all:experiments permission to the users prodile using the "ASSIGN PERMISSIONS" button.

You can now use various commands to manipulate the APIs outside of what will exist in the application code, this is a distinct advantage over directly using enterprise tools such as Okta.  Should you wish to use Okta as an Identity provider, or backend, to Auth0 then this can be done however you will need help from our Tech Ops department to do this and is an expensive option.  At this time the user and passwords being used for securing APIs can be managed through the Auth0 dashboard including the ability to invite users to become admins.

<pre><code><b>
export AUTH0_DOMAIN=karlmutch.auth0.com
export AUTH0_CLIENT_ID=pL3iSUmOB7EPiXae4gPfuEasccV7PATs
export AUTH0_CLIENT_SECRET=KHSCFuFumudWGKISCYD79ZkwF2YFCiQYurhjik0x6OKYyOb7TkfGKJrHKXXADzqG
export AUTH0_REQUEST=$(printf '{"client_id": "%s", "client_secret": "%s", "audience":"http://api.karlmutch.com/experimentsrv","grant_type":"password", "username": "karlmutch@gmail.com", "password": Ap9ss2345f"", "scope": "all:experiments", "realm": "Username-Password-Authentication" }' "$AUTH0_CLIENT_ID" "$AUTH0_CLIENT_SECRET")
export AUTH0_TOKEN=$(curl -s --request POST --url https://karlmutch.auth0.com/oauth/token --header 'content-type: application/json' --data "$AUTH0_REQUEST" | jq -r '"\(.access_token)"')

</b>
c.f. https://auth0.com/docs/quickstart/backend/golang/02-using#obtaining-an-access-token-for-testing.
</code></pre>

If you are using the test API and you are either running a kubectl port-forward or have a local instance of the postgres DB, you can do something like:

<pre><code><b>kubectl port-forward --namespace default svc/$PGRELEASE-postgresql 5432:5432 &
cd cmd/downstream
go run . --ip-port=":30008" &
cd ../..
cd cmd/experimentsrv
export AUTH0_DOMAIN=karlmutch.auth0.com
export AUTH0_CLIENT_ID=pL3iSUmOB7EPiXae4gPfuEasccV7PATs
export AUTH0_CLIENT_SECRET=KHSCFuFumudWGKISCYD79ZkwF2YFCiQYurhjik0x6OKYyOb7TkfGKJrHKXXADzqG
export AUTH0_REQUEST=$(printf '{"client_id": "%s", "client_secret": "%s", "audience":"http://api.karlmutch.com/experimentsrv","grant_type":"password", "username": "karlmutch@gmail.com", "password": "Passw0rd!", "scope": "all:experiments", "realm": "Username-Password-Authentication" }' "$AUTH0_CLIENT_ID" "$AUTH0_CLIENT_SECRET")
export AUTH0_TOKEN=$(curl -s --request POST --url https://karlmutch.auth0.com/oauth/token --header 'content-type: application/json' --data "$AUTH0_REQUEST" | jq -r '"\(.access_token)"')
go test -v --dbaddr=localhost:5432 -ip-port="[::]:30007" -dbname=platform -downstream="[::]:30008"
</b></code></pre>

## Auth0 claims extensibility

Auth0 can be configured to include additional headers with user metadata such as email addresses etc using custom rules in the Auth0 rules configuration.  Header that are added can be queried and extracted from gRPC HTTP authorization header meta data as shown in the experimentsrv server.go file. An example of a rule is as follows:

<pre><code>
function (user, context, callback) {
  context.accessToken["http://karlmutch.com/user"] = user.email;
  callback(null, user, context);
 }</code></pre>

 An example of extracting this item on the gRPC client side can be found in cmd/experimentsrv/server.go in the GetUserFromClaims function.

# Manually invoking and using production services with TLS

When using the gRPC services within a secured cluster these instructions can be used to access and exercise the services.

An example client for running a simple ping style test against the cluster is provided in the cmd/cli-experiment directory.  This client acts as a test for the presence of the service.  If the commands to obtain a JWT have been followed then this command can be run against the cluster as follows:

<pre><code><b>
cd cmd/cli-experiment
go run . --server-addr=platform-services.karlmutch.com:443 --auth0-token="$AUTH0_TOKEN"</b>
(*dev_karlmutch_experiment.CheckResponse)(0xc00003c280)(modules:"downstream" )
</code></pre>

Once valid transactions are being performed you should go back to the section on Honeycomb and configure the relevant fields inside the definitions panel for your Dataset.

# Manually invoking and using services without TLS

When using the gRPC services within an unsecured cluster these instructions can be used to access and exercise the services.

A pre-requiste of manually invoking GRPC servcies is that the grpc_cli tooling is installed.  The instructions for doing this can be found within the grpc source code repository at, https://github.com/grpc/grpc/blob/master/doc/command_line_tool.md.

The following instructions identify a $INGRESS_HOST value for cases where a LoadBalancer is being used.  If you are using Kubernetes in Docker and the cluster is hosted locally then the INGRESS_HOST value should be 127.0.0.1 for the following instructions.

Services used within the platform require that not only is the link integrity and security is maintained using mTLS but that an authorization block is also supplied to verify the user requesting a service.  The authorization can be supplied when using the gRPC command line tool using the metadata options.  First we retrieve a token using curl and then make a request against the service, run in this case as a local docker container, as follows:

<pre><code><b>grpc_cli call localhost:30001 dev.karlmutch.experiment.Service.Get "id: 'test'" --metadata authorization:"Bearer $AUTH0_TOKEN"
</b></code></pre>

The services used within the platform all support reflection when using gRPC.  To examine calls available for a server you should first identify the endpoint through which the gateway is being routed, in this case as part of an Istio cluster on AWS with TLS enabled, for example:

<pre><code><b>export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
export INGRESS_HOST=$(kubectl get po -l istio=ingressgateway -n istio-system -o jsonpath='{.items[0].status.hostIP}')
export CLUSTER_INGRESS=$INGRESS_HOST:$SECURE_INGRESS_PORT
<b>grpc_cli ls $CLUSTER_INGRESS -l --channel_creds_type=ssl</b>
filename: experimentsrv.proto
package: dev.karlmutch.experiment;
service Service {
  rpc Create(dev.karlmutch.experiment.CreateRequest) returns (dev.karlmutch.experiment.CreateResponse) {}
  rpc Get(dev.karlmutch.experiment.GetRequest) returns (dev.karlmutch.experiment.GetResponse) {}
  rpc MeshCheck(dev.karlmutch.experiment.CheckRequest) returns (dev.karlmutch.experiment.CheckResponse) {}
}

filename: grpc/health/v1/health.proto
package: grpc.health.v1;
service Health {
  rpc Check(grpc.health.v1.HealthCheckRequest) returns (grpc.health.v1.HealthCheckResponse) {}
  rpc Watch(grpc.health.v1.HealthCheckRequest) returns (stream grpc.health.v1.HealthCheckResponse) {}
}

filename: grpc_reflection_v1alpha/reflection.proto
package: grpc.reflection.v1alpha;
service ServerReflection {
  rpc ServerReflectionInfo(stream grpc.reflection.v1alpha.ServerReflectionRequest) returns (stream grpc.reflection.v1alpha.ServerReflectionResponse) {}
}
</code></pre>


In circumstances where you have the CN name validation enabled then the host name should reflect the common name for the host, for example:

<pre><code><b>grpc_cli call platform-services.karlmutch.com:$SECURE_INGRESS_PORT dev.karlmutch.experiment.Service.Get "id: 'test'" --metadata authorization:"Bearer $AUTH0_TOKEN" --channel_creds_type=ssl
</b></code></pre>


To drill further into interfaces and examine the types being used within calls you can perform commands such as:

<pre><code><b>grpc_cli type $CLUSTER_INGRESS dev.karlmutch.experiment.CreateRequest -l --channel_creds_type=ssl</b>
message CreateRequest {
.dev.karlmutch.experiment.Experiment experiment = 1[json_name = "experiment"];
}
<b>grpc_cli type $CLUSTER_INGRESS dev.karlmutch.experiment.Experiment -l --channel_creds_type=ssl</b>
message Experiment {
string uid = 1[json_name = "uid"];
string name = 2[json_name = "name"];
string description = 3[json_name = "description"];
.google.protobuf.Timestamp created = 4[json_name = "created"];
map&lt;uint32, .dev.karlmutch.experiment.InputLayer&gt; inputLayers = 5[json_name = "inputLayers"];
map&lt;uint32, .dev.karlmutch.experiment.OutputLayer&gt; outputLayers = 6[json_name = "outputLayers"];
}
<b>grpc_cli type $CLUSTER_INGRESS dev.karlmutch.experiment.InputLayer -l --channel_creds_type=ssl</b>
message InputLayer {
enum Type {
	Unknown = 0;
	Enumeration = 1;
	Time = 2;
	Raw = 3;
}
string name = 1[json_name = "name"];
.dev.karlmutch.experiment.InputLayer.Type type = 2[json_name = "type"];
repeated string values = 3[json_name = "values"];
}
</code></pre>


# Shutting down a service, or cluster

<pre><code><b>kubectl delete -f experimentsrv.yaml
</b></code></pre>

<pre><code><b>kops delete cluster $CLUSTER_NAME --yes
</b></code></pre>
