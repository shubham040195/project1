package main

import (
    "fmt"
    "os"
    "log"
    "helm.sh/helm/v3/pkg/action"
    "helm.sh/helm/v3/pkg/chart/loader"
    "helm.sh/helm/v3/pkg/kube"
    "helm.sh/helm/v3/pkg/cli"

     _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)


var (
        chartPath string = "/home/shubham_khatri1995/work/helm/charts/mysql"
        kubeconfigPath string = "/home/shubham_khatri1995/.kube/config"
        chartName string = "mysql1"
        namespaceName string = "rakuten"
        kubeContext string = "minikube"
//      kubeContext string = "gke_thermal-talon-297909_us-central1-c_gocluster-1"
)


func helmK8sConfiguration() *action.Configuration {
    actionConfig := new(action.Configuration)
    if err := actionConfig.Init(kube.GetConfig(kubeconfigPath, kubeContext, ""), namespaceName, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
        fmt.Sprintf(format, v)
    }); err != nil {
        panic(err)
    }
    return actionConfig
}

func install_chart() {
    chart, err := loader.Load(chartPath)
    if err != nil {
        panic(err)
    }

    actionConfig := helmK8sConfiguration()
    client := action.NewInstall(actionConfig)
    client.CreateNamespace = true
    client.GenerateName = true
    client.Namespace = namespaceName
    client.ReleaseName = chartName
    rel, err := client.Run(chart, nil)
    if err != nil {
        panic(err)
    }
    fmt.Println("Successfully installed helm chart:",rel.Name)
}

func uninstall_chart(){
        actionConfig := helmK8sConfiguration()
        client := action.NewUninstall(actionConfig)
        chart_list := list_chart()
        for _, chart := range chart_list{
                _, err := client.Run(chart)
                if err != nil {
                        log.Printf("%+v",err)
                        os.Exit(1)
                }
        }
        fmt.Println("Successfully uninstalled helm charts")

}

func list_chart() []string{

    actionConfig := helmK8sConfiguration()
    client := action.NewList(actionConfig)
    // Only list deployed
    client.Deployed = true
    results, err := client.Run()
    if err != nil {
        log.Printf("%+v", err)
        os.Exit(1)
    }
    chartlist:=[]string{}
    // chartlist := make([]string,1)
    for _, rel := range results {
            chartlist= append(chartlist,rel.Name)
    }
    return chartlist
}


func pull_chart() {
        chartRef := "nginx"
        client := action.NewPull()
        client.Settings = &cli.EnvSettings{}
        client.ChartPathOptions.RepoURL = "https://charts.bitnami.com/bitnami"
        client.DestDir = chartPath
        client.UntarDir = chartPath
        client.Untar = true
        res, err := client.Run(chartRef)
        if err != nil {
                log.Printf("%+v",err)
                os.Exit(1)
        }
        fmt.Println("Successfully pulled chart:",res)
}

func chartCreate(){
        res, err := chartutil.Create(chartName, chartPath)
        if err != nil{
                panic(err)
        }
        fmt.Println("Successfully created helm chart:",res)
}



func main() {
     install_chart()
     uninstall_chart()
     pull_chart()
     results :=list_chart()
     fmt.Println("list of Installed helm charts",results)
}
