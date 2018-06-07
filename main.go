/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 17-11-2017
 * |
 * | File Name:     main.go
 * +===============================================
 */

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/aiotrc/pm/plog"
	"github.com/aiotrc/pm/project"
	"github.com/aiotrc/pm/runner"
	"github.com/aiotrc/pm/thing"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/core/options"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// Config represents main configuration
var Config = struct {
	DB struct {
		URL string `default:"mongodb://172.18.0.1:27017" env:"db_url"`
	}
}{}

// ISRC database
var isrcDB *mgo.Database

// prometheus monitoring
var (
	numberOfCreatedProjects = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "pm",
		Name:      "project_created_total",
		Help:      "Number of created projects.",
	})

	requestsDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "pm",
			Name:      "request_duration_seconds",
			Help:      "A histogram of latencies for requests.",
			Buckets:   []float64{.25, .5, 1, 2.5, 5, 10},
		},
		[]string{"path", "method", "code"},
	)
)

// init initiates global variables
func init() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = log.WithFields(log.Fields{
		"component": "pm",
	}).Writer()

	// Load configuration
	if err := configor.Load(&Config, "config.yml"); err != nil {
		panic(err)
	}

	// Prometheus
	prometheus.MustRegister(numberOfCreatedProjects, requestsDuration)
}

// handle registers apis and create http handler
func handle() http.Handler {
	r := gin.Default()

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 Not Found"})
	})

	r.Use(gin.ErrorLogger(), func(c *gin.Context) {
		gin.WrapH(promhttp.InstrumentHandlerDuration(
			requestsDuration.MustCurryWith(prometheus.Labels{"path": c.Request.URL.Path}),
			http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { c.Next() }),
		))(c)
	})

	api := r.Group("/api")
	{
		api.GET("/about", aboutHandler)

		api.POST("/project", projectNewHandler)
		api.GET("/project", projectListHandler)
		api.DELETE("/project/:name", projectRemoveHandler)
		api.POST("/project/:project/things", thingAddHandler)
		api.GET("/project/:project/logs", projectLogHandler)
		api.GET("/project/:project", projectDetailHandler)
		api.GET("/project/:project/activate", projectActivateHandler)
		api.GET("/project/:project/deactivate", projectDeactivateHandler)

		api.GET("/things/:name", thingGetHandler)
		api.GET("/things/:name/activate", thingActivateHandler)
		api.GET("/things/:name/deactivate", thingDeactivateHandler)
		api.DELETE("/things/:name", thingRemoveHandler)
		api.GET("/things", thingListHandler)
	}
	r.Any("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}

func setupDB() {
	// Create a Mongo Session
	client, err := mgo.Connect(context.Background(), Config.DB.URL, nil)
	if err != nil {
		log.Fatalf("Mongo session %s: %v", Config.DB.URL, err)
	}
	isrcDB = client.Database("isrc")
	// TODO removes logs on sepcific period (Time field)

	if _, err := isrcDB.Collection("pm").Indexes().CreateMany(
		context.Background(),
		mgo.IndexModel{
			Keys: bson.NewDocument(
				bson.EC.Int32("name", 1),
			),
			Options: bson.NewDocument(
				bson.EC.Boolean("unique", true),
			),
		},
		mgo.IndexModel{
			Keys: bson.NewDocument(
				bson.EC.Int32("things.id", 1),
			),
			/*Options: bson.NewDocument(
				bson.EC.Boolean("unique", true),
			),*/
		},
	); err != nil {
		log.Fatalf("Create index %v", err)
	}
}

func main() {
	setupDB()

	fmt.Println("PM AIoTRC @ 2018")

	r := handle()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		fmt.Printf("PM Listen: %s\n", srv.Addr)
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("Listen Error:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("PM Shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown Error:", err)
	}
}

func aboutHandler(c *gin.Context) {
	c.String(http.StatusOK, "18.20 is leaving us")
}

func projectNewHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var json projectReq
	if err := c.BindJSON(&json); err != nil {
		return
	}

	name := json.Name

	p, err := project.New(name, []runner.Env{
		{Name: "MONGO_URL", Value: Config.DB.URL},
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if _, err := isrcDB.Collection("pm").InsertOne(context.Background(), p); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	numberOfCreatedProjects.Inc()

	c.JSON(http.StatusOK, p)
}

func projectDetailHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	name := c.Param("project")

	var p project.Project

	dr := isrcDB.Collection("pm").FindOne(context.Background(), bson.NewDocument(
		bson.EC.String("name", name),
	))

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Project %s not found", name))
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.JSON(http.StatusOK, p)
}

func projectRemoveHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	name := c.Param("name")

	var p project.Project

	dr := isrcDB.Collection("pm").FindOne(context.Background(), bson.NewDocument(
		bson.EC.String("name", name),
	))

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Project %s not found", name))
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	if err := p.Runner.Remove(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if _, err := isrcDB.Collection("pm").DeleteOne(context.Background(), bson.NewDocument(
		bson.EC.String("name", name),
	)); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, p)
}

func projectListHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	ps := make([]project.Project, 0)

	cur, err := isrcDB.Collection("pm").Find(context.Background(), bson.NewDocument())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	for cur.Next(context.Background()) {
		var p project.Project

		if err := cur.Decode(&p); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ps = append(ps, p)
	}
	if err := cur.Close(context.Background()); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, ps)
}

func projectActivateHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	name := c.Param("project")

	dr := isrcDB.Collection("pm").FindOneAndUpdate(context.Background(), bson.NewDocument(
		bson.EC.String("name", name),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$set", bson.EC.Boolean("status", true)),
	), mgo.Opt.ReturnDocument(options.After))

	var p project.Project

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Thing %s not found", name))
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.JSON(http.StatusOK, p)
}

func projectDeactivateHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	name := c.Param("project")

	dr := isrcDB.Collection("pm").FindOneAndUpdate(context.Background(), bson.NewDocument(
		bson.EC.String("name", name),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$set", bson.EC.Boolean("status", false)),
	), mgo.Opt.ReturnDocument(options.After))

	var p project.Project

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Project %s not found", name))
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.JSON(http.StatusOK, p)
}

func thingAddHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	name := c.Param("project")

	var json thingReq
	if err := c.BindJSON(&json); err != nil {
		return
	}

	t := thing.Thing{
		ID:     json.Name,
		Status: true,
	}

	dr := isrcDB.Collection("pm").FindOneAndUpdate(context.Background(), bson.NewDocument(
		bson.EC.String("name", name),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$addToSet", bson.EC.Interface("things", t)),
	), mgo.Opt.ReturnDocument(options.After))

	var p project.Project

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Project %s not found", name))
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.JSON(http.StatusOK, p)
}

func thingGetHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	name := c.Param("name")

	var p project.Project

	dr := isrcDB.Collection("pm").FindOne(context.Background(), bson.NewDocument(
		bson.EC.Boolean("status", true),
		bson.EC.SubDocumentFromElements("things", bson.EC.SubDocumentFromElements("$elemMatch",
			bson.EC.String("id", name),
		)),
	))

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Thing %s not found", name))
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.JSON(http.StatusOK, p)
}

func thingActivateHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	name := c.Param("name")

	dr := isrcDB.Collection("pm").FindOneAndUpdate(context.Background(), bson.NewDocument(
		bson.EC.SubDocumentFromElements("things", bson.EC.SubDocumentFromElements(
			"$elemMatch", bson.EC.String("id", name),
		)),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$set", bson.EC.Boolean("things.$.status", true)),
	), mgo.Opt.ReturnDocument(options.After))

	var p project.Project

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Thing %s not found", name))
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.JSON(http.StatusOK, p)
}

func thingDeactivateHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	name := c.Param("name")

	dr := isrcDB.Collection("pm").FindOneAndUpdate(context.Background(), bson.NewDocument(
		bson.EC.SubDocumentFromElements("things", bson.EC.SubDocumentFromElements(
			"$elemMatch", bson.EC.String("id", name),
		)),
	), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$set", bson.EC.Boolean("things.$.status", false)),
	), mgo.Opt.ReturnDocument(options.After))

	var p project.Project

	if err := dr.Decode(&p); err != nil {
		if err == mgo.ErrNoDocuments {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("Thing %s not found", name))
		} else {
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	c.JSON(http.StatusOK, p)
}

func projectLogHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var pls = make([]plog.ProjectLog, 0)

	id := c.Param("project")

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cur, err := isrcDB.Collection("errors").Aggregate(context.Background(), bson.NewArray(
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$match", bson.EC.String("project", id)),
		),
		bson.VC.DocumentFromElements(
			bson.EC.Int32("$limit", int32(limit)),
		),
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$sort", bson.EC.Int32("Time", -1)),
		),
	))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	for cur.Next(context.Background()) {
		var pl plog.ProjectLog

		if err := cur.Decode(&pl); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		pls = append(pls, pl)
	}
	if err := cur.Close(context.Background()); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, pls)
}

func thingRemoveHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	name := c.Param("name")

	dr := isrcDB.Collection("pm").FindOneAndUpdate(context.Background(), bson.NewDocument(), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$pull", bson.EC.SubDocumentFromElements(
			"things", bson.EC.String("id", name)),
		),
	), mgo.Opt.ReturnDocument(options.After))

	var p project.Project

	if err := dr.Decode(&p); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, p)
}

func thingListHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	results := make([]thing.Thing, 0)

	cur, err := isrcDB.Collection("pm").Aggregate(context.Background(), bson.NewArray(
		bson.VC.DocumentFromElements(
			bson.EC.String("$unwind", "$things"),
		),
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$replaceRoot", bson.EC.String("newRoot", "$things")),
		),
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements("$sort", bson.EC.Int32("Time", -1)),
		),
	))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for cur.Next(context.Background()) {
		var result thing.Thing

		if err := cur.Decode(&result); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		results = append(results, result)
	}
	if err := cur.Close(context.Background()); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, results)

}
