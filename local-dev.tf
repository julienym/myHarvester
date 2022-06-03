resource "local_file" "harv_kubeconfig" {
    content  = rancher2_cluster_v2.this.kube_config
    filename = pathexpand("~/.kube/clusters/${local.rke2_name}")
}
