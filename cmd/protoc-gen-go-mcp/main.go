// Copyright 2025 Redpanda Data, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"flag"
	"fmt"

	"github.com/redpanda-data/protoc-gen-go-mcp/pkg/generator"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	var flagSet flag.FlagSet
	packageSuffix := flagSet.String(
		"package_suffix",
		"mcp",
		"Generate files into a sub-package of the package containing the base .pb.go files using the given suffix. An empty suffix denotes to generate into the same package as the base pb.go files.",
	)
	debug := flagSet.Bool(
		"debug",
		false,
		"Enable debug logging to stderr",
	)
	maxRecursionDepth := flagSet.Int(
		"max_recursion_depth",
		3,
		"Maximum depth for recursive message expansion in schema generation (default 3). Use 1 to detect circular references immediately and reduce schema size.",
	)
	generateStandard := flagSet.Bool(
		"generate_standard",
		true,
		"Generate standard MCP tool variants (default true)",
	)
	generateOpenAI := flagSet.Bool(
		"generate_openai",
		true,
		"Generate OpenAI-compatible tool variants (default true)",
	)

	protogen.Options{
		ParamFunc: flagSet.Set,
	}.Run(func(gen *protogen.Plugin) error {
		if !*generateStandard && !*generateOpenAI {
			return fmt.Errorf("at least one of --generate_standard or --generate_openai must be true")
		}
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generator.NewFileGenerator(f, gen, *debug, *maxRecursionDepth, *generateStandard, *generateOpenAI).Generate(*packageSuffix)
		}
		return nil
	})
}
