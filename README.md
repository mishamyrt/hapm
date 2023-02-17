# Home Assistant Package Manager

Package manager for Home Assistant which allows you to prepare components for a Stateless Docker image.

## Installation

```sh
pip install hapm
```

## Configuration

A manifest (`hapm.yaml`) is a list of links, divided into several categories. The links are written in a special format and necessarily have a version with `@` in it.

For example: `github.com/user/integration@v1.0.0`.

The version specifies an existing tag or branch in the specified repository.

A link may not have an https prefix, in which case it will be inserted automatically during the reading phase.

```yaml
---
integrations:
  - github.com/mishamyrt/dohome_rgb@v0.3.0
  - github.com/mishamyrt/myrt_desk_hass@master
```

## Sync remote packages

```sh
hapm
```

## Export 

```sh
hapm put -t integrations <path>
```

## List 

```sh
hapm list
```

## Add new package

**not yet implemented**

```sh
hapm get -t integration https://github.com/AlexxIT/XiaomiGateway3
```

## Update

**not yet implemented**

```sh
# Update all packages to latest version
hapm update
# Update all packages to latest version
hapm update <url>
# Only print
hapm update --dry
```
