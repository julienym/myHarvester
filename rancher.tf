resource "rancher2_cluster_v2" "this" {
  name = local.rke2_name
  
  # Keep compatible with Harvester version
  # https://github.com/harvester/harvester/releases
  kubernetes_version = var.harvester.rke2_version

  rke_config {
    machine_global_config = <<EOF
cni: "multus,canal"
EOF
  }
}

resource "null_resource" "healthcheck" {
  depends_on = [
    maas_instance.this
  ]
  provisioner "local-exec" {
    command = "bash scripts/healthcheck.sh"

    environment = {
      KUBE64 = base64encode(rancher2_cluster_v2.this.kube_config)
    }
  }
}

module "rancher-harvester" {
  depends_on = [
    # null_resource.healthcheck
    null_resource.provisionner
  ]
  source = "git::https://github.com/julienym/myTerraformModules.git//helm?ref=1.0.0"
  name = "rancher"
  repository = "https://releases.rancher.com/server-charts/stable"
  namespace = "cattle-system"
  chart = "rancher"
  chart_version = "2.6.4"
  values = {
    "tls" = "external"
    "rancherImage" = "rancher/rancher"
    "rancherImageTag" = "v2.6.4-harvester3"
    "noDefaultAdmin" = false
    "features" = "multi-cluster-management=false\\,multi-cluster-management-agent=false"
    "useBundledSystemChart" = true
    "bootstrapPassword" = var.harvester_bootstrap_passwd
  }
}

# resource "null_resource" "provisionner" {
#   depends_on = [
#     null_resource.healthcheck
#     # module.rancher-harvester
#   ]
#   provisioner "local-exec" {
#     command = <<EOF
# echo -n $KUBE64 | base64 -d > /tmp/.cluster
# export KUBECONFIG="/tmp/.cluster"
# kubectl patch clusters.management.cattle.io/${rancher2_cluster_v2.this.cluster_v1_id} --type='merge' --patch '{"status": {"provider": "harvester"}}'
# EOF

#     environment = {
#       KUBE64 = base64encode(rancher2_cluster_v2.this.kube_config)
#     }
#   }
# }

resource "null_resource" "provisionner" {
  depends_on = [
    null_resource.healthcheck
    # module.rancher-harvester
  ]
  provisioner "local-exec" {
    command = <<EOF
export KUBECONFIG="/home/julien/.kube/clusters/rancher"
kubectl patch clusters.management.cattle.io/${rancher2_cluster_v2.this.cluster_v1_id} --type='merge' --patch '{"status": {"provider": "harvester"}}'
EOF
  }
}
