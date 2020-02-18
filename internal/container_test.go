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

func TestCreate(t *testing.T) {
	config := miniConfig("create")

	err := config.Create(true)
	if err != nil {
		t.Errorf("Error creating container: %s", err)
	}
	status, _ := config.Status()
	if status != 3 {
		t.Errorf("Create did not leave container running (container is %s)", ParseStatus(status))
	}
	config.Down()

	err = config.Create(false)
	if err != nil {
		t.Errorf("Error creating container: %s", err)
	}
	status, _ = config.Status()
	if status == 3 {
		t.Errorf("Create left container running")
	}
	config.Down()
}

func TestName(t *testing.T) {
	config := miniConfig("name")
	name := config.Name()
	expectStrEq("test_name", name, t)

	config.Fields.Name = "named_container"
	name = config.Name()
	expectStrEq("named_container", name, t)
}
