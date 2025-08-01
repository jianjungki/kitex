// Copyright 2021 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/kitex/tool/cmd/kitex/utils"
	"github.com/cloudwego/kitex/tool/internal_pkg/util/env"

	"github.com/cloudwego/kitex/tool/cmd/kitex/sdk"

	"github.com/cloudwego/kitex"
	kargs "github.com/cloudwego/kitex/tool/cmd/kitex/args"
	"github.com/cloudwego/kitex/tool/cmd/kitex/versions"
	"github.com/cloudwego/kitex/tool/internal_pkg/log"
	"github.com/cloudwego/kitex/tool/internal_pkg/pluginmode/protoc"
	"github.com/cloudwego/kitex/tool/internal_pkg/pluginmode/thriftgo"
	"github.com/cloudwego/kitex/tool/internal_pkg/prutal"
)

var args kargs.Arguments

func init() {
	var queryVersion bool
	args.AddExtraFlag(&kargs.ExtraFlag{
		Apply: func(f *flag.FlagSet) {
			f.BoolVar(&queryVersion, "version", false,
				"Show the version of kitex")
		},
		Check: func(a *kargs.Arguments) error {
			if queryVersion {
				println(a.Version)
				os.Exit(0)
			}
			return nil
		},
	})
	if err := versions.RegisterMinDepVersion(
		&versions.MinDepVersion{
			RefPath: "github.com/cloudwego/kitex",
			Version: "v0.11.0",
		},
	); err != nil {
		log.Error(err)
		os.Exit(versions.CompatibilityCheckExitCode)
	}
}

func main() {
	mode := os.Getenv(kargs.EnvPluginMode)
	if len(os.Args) <= 1 && mode != "" {
		// run as a plugin
		switch mode {
		case thriftgo.PluginName:
			os.Exit(thriftgo.Run())
		case protoc.PluginName:
			os.Exit(protoc.Run())
		}
		panic(mode)
	}

	curpath, err := filepath.Abs(".")
	if err != nil {
		log.Errorf("Get current path failed: %s", err)
		os.Exit(1)
	}
	// run as kitex
	err = args.ParseArgs(kitex.Version, curpath, os.Args[1:])
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		log.Error(err)
		os.Exit(2)
	}
	if !args.NoDependencyCheck {
		// check dependency compatibility between kitex cmd tool and dependency in go.mod
		if err := versions.DefaultCheckDependencyAndProcess(); err != nil {
			os.Exit(versions.CompatibilityCheckExitCode)
		}
	}
	if args.IsProtobuf() {
		// Whether using protoc or prutal, no longer generate the fast api for protobuf
		args.Config.NoFastAPI = true
		if !env.UseProtoc() {
			g := prutal.NewPrutalGen(args.Config)
			if err := g.Process(); err != nil {
				log.Errorf("%s", err)
				os.Exit(1)
			}
			return
		}
	}

	out := new(bytes.Buffer)
	cmd, err := args.BuildCmd(out)
	if err != nil {
		log.Warn(err)
		os.Exit(1)
	}

	if args.IsThrift() && !args.LocalThriftgo {
		if err = sdk.InvokeThriftgoBySDK(curpath, cmd); err != nil {
			// todo: optimize -use and remove error returned from thriftgo
			out.WriteString(err.Error())
		}
	} else {
		err = kargs.ValidateCMD(cmd.Path, args.IDLType)
		if err != nil {
			log.Warn(err)
			os.Exit(1)
		}
		err = cmd.Run()
	}

	if err != nil {
		if args.Use != "" {
			out := strings.TrimSpace(out.String())
			if strings.HasSuffix(out, thriftgo.TheUseOptionMessage) {
				goto NormalExit
			}
		}
		log.Warn(err)
		os.Exit(1)
	}
NormalExit:
	utils.OnKitexToolNormalExit(args)
}
