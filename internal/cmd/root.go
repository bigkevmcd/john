package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mhale/smtpd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	gomaildir "github.com/emersion/go-maildir"

	"github.com/bigkevmcd/john/pkg/handler"
	"github.com/bigkevmcd/john/pkg/httsmtp"
	"github.com/bigkevmcd/john/pkg/mailet/maildir"
)

const (
	portFlag        = "port"
	maildirFlag     = "maildir"
	initMaildirFlag = "init-maildir"
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
	cmd.AddCommand(makeHTTPServerCmd())
	return cmd
}

func makeSMTPServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "smtp",
		Short: "Act as an SMTP server and process incoming emails",
		Run: func(cmd *cobra.Command, args []string) {
			listen := fmt.Sprintf(":%d", viper.GetInt(portFlag))
			log.Printf("listening on %s", listen)

			if viper.GetBool(initMaildirFlag) {
				cobra.CheckErr(gomaildir.Dir(viper.GetString(maildirFlag)).Init())
			}

			cobra.CheckErr(smtpd.ListenAndServe(listen,
				handler.MakeHandler(maildir.New(viper.GetString(maildirFlag))), "John SMTP", ""))
		},
	}

	cmd.Flags().Int(
		portFlag,
		2525,
		"port to receive mail on",
	)
	cobra.CheckErr(viper.BindPFlag(portFlag, cmd.Flags().Lookup(portFlag)))

	cmd.Flags().Bool(
		initMaildirFlag,
		false,
		"If true, initialise the directory as a Maildir (create cur,new,tmp)",
	)
	cobra.CheckErr(viper.BindPFlag(initMaildirFlag, cmd.Flags().Lookup(initMaildirFlag)))

	cmd.Flags().String(
		maildirFlag,
		"./tmp",
		"Path to store Maildir mail",
	)
	cobra.CheckErr(viper.BindPFlag(maildirFlag, cmd.Flags().Lookup(maildirFlag)))

	return cmd
}

func makeHTTPServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "http",
		Short: "Accepts SMTP over HTTP - for testing purposes",
		Run: func(cmd *cobra.Command, args []string) {
			listen := fmt.Sprintf(":%d", viper.GetInt(portFlag))
			log.Printf("listening on %s", listen)

			if viper.GetBool(initMaildirFlag) {
				cobra.CheckErr(gomaildir.Dir(viper.GetString(maildirFlag)).Init())
			}

			cobra.CheckErr(http.ListenAndServe(listen, httsmtp.MakeHandler(maildir.New(viper.GetString(maildirFlag)))))
		},
	}

	cmd.Flags().Int(
		portFlag,
		8080,
		"port to serve HTTP",
	)
	cobra.CheckErr(viper.BindPFlag(portFlag, cmd.Flags().Lookup(portFlag)))

	cmd.Flags().Bool(
		initMaildirFlag,
		false,
		"If true, initialise the directory as a Maildir (create cur,new,tmp)",
	)
	cobra.CheckErr(viper.BindPFlag(initMaildirFlag, cmd.Flags().Lookup(initMaildirFlag)))

	cmd.Flags().String(
		maildirFlag,
		"./tmp",
		"Path to store Maildir mail",
	)
	cobra.CheckErr(viper.BindPFlag(maildirFlag, cmd.Flags().Lookup(maildirFlag)))

	return cmd
}

// Execute is the main entry point into this component.
func Execute() {
	cobra.CheckErr(makeRootCmd().Execute())
}
