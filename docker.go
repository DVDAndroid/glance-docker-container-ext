package main

import (
	docker "github.com/fsouza/go-dockerclient"
	"sort"
	"strconv"
	"strings"
)

type DockerContainer struct {
	Name        string
	Description string
	State       string // created|restarting|running|removing|paused|exited|dead
	Status      string
	Icon        string
	IsSvgIcon   bool
	URL         string
	SameTab     bool
}

type GlanceLabel struct {
	Enable      bool
	Name        string
	Description string
	Url         string
	Icon        string
	Group       string
	SameTab     bool
}

func LoadContainers(dockerClient *docker.Client, p params) ([]DockerContainer, error) {
	containerList, err := dockerClient.ListContainers(docker.ListContainersOptions{
		All: p.AllContainers,
	})

	if err != nil {
		return nil, err
	}

	var containers []DockerContainer
	for _, container := range containerList {
		glanceLabels := make(map[int]GlanceLabel)

		for label, value := range container.Labels {
			if !strings.HasPrefix(label, "glance.") {
				continue
			}

			parts := strings.Split(label, ".")
			if len(parts) != 3 {
				continue
			}

			index, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
			}

			gl, exists := glanceLabels[index]
			if !exists {
				glanceLabels[index] = GlanceLabel{}
				gl = glanceLabels[index]
			}
			switch parts[2] {
			case "enable":
				gl.Enable = value == "true"
			case "name":
				gl.Name = value
			case "description":
				gl.Description = value
			case "group":
				gl.Group = value
			case "icon":
				gl.Icon = value
				if strings.HasPrefix(value, "si:") {
					gl.Icon = strings.TrimPrefix(value, "si:")
					gl.Icon = "https://cdnjs.cloudflare.com/ajax/libs/simple-icons/11.14.0/" + gl.Icon + ".svg"
				}
			case "url":
				gl.Url = value
			case "same-tab":
				gl.SameTab = value == "true"
			}
			glanceLabels[index] = gl
		}

		for _, gl := range glanceLabels {
			if !gl.Enable {
				continue
			}
			if gl.Group != p.Group {
				continue
			}

			state := container.State
			if p.IgnoreStatus {
				state = ""
			}

			if gl.Name == "" {
				gl.Name = container.Names[0][1:]
			}

			containers = append(containers, DockerContainer{
				Name:        gl.Name,
				Status:      container.Status,
				State:       state,
				Description: gl.Description,
				Icon:        gl.Icon,
				IsSvgIcon:   strings.Contains(gl.Icon, "/simple-icons/") || strings.HasSuffix(gl.Icon, ".svg"),
				URL:         gl.Url,
				SameTab:     p.SameTab || gl.SameTab,
			})
		}
	}

	sortContainers(containers, strings.Split(p.Order, ","))

	return containers, nil
}

func sortContainers(containers []DockerContainer, order []string) {
	sort.Slice(containers, func(i, j int) bool {
		for _, field := range order {
			switch field {
			case "name":
				name1 := strings.ToLower(containers[i].Name)
				name2 := strings.ToLower(containers[j].Name)
				if name1 != name2 {
					return name1 < name2
				}
				description1 := strings.ToLower(containers[i].Description)
				description2 := strings.ToLower(containers[j].Description)
				if description1 != description2 {
					return description1 < description2
				}
			case "status":
				if containers[i].State != containers[j].State {
					return containers[i].State < containers[j].State
				}
			}
		}
		return false
	})
}
