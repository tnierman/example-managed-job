# Example Managed Jobs

## Introduction
The following tutorial creates a basic client in Golang using Kubernetes' Go API. It's recommended you complete this at least once, then use this as a base branch for future steps.

Since this and future sections place a heavier emphasis on writing Go, the following steps involve editing files directly instead of using heredocs to generate them from the command line. The final product created in this section is included, for reference.

## Instructions

### Create the client

Open `main.go` with your favorite text editor. The file should look similar to this

```golang
package main

import "fmt"

func main() {
   fmt.Println("hello world")
}
```

To break this down:
- `package main` - Declares a new go package named `main` for the [previously initialized module](https://github.com/tnierman/example-managed-job/tree/main#go). Basically, [modules](https://go.dev/doc/modules/developing) are collections of packages, and [packages](https://go.dev/doc/code) are collections of files. Every go file we add needs a package declaration as the first statement at the top; only comments can preceed the package declaration.
- `import "fmt"` - Imports the `fmt` package, which is a built-in package included with every go release. Imports let you use code defined elsewhere in your code.
- `func main() {...}` - Represents the `main` function. The `main` function is where every Go program starts executing from.
- `fmt.Println("hello world")` - Uses the standard `fmt` package to print "hello world" to your console.

Since this tutorial is geared specifically toward creating managed-jobs, we won't spend much time covering the basics of Go. If at any time the syntax or behavior of the Go code below feels unfamiliar, _stop_! In order to build effective Go programs, you need, at a minimum, a basic understanding of the language and the tools it provides. There are numerous tutorials online that build on simple programs like the one above to gradually cover everything from the most basic syntax to the most advanced concepts.

If you need somewhere to start, a (very) short list:
- [A Tour of Go](https://go.dev/tour/welcome/1) - Walks you step-by-step through the Go documentation
- [W3 school's Go tutorial](https://www.w3schools.com/go/) - Similar to "A Tour of Go": walks you step-by-step through Go concepts
- [Effective Go](https://go.dev/doc/effective_go) - A book that covers best practices and effective design

One tool worth calling out in particular is https://pkg.go.dev: this site acts as a reference for many of the Go packages we'll be importing, and can be invaluable when working with unfamiliar methods or types.

---

As a first step, lets modify the existing `Println` in a few ways:

```golang
package main

import "log"

func main() {
    log.Println("Building client")
}
```

The changes here include:
- Updating `import "fmt"` to `import "log"`
- Changing `fmt.Println("hello world")` to `log.Println("Building client")`

This achieves a few things: 1) importing `log` instead of `fmt` lets us call the functions present in the `log` package, and 2) using `log`'s version of `Println` automatically prepends the date and time to the message you've passed to it. This will be useful to ensure the code we write remains performant. Finally, `"Building client"` is a more meaningful log when debugging than `"hello world"`

> **Note**
> If you want to know more about the methods supplied by the log or fmt libraries, refer to their respective go.pkg.dev pages [here](https://pkg.go.dev/log) and [here](https://pkg.go.dev/fmt)

> **Note**
> It's wise to periodically run your code as you follow along. The easiest way to do this is invoking `go run main.go` from the root directory (same directory as your `go.mod`) of your project


Next, lets create our client:

```golang
package main

import (
        "fmt"
        "log"

        ctrl "sigs.k8s.io/controller-runtime"
        kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
        log.Println("Creating the client")

        client, err := kclient.New(ctrl.GetConfigOrDie(), kclient.Options{})
        if err != nil {
                log.Fatalf("Couldn't create client: %v", err)
        }
        fmt.Printf("client: %#v\n", client)
}
```

A few things had to happen to instantiate our client:
- The `fmt` package is re-imported - we'll talk about why in a bit
- The first kubernetes packages have been imported as well! You'll notice we named these using [import declarations](https://go.dev/ref/spec#Import_declarations) - this lets us call the packages using a different name than normal. We did this for `controller-runtime`, because typing out that name every time you want something from it is a good way to get carpal tunnel syndrome. `client`, on the other hand, is a name we want to use for _our_ client object, not the package, so we rename it `kclient` here (short for "kubernetes client").
- Using the `kclient` (really the [controller-runtime/pkg/client](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/client) package), we call `New()` to create our client, and provide it with a few arguments:
  - The `ctrl.GetConfigOrDie()` is a `controller-runtime` function that retrieves a kubeconfig from the cluster. The exact permissions granted depend on who's running this code (we'll get to that later), but for now this _should_ be sufficient to perform some basic actions.
  - The `kclient.Options{}` specifies that we don't need any special options associated with this client. If you're curious, the options we could request are [here](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/client#Options).
- We set two new variables, `client` and `err`, as the output from this function:
  - `client` is now an object we can use to interact with the cluster (assuming the `New()` function succeeded). Using `client`, we can `Create()`, `List()`, `Delete()`, etc objects from the cluster, like you would with `oc`.
  - `err` determines whether the `New()` function actually succeeded or not. By checking if `err != nil`, we're checking if an error was returned, and, if so, we're logging it and exiting.
- Finally, in order for this Go code to be valid, every variable we declare has to be used in some way. Usually, the easiest way to ensure this happens during development is to print its value - thus `Printf()`. We're using `Printf` instead of `Println()` so that we can use [formatting verbs](https://pkg.go.dev/fmt#hdr-Printing) to see what's currently defined within the `client` object. (This is also why `fmt` had to be re-imported)

> **Note**
> Before running the above code, make sure you're logged into a cluster you own. You may also have to run `go mod tidy` to download the dependencies used by this code.

Finally, lets use the client to retrieve some data from the cluster:

```golang
package main

import (
        "context"
        "fmt"
        "log"

        corev1 "k8s.io/api/core/v1"
        ctrl "sigs.k8s.io/controller-runtime"
        kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
        log.Println("Creating the client")

        client, err := kclient.New(ctrl.GetConfigOrDie(), kclient.Options{})
        if err != nil {
                log.Fatalf("Couldn't create client: %v", err)
        }

        log.Println("Looking for kube-apiserver pods...")
        pods := corev1.PodList{}
        err = client.List(context.TODO(), &pods, &kclient.ListOptions{Namespace: "openshift-kube-apiserver"})
        if err != nil {
                log.Fatalf("Couldn't list pods in 'openshift-kube-apiserver': %v", err)
        }

        log.Println("Pods in 'openshift-kube-apiserver': ")
        for _, pod := range pods.Items {
                fmt.Println("- ", pod.Name)
        }
}
```

A lot more has changed from before, so to break things down one last time:
- Two new imports: `context` and `corev1`:
  - `context` is yet another standard package (notice the pattern: these built-in libraries are called by name, no url needed) which "carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes" ([reference](https://pkg.go.dev/context))
  - `corev1` is a Kubernetes API most people are familiar with - this is where pods, volumes, services, etc are [defined](https://pkg.go.dev/k8s.io/api/core/v1)
- A log message was added, indicating that the `client` was created successfully.
- To keep our output concise, the `fmt.Printf()` line was removed. This is allowed, because `client` is now used later to `List()` some stuff.
- Now, the good stuff: first we're initializing a new variable, `pods`, which is a [PodList](https://pkg.go.dev/k8s.io/api/core/v1#PodList). A `PodList` isn't something you've probably seen if you haven't worked with Kubernetes APIs before - it's how our `client` is able work with multiple `Pod` objects at once.
- Next, we're using our `client` to `List()` pods: in order to do that we need to pass it a `context` (we just use `TODO()` here), the `PodList` we just created, and some options. This time, we do specify some options: specifically that we want to see pods in the `openshift-kube-apiserver` namespace. All the different options we could've also specified are [here](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/client#ListOptions).
- Then, we're checking our error again, to make sure the `List()` succeeded. If it doesn't, we log why and exit.
- Finally, we print out the pods we found in `openshift-kube-apiserver`

Hopefully, when you run this file, the output looks similar to:

```bash
$ go run main.go
2023/04/19 20:28:08 Creating the client
I0419 20:28:10.537290  192060 request.go:690] Waited for 1.046967266s due to client-side throttling, not priority and fairness, request: GET:https://api-backplane.apps.hive-stage-01.n1u3.p1.openshiftapps.com/backplane/cluster/2371i7699p060unprp448283ibe1m6mi/apis/events.k8s.io/v1beta1?timeout=32s
2023/04/19 20:28:12 Looking for kube-apiserver pods...
2023/04/19 20:28:12 Pods in 'openshift-kube-apiserver':
-  installer-11-ip-10-0-117-108.ec2.internal
-  installer-11-ip-10-0-157-51.ec2.internal
-  installer-12-ip-10-0-117-108.ec2.internal
-  installer-12-ip-10-0-157-51.ec2.internal
-  installer-12-ip-10-0-166-180.ec2.internal
-  kube-apiserver-guard-ip-10-0-117-108.ec2.internal
-  kube-apiserver-guard-ip-10-0-157-51.ec2.internal
-  kube-apiserver-guard-ip-10-0-166-180.ec2.internal
-  kube-apiserver-ip-10-0-117-108.ec2.internal
-  kube-apiserver-ip-10-0-157-51.ec2.internal
-  kube-apiserver-ip-10-0-166-180.ec2.internal
-  revision-pruner-10-ip-10-0-117-108.ec2.internal
-  revision-pruner-11-ip-10-0-117-108.ec2.internal
-  revision-pruner-11-ip-10-0-157-51.ec2.internal
-  revision-pruner-12-ip-10-0-117-108.ec2.internal
-  revision-pruner-12-ip-10-0-157-51.ec2.internal
-  revision-pruner-12-ip-10-0-166-180.ec2.internal
-  revision-pruner-9-ip-10-0-117-108.ec2.internal
```

Now, if you run `oc get pod -n openshift-kube-apiserver`, the output should match:

```bash
$ oc get po -n openshift-kube-apiserver
NAME                                                READY   STATUS      RESTARTS   AGE
installer-11-ip-10-0-117-108.ec2.internal           0/1     Completed   0          3h51m
installer-11-ip-10-0-157-51.ec2.internal            0/1     Completed   0          3h47m
installer-12-ip-10-0-117-108.ec2.internal           0/1     Completed   0          3h36m
installer-12-ip-10-0-157-51.ec2.internal            0/1     Completed   0          3h44m
installer-12-ip-10-0-166-180.ec2.internal           0/1     Completed   0          3h40m
kube-apiserver-guard-ip-10-0-117-108.ec2.internal   1/1     Running     0          3h51m
kube-apiserver-guard-ip-10-0-157-51.ec2.internal    1/1     Running     0          3h47m
kube-apiserver-guard-ip-10-0-166-180.ec2.internal   1/1     Running     0          3h42m
kube-apiserver-ip-10-0-117-108.ec2.internal         5/5     Running     0          3h33m
kube-apiserver-ip-10-0-157-51.ec2.internal          5/5     Running     0          3h41m
kube-apiserver-ip-10-0-166-180.ec2.internal         5/5     Running     0          3h37m
revision-pruner-10-ip-10-0-117-108.ec2.internal     0/1     Completed   0          3h54m
revision-pruner-11-ip-10-0-117-108.ec2.internal     0/1     Completed   0          3h52m
revision-pruner-11-ip-10-0-157-51.ec2.internal      0/1     Completed   0          3h49m
revision-pruner-12-ip-10-0-117-108.ec2.internal     0/1     Completed   0          3h44m
revision-pruner-12-ip-10-0-157-51.ec2.internal      0/1     Completed   0          3h44m
revision-pruner-12-ip-10-0-166-180.ec2.internal     0/1     Completed   0          3h44m
revision-pruner-9-ip-10-0-117-108.ec2.internal      0/1     Completed   0          3h54m
```

If so, this indicates that the Go worked, and that you've successfully created and used the Kubernetes client!

### Run as a job
Since we're making managed-jobs, we need to make sure this code runs properly _in_ the cluster. Using the resources created in the previous step, this shouldn't be too hard:


First, lets rebuild our docker image. To do this, we need to add `COPY go.sum go.sum` to the previously generated Dockerfile. Refer to the current [Dockerfile](./Dockerfile#L7) for clarity. This wasn't included in the previous step, because go.sum files are only created after an [external dependency is used](https://go.dev/blog/using-go-modules#adding-a-dependency). After doing that, run the following:

> **Note**
> Be sure to replace the quay.io user in the below commands to reference yourself

```bash
podman build . -t quay.io/tnierman_openshift/example-managed-job:latest
```

This tags the container correctly the first time, so we no longer need to run `podman tag` before pushing.

Now push the newly built image:

```bash
podman push quay.io/tnierman_openshift/example-managed-job:latest
```

Next, lets recreate the job from before, and see if the code runs:

```bash
oc create -f deploy/job.yaml
oc get pod -n openshift-backplane-cee
```

Hopefully, the output looks something like below:
```bash
$ oc logs example-managed-job-pm6t5 -n openshift-backplane-cee
2023/04/20 03:00:11 Creating the client
I0420 03:00:12.106670       1 request.go:690] Waited for 1.035043861s due to client-side throttling, not priority and fairness, request: GET:https://172.30.0.1:443/apis/console.openshift.io/v1?timeout=32s
2023/04/20 03:00:14 Looking for kube-apiserver pods...
2023/04/20 03:00:14 Pods in 'openshift-kube-apiserver':
-  installer-11-ip-10-0-117-108.ec2.internal
-  installer-11-ip-10-0-157-51.ec2.internal
-  installer-12-ip-10-0-117-108.ec2.internal
-  installer-12-ip-10-0-157-51.ec2.internal
-  installer-12-ip-10-0-166-180.ec2.internal
-  kube-apiserver-guard-ip-10-0-117-108.ec2.internal
-  kube-apiserver-guard-ip-10-0-157-51.ec2.internal
-  kube-apiserver-guard-ip-10-0-166-180.ec2.internal
-  kube-apiserver-ip-10-0-117-108.ec2.internal
-  kube-apiserver-ip-10-0-157-51.ec2.internal
-  kube-apiserver-ip-10-0-166-180.ec2.internal
-  revision-pruner-10-ip-10-0-117-108.ec2.internal
-  revision-pruner-11-ip-10-0-117-108.ec2.internal
-  revision-pruner-11-ip-10-0-157-51.ec2.internal
-  revision-pruner-12-ip-10-0-117-108.ec2.internal
-  revision-pruner-12-ip-10-0-157-51.ec2.internal
-  revision-pruner-12-ip-10-0-166-180.ec2.internal
-  revision-pruner-9-ip-10-0-117-108.ec2.internal
```

If it doesn't, make sure the job was deployed to the `openshift-backplane-cee` (or similar) namespace. Deploying to the `default` namespace will attach the `default` serviceaccount to the pod, which doesn't have sufficient permissions to view pods in the `openshift-kube-apiserver` namespace.

## Next steps
If you were able to retrieve similar output from your test job, then use the following link to reach the next step.

- [ ] [Roll a deployment](https://github.com/tnierman/example-managed-job/blob/roll_deployment/README.md)

**If you were not**, you likely will not be able to complete the next step, since it iterates on the existing codebase. Keep trying! Sometimes the best way to learn something is to bang your head against your desk :)
