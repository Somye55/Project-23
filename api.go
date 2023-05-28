package main

import (
	"log"
	"os"
	"time"

	"github.com/Somye55/Project-23/config"
	"github.com/Somye55/Project-23/docker"
	"github.com/Somye55/Project-23/event"
	"github.com/Somye55/Project-23/handlers"
	"github.com/Somye55/Project-23/id"
	"github.com/Somye55/Project-23/k8s"
	"github.com/Somye55/Project-23/provisioner"
	"github.com/Somye55/Project-23/pwd"
	"github.com/Somye55/Project-23/pwd/types"
	"github.com/Somye55/Project-23/scheduler"
	"github.com/Somye55/Project-23/scheduler/task"
	"github.com/Somye55/Project-23/storage"
)

func main() {
	config.ParseFlags()

	e := initEvent()
	s := initStorage()
	df := initDockerFactory(s)
	kf := initK8sFactory(s)

	ipf := provisioner.NewInstanceProvisionerFactory(provisioner.NewWindowsASG(df, s), provisioner.NewDinD(id.XIDGenerator{}, df, s))
	sp := provisioner.NewOverlaySessionProvisioner(df)

	core := pwd.NewPWD(df, e, s, sp, ipf)

	tasks := []scheduler.Task{
		task.NewCheckPorts(e, df),
		task.NewCheckSwarmPorts(e, df),
		task.NewCheckSwarmStatus(e, df),
		task.NewCollectStats(e, df, s),
		task.NewCheckK8sClusterStatus(e, kf),
		task.NewCheckK8sClusterExposedPorts(e, kf),
	}
	sch, err := scheduler.NewScheduler(tasks, s, e, core)
	if err != nil {
		log.Fatal("Error initializing the scheduler: ", err)
	}

	sch.Start()

	d, err := time.ParseDuration("4h")
	if err != nil {
		log.Fatalf("Cannot parse duration Got: %v", err)
	}

	playground := types.Playground{Domain: config.PlaygroundDomain, DefaultDinDInstanceImage: "franela/dind", AvailableDinDInstanceImages: []string{"franela/dind"}, AllowWindowsInstances: config.NoWindows, DefaultSessionDuration: d, Extras: map[string]interface{}{"LoginRedirect": "http://localhost:3000"}, Privileged: true}
	if _, err := core.PlaygroundNew(playground); err != nil {
		log.Fatalf("Cannot create default playground. Got: %v", err)
	}

	handlers.Bootstrap(core, e)
	handlers.Register(nil)
}

func initStorage() storage.StorageApi {
	s, err := storage.NewFileStorage(config.SessionsFile)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Error initializing StorageAPI: ", err)
	}
	return s
}

func initEvent() event.EventApi {
	return event.NewLocalBroker()
}

func initDockerFactory(s storage.StorageApi) docker.FactoryApi {
	return docker.NewLocalCachedFactory(s)
}

func initK8sFactory(s storage.StorageApi) k8s.FactoryApi {
	return k8s.NewLocalCachedFactory(s)
}
