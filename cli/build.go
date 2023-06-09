package cli

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

func BuildCmd() *cobra.Command {
	var exportFile string

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Use yaml to build a docker image",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Get filename from cli
			// TODO: Read yaml file
			var yamlFile []byte
			var err error
			if args[0] == "-" {
				yamlFile, err = ioutil.ReadAll(os.Stdin)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				yamlFile, err = ioutil.ReadFile(args[0])
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			var v struct {
				Gpu     bool     `yaml:"gpu"`
				Src     []string `yaml:"src"`
				Workdir string   `yaml:"workdir"`
				Python  struct {
					Version       string   `yaml:"version"`
					ProjectType   string   `yaml:"project_type"`
					ProjectSource []string `yaml:"project_source"`
				} `yaml:"python"`
			}

			if err := yaml.Unmarshal(yamlFile, &v); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// --------------------

			// TODO: Construct docker file

			var dockerFile strings.Builder

			if v.Gpu == true {
				dockerFile.WriteString("FROM nvidia-gpu\n")
			} else {
				dockerFile.WriteString("FROM ubuntu:22.04\n")
			}
			dockerFile.WriteString("WORKDIR " + v.Workdir + "\n")
			for _, item := range v.Python.ProjectSource {
				dockerFile.WriteString("COPY " + item + " " + item + "\n")
			}

			// --------------------
			finalizedFile := dockerFile.String()

			switch exportFile {
			case "":
				return
			case "-":
				writer := bufio.NewWriter(os.Stdout)
				_, err = writer.WriteString(finalizedFile)
				if err != nil {
					fmt.Println("Error writing to buffer")
					return
				}
				writer.Flush()
				return
			default:
				f, err := os.Create(exportFile)
				if err != nil {
					fmt.Println(err)
					return
				}
				writer := bufio.NewWriter(f)
				_, err = writer.WriteString(finalizedFile)
				if err != nil {
					fmt.Println("Error writing to buffer")
					return
				}
				writer.Flush()
			}

			// TODO: Build container

			// TODO: Export Dockerfile to file for stdout
		},
	}

	cmd.Flags().StringVar(&exportFile, "export", "", "File to export dockerfile to")
	return cmd
}
