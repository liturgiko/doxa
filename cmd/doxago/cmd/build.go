// Copyright © 2019 The Orthodox Christian Mission Center (ocmc.org)
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
// TODO: Currently this command calls lt.build, which was a go template based solution for liturgical templates.  Replace it with site.build.
package cmd

import (
	"fmt"
	"github.com/liturgiko/doxa/pkg/lt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)


var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build liturgical website",
	Long: `build a liturgical website from templates based on settings in the config file`,
	Run: func(cmd *cobra.Command, args []string) {

		start := time.Now()

		// setRecord popPath the logger, which will be passed to the functions that do the processing
		LogFile, err := os.OpenFile(LogFilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		defer LogFile.Close()

		Logger.SetOutput(LogFile)
		Logger.SetFlags(log.Ldate + log.Ltime + log.Lshortfile)

		// get the flags
		domains := viper.GetStringSlice("generate.domains")
		patterns := viper.GetStringSlice("generate.template.filename.patterns")
		extension := viper.GetString("generate.template.extension")
//		types := viper.GetStringSlice("generate.output.types")

		msg := fmt.Sprintf("building...\n")
		fmt.Println(msg)
		Logger.Println(msg)
		if err = 	lt.Build(Paths.TemplatesPath,
			Paths.DbPath,
			Paths.SitePath,
			patterns,
			extension,
			domains,
			); err != nil {
			Logger.Println(err.Error())
		}
		Elapsed(start)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

