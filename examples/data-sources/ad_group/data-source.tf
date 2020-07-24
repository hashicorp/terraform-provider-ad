data "ad_group" "g" {
    guid = "DC3E5929-71C0-4232-9C32-9C7AFAABF0BB"
}

output "groupname" {
    value = data.ad_group.g.name
}
