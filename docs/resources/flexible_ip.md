---
page_title: "Scaleway: scaleway_flexible_ip"
description: |-
  Manages Scaleway Flexible IPs.
---

# scaleway_flexible_ip

Creates and manages Scaleway flexible IPs.
For more information, see [the documentation](https://developers.scaleway.com/en/products/flexible-ip/api).

## Examples

### Basic

```hcl
resource "scaleway_flexible_ip" "main" {
    reverse = "my-reverse.com"
}
```

### With zone

```hcl
resource "scaleway_flexible_ip" "main" {
    zone = "fr-par-2"
}
```

### With baremetal server

```hcl
resource "scaleway_account_ssh_key" "main" {
    name 	   = "main"
    public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAILHy/M5FVm5ydLGcal3e5LNcfTalbeN7QL/ZGCvDEdqJ foobar@example.com"
}

data "scaleway_baremetal_os" "by_id" {
    zone = "fr-par-2"
    name = "Ubuntu"
    version = "20.04 LTS (Focal Fossa)"						
}

data "scaleway_baremetal_offer" "my_offer" {
    zone = "fr-par-2"
    name = "EM-A210R-HDD"
}	

resource "scaleway_baremetal_server" "base" {
    zone = "fr-par-2"
    offer = data.scaleway_baremetal_offer.my_offer.offer_id
    os = data.scaleway_baremetal_os.by_id.os_id
    ssh_key_ids = scaleway_account_ssh_key.main.id
}

resource "scaleway_flexible_ip" "main" {
	server_id = scaleway_baremetal_server.base.id
	zone = "fr-par-2"
}
```

## Arguments Reference

The following arguments are supported:

- `description`: (Optional) A description of the flexible IP.
- `tags`: (Optional) A list of tags to apply to the flexible IP.
- `reverse` - (Optional) The reverse domain associated with this flexible IP.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Flexible IP

~> **Important:** Flexible IPs' IDs are [zoned](../guides/regions_and_zones.md#resource-ids), which means they are of the form `{zone}/{id}`, e.g. `fr-par-1/11111111-1111-1111-1111-111111111111`

- `ip_address` -  The IPv4 address of the Flexible IP
- `zone` - The zone of the Flexible IP
- `organization_id` - The organization of the Flexible IP
- `project_id` - The project of the Flexible IP
- `server_id` - The ID of the associated server

## Import

Flexible IPs can be imported using the `{zone}/{id}`, e.g.

```bash
$ terraform import scaleway_flexible_ip.main fr-par-1/11111111-1111-1111-1111-111111111111
```
