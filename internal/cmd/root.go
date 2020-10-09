package cmd

import (
	"fmt"
	"log"

	"github.com/mhale/smtpd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bigkevmcd/john/pkg/handler"
	"github.com/bigkevmcd/john/pkg/mailet/maildir"
)

const (
	portFlag    = "port"
	maildirFlag = "maildir"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.AutomaticEnv()
}

func makeRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "john",
		Short: "Simple Mailet processing service.",
	}
	cmd.AddCommand(makeSMTPServerCmd())
	return cmd
}

func makeSMTPServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "smtp",
		Short: "Act as an SMTP server and process incoming emails",
		Run: func(cmd *cobra.Command, args []string) {
			listen := fmt.Sprintf(":%d", viper.GetInt(portFlag))
			log.Printf("listening on %s", listen)
			h, err := maildir.New(viper.GetString(maildirFlag))
			logIfError(err)
			logIfError(smtpd.ListenAndServe(listen, handler.MakeHandler(h), "John SMTP", ""))
		},
	}

	cmd.Flags().Int(
		portFlag,
		2525,
		"port to receive mail on",
	)
	logIfError(viper.BindPFlag(portFlag, cmd.Flags().Lookup(portFlag)))

	cmd.Flags().String(
		maildirFlag,
		"./tmp",
		"Path to store Maildir mail",
	)
	logIfError(viper.BindPFlag(maildirFlag, cmd.Flags().Lookup(maildirFlag)))
	return cmd
}

// Execute is the main entry point into this component.
func Execute() {
	logIfError(makeRootCmd().Execute())
}

func logIfError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
