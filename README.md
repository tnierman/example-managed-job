# Example Managed Jobs

## Introduction
This repository is intended to serve as a basic tutorial to creating Golang applications, specifically managed-jobs, that interact with Kubernetes clusters.

In order to follow along:
- Create a new repo in GitHub
- Clone it to your local workstation
- Follow the instructions below. Try to complete as many as you can without referencing the files provided here; but if you do get stuck, they can serve as a reference. The link to the next step is given [at the bottom of this page](#next-steps), if you'd like to skip ahead.

## Getting started
These initial steps establish a containerized go project. Except where mentioned, most of these instructions should be copy/paste-able

### Go
The following steps create a basic Golang project to build off of:

Initialize your go module:
```bash
go mod init github.com/<user>/example-managed-job
```

Create a main.go:
```bash
cat << EOF > main.go
package main

import "fmt"

func main() {
   fmt.Println("hello world")
}
EOF
```

Test your program:
```bash
go run main.go
```

### Docker
The following steps containerize the project created above, and store it in a remote registry:

Create a `Dockerfile`:
```bash
cat << EOF > Dockerfile
FROM golang:1.20 as builder

WORKDIR /workspace

# Copy Go Module manifest & dependency files
COPY go.mod go.mod

# Install deps
RUN go mod download

# Copy source files
COPY main.go

# Build the thing
RUN CGO_ENABLED=0 go build -o job main.go

# Use distroless as minimal base image to package the binary
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/job .
USER 65532:65532

ENTRYPOINT ["/job"]
EOF
```

Build the container image:
```bash
podman build . -t example-managed-job:latest
```
(Substitute `podman` for `docker` if your container engine differs)

Test the container:
```bash
podman run example-managed-job:latest
```

Push the container (create a quay.io user and `example-managed-job` repo first!):
```bash
# Modify the quay url in the following commands to point to your own user and repo
# The 'tag' command isn't needed if the original 'build' is invoked with the correct '-t' argument
podman tag example-managed-job:latest quay.io/tnierman_openshift/example-managed-job:latest
podman push quay.io/tnierman_openshift/example-managed-job:latest
```

### Kubernetes/OpenShift
The following steps deploy the containerized project to a Kubernetes or OpenShift cluster:

Create a job manifest:
```bash
# Most SRE-P repositories define a dedicated directory, like 'deploy', for k8s files
mkdir deploy/
cd deploy/
cat << EOF > job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: example-managed-job
spec:
  completions: 1
  template:
    metadata:
      labels:
        app: example-managed-job
    spec:
      containers:
      - command:
        - /job
        image: quay.io/tnierman_openshift/example-managed-job:latest
        imagePullPolicy: Always
        name: job
      restartPolicy: OnFailure
EOF
```
Be sure to modify the `image` url to match your own registry

Create and login to a cluster. Deploy the job with
```bash
# Clean up any previous jobs prior to deploying a new one
oc delete -f job.yaml --ignore-not-found
oc create -f job.yaml
```
NOTE: If your job encounters an `ErrImagePull`, make sure the registry you pushed to prior is public!

Verify the job ran successfully. For example
```bash
$ oc get po -n openshift-backplane-cee
NAME                                                  READY   STATUS      RESTARTS   AGE
example-managed-job-pk4kj                             0/1     Completed   0          1m

# Copy the pod name - it will be unique every time a new job is created
$ oc logs example-managed-job-pk4kj -n openshift-backplane-cee
hello world
```

## Next Steps
After completing the basics here, it's recommended that you commit your changes and continue to the next step.

- [ ] Create a [basic Kubernets client](https://github.com/tnierman/example-managed-job/tree/basic_client)
