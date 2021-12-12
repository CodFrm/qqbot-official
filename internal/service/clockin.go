package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/CodFrm/qqbot-official/internal/db"
	"github.com/CodFrm/qqbot/utils"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type ClockIn interface {
	ClockIn(guild, user string) (string, error)
	SleepClockIn(guild, user string) (string, error)
	GetUpClockIn(guild, user string) (string, error)
	GetUpList(guild string) ([]string, error)
}

type clockIn struct {
	user User
}

func NewClockIn(c *cron.Cron, user User) ClockIn {
	ret := &clockIn{user: user}
	_, err := c.AddFunc("0 8 * * *", ret.notGetup)
	if err != nil {
		logrus.Fatalf("clockin task: %v", err)
	}
	return ret
}

func (c *clockIn) notGetup() {
	now := time.Now().Add(-time.Hour * 24)
	guildList, err := db.SGetAll("clockin:guild:" + now.Format("2006/01/02"))
	if err != nil {
		logrus.Errorf("clockin guild list: %v", err)
		return
	}
	for _, guild := range guildList {
		users, err := db.SGetAll("clockin:getup:member:" + now.Format("2006/01/02") + ":" + string(guild))
		if err != nil {
			logrus.Errorf("clockin getup member %s: %v", guild, err)
			continue
		}
		for _, user := range users {
			if _, err := c.user.AddIntegral(string(guild), string(user), -3); err != nil {
				logrus.Errorf("guild %s user %s add -3 integral: %v", guild, user, err)
			}
		}
	}
}

func (c *clockIn) GetUpList(guild string) ([]string, error) {
	now := time.Now()
	if now.Hour() < 8 {
		return nil, errors.New("榜单还未生成")
	}
	users, err := db.SGetAll("clockin:getup:member:" + now.Add(-time.Hour*24).Format("2006/01/02") + ":" + guild)
	if err != nil {
		return nil, err
	}
	var list []string
	for _, user := range users {
		list = append(list, string(user))
	}
	return list, nil
}

func (c *clockIn) key(t time.Time, guild, user string) string {
	return fmt.Sprintf("clockin:sleep:%s:%s:%s", t.Format("2006/01/02"), guild, user)
}

func (c *clockIn) SleepClockIn(guild, user string) (string, error) {
	now := time.Now()
	if now.Hour() > 22 {
		return "", errors.New("已经过点了,请在晚上10点前打卡")
	}
	if now.Hour() < 19 {
		return "", errors.New("天还没暗呢")
	}
	if err := db.Put(c.key(now, guild, user), []byte(strconv.FormatInt(time.Now().Unix(), 10)), 86400*2); err != nil {
		return "", err
	}
	if err := db.SAdd("clockin:getup:member:"+now.Format("2006/01/02")+":"+guild, []byte(user)); err != nil {
		logrus.Errorf("clockin add sleep %s: %v", user, err)
	}
	if err := db.SAdd("clockin:guild:"+now.Format("2006/01/02"), []byte(guild)); err != nil {
		logrus.Errorf("clockin add guild %s: %v", guild, err)
	}
	return fmt.Sprintf("现在是%v,早早进入梦乡吧,起床之后记得艾特猫猫早起打卡哦", now.Format("15:04")), nil
}

func (c *clockIn) GetUpClockIn(guild, user string) (string, error) {
	now := time.Now()
	val, err := db.Get(fmt.Sprintf("clockin:getup:days:%s:%s:%s", now.Format("2006/01/02"), guild, user))
	if err != nil {
		return "", err
	}
	if val != nil {
		return "", errors.New("已经打过卡了")
	}
	val, err = db.Get(c.key(now.Add(-time.Hour*24), guild, user))
	if err != nil {
		return "", err
	}
	if val == nil {
		return "", errors.New("昨天没有参与打卡")
	}
	if now.Hour() > 8 {
		if err := db.SRem("clockin:getup:member:"+now.Add(-time.Hour*24).Format("2006/01/02")+":"+guild, []byte(user)); err != nil {
			logrus.Errorf("clockin remove sleep %s: %v", user, err)
		}
		return "", errors.New("已经过点了,扣除3积分!")
	}
	err = db.Put(fmt.Sprintf("clockin:getup:days:%s:%s:%s", time.Now().Format("2006/01/02"), guild, user), []byte("1"), 86400*2)
	if err != nil {
		return "", err
	}
	num, _ := db.Incr("clockin:getup:"+guild, 1, 86400*2)
	if err := db.SRem("clockin:getup:member:"+now.Add(-time.Hour*24).Format("2006/01/02")+":"+guild, []byte(user)); err != nil {
		logrus.Errorf("clockin remove sleep %s: %v", user, err)
	}
	if _, err := c.user.AddIntegral(guild, user, 6); err != nil {
		logrus.Errorf("guild %s user %s add 6 integral: %v", guild, user, err)
	}
	return fmt.Sprintf("新的一天开始啦~您是第%d位起床者,昨晚睡了%d小时,奖励%d积分", num, (now.Unix()-utils.StringToInt64(string(val)))/3600, 6), nil
}

func (c *clockIn) ClockIn(guild, user string) (string, error) {
	val, err := db.Get(fmt.Sprintf("clockin:days:%s:%s:%s", time.Now().Format("2006/01/02"), guild, user))
	if err != nil {
		return "", err
	}
	if val != nil {
		return "", errors.New("今天已经打过卡了")
	}
	err = db.Put(fmt.Sprintf("clockin:days:%s:%s:%s", time.Now().Format("2006/01/02"), guild, user), []byte("1"), 86400*2)
	if err != nil {
		return "", err
	}
	if _, err := c.user.AddIntegral(guild, user, 2); err != nil {
		logrus.Errorf("guild %s user %s add 2 integral: %v", guild, user, err)
	}
	return "打卡成功,增加积分2", nil
}
