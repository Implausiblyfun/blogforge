package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/braintree/manners"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v2"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	"github.com/implausiblyfun/blogforge/internal/apid"
	"github.com/implausiblyfun/blogforge/internal/auth"
	"github.com/implausiblyfun/blogforge/internal/buildartifacts"
	"github.com/implausiblyfun/blogforge/internal/database"
	"github.com/implausiblyfun/blogforge/internal/encrypt"
	"github.com/implausiblyfun/blogforge/internal/standardroutes"
	"github.com/implausiblyfun/blogforge/internal/tracing"
)

// Extract some interesting terms for easier reuse.
// While it does pull them away from the code using them the name is fairly silly may crop up a few times...
// Therefore lets place it and a few other similar things here.
const (
	appName = "blogforge-apid"
)

func main() {
	// lets pretend someone implements flags and doesnt put something heavier weight like cobra on this.
	flag.Parse()

	// Lets do a crazy getter because of bad docker compose and soem shenanagins on my local windows.
	// Side note: windows and docker works but... things break in silly ways
	os.Chdir(os.Getenv("BASEPATH"))

	buildartifacts.LoadBuildInfo()

	// using fmt for now rather than emitting up or using a full fledged logger due to ease of impl.
	fmt.Printf("Starting run of %s version %s\n ", appName, buildartifacts.Build.Version)
	err := run(context.Background())

	exitMsg := fmt.Sprintf("Closing down %s version %s", appName, buildartifacts.Build.Version)
	if err != nil {
		exitMsg += " with error of " + err.Error()
	}
	fmt.Println(exitMsg)
}

// run is a nice encapsulation of how we kick off the main processes.
// This is pulled out from main so if we want to call it from a different location its not entirely annoying.
func run(c context.Context) error {

	// add a cancellation onto our context so that we can correctly clean up if need be
	ctx, cancelCaller := context.WithCancel(c)
	defer cancelCaller()

	// load some configuration options probably things like what ports to expose, databases, admin cred locations that thing
	cfg, err := loadConfigs()
	if err != nil {
		// too lazy to get wrap working here but that would be cleaner
		return fmt.Errorf("failed to load config: %w", err)
	}
	fmt.Printf("Configs loaded where we specified %d fields at the top level\n", len(cfg))

	if secKey, ok := cfg["secret"]; ok {
		fmt.Println("Setting our insecurely stored key")
		encrypt.SetKey(fmt.Sprintf("%v", secKey))

	}

	// Load tracing infos and the like
	if tracez, ok := cfg["tracing"]; ok {
		fmt.Println("Setting tracing settings of", tracez)
		shouldEnable := checkConfigTruth("enabled", tracez)
		fmt.Println("Tracing Enabled:", shouldEnable)
		if shouldEnable {
			tp, err := tracing.NewBaseTracer("http://jeager:14268/api/traces", appName)
			if err != nil {
				return err
			}
			defer tp.ForceFlush(ctx)
		}
	}

	// Any configuration options needed for the router and the like
	if ginEn, ok := cfg["gin"]; ok {
		fmt.Println("Setting Gin settings of", ginEn)

		shouldSet := checkConfigTruth("productionmode", ginEn)
		fmt.Println("Gin Enabled:", shouldSet)
		if shouldSet {
			gin.SetMode(gin.ReleaseMode)
		}
	}

	// attempt to set up the db connections per our config
	connectionMap, err := initDBConnections(cfg["databases"])
	if err != nil {
		return fmt.Errorf("failed to establish database connections: %w", err)
	}
	fmt.Printf("Db Connections Secured Stood up %d connections \n", len(connectionMap))

	// actually start serving
	router := gin.New()
	router.Use(otelgin.Middleware("apid"))
	router.NoRoute(standardroutes.NotFoundRouteHandler())
	router.Use(gin.Recovery())
	router.GET("/ping", standardroutes.EmptyGoodHandler())
	router.GET("/index", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "instructions": fmt.Sprintf("check out the routes %v", router.Routes())})
	})
	router.POST("/newuser", apid.NewUserHandler(connectionMap["basicDB"]))
	// router.GET("/metrics",  standardroutes.EmptyGoodHandler())

	// unprivileged routes

	// logged in routes
	loggedIn := router.Group("/")
	loggedIn.Use(auth.BasicAuther(connectionMap["basicDB"]))
	loggedIn.Any("/fire", apid.WithFire(cancelCaller))
	loggedIn.Any("/brimstone", apid.WithFireAndBrimstone())
	loggedIn.GET("/users/all", apid.ListUsers(connectionMap["basicDB"]))
	loggedIn.GET("/configs", apid.AdminConfig(cfg))

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		errCh := make(chan error, 1)
		var e error
		go func() {
			errCh <- manners.ListenAndServe(":8080", router)
		}()
		select {
		case err := <-errCh:
			e = err
		case <-ctx.Done():
			fmt.Println("Closing Gin Router early and non-gracefully")
		}
		cancelCaller()

		return e
	})
	g.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case <-signalChannel:
			fmt.Println("Detected Shutdown request")
		case <-ctx.Done():
		}
		cancelCaller()
		return nil
	})
	err = g.Wait()
	return err
}

// loadConfigs is a placeholder function with what shall likely become a stricter signature.
// For now it returns a structure that `could` hold our info and an error to describe anything that went wrong.
func loadConfigs() (map[string]interface{}, error) {
	cfgDat := map[string]interface{}{}
	buf, err := ioutil.ReadFile("./artifacts/configs/apid.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buf, cfgDat)
	return cfgDat, err
}

// initDBConnections to see if we can get the app to have basic accesses.
func initDBConnections(dbCfg interface{}) (map[string]*sqlx.DB, error) {

	cfg, ok := dbCfg.([]interface{})
	if !ok {
		return nil, errors.New("No databases defined")
	}

	databases := map[string]*sqlx.DB{}
	for _, dbparams := range cfg {
		dbName, con, err := database.MySqlConnection(dbparams)
		if err != nil {
			return databases, err
		}
		databases[dbName] = con
	}

	return databases, nil
}

func checkConfigTruth(key string, data interface{}) bool {
	datMap, alright := data.(map[interface{}]interface{})
	if !alright {
		return false
	}

	dat, alright := datMap[key]
	if !alright {
		return false
	}
	datStr := fmt.Sprintf("%v", dat)
	if len(datStr) > 0 && strings.HasPrefix(strings.ToLower(datStr), "t") {
		return true
	}

	return false

}
