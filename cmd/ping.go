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
	Connexion int
	Port      int
	Duration  int
	Host      string
	Password  string
}

var rdconnection = RedisConnection{}

var getCmd = &cobra.Command{
	Use:   "ping",
	Short: "ping pong redis",
	RunE: func(cmd *cobra.Command, args []string) error {
		return chaos(rdconnection)

	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.PersistentFlags().IntVarP(&rdconnection.Connexion, "connection", "c", 1, "number of simulated connections")
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
		MaxActive:       r.Connexion,
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
	for i <= r.Connexion {
		go ping()
		//fmt.Printf("number in for %v \n", i)
		//time.Sleep(time.Duration(r.Duration) * time.Second)
		i += 1

	}
	time.Sleep(1 * time.Second)
	fmt.Printf("active connection %v \n", pool.ActiveCount())
	time.Sleep(time.Duration(r.Duration) * time.Second)
	return nil
}

func ping() error {
	conn := pool.Get()
	//defer conn.Close()
	//fmt.Printf("active connection %v \n", pool.ActiveCount())
	for {
		s, err := redis.String(conn.Do("PING"))
		fmt.Printf("%v", s)
		if err != nil {
			log.Printf("ERROR: fail ping redis, error %s", err.Error())
			return err
		}
	}

}
