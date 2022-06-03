terraform {
  required_providers {
    rancher2 = {
      source = "rancher/rancher2"
    }
    kubectl = {
      source  = "gavinbunney/kubectl"
    }
    maas = {
      source  = "ionutbalutoiu/maas"
      version = "1.0.1"
    }
    harvester = {
      source  = "harvester/harvester"
    }
  }
  required_version = ">= 0.13"
}
