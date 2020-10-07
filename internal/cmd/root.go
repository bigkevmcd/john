package cmd

import (
	"fmt"
	"log"

	"github.com/bigkevmcd/john/pkg/handler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mhale/smtpd"
)

const (
	portFlag = "port"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.AutomaticEnv()
}

func makeHTTPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "john",
		Short: "Simple SMTP server with Mailets",
		Run: func(cmd *cobra.Command, args []string) {
			listen := fmt.Sprintf(":%d", viper.GetInt(portFlag))
			log.Printf("listening on %s", listen)
			logIfError(smtpd.ListenAndServe(listen, handler.MailetsHandler, "John SMTP", ""))
		},
	}

	cmd.Flags().Int(
		portFlag,
		2525,
		"port to receive mail on",
	)
	logIfError(viper.BindPFlag(portFlag, cmd.Flags().Lookup(portFlag)))
	return cmd
}

// Execute is the main entry point into this component.
func Execute() {
	if err := makeHTTPCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}

func logIfError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}