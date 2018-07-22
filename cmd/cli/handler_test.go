/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package cli

import (
	"testing"

	"github.com/ory/hydra/config"
)

//var db *sqlx.DB
//
//func TestMain(m *testing.M) {
//	runner := dockertest.Register()
//
//	flag.Parse()
//	if !testing.Short() {
//		dockertest.Parallel([]func(){
//			func() {
//				var err error
//				db, err = dockertest.ConnectToTestPostgreSQL()
//				if err != nil {
//					log.Fatalf("Unable to connect to database: %s", err)
//				}
//			},
//		})
//	}
//
//	runner.Exit(m.Run())
//}

func TestNewHandler(t *testing.T) {
	_ = NewHandler(&config.Config{})
}
