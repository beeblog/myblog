package models

import (
	// "github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	// 设置数据库路径
	_DB_NAME = "data/beeblog.db"
	// 设置数据库名称
	_SQLITE3_DRIVER = "sqlite3"
)

// 分类
type Category struct {
	Id              int64
	Title           string
	Created         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	TopicTime       time.Time `orm:"index"`
	TopicCount      int64
	TopicLastUserId int64
}

// 文章
type Topic struct {
	Id              int64
	Uid             int64
	Title           string
	Category        string
	Content         string `orm:"size(5000)"`
	Attachment      string
	Created         time.Time `orm:"index"`
	Updated         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	Author          string
	ReplyCount      int64
	ReplyLastUserId int64
}

type Reply struct {
	Id      int64
	Tid     int64
	Name    string
	Content string
	Created time.Time
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func RegisterDB() {
	// 检查数据库文件
	if !Exist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}

	// 注册驱动（“sqlite3” 属于默认注册，此处代码可省略）
	orm.RegisterDriver(_SQLITE3_DRIVER, orm.DRSqlite)
	// 注册默认数据库
	orm.RegisterDataBase("default", _SQLITE3_DRIVER, _DB_NAME, 10)
	// 注册模型
	orm.RegisterModel(new(Category), new(Topic), new(Reply))
}

func AddCategory(name string, flag int) error {
	o := orm.NewOrm()
	var cate *Category
	if flag == 0 {
		cate = &Category{
			Title:     name,
			Created:   time.Now(),
			TopicTime: time.Now(),
		}
	} else {
		cate = &Category{
			Title:      name,
			TopicCount: 1,
			Created:    time.Now(),
			TopicTime:  time.Now(),
		}
	}

	// 查询数据
	qs := o.QueryTable("category")
	err := qs.Filter("title", name).One(cate)
	if err == nil {
		return err
	}

	// 插入数据
	_, err = o.Insert(cate)
	if err != nil {
		return err
	}

	return nil
}

func DeleteCategory(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	cate := &Category{Id: cid}
	_, err = o.Delete(cate)
	return err
}

func GetAllCategories() ([]*Category, error) {
	o := orm.NewOrm()

	cates := make([]*Category, 0)

	qs := o.QueryTable("category")
	_, err := qs.All(&cates)
	return cates, err
}

func AddTopic(title, content, category string) error {
	o := orm.NewOrm()

	topic := &Topic{
		Title:    title,
		Category: category,
		Content:  content,
		Created:  time.Now(),
		Updated:  time.Now(),
	}
	_, err := o.Insert(topic)
	if err != nil {
		return err
	}
	cate := new(Category)

	qs := o.QueryTable("category")

	err = qs.Filter("title", category).One(cate)

	if err == nil {
		cate.TopicCount++
		_, err = o.Update(cate)
	} else {
		AddCategory(category, 1)
		return nil
	}
	return nil
}

func GetAllTopics() ([]*Topic, error) {
	o := orm.NewOrm()

	topics := make([]*Topic, 0)
	qs := o.QueryTable("topic")
	_, err := qs.All(&topics)
	return topics, err
}

func GetTopic(id string) (*Topic, error) {
	tid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()

	// topic := new(Topic)
	// qs := o.QueryTable("topic")
	// err = qs.Filter("id", tid).One(topic)
	// if err != nil {
	// 	return nil, err
	// }
	topic := &Topic{Id: tid}
	if o.Read(topic) == nil {
		topic.Updated = time.Now()
		topic.Views++
		o.Update(topic)
	}

	// topic.Views++
	// _, err = o.Update(topic)
	// if err != nil {
	// 	return nil, err
	// }
	return topic, nil
}

func ModifyTopic(tid, title, category, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	var oldcate string
	o := orm.NewOrm()
	topic := &Topic{Id: tidNum}
	if o.Read(topic) == nil {
		topic.Title = title
		oldcate = topic.Category
		topic.Category = category
		topic.Content = content
		topic.Updated = time.Now()
		o.Update(topic)
	}

	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil {
		cate.TopicCount++
		o.Update(cate)
	} else {
		AddCategory(category, 1)
	}

	err = qs.Filter("title", oldcate).One(cate)
	if err == nil {
		cate.TopicCount--
		o.Update(cate)
	}

	return nil
}

func DeleteTopic(tid string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()

	topic := &Topic{Id: tidNum}
	var oldcate string

	if o.Read(topic) == nil {
		oldcate = topic.Category
	}

	cate := new(Category)

	qs := o.QueryTable("category")
	err = qs.Filter("title", oldcate).One(cate)

	if err == nil {
		cate.TopicCount--
		o.Update(cate)
	}

	_, err = o.Delete(topic)
	return err
}

func AddReply(tid, name, content string) error {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	o := orm.NewOrm()
	reply := &Reply{
		Tid:     tidNum,
		Name:    name,
		Content: content,
		Created: time.Now(),
	}

	_, err = o.Insert(reply)

	if err != nil {
		return err
	}

	topic := new(Topic)

	qs := o.QueryTable("topic")
	err = qs.Filter("id", tidNum).One(topic)

	if err == nil {
		topic.ReplyCount++
		_, err = o.Update(topic)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetAllReplies(tid string) ([]*Reply, error) {
	tidNum, err := strconv.ParseInt(tid, 10, 64)
	replies := make([]*Reply, 0)
	o := orm.NewOrm()

	qs := o.QueryTable("reply")
	_, err = qs.Filter("tid", tidNum).All(&replies)

	if err == nil {
		return replies, nil
	}

	return nil, nil

}

func DelReply(tid, rid string) error {
	tidNum, _ := strconv.ParseInt(tid, 10, 64)
	ridNum, _ := strconv.ParseInt(rid, 10, 64)

	o := orm.NewOrm()

	reply := &Reply{Id: ridNum}

	o.Delete(reply)

	topic := &Topic{Id: tidNum}

	err := o.Read(topic)

	if err == nil {
		topic.ReplyCount--
		o.Update(topic)
	}

	return nil

}
