resource "ad_ou" "o" { 
    name = "gplinktestOU"
    path = "dc=yourdomain,dc=com"
    description = "OU for gplink tests"
    protected = false
}
