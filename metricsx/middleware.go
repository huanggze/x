package metricsx

import (
	"math"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"

	"github.com/huanggze/x/cmdx"
	"github.com/huanggze/x/configx"
	"github.com/huanggze/x/logrusx"
	"github.com/huanggze/x/resilience"
	"github.com/ory/analytics-go/v5"
)

var instance *Service
var lock sync.Mutex

// Service helps with providing context on metrics.
type Service struct {
	optOut     bool
	instanceId string

	o *Options

	c analytics.Client
	l *logrusx.Logger

	mem *MemoryStatistics
}

// Options configures the metrics service.
type Options struct {
	// Service represents the service name, for example "ory-hydra".
	Service string

	// DeploymentId represents the cluster id, typically a hash of some unique configuration properties.
	DeploymentId string

	DBDialect string

	// When this instance was started
	StartTime time.Time

	// IsDevelopment should be true if we assume that we're in a development environment.
	IsDevelopment bool

	// WriteKey is the segment API key.
	WriteKey string

	// BuildVersion represents the build version.
	BuildVersion string

	// BuildHash represents the build git hash.
	BuildHash string

	// BuildTime represents the build time.
	BuildTime string

	// Config overrides the analytics.Config. If nil, sensible defaults will be used.
	Config *analytics.Config

	// MemoryInterval sets how often memory statistics should be transmitted. Defaults to every 12 hours.
	MemoryInterval time.Duration
}

type void struct {
}

func (v *void) Logf(format string, args ...interface{}) {
}

func (v *void) Errorf(format string, args ...interface{}) {
}

// New returns a new metrics service. If one has been instantiated already, no new instance will be created.
func New(
	cmd *cobra.Command,
	l *logrusx.Logger,
	c *configx.Provider,
	o *Options,
) *Service {
	lock.Lock()
	defer lock.Unlock()

	if instance != nil {
		return instance
	}

	o.StartTime = time.Now()

	if o.BuildTime == "" {
		o.BuildTime = "unknown"
	}

	if o.BuildVersion == "" {
		o.BuildVersion = "unknown"
	}

	if o.BuildHash == "" {
		o.BuildHash = "unknown"
	}

	if o.Config == nil {
		o.Config = &analytics.Config{
			Interval: time.Hour * 6,
		}
	}

	o.Config.Logger = new(void)

	if o.MemoryInterval < time.Minute {
		o.MemoryInterval = time.Hour * 12
	}

	segment, err := analytics.NewWithConfig(o.WriteKey, *o.Config)
	if err != nil {
		l.WithError(err).Fatalf("Unable to initialise software quality assurance features.")
		return nil
	}

	optOut, err := cmd.Flags().GetBool("sqa-opt-out")
	if err != nil {
		cmdx.Must(err, `Unable to get command line flag "sqa-opt-out": %s`, err)
	}

	if !optOut {
		optOut = c.Bool("sqa.opt_out")
	}

	if !optOut {
		optOut = c.Bool("sqa_opt_out")
	}

	if !optOut {
		optOut, _ = strconv.ParseBool(os.Getenv("SQA_OPT_OUT"))
	}

	if !optOut {
		optOut, _ = strconv.ParseBool(os.Getenv("SQA-OPT-OUT"))
	}

	if !optOut {
		l.Info("Software quality assurance features are enabled. Learn more at: https://www.ory.sh/docs/ecosystem/sqa")
	}

	m := &Service{
		optOut:     optOut,
		instanceId: uuid.Must(uuid.NewV4()).String(),
		o:          o,
		c:          segment,
		l:          l,
		mem:        new(MemoryStatistics),
	}

	instance = m

	go m.Identify()
	go m.Track()

	return m
}

// Identify enables reporting to segment.
func (sw *Service) Identify() {
	IdentifySend(sw, true)

	// User has not opt-out then make identify to be sent every 6 hours
	if !sw.optOut {
		for range time.Tick(time.Hour * 6) {
			IdentifySend(sw, false)
		}
	}
}

func IdentifySend(sw *Service, startup bool) {
	if err := resilience.Retry(sw.l, time.Minute*5, time.Hour*6, func() error {
		return sw.c.Enqueue(analytics.Identify{
			InstanceId:   sw.instanceId,
			DeploymentId: sw.o.DeploymentId,
			Project:      sw.o.Service,

			DatabaseDialect:  sw.o.DBDialect,
			ProductVersion:   sw.o.BuildVersion,
			ProductBuild:     sw.o.BuildHash,
			UptimeDeployment: 0,
			UptimeInstance:   math.Round(time.Since(sw.o.StartTime).Seconds()),
			IsDevelopment:    sw.o.IsDevelopment,
			IsOptOut:         sw.optOut,
			Startup:          startup,
		})
	}); err != nil {
		sw.l.WithError(err).Debug("Could not commit anonymized environment information")
	}
}

// Track commits memory statistics to segment.
func (sw *Service) Track() {
	if sw.optOut {
		return
	}

	for {
		sw.mem.Update()
		if err := sw.c.Enqueue(analytics.Track{
			InstanceId:   sw.instanceId,
			DeploymentId: sw.o.DeploymentId,
			Project:      sw.o.Service,

			CPU:            runtime.NumCPU(),
			OsName:         runtime.GOOS,
			OsArchitecture: runtime.GOARCH,
			Alloc:          sw.mem.Alloc,
			TotalAlloc:     sw.mem.TotalAlloc,
			Frees:          sw.mem.Frees,
			Mallocs:        sw.mem.Mallocs,
			Lookups:        sw.mem.Lookups,
			Sys:            sw.mem.Sys,
			NumGC:          sw.mem.NumGC,
			HeapAlloc:      sw.mem.HeapAlloc,
			HeapInuse:      sw.mem.HeapInuse,
			HeapIdle:       sw.mem.HeapIdle,
			HeapObjects:    sw.mem.HeapObjects,
			HeapReleased:   sw.mem.HeapReleased,
			HeapSys:        sw.mem.HeapSys,
		}); err != nil {
			sw.l.WithError(err).Debug("Could not commit anonymized telemetry data")
		}
		time.Sleep(sw.o.MemoryInterval)
	}
}
