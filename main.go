package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/authn/k8schain"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	ctx := context.Background()

	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, nil).ClientConfig()
	if err != nil {
		log.Fatalf("Failed to load kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create kubernetes client: %v", err)
	}

	pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list pods: %v", err)
	}

	// Collect unique images across all pods.
	images := make(map[string]struct{})
	for _, pod := range pods.Items {
		for _, c := range pod.Spec.InitContainers {
			images[c.Image] = struct{}{}
		}
		for _, c := range pod.Spec.Containers {
			images[c.Image] = struct{}{}
		}
	}

	kc, err := k8schain.NewInCluster(ctx, k8schain.Options{})
	if err != nil {
		log.Printf("Warning: k8schain in-cluster auth unavailable (%v), falling back to default keychain", err)
		kc = authn.DefaultKeychain
	}

	exitCode := 0
	for img := range images {
		ref, err := name.ParseReference(img)
		if err != nil {
			log.Printf("ERROR: failed to parse image ref %q: %v", img, err)
			exitCode = 1
			continue
		}

		// Pulling the manifest is enough to count as a "pull" for ECR lifecycle policies.
		desc, err := remote.Get(ref, remote.WithAuthFromKeychain(kc), remote.WithContext(ctx))
		if err != nil {
			log.Printf("ERROR: failed to fetch manifest for %s: %v", ref, err)
			exitCode = 1
			continue
		}

		fmt.Printf("OK: %s (digest: %s)\n", ref, desc.Digest)
	}

	if exitCode != 0 {
		os.Exit(exitCode)
	}
}
