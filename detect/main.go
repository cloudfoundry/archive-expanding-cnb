/*
 * Copyright 2018-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/archive-expanding-cnb/expand"
	"github.com/cloudfoundry/jvm-application-cnb/jvmapplication"
	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/libcfbuildpack/helper"
)

var archiveTypes = []*regexp.Regexp{
	regexp.MustCompile(".*\\.jar$"),
	regexp.MustCompile(".*\\.war$"),
	regexp.MustCompile(".*\\.tar$"),
	regexp.MustCompile(".*\\.tar\\.gz$"),
	regexp.MustCompile(".*\\.tgz$"),
	regexp.MustCompile(".*\\.zip$"),
}

func main() {
	detect, err := detect.DefaultDetect()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialize Detect: %s\n", err)
		os.Exit(101)
	}

	if code, err := d(detect); err != nil {
		detect.Logger.Info(err.Error())
		os.Exit(code)
	} else {
		os.Exit(code)
	}
}

func d(detect detect.Detect) (int, error) {
	var c []string

	for _, r := range archiveTypes {
		if f, err := helper.FindFiles(detect.Application.Root, r); err != nil {
			return -1, err
		} else {
			c = append(c, f...)
		}
	}

	if len(c) != 1 {
		return detect.Fail(), nil
	}

	bp := detect.BuildPlan[expand.Dependency]
	if bp.Metadata == nil {
		bp.Metadata = make(buildplan.Metadata)
	}
	bp.Metadata[expand.Archive] = c[0]

	return detect.Pass(buildplan.BuildPlan{
		expand.Dependency:         bp,
		jvmapplication.Dependency: detect.BuildPlan[jvmapplication.Dependency],
	})
}
