### Notes and Findings

* Usage of `NewYAMLOrJSONDecoder` from package `"k8s.io/apimachinery/pkg/util/yaml"` is a helpful tool to parse resources directly from yaml files instead of writing a structured code for a resource.

* `client-go` requires you to create a client separately for each of your k8s resources to talk to the API server, whereas `controller-runtime` library provides you a single dynamic client that can be used with all resources.

* Robust config fetcher `GetConfigOrDie()` in `controller-runtime` library that tries to get config with multi level precedence from different sources ( --kubeconfig flag -> KUBECONFIG env varibale -> in-cluster config if inside cluster -> $HOME/.kube/config) before failing. While for `client-go`, you'd have to fetch config separately.

