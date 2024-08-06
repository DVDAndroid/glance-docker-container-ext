package main

import (
	docker "github.com/fsouza/go-dockerclient"
	"sort"
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

func LoadContainers(dockerClient *docker.Client, p params) ([]DockerContainer, error) {
	containerList, err := dockerClient.ListContainers(docker.ListContainersOptions{
		All: p.AllContainers,
		Filters: map[string][]string{
			"label": {"glance.enable=true"},
		},
	})

	if err != nil {
		return nil, err
	}

	containers := make([]DockerContainer, len(containerList))
	for i, container := range containerList {
		name := container.Names[0][1:]
		description := ""
		url := ""
		icon := ""
		group := ""
		sameTab := false

		for label, value := range container.Labels {
			switch label {
			case "glance.name":
				name = value
			case "glance.description":
				description = value
			case "glance.group":
				group = value
			case "glance.icon":
				icon = value
				if strings.HasPrefix(value, "si:") {
					icon = strings.TrimPrefix(value, "si:")
					icon = "https://cdnjs.cloudflare.com/ajax/libs/simple-icons/11.14.0/" + icon + ".svg"
				}
			case "glance.url":
				url = value
			case "glance.same-tab":
				sameTab = value == "true"
			}
		}

		if group != p.Group {
			continue
		}

		state := container.State
		if p.IgnoreStatus {
			state = ""
		}

		containers[i] = DockerContainer{
			Name:        name,
			Description: description,
			State:       state,
			Status:      container.Status,
			Icon:        icon,
			IsSvgIcon:   strings.Contains(icon, "/simple-icons/") || strings.HasSuffix(icon, ".svg"),
			URL:         url,
			SameTab:     p.SameTab || sameTab,
		}
	}

	for i := 0; i < len(containers); i++ {
		// happens if container group is different than the requested group
		if containers[i].Name == "" {
			containers = append(containers[:i], containers[i+1:]...)
			i--
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
