data "ad_gpo" "g" {
    name = "Some GPO"
}

output "gpo_uuid" {
    value = data.ad_gpo.g.guid
}
