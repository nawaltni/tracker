package cmd

import (
	"fmt"
	"log"

	"github.com/nawaltni/tracker/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "places",
	Short: "Places is the service that manages places data",
	Run:   run,
}

func RootCommand() *cobra.Command {
	return rootCmd
}

func run(cmd *cobra.Command, args []string) {
	// 1. Read Config
	conf, err := config.LoadConfig(cmd)
	if err != nil {
		log.Fatal("Failed to load config: " + err.Error())
	}

	fmt.Println(conf)

	// // 2. Create logger instance
	// logger, err := config.NewZeroLog("debug", conf.Environment)
	// if err != nil {
	// 	logger.FatalGCP(err, "Could not config logger")
	// }

	// // 3. Create BigQuery Client
	// bqClient, err := bigquery.NewClient(conf.BigQuery, logger)
	// if err != nil {
	// 	logger.FatalGCP(err, "Could not establish connection to BigQuery after 10 tries.")
	// }

	// // 3. Create CBS Client
	// cbsClient, err := cbs.NewCBSClient(conf.CBS)
	// if err != nil {
	// 	logger.FatalGCP(err, "Could not establish connection to cbs.")
	// }

	// recaptchaClient, err := recaptcha.NewClient(conf.Recaptcha2.Secret, *logger)
	// if err != nil {
	// 	logger.FatalGCP(err, "Could not connect to reCAPTCHA v2")
	// 	panic(err)
	// }

	// recaptchaV3Client, err := recaptcha.NewClient(conf.Recaptcha3.Secret, *logger)
	// if err != nil {
	// 	logger.FatalGCP(err, "Could not connect to reCAPTCHA v3")
	// 	panic(err)
	// }

	// // 4. Setup encoder

	// encoder, err := urldecoder.NewCryptoAES(conf.URLEncoding.Key, conf.URLEncoding.Num)
	// if err != nil {
	// 	logger.FatalGCP(err, "Could not start AES encoder")
	// 	return
	// }

	// // 5. Setup URL Generator
	// urlEncoder := urldecoder.NewURLBuilder(conf.URLEncoding.Host, conf.URLEncoding.Scheme, encoder)

	// // 6.
	// ipRegistryClient := ipregistry.NewIPRegistryHandler(conf.IPRegistry)

	// // 7
	// redisCacher := redis.NewCache(*conf)

	// // 8 Create router evaluator

	// rEvaluator := rules.NewRuleRoutingManager(conf.ExpiredClickOffset)

	// // 9. Create Kafka Client

	// var opts []kafka.ClientOptionFunc
	// if conf.Kafka.Auth {
	// 	opts = append(opts, kafka.SASLPlain(conf.Kafka.Key, conf.Kafka.Secret))
	// }

	// if conf.Kafka.SSL {
	// 	opts = append(opts, kafka.DefaultSSL())
	// }

	// kafkaClient, err := kafka.NewClient(conf.Kafka.Address, opts...)
	// if err != nil {
	// 	logger.FatalGCP(err, "Could not start kafka client")
	// 	return
	// }

	// // 10. Create backend requests helper
	// backendClient := httpclient.NewBackendClient(&http.Client{}, conf, logger)

	// // 11. Prepare Services, injecting dependencies.
	// svcs := services.NewServices(
	// 	logger, *conf, bqClient, recaptchaClient, recaptchaV3Client, cbsClient,
	// 	urlEncoder, ipRegistryClient, redisCacher, rEvaluator, kafkaClient, backendClient,
	// )

	// // 12. Create a API Service
	// apiSVC := api.New(logger, conf, svcs)

	// // 14. Prepare gRPC Server
	// srvGRPC, err := grpco.New(*conf, svcs, *logger)
	// if err != nil {
	// 	logger.FatalGCP(err, "could not create grpc server")
	// 	return
	// }

	// go func() {
	// 	// 15. Start gRPC Service
	// 	err = srvGRPC.Start()
	// 	if err != nil {
	// 		logger.FatalGCP(err, "could not start grpc server")
	// 		return
	// 	}
	// }()

	// // 13. Execute services
	// apiSVC.Start()
}
