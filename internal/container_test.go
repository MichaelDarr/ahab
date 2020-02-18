package internal

import (
	"os"
	"testing"
)

func TestGetContainer(t *testing.T) {
	os.Chdir(noConfDir)
	container, err := GetContainer()
	if err == nil {
		t.Errorf("Unexpected success finding project config: %s", err)
	}

	os.Chdir(exampleConfDir)
	container, err = GetContainer()
	if err != nil {
		t.Errorf("Error finding project config: %s", err)
	} else if container.FilePath != exampleConfPath {
		t.Errorf("Error finding project config path. Expected %s, Found %s", exampleConfPath, container.FilePath)
	}

	os.Chdir(exampleConfChildDir)
	container, err = GetContainer()
	if err != nil {
		t.Errorf("Error finding project config from child dir: %s", err)
	} else if container.FilePath != exampleConfPath {
		t.Errorf("Error finding project config path from child dir. Expected %s, Found %s", exampleConfPath, container.FilePath)
	}
}

func TestCmd(t *testing.T) {
	container := miniContainer("cmd")
	container.Create(true)
	expectContainerStatus(3, container, t)

	err := container.Cmd("pause")
	if err != nil {
		t.Errorf("Error executing pause cmd on container: %s", err)
	}
	expectContainerStatus(5, container, t)

	err = container.Cmd("unpause")
	if err != nil {
		t.Errorf("Error executing unpause cmd on container: %s", err)
	}
	expectContainerStatus(3, container, t)

	err = container.Cmd("stop")
	if err != nil {
		t.Errorf("Error executing stop cmd on container: %s", err)
	}
	expectContainerStatus(6, container, t)

	err = container.Cmd("start")
	if err != nil {
		t.Errorf("Error executing start cmd on container: %s", err)
	}
	expectContainerStatus(3, container, t)
	container.Down()
}

func TestCreate(t *testing.T) {
	container := miniContainer("create")
	err := container.Create(true)
	if err != nil {
		t.Errorf("Error creating container: %s", err)
	}
	expectContainerStatus(3, container, t)
	container.Down()

	err = container.Create(false)
	if err != nil {
		t.Errorf("Error creating container: %s", err)
	}
	prohibitContainerStatus(3, container, t)
	container.Down()
}

func TestDown(t *testing.T) {
	container := miniContainer("down")
	container.Create(true)
	expectContainerStatus(3, container, t)
	err := container.Down()
	if err != nil {
		t.Errorf("Error down-ing container: %s", err)
	}
	expectContainerStatus(0, container, t)

	container.Create(false)
	expectContainerStatus(6, container, t)
	err = container.Down()
	if err != nil {
		t.Errorf("Error down-ing container: %s", err)
	}
	expectContainerStatus(0, container, t)
}

func TestName(t *testing.T) {
	container := miniContainer("name")
	name := container.Name()
	expectStrEq("test_name", name, t)

	container.Fields.Name = "named_container"
	name = container.Name()
	expectStrEq("named_container", name, t)
}

func TestProp(t *testing.T) {
	container := miniContainer("prop")
	status, err := container.Prop("State.Status")
	if err != nil {
		t.Error("Unexpected error on prop check for nonexistent container")
	} else if status != "" {
		t.Errorf("Non-empty string on prop check for nonexistent container. Received: %s", status)
	}

	container.Create(true)
	expectContainerStatus(3, container, t)
	status, err = container.Prop("State.Status")
	if err != nil {
		t.Errorf("Unexpected an error on prop check: %s", err)
	} else if status == "" {
		t.Error("Empty string receieved on prop check for existing container")
	}
	container.Down()
}
