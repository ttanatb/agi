// Copyright (C) 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"google.golang.org/grpc/grpclog"

	"github.com/google/gapid/core/app"
	"github.com/google/gapid/core/app/auth"
	"github.com/google/gapid/core/app/crash"
	"github.com/google/gapid/core/event/task"
	"github.com/google/gapid/core/log"
	"github.com/google/gapid/core/os/android/adb"
	"github.com/google/gapid/core/os/device/bind"
	"github.com/google/gapid/core/os/device/remotessh"
	"github.com/google/gapid/core/os/file"
	"github.com/google/gapid/core/os/fuchsia/ffx"
	"github.com/google/gapid/core/text"
	"github.com/google/gapid/gapir/client"
	"github.com/google/gapid/gapis/database"
	"github.com/google/gapid/gapis/replay"
	"github.com/google/gapid/gapis/server"
	"github.com/google/gapid/gapis/service"
	"github.com/google/gapid/gapis/service/path"
	"github.com/google/gapid/gapis/stringtable"
	"github.com/google/gapid/gapis/trace"
)

var (
	rpc              = flag.String("rpc", "localhost:0", "TCP host:port of the server's RPC listener")
	stringsPath      = flag.String("strings", "strings", "_Directory containing string table packages")
	persist          = flag.Bool("persist", false, "Server will keep running even when no connections remain")
	gapisAuthToken   = flag.String("gapis-auth-token", "", "_The connection authorization token for gapis")
	gapirAuthToken   = flag.String("gapir-auth-token", "", "_The connection authorization token for gapir")
	gapirArgStr      = flag.String("gapir-args", "", "_The arguments to be passed to the host-run gapir")
	scanAndroidDevs  = flag.Bool("monitor-android-devices", true, "Server will scan for locally connected Android devices")
	scanFuchsiaDevs  = flag.Bool("monitor-fuchsia-devices", true, "Server will scan for locally connected Fuchsia devices")
	androidSerial    = flag.String("android-serial", "", "Server will only consider the Android device with this serial id")
	addLocalDevice   = flag.Bool("add-local-device", true, "Server can trace and replay locally")
	idleTimeout      = flag.Duration("idle-timeout", 0, "_Closes GAPIS if the server is not repeatedly pinged within this duration (e.g. '30s', '2m'). Default: 0 (no timeout).")
	adbPath          = flag.String("adb", "", "Path to the adb executable; leave empty to search the environment")
	enableLocalFiles = flag.Bool("enable-local-files", false, "Allow clients to access local .gfxtrace files by path")
	remoteSSHConfig  = flag.String("ssh-config", "", "_Path to an ssh config file for remote devices")
	preloadDepGraph  = flag.Bool("preload-dep-graph", true, "_Preload the dependency graph when loading captures")
)

func main() {
	app.ShortHelp = "GAPIS is the graphics API server"
	app.Name = "GAPIS" // Has to be this for version parsing compatability
	app.Run(run)
}

type testBigStruct struct {
    a1  int64
    a2  int64
    a3  int64
    a4  int64
    a5  int64
    a6  int64
    a7  int64
    a8  int64
    a9  int64
    a0  int64
	aa1  int64
    aa2  int64
    aa3  int64
    aa4  int64
    aa5  int64
    aa6  int64
    aa7  int64
    aa8  int64
    aa9  int64
    aa0  int64
	ba1  int64
    ba2  int64
    ba3  int64
    ba4  int64
    ba5  int64
    ba6  int64
    ba7  int64
    ba8  int64
    ba9  int64
    ba0  int64
	baa1  int64
    baa2  int64
    baa3  int64
    baa4  int64
    baa5  int64
    baa6  int64
    baa7  int64
    baa8  int64
    baa9  int64
    baa0  int64
	ca1  int64
    ca2  int64
    ca3  int64
    ca4  int64
    ca5  int64
    ca6  int64
    ca7  int64
    ca8  int64
    ca9  int64
    ca0  int64
	caa1  int64
    caa2  int64
    caa3  int64
    caa4  int64
    caa5  int64
    caa6  int64
    caa7  int64
    caa8  int64
    caa9  int64
    caa0  int64
	cba1  int64
    cba2  int64
    cba3  int64
    cba4  int64
    cba5  int64
    cba6  int64
    cba7  int64
    cba8  int64
    cba9  int64
    cba0  int64
	cbaa1  int64
    cbaa2  int64
    cbaa3  int64
    cbaa4  int64
    cbaa5  int64
    cbaa6  int64
    cbaa7  int64
    cbaa8  int64
    cbaa9  int64
    cbaa0  int64
	
	da1  int64
    da2  int64
    da3  int64
    da4  int64
    da5  int64
    da6  int64
    da7  int64
    da8  int64
    da9  int64
    da0  int64
	daa1  int64
    daa2  int64
    daa3  int64
    daa4  int64
    daa5  int64
    daa6  int64
    daa7  int64
    daa8  int64
    daa9  int64
    daa0  int64
	dba1  int64
    dba2  int64
    dba3  int64
    dba4  int64
    dba5  int64
    dba6  int64
    dba7  int64
    dba8  int64
    dba9  int64
    dba0  int64
	dbaa1  int64
    dbaa2  int64
    dbaa3  int64
    dbaa4  int64
    dbaa5  int64
    dbaa6  int64
    dbaa7  int64
    dbaa8  int64
    dbaa9  int64
    dbaa0  int64
	dca1  int64
    dca2  int64
    dca3  int64
    dca4  int64
    dca5  int64
    dca6  int64
    dca7  int64
    dca8  int64
    dca9  int64
    dca0  int64
	dcaa1  int64
    dcaa2  int64
    dcaa3  int64
    dcaa4  int64
    dcaa5  int64
    dcaa6  int64
    dcaa7  int64
    dcaa8  int64
    dcaa9  int64
    dcaa0  int64
	dcba1  int64
    dcba2  int64
    dcba3  int64
    dcba4  int64
    dcba5  int64
    dcba6  int64
    dcba7  int64
    dcba8  int64
    dcba9  int64
    dcba0  int64
	dcbaa1  int64
    dcbaa2  int64
    dcbaa3  int64
    dcbaa4  int64
    dcbaa5  int64
    dcbaa6  int64
    dcbaa7  int64
    dcbaa8  int64
    dcbaa9  int64
    dcbaa0  int64
}

// features is the reported list of features supported by the server.
// This feature list can be used by the client to determine what new RPCs can be
// called.
var features = []string{}

// addFallbackLogHandler adds a handler to b that calls to fallback, only when
// nothing else is listening, i.e. nothing is listening to the log stream RPC.
func addFallbackLogHandler(b *log.Broadcaster, fallback log.Handler) {
	b.Listen(log.NewHandler(func(m *log.Message) {
		if b.Count() == 1 {
			fallback.Handle(m)
		}
	}, fallback.Close))
}

func run(ctx context.Context) error {
	log.I(ctx, "%+v", testBigStruct{})

	logBroadcaster := log.Broadcast()
	if oldHandler, oldWasDefault := app.LogHandler.SetTarget(logBroadcaster, false); oldWasDefault {
		addFallbackLogHandler(logBroadcaster, oldHandler)
	} else {
		logBroadcaster.Listen(oldHandler)
	}

	if *adbPath != "" {
		adb.ADB = file.Abs(*adbPath)
	}

	r := bind.NewRegistry()
	ctx = bind.PutRegistry(ctx, r)
	m := replay.New(ctx)
	ctx = replay.PutManager(ctx, m)
	ctx = trace.PutManager(ctx, trace.New(ctx))
	ctx = database.Put(ctx, database.NewInMemory(ctx))

	// Grpc is very verbose, turn that down
	grpclog.SetLogger(log.From(ctx).SetFilter(log.SeverityFilter(log.Error)))

	var hostDevice *path.Device

	if *addLocalDevice {
		host := bind.Host(ctx)
		hostDevice = path.NewDevice(host.Instance().ID.ID())
		r.AddDevice(ctx, host)
		r.SetDeviceProperty(ctx, host, client.LaunchArgsKey, text.SplitArgs(*gapirArgStr))
	}

	wg := sync.WaitGroup{}

	if *androidSerial != "" {
		adb.LimitToSerial(*androidSerial)
	}

	if *scanAndroidDevs {
		wg.Add(1)
		crash.Go(func() { monitorAndroidDevices(ctx, r, wg.Done) })
	}

	if *scanFuchsiaDevs {
		wg.Add(1)
		crash.Go(func() { monitorFuchsiaDevices(ctx, r, wg.Done) })
	}

	if *remoteSSHConfig != "" {
		wg.Add(1)
		crash.Go(func() { monitorRemoteSSHDevices(ctx, r, wg.Done) })
	}

	deviceScanDone, onDeviceScanDone := task.NewSignal()
	crash.Go(func() {
		wg.Wait()
		onDeviceScanDone(ctx)
	})

	hostname, err := os.Hostname()
	if err != nil {
		return log.Err(ctx, err, "Failed to retrieve hostname")
	}

	return server.Listen(ctx, *rpc, server.Config{
		Info: &service.ServerInfo{
			Name:              hostname,
			VersionMajor:      uint32(app.Version.Major),
			VersionMinor:      uint32(app.Version.Minor),
			VersionPoint:      uint32(app.Version.Point),
			Features:          features,
			ServerLocalDevice: hostDevice,
		},
		StringTables:     loadStrings(ctx),
		EnableLocalFiles: *enableLocalFiles,
		PreloadDepGraph:  *preloadDepGraph,
		AuthToken:        auth.Token(*gapisAuthToken),
		DeviceScanDone:   deviceScanDone,
		LogBroadcaster:   logBroadcaster,
		IdleTimeout:      *idleTimeout,
	})
}

func monitorAndroidDevices(ctx context.Context, r *bind.Registry, scanDone func()) {
	// Populate the registry with all the existing devices.
	func() {
		defer scanDone() // Signal that we have a primed registry.

		if devs, err := adb.Devices(ctx); err == nil {
			for _, d := range devs {
				r.AddDevice(ctx, d)
				r.SetDeviceProperty(ctx, d, client.LaunchArgsKey, text.SplitArgs(*gapirArgStr))
			}
		}
	}()

	if err := adb.Monitor(ctx, r, time.Second*3); err != nil {
		log.W(ctx, "Could not scan for local Android devices. Error: %v", err)
	}
}

func monitorFuchsiaDevices(ctx context.Context, r *bind.Registry, scanDone func()) {
	// Populate the registry with all the existing devices.
	func() {
		defer scanDone() // Signal that we have a primed registry.

		if devs, err := ffx.Devices(ctx); err == nil {
			for _, d := range devs {
				r.AddDevice(ctx, d)
			}
		}
	}()

	if err := ffx.Monitor(ctx, r, time.Second*3); err != nil {
		log.W(ctx, "Could not scan for local Fuchsia devices. Error: %v", err)
	}
}

func monitorRemoteSSHDevices(ctx context.Context, r *bind.Registry, scanDone func()) {
	getRemoteSSHConfig := func() ([]io.ReadCloser, error) {
		f, err := os.Open(*remoteSSHConfig)
		if err != nil {
			return nil, err
		}
		return []io.ReadCloser{f}, nil
	}

	func() {
		// Populate the registry with all the existing devices.
		defer scanDone() // Signal that we have a primed registry.

		f, err := getRemoteSSHConfig()
		if err != nil {
			log.E(ctx, "Could not open remote ssh config")
			return
		}

		if devs, err := remotessh.Devices(ctx, f); err == nil {
			for _, d := range devs {
				r.AddDevice(ctx, d)
				r.SetDeviceProperty(ctx, d, client.LaunchArgsKey, text.SplitArgs(*gapirArgStr))
			}
		}
	}()

	if err := remotessh.Monitor(ctx, r, time.Second*15, getRemoteSSHConfig); err != nil {
		log.W(ctx, "Could not scan for remote SSH devices. Error: %v", err)
	}
}

func loadStrings(ctx context.Context) []*stringtable.StringTable {
	files, err := filepath.Glob(filepath.Join(*stringsPath, "*.stb"))
	if err != nil {
		log.E(ctx, "Couldn't scan for stringtables. Error: %v", err)
		return nil
	}

	out := make([]*stringtable.StringTable, 0, len(files))

	for _, path := range files {
		ctx := log.V{"path": path}.Bind(ctx)
		st, err := stringtable.Load(path)
		if err != nil {
			log.E(ctx, "Couldn't load stringtable file. Error: %v", err)
			continue
		}
		out = append(out, st)
	}

	return out
}
