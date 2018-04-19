package main

type mysqlConf struct {
	Host     string
	Port     string
	User     string
	Pswd     string
	Charset  string
	Database string
}

type redisConf struct {
	Host           string
	Port           string
	Pswd           string
	Database       string
	Queue          string
	HaoServiceCity string
	Timeout        int
}

type logConf struct {
	File   string
	SqlLog bool
	Debug  bool
}

type antConf struct {
	AntNestUrl string
	AppKey     string
	SecretKey  string
}

type config struct {
	Mysql mysqlConf
	Redis redisConf
	Log   logConf
	Ant   antConf
}

type AuArticle struct {
	ArticleId      int64 `xorm:"article_id pk not null autoincr"`
	CategoryId     int64
	ChildId        int64
	ArticleTitle   string
	ArticleFocus   string
	ArticleThumb   string
	ArticleDesc    string
	ArticleContent string `xorm:"text"`
	Author         string
	IsTop          int64
	Click          int64
	Share          int64
	Status         int64
	CreateTime     int64
	UpdateTime     int64
	PublishTime    int64
	Type           int64
	FromName       string `xorm:"from_name varchar(10)"`
	FromUrl        string
}

