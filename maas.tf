resource "maas_instance" "this" {
  for_each = var.maas.instances

  allocate_params {
    hostname = "${each.key}.${var.maas.domain}"
  }
  deploy_params {
    distro_series = "bionic"
    user_data = templatefile("cloud-inits/ubuntu18_harv.yml",
    {
      rancher_join_command = rancher2_cluster_v2.this.cluster_registration_token[0].node_command
      roles = [
        "etcd",
        "controlplane",
        "worker"
      ]
    })
  }
  network_interfaces {
    name = each.value.prov_net_int
    subnet_cidr = var.maas.prov_net_subnet
  }
}