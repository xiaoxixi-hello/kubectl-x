/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
    "context"
    "fmt"
    "log"

    "github.com/spf13/cobra"
    "github.com/ylinyang/kubectl-x/pkg/kube"
    "github.com/ylinyang/kubectl-x/pkg/mtable"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// xCmd represents the x command
var xCmd = &cobra.Command{
    Use:   "x",
    Short: "x",
    Long:  "x",
}

// hxCmd is Subcommand for xCmd
var hxCmd = &cobra.Command{
    Use:   "hx",
    Short: "Show the relationship between deployment server and ingress",
    Long:  "Show the relationship between deployment server and ingress",
    Run:   hx,
}

func init() {
    rootCmd.AddCommand(xCmd)
    xCmd.AddCommand(hxCmd)
}

func hx(cmd *cobra.Command, args []string) {
    clientset := kube.ClientSet(KubernetesConfigFlags)
    ns, _ := rootCmd.Flags().GetString("namespace")

    m := make([]map[string]string, 0)
    // if flag, _ := cmd.Flags().GetBool("hx"); flag {

    // display deploymentName status  ip:port svcIp  ingressUrl
    podList, err := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        log.Fatalln(err)
    }

    sMap := make(map[string]string)
    serviceList, err := clientset.CoreV1().Services(ns).List(context.TODO(), metav1.ListOptions{})
    if err == nil {
        for _, service := range serviceList.Items {
            if service.Spec.ClusterIP == "" {
                sMap[service.ObjectMeta.Name] = ""
                continue
            }
            sMap[service.ObjectMeta.Name] = service.Spec.ClusterIP
        }
    }

    iMap := make(map[string]string, 100)
    ingressList, err := clientset.NetworkingV1beta1().Ingresses(ns).List(context.TODO(), metav1.ListOptions{})
    if err == nil {
        for _, ingress := range ingressList.Items {
            if ingress.Spec.Rules[0].Host == "" {
                iMap[ingress.ObjectMeta.Name] = ""
                continue
            }
            iMap[ingress.ObjectMeta.Name] = ingress.Spec.Rules[0].Host
        }
    }

    for _, pod := range podList.Items {
        pMap := make(map[string]string)
        name := pod.Spec.Containers[0].Name
        pMap["NAME"] = pod.Name
        pMap["IP"] = pod.Status.PodIP
        pMap["SVC"] = sMap[name]
        pMap["INGRESS"] = iMap[name]
        if pod.Status.ContainerStatuses != nil {
            //  fmt.Println(pod.Status.ContainerStatuses)
            pMap["STATUS"] = func(p corev1.Pod) string {
                if p.Status.ContainerStatuses[0].State.Terminated != nil {
                    switch p.Status.ContainerStatuses[0].State.Terminated.ExitCode {
                    case 127:
                        return "Error"
                    default:
                        return "Terminated"
                    }
                } else if p.Status.ContainerStatuses[0].State.Running != nil {
                    return "Running"
                } else if p.Status.ContainerStatuses[0].State.Waiting.Reason != "" {
                    switch p.Status.ContainerStatuses[0].State.Waiting.Reason {
                    case "CrashLoopBackOff":
                        return "CrashLoopBackOff"
                    case "ContainerCreating":
                        return "ContainerCreating"
                    case "ImagePullBackOff":
                        return "ImagePullBackOff"
                    default:
                        return "Waiting"
                    }
                }
                return ""
            }(pod)
        } else {
            pMap["STATUS"] = string(pod.Status.Phase)
        }
        m = append(m, pMap)
    }
    // }
    // gen table
    table := mtable.GenTable(m)
    fmt.Println(table)
}
