data "sshconfig_host" "github_ssh_config" {
  path = pathexpand("~/.ssh/config")
  host = "github.com"
}

output "github_config" {
  value = "${data.sshconfig_host.github_ssh_config.rendered}"
}

output "host_map" {
  value = "${data.sshconfig_host.github_ssh_config.host_map}"
}