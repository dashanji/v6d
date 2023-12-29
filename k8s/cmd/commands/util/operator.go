/*
* Copyright 2020-2023 Alibaba Group Holding Limited.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package util

import (
	"log"

	"sigs.k8s.io/kustomize/api/filesys"
	"sigs.k8s.io/kustomize/api/krusty"
)

//"sigs.k8s.io/kustomize/kyaml/filesys"

/*
	"github.com/v6d-io/v6d/k8s/config"
	"sigs.k8s.io/kustomize/k8sdeps"
	"sigs.k8s.io/kustomize/pkg/commands/build"
	kustomizefs "sigs.k8s.io/kustomize/pkg/fs"*/
//"sigs.k8s.io/kustomize/api/filesys"
//"sigs.k8s.io/kustomize/kustomize/v4/commands/build"

func BuildKustomizeInEmbedDir() (Manifests, error) {
	/*
		// Set the directory you want to build.
		kustomizationDirectory := "/opt/caoye/v6d/k8s/config/default/kustomization.yaml"
		// Prepare the arguments for the command.
		args := []string{"build", kustomizationDirectory}

		// Capture the output in a buffer.
		var outBuffer bytes.Buffer

		stdOut := os.Stdout
		// Create a temporary directory to extract the embedded config files
		tmpDir, err := os.MkdirTemp("", "v6d-operator-manifests-")
		if err != nil {
			return nil, err
		}
		defer os.RemoveAll(tmpDir)

		// Extract the embedded config files to the temporary directory
		if err := fs.WalkDir(config.Manifests, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				data, err := config.Manifests.ReadFile(path)
				if err != nil {
					return err
				}
				destPath := filepath.Join(tmpDir, path)
				if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
					return err
				}
				if err := os.WriteFile(destPath, data, os.ModePerm); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return nil, err
		}

		ftory := k8sdeps.NewFactory()
		fsys := kustomizefs.MakeFakeFS()
		fmt.Print("ftory: ", ftory, fsys, stdOut, args, outBuffer)
		// Create the build command.
		cmd := build.NewCmdBuild(stdOut, fsys, ftory.ResmapF, ftory.TransformerF)
		cmd.SetArgs(args)
		cmd.SetOutput(&outBuffer)

		if err := cmd.RunE(cmd, []string{tmpDir + "/default/"}); err != nil {
			return nil, err
		}
		fmt.Println(outBuffer.String())
		//return ParseManifestsToObjects(buffy.Bytes())
		/*cmd.RunE(cmd, []string{tmpDir + "/default"})
		// Run the build command.
		if err := cmd.Execute(); err != nil {
			log.Fatalf("Failed to execute kustomize build command: %v", err)
		}

		// Output the result.
		fmt.Println(outBuffer.String())*/
	/*
		fSys := fs.MakeFakeFS()
		buffy := new(bytes.Buffer)
		cmd := build.NewCmdBuild(fSys, build.MakeHelp("", ""), buffy)

		// Create a temporary directory to extract the embedded config files
		tmpDir, err := os.MkdirTemp("", "v6d-operator-manifests-")

		if err != nil {
			return nil, err
		}

		defer os.RemoveAll(tmpDir)

		// Extract the embedded config files to the temporary directory

		if err := fs.WalkDir(config.Manifests, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				data, err := config.Manifests.ReadFile(path)
				if err != nil {
					return err
				}
				destPath := filepath.Join(tmpDir, path)
				if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
					return err
				}
				if err := os.WriteFile(destPath, data, os.ModePerm); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {

			return nil, err
		}

		if err := cmd.RunE(cmd, []string{tmpDir + "/default"}); err != nil {
			return nil, err
		}

		return ParseManifestsToObjects(buffy.Bytes())
		return nil, nil
	*/
	/*
		fSys := filesys.MakeFsOnDisk()
		buffy := new(bytes.Buffer)
		cmd := build.NewCmdBuild(fSys, build.MakeHelp("", ""), buffy)

		// Create a temporary directory to extract the embedded config files
		tmpDir, err := os.MkdirTemp("", "v6d-operator-manifests-")

		if err != nil {
			return nil, err
		}

		defer os.RemoveAll(tmpDir)

		// Extract the embedded config files to the temporary directory

		if err := fs.WalkDir(config.Manifests, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				data, err := config.Manifests.ReadFile(path)
				if err != nil {
					return err
				}
				destPath := filepath.Join(tmpDir, path)
				if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
					return err
				}
				if err := os.WriteFile(destPath, data, os.ModePerm); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {

			return nil, err
		}

		if err := cmd.RunE(cmd, []string{tmpDir + "/default"}); err != nil {
			return nil, err
		}

		return ParseManifestsToObjects(buffy.Bytes())*/
	// 设置要构建的 kustomization 目录路径。
	kustomizationDir := "/opt/caoye/v6d/k8s/config/default"

	// 创建一个 FileSystem
	fSys := filesys.MakeFsOnDisk()

	// 创建一个 kustomize 构建对象。
	k := krusty.MakeKustomizer(fSys, krusty.MakeDefaultOptions())

	// 运行 kustomize 构建
	resMap, err := k.Run(kustomizationDir)
	if err != nil {
		log.Fatalf("failed to run kustomize build: %v", err)
	}

	// 将结果转换为YAML格式
	yaml, err := resMap.AsYaml()
	if err != nil {
		log.Fatalf("failed to convert resources to YAML: %v", err)
	}

	return ParseManifestsToObjects(yaml)
}
