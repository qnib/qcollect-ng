// Copyright © 2017 Christian Kniep <christian@qnib.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
"os"

"github.com/codegangsta/cli"
"github.com/qnib/qframe/server"
)

func main() {
	app := cli.NewApp()
	app.Name = "ETC metrics collector based on qframe, inspired by qcollect,logstash and fullerite"
	app.Usage = "qcollect-ng [options]"
	app.Version = "0.3.2.4"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "qcollect-ng.yml",
			Usage: "Config file, will overwrite flag default if present.",
		},
		cli.StringFlag{
			Name:  "ld-path",
			Value: "lib/",
			Usage: "Library base path for golang plugins",
		},
	}
	app.Action = qserver.Run
	app.Run(os.Args)
}

