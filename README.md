glance-docker-container-ext
===

![Docker Image Size](https://img.shields.io/docker/image-size/dvdandroid/glance-docker-container-ext)
![Docker Image Version](https://img.shields.io/docker/v/dvdandroid/glance-docker-container-ext)

[Glance](https://github.com/glanceapp/glance) extension that creates a widget that displays the running Docker containers.

![Sample Screenshot](./assets/screen.png)

## Installation

Assuming you are using Docker compose, add the following to your `docker-compose.yml` file containing Glance container:

```yaml
glance-docker-container-ext:
  image: dvdandroid/glance-docker-container-ext
  container_name: glance-docker-container-ext
  restart: unless-stopped
  environment:
    - DOCKER_HOST=unix:///var/run/docker.sock
    - PORT=8081 # Optional, default is 8081
  volumes:
    - /var/run/docker.sock:/var/run/docker.sock
```

then in your `glance.yml` config file, add the following:

```yaml
- type: extension
  allow-potentially-dangerous-html: true
  url: http://glance-docker-container-ext:8081
  cache: 5m
  parameters:
    title: Docker Containers
    all: true
    order: name,status
```

### Parameters

| Parameter       | Description                                                                                                              | Default             |
|-----------------|--------------------------------------------------------------------------------------------------------------------------|---------------------|
| `title`         | Title of the widget                                                                                                      | "Docker Containers" |
| `all`           | Show all containers or only running ones                                                                                 | `true`              |
| `order`         | Order of the containers, comma separated **string** of `name`, `status`<br>(`name`,`status`,`name,status`,`status,name`) | `name`              |
| `group`         | Identifier for the group of containers. If set, only containers with the same group will be displayed.                   |                     |
| `same-tab`      | Open the URL in the same tab. Value customizable per container                                                           | `false`             |
| `ignore-status` | Status of the containers will not be displayed                                                                           | `false`             |

## Configuration

Then, for every container you want to monitor, add the following labels to its compose file:

```yaml
labels:
  glance.0.enable: true
  glance.0.name: Sonarr
  glance.0.description: TV show search
  glance.0.group: media
  glance.0.url: http://sonarr.lan
  glance.0.icon: ./assets/imgs/television-classic.svg
```

:warning: Multiple labels can be added to the same container. Read below

| Label                  | Description                                                                                                | Default        |
|------------------------|------------------------------------------------------------------------------------------------------------|----------------|
| `glance.X.enable`      | Enable monitoring for this container                                                                       |                |
| `glance.X.name`        | Name of the container                                                                                      | container name |
| `glance.X.description` | Description of the container                                                                               |                |
| `glance.X.group`       | Identifier for the group of containers, used in combination with parameter `group` in glance configuration |                |
| `glance.X.url`         | URL to open when clicking on the container                                                                 |                |
| `glance.X.icon`        | Icon to display, pointing to assets or Simple Icon (`si:` prefix)                                          |                |
| `glance.X.same-tab`    | Open the URL in the same tab                                                                               | `false`        |

Value of `X` must be replaced with a number starting from 0: this allows to add multiple widget referring to the same container.

For example, if you want to define two widgets for the same container, but with different labels, you can do it like this (in the compose file):

```yaml
labels:
  glance.0.enable: true
  glance.0.name: Container (ADMIN)
  glance.0.url: http://website.lan/admin
  glance.1.enable: true
  glance.1.group: User
  glance.1.name: Container (USER)
  glance.1.description: User access
  glance.1.url: http://website.lan/user
```
