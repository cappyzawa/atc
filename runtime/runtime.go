package runtime

import (
	"context"
	"io"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/resource"
	"github.com/concourse/atc/worker"
)

//go:generate counterfeiter . Resource

type Resource interface {
	Get(context.Context, worker.Volume, IOConfig, atc.Source, atc.Params, atc.Version) (resource.VersionedSource, error)
	Put(context.Context, IOConfig, atc.Source, atc.Params) (resource.VersionedSource, error)
	Check(atc.Source, atc.Version) ([]atc.Version, error)
}

//go:generate counterfeiter . Task

type Task interface {
	Run(context.Context, IOConfig, atc.TaskConfig) (chan TaskResult, []worker.VolumeMount, error)
}

// ExitStatus is the resulting exit code from the process that ran.
// Typically if the ExitStatus result is 0, the Success result is true.
type ExitStatus int

// kinda nasty; but we need to avoid an import cycle between exec and runtime
// until the TaskDelegate can be removed from there.
type TaskExecutionDelegate interface {
	worker.ImageFetchingDelegate

	Initializing(lager.Logger, atc.TaskConfig)
	Starting(lager.Logger, atc.TaskConfig)
	Finished(lager.Logger, ExitStatus)
}

//go:generate counterfeiter . Orchestrator
type Orchestrator interface {
	RunTask(
		context.Context,
		TaskExecutionDelegate,

		//TODO : See if this can be discerned from the TaskConfig
		db.ContainerOwner,
		db.ContainerMetadata,
		worker.ContainerSpec,
		//TODO

		creds.VersionedResourceTypes,
		IOConfig,
		atc.TaskConfig,
	) (chan TaskResult, []worker.VolumeMount, error)
	// Resource(atc.Source, atc.Params) Resource
}

type IOConfig struct {
	Stdout io.Writer
	Stderr io.Writer
}

type Process interface {
	Wait() (int, error)
}