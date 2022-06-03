module "harvester-crd" {
  depends_on = [
    null_resource.provisionner,
    # null_resource.healthcheck,
    module.rancher-harvester
  ]
  source = "git::https://github.com/julienym/myTerraformModules.git//helm?ref=1.0.0"
  name = "harvester-crd"
  chart = "charts/harvester/deploy/charts/harvester-crd"
  namespace = "harvester-system"
}

module "harvester" {
  depends_on = [
    module.harvester-crd
  ]
  source = "git::https://github.com/julienym/myTerraformModules.git//helm?ref=1.0.0"
  name = "harvester"
  #TODO get remote chart
  chart = "charts/harvester/deploy/charts/harvester"
  namespace = "harvester-system"
  # values_file = "../../charts/harvester-install/iso-values.yml"
  values_file = "../../charts/harvester-install/default.yml"
  values = { #Does this means all node need the same interface name ????
    "kube-vip.config.vip_interface" = var.maas.instances.chronos.prov_net_int
    "kube-vip.config.vip_ip" = var.harvester.vip
    "service.vip.ip" = var.harvester.vip
  }
}

resource "null_resource" "expose_harvester" {
  depends_on = [
    module.harvester
  ]
  provisioner "local-exec" {
    command = <<EOF
sleep 5m
echo -n $KUBE64 | base64 -d > /tmp/.cluster
export KUBECONFIG="/tmp/.cluster"
kubectl create configmap -n kube-system kubevip --from-literal cidr-global=${var.harvester.vip}/32
kubectl patch svc/rancher -n cattle-system --type='merge' --patch '{"spec": {"type": "LoadBalancer"}}'
EOF

    environment = {
      KUBE64 = base64encode(rancher2_cluster_v2.this.kube_config)
    }
  }
}

# resource "harvester_image" "ubuntu20" {
#   depends_on = [
#     module.harvester
#   ]
#   name      = "ubuntu20"
#   namespace = "harvester-public"

#   display_name = "ubuntu20"
#   source_type  = "download"
#   url          = "https://cloud-images.ubuntu.com/focal/current/focal-server-cloudimg-amd64.img"
# }

# resource "harvester_image" "ubuntu18" {
#   depends_on = [
#     module.harvester
#   ]
#   name      = "ubuntu18"
#   namespace = "harvester-public"

#   display_name = "ubuntu18"
#   source_type  = "download"
#   url          = "https://cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64.img"
# }