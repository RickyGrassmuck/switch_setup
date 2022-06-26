/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"github.com/RickyGrassmuck/switch_setup/utils"
	"github.com/google/go-github/v45/github"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"path/filepath"
)

type Package struct {
	Owner    string
	Repo     string
	AssetExt string
}

var downloadDirectory string

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download latest releases of neccessary CFW components",
	Long:  "Download latest releases of neccessary CFW components",
	Run:   downloadFunc,
}

func downloadFunc(cmd *cobra.Command, args []string) {
	var packages = map[string]Package{
		"hb-menu": {
			Owner:    "switchbrew",
			Repo:     "nx-hbmenu",
			AssetExt: ".zip",
		},
		"tesla-menu": {
			Owner:    "WerWolv",
			Repo:     "Tesla-Menu",
			AssetExt: ".zip",
		},
		"ldn_mitm": {
			Owner:    "spacemeowx2",
			Repo:     "ldn_mitm",
			AssetExt: ".zip",
		},
		"sys-botbase": {
			Owner:    "Koi-3088",
			Repo:     "sys-usb-botbase",
			AssetExt: ".zip",
		},

		"atmosphere": {
			Owner:    "Atmosphere-NX",
			Repo:     "Atmosphere",
			AssetExt: ".zip",
		},
	}
	fmt.Printf("Creating Directory: %s\n", downloadDirectory)
	utils.CreatePathIfNotExist(downloadDirectory)

	for _, p := range packages {
		p.DownloadLatestRelease()
	}
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	defaultDownloadDirectory, err := homedir.Expand("~/.switch_setup")
	if err != nil {
		fmt.Println(err)
		return
	}

	rootCmd.PersistentFlags().StringVar(&downloadDirectory, "download-dir", defaultDownloadDirectory, "")
}

func (p *Package) DownloadLatestRelease() error {
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(context.TODO(), p.Owner, p.Repo)
	if err != nil {
		return err
	}
	for _, asset := range release.Assets {
		if filepath.Ext(asset.GetName()) == p.AssetExt {
			fmt.Printf("Downloading %s\n", asset.GetName())
			err := utils.DownloadFile(asset.GetBrowserDownloadURL(), filepath.Join(downloadDirectory, asset.GetName()))
			if err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}
