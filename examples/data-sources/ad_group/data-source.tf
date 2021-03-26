data "ad_group" "g" {
    group_id = "DC3E5929-71C0-4232-9C32-9C7AFAABF0BB"
}

output "groupname" {
    value = data.ad_group.g.name
}

data "ad_group" "g2" {
    group_id = "some_group_sam_account_name"
}

output "g2_guid" {
    value = data.ad_group.g2.id
}

output "g2_description" {
    value = data.ad_group.g2.description
}