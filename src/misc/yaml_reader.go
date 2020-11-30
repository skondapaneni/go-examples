package main

import (
    "fmt"
    "io/ioutil"
    //"path/filepath"
    "gopkg.in/yaml.v2"
    "log"
)

type Service struct {
    APIVersion string `yaml:"apiVersion"`
    Kind       string `yaml:"kind"`
    Metadata   struct {
        Name      string `yaml:"name"`
        Namespace string `yaml:"namespace"`
    } `yaml:"metadata"`
    Spec struct {
        Type     string `yaml:"type"`
        Containers []struct {
            Name       string `yaml:"name"`
            IpAddr     string `yaml:"ip"`
            BridgeIp   string `yaml:"bridge_ip"`
            GwIp       string `yaml:"gw_ip"`
        } `yaml:"containers"`
    } `yaml:"spec"`
}

func parse_config(filename string) {
    var service Service
    data, _ := ioutil.ReadFile(filename)

    err := yaml.Unmarshal(data, &service)
    if err != nil {
       log.Fatal("%v", err)
    }

    fmt.Println("serviceName : ", service.Metadata.Name)
    fmt.Println("service: ", service)
}


func main() {
    parse_config("./service.yaml")
}
