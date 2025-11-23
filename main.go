package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/pkhamre/kube-waste/pkg/analyzer"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
)

func main() {
	// 0. Global Silence: Turn off klog completely
	klog.SetOutput(io.Discard)
	klog.LogToStderr(false)
	// This prevents klog from hijacking flags
	flag.Set("logtostderr", "false")
	flag.Parse()
	// 1. Resolve Kubeconfig Path
	// Priority: ./kubeconfig -> ./.kubeconfig -> ~/.kube/config -> KUBECONFIG env
	cwd, _ := os.Getwd()
	localConfig := filepath.Join(cwd, "kubeconfig")
	localDotConfig := filepath.Join(cwd, ".kubeconfig")

	var kubeconfig string

	if _, err := os.Stat(localConfig); err == nil {
		kubeconfig = localConfig
		fmt.Printf("Using config from current directory: %s\n", kubeconfig)
	} else if _, err := os.Stat(localDotConfig); err == nil {
		kubeconfig = localDotConfig
		fmt.Printf("Using config from current directory: %s\n", kubeconfig)
	} else if home := homedir.HomeDir(); home != "" {
		// Check home dir, but don't crash if missing, clientcmd handles it
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		kubeconfig = os.Getenv("KUBECONFIG")
	}

	// 2. Build Config
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %v\n", err)
		os.Exit(1)
	}

	config.WarningHandler = rest.NoWarnings{}

	// 3. Create Client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Scanning cluster for waste (Go Version)...")

	// 4. Run Analyzers
	var allWaste []analyzer.WasteItem

	// Check PVCs
	pvcWaste, err := analyzer.DetectPVCWaste(clientset)
	if err != nil {
		fmt.Printf("Error scanning PVCs: %v\n", err)
	}
	allWaste = append(allWaste, pvcWaste...)

	// Check Deployments
	deployWaste, err := analyzer.DetectZombieDeployments(clientset)
	if err != nil {
		fmt.Printf("Error scanning Deployments: %v\n", err)
	} else {
		allWaste = append(allWaste, deployWaste...)
	}

	// Check Services
	svcWaste, err := analyzer.DetectUnusedServices(clientset)
	if err != nil {
		fmt.Printf("Error scanning Services: %v\n", err)
	} else {
		allWaste = append(allWaste, svcWaste...)
	}

	// Check Orphaned Pods
	podWaste, err := analyzer.DetectOrphanedPods(clientset)
	if err != nil {
		fmt.Printf("Error scanning Pods: %v\n", err)
	} else {
		allWaste = append(allWaste, podWaste...)
	}

	// 5. Print Table
	printTable(allWaste)
}

func printTable(items []analyzer.WasteItem) {
	totalSavings := 0.0

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, "TYPE\tNAMESPACE\tNAME\tDETAILS\tEST. SAVINGS")
	fmt.Fprintln(w, "----\t---------\t----\t-------\t------------")

	for _, item := range items {
		costStr := fmt.Sprintf("$%.2f/mo", item.Cost)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			item.Type,
			item.Namespace,
			item.Name,
			item.Details,
			costStr,
		)
		totalSavings += item.Cost
	}

	w.Flush()

	fmt.Println("-------------------------------------------------------------------------")
	fmt.Printf("TOTAL POTENTIAL SAVINGS: $%.2f / month\n", totalSavings)
}
