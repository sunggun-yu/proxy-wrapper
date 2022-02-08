package cmd

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sunggun-yu/envp/internal/config"
)

// flags struct for edit command
type editFlags struct {
	desc string
	env  []string
}

func init() {
	rootCmd.AddCommand(editCommand())
}

// example of edit command
func cmdExampleEdit() string {
	return `
  envp edit my-proxy \
    -d 'updated profile desc' \
    -e 'NO_PROXY=127.0.0.1,localhost'
  `
}

// editCommand edit/update environment variable profile and it's envionment variables in the config file
func editCommand() *cobra.Command {
	var flags editFlags

	cmd := &cobra.Command{
		Use:          "edit profile-name [flags]",
		Short:        "Edit environment variable profile",
		SilenceUsage: true,
		Example:      cmdExampleEdit(),
		Args: cobra.MatchAll(
			Arg0AsProfileName(),
			Arg0NotExistingProfile(),
		),
		ValidArgsFunction: ValidArgsProfileList,
		RunE: func(cmd *cobra.Command, args []string) error {

			profileName := args[0]
			var profile config.Profile

			// validate selected profile
			selected := configProfiles.Sub(profileName)
			// unmarshal into Profile
			err := selected.Unmarshal(&profile)
			if err != nil {
				return fmt.Errorf("profile '%v' malformed configuration %e", profile, err)
			}

			// update desc if input is not empty
			if flags.desc != "" {
				profile.Desc = flags.desc
			}

			// update env
			// parse flag.env into a map for easy checking
			menv := config.ParseEnvFlagToMap(flags.env)
			if menv != nil {
				// loop profile.Env and check if flag.env has updated value(exist)
				for _, e := range profile.Env {
					if _, exist := menv[e.Name]; !exist {
						menv[e.Name] = e.Value
					}
				}
				profile.Env = config.MapToEnv(menv)
			}

			// set updated profile
			configProfiles.Set(profileName, profile)

			// overwrite the profile
			viper.Set(ConfigKeyProfile, configProfiles.AllSettings())

			// wait for the config file update and verify profile is added or not
			rc := make(chan error, 1)

			// it's being watched in root initConfig - viper.WatchConfig()
			viper.OnConfigChange(func(e fsnotify.Event) {
				// assuming
				if configProfiles.Get(profileName) == nil {
					rc <- fmt.Errorf("profile %v not added", profileName)
					return
				}
				fmt.Println("profile", profileName, "updated successfully:", e.Name)
				rc <- nil
			})

			if err := viper.WriteConfig(); err != nil {
				return err
			}
			// wait for profile validation channel
			err = <-rc
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.desc, "desc", "d", "", "description of profile")
	cmd.Flags().StringSliceVarP(&flags.env, "env", "e", []string{}, "'VAR=VAL' format of string")
	cmd.MarkFlagRequired("env")

	return cmd
}
