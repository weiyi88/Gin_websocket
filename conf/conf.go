package conf

import (
	"chat/model"
	"context"
	"fmt"
	logging "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/ini.v1"
	"strings"
)

var (
	MongoDBClinet *mongo.Client
	AppMode       string
	HttpPort      string

	Db     string
	DbHost string
	DbPort string

	DbUser      string
	DbPassWord  string
	DbName      string
	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string

	MongoDBName string
	MongoDBAddr string
	MongoDBPwd  string
	MongoDBPort string
)

// 本地环境读取
func Init() {
	// 读取配置信息
	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		fmt.Println("ini load failed", err)
	}
	// 加载配置信息
	LoadServer(file)
	LoadMysql(file)
	LoadMongoDb(file)
	// 数据库链接
	MongoDB()

	// mysql 链接
	path := strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8&parseTime=true"}, "")
	model.Database(path)
}

func LoadMysql(file *ini.File) {
	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
}

func LoadMongoDb(file *ini.File) {
	MongoDBName = file.Section("mongoDB").Key("MongoDBPwd").String()
	MongoDBAddr = file.Section("mongoDB").Key("MongoDBPwd").String()
	MongoDBPwd = file.Section("mongoDB").Key("MongoDBPwd").String()
	MongoDBPort = file.Section("mongoDB").Key("MongoDBPwd").String()
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}

func MongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://" + MongoDBAddr + ":" + MongoDBPort)
	var err error
	MongoDBClinet, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logging.Info(err)
		panic(err)
	}
	logging.Info("MongoDB Connnect Successfully")
}
