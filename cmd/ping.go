package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/cobra"
)

type RedisConnection struct {
	MaxConnection int
	MaxParallel   int
	Port          int
	Duration      int
	Host          string
	Password      string
}

var rdconnection = RedisConnection{}

var getCmd = &cobra.Command{
	Use:   "ping",
	Short: "open connection to redis",
	RunE: func(cmd *cobra.Command, args []string) error {
		return chaos(rdconnection)

	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.PersistentFlags().IntVarP(&rdconnection.MaxConnection, "max-connection", "c", 0, "Max Active connection, by default unlimited")
	getCmd.PersistentFlags().IntVarP(&rdconnection.MaxParallel, "max-parallel", "w", 10, "Max Go routine in parallel")
	getCmd.PersistentFlags().IntVarP(&rdconnection.Port, "port", "P", 6379, "redis port")
	getCmd.PersistentFlags().IntVarP(&rdconnection.Duration, "duration", "d", 300, "chaos duration")
	getCmd.PersistentFlags().StringVar(&rdconnection.Host, "host", "localhost", "redis host")
	getCmd.PersistentFlags().StringVar(&rdconnection.Password, "password", "", "redis password")
}

var pool *redis.Pool

func initPool(r RedisConnection) {
	adresse := fmt.Sprintf("%s:%v", r.Host, r.Port)
	pool = &redis.Pool{
		MaxIdle:         80,
		MaxActive:       r.MaxConnection,
		MaxConnLifetime: 0,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", adresse)
			if err != nil {
				log.Printf("ERROR: fail init redis pool: %s", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}
}

func chaos(r RedisConnection) error {
	initPool(r)
	i := 1
	for i <= r.MaxParallel {
		go getPool()
		i += 1

	}
	time.Sleep(1 * time.Second)
	go countConexion()
	time.Sleep(time.Duration(r.Duration) * time.Second)
	return nil
}

// Print every 30s number of active connection
func countConexion() {
	for {
		s := pool.Stats()
		fmt.Printf("active connection %v \n",s.ActiveCount)
		fmt.Printf("idle connection %v \n", s.IdleCount)
		fmt.Printf("wait connection %v \n", s.WaitCount)
		time.Sleep(30 * time.Second)
	}
}

// Open a new connection
func getPool() error {
	for {
		pool.Get()
	}
}
