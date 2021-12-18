package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/CodFrm/qqbot-official/internal/db"
	"github.com/CodFrm/qqbot/utils"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi"
)

var ErrClockExist = errors.New("今天已经打过卡了")

type ChannelClockIn interface {
	IsClock(guild, user string) (bool, error)
	ClockIn(guild, user string) error
	SetNotice(guild, channel, cron, title, content string) error
	SetClock(guild, channel string) error
}

type ClockIn interface {
	ClockIn(guild, user string) (string, error)
	SleepClockIn(guild, user string) (string, error)
	GetUpClockIn(guild, user string) (string, error)
	GetUpList(guild string) ([]string, error)
}

type clockIn struct {
	user User
	api  openapi.OpenAPI
}

func NewClockIn(c *cron.Cron, user User, api openapi.OpenAPI) ClockIn {
	ret := &clockIn{user: user, api: api}
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
	// 这与9点前打卡也成功是一个feature
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
	if now.Hour() < 4 {
		return "", errors.New("还没到点呢")
	}
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
	num, _ := db.Incr("clockin:getup:"+now.Format("2006/01/02")+":"+guild, 1, 86400*2)
	if err := db.SRem("clockin:getup:member:"+now.Add(-time.Hour*24).Format("2006/01/02")+":"+guild, []byte(user)); err != nil {
		logrus.Errorf("clockin remove sleep %s: %v", user, err)
	}
	if _, err := c.user.AddIntegral(guild, user, 6); err != nil {
		logrus.Errorf("guild %s user %s add 6 integral: %v", guild, user, err)
	}
	return fmt.Sprintf("新的一天开始啦~您是第%d位起床者,昨晚睡了%d小时,奖励%d积分", num, (now.Unix()-utils.StringToInt64(string(val)))/3600, 6), nil
}

func (c *clockIn) ClockIn(guild, user string) (string, error) {
	if err := c.clockIn("days", guild, user, 2); err != nil {
		return "", err
	}
	return "打卡成功,增加积分2", nil
}

func (c *clockIn) clockIn(feat, guild, user string, integral int) error {
	val, err := db.Get(fmt.Sprintf("clockin:days:%s:%s:%s:%s", feat, time.Now().Format("2006/01/02"), guild, user))
	if err != nil {
		return err
	}
	if val != nil {
		return ErrClockExist
	}
	err = db.Put(fmt.Sprintf("clockin:days:%s:%s:%s:%s", feat, time.Now().Format("2006/01/02"), guild, user), []byte("1"), 86400*2)
	if err != nil {
		return err
	}
	if _, err := c.user.AddIntegral(guild, user, int64(integral)); err != nil {
		logrus.Errorf("guild %s user %s %s add %d integral: %v", guild, user, feat, integral, err)
	}
	return nil
}

type channelClockIn struct {
	opts *ChannelClockInOptions

	api     openapi.OpenAPI
	cron    *cron.Cron
	user    User
	feat    string
	oldCron []cron.EntryID
}

type ChannelClockInOptions struct {
	Integral int64
	Limit    func() error
}

func NewChannelClockIn(c *cron.Cron, user User, api openapi.OpenAPI, feat string, opts *ChannelClockInOptions) ChannelClockIn {
	ret := &channelClockIn{
		api:  api,
		cron: c,
		user: user,
		feat: feat,
		opts: opts,
	}
	if err := ret.cronNotice(); err != nil {
		logrus.Errorf("cron notice: %v", err)
	}
	return ret
}

func (c *channelClockIn) cronNotice() error {
	list, err := db.SGetAll("clockin:notice:guild")
	if err != nil {
		logrus.Errorf("notice list: %v", err)
		return err
	}
	for _, v := range c.oldCron {
		c.cron.Remove(v)
	}
	c.oldCron = make([]cron.EntryID, 0)
	for _, v := range list {
		channel, err := db.Get("clockin:notice:channel:" + c.feat + ":" + string(v))
		if err != nil {
			logrus.Errorf("notice channel: %v", err)
			continue
		}
		clock, err := db.Get("clockin:notice:clock:" + c.feat + ":" + string(v))
		if err != nil {
			logrus.Errorf("notice channel: %v", err)
			continue
		}
		if clock == nil {
			clock = channel
		}
		cron, err := db.Get("clockin:notice:cron:" + c.feat + ":" + string(v))
		if err != nil {
			logrus.Errorf("notice cron: %v", err)
			continue
		}
		title, err := db.Get("clockin:notice:title:" + c.feat + ":" + string(v))
		if err != nil {
			logrus.Errorf("notice title: %v", err)
			continue
		}
		content, err := db.Get("clockin:notice:content:" + c.feat + ":" + string(v))
		if err != nil {
			logrus.Errorf("notice content: %v", err)
			continue
		}
		if channel == nil || cron == nil || content == nil || title == nil || clock == nil {
			continue
		}
		e, err := c.cron.AddFunc(string(cron), func() {
			_, err = c.api.PostMessage(context.Background(), string(channel), &dto.MessageToCreate{
				Ark: &dto.Ark{
					TemplateID: 24,
					KV: []*dto.ArkKV{{
						Key:   "#DESC#",
						Value: string(title),
					}, {
						Key:   "#PROMPT#",
						Value: string(title),
					}, {
						Key:   "#TITLE#",
						Value: string(title),
					}, {
						Key:   "#METADESC#",
						Value: strings.ReplaceAll(string(content), "{channel}", "<#"+string(clock)+">"),
					}},
				},
			})
			if err != nil {
				logrus.Errorf("notice send: %v", err)
			}
		})
		//"新的一天开始啦,快前往{channel}进行打卡签到吧:早起打卡(+6),早八人打卡(+3),使用学习软件分享进行打卡(+4)"
		if err != nil {
			return err
		}
		c.oldCron = append(c.oldCron, e)
	}
	return nil
}

func (c *channelClockIn) IsClock(guild, user string) (bool, error) {
	val, err := db.Get(fmt.Sprintf("clockin:days:%s:%s:%s:%s", c.feat, time.Now().Format("2006/01/02"), guild, user))
	if err != nil {
		return false, err
	}
	if val != nil {
		return true, nil
	}
	return false, nil
}

func (c *channelClockIn) ClockIn(guild, user string) error {
	ok, err := c.IsClock(guild, user)
	if err != nil {
		return err
	}
	if ok {
		return ErrClockExist
	}
	if c.opts.Limit != nil {
		if err := c.opts.Limit(); err != nil {
			return err
		}
	}
	err = db.Put(fmt.Sprintf("clockin:days:%s:%s:%s:%s", c.feat, time.Now().Format("2006/01/02"), guild, user), []byte("1"), 86400*2)
	if err != nil {
		return err
	}
	if _, err := c.user.AddIntegral(guild, user, c.opts.Integral); err != nil {
		logrus.Errorf("guild %s user %s %s add %d integral: %v", guild, user, c.feat, c.opts.Integral, err)
	}
	return nil
}

func (c *channelClockIn) SetNotice(guild, channel, cron, title, content string) error {
	if err := db.SAdd("clockin:notice:guild", []byte(guild)); err != nil {
		return err
	}
	if err := db.Put("clockin:notice:channel:"+c.feat+":"+guild, []byte(channel), 0); err != nil {
		return err
	}
	if err := db.Put("clockin:notice:cron:"+c.feat+":"+guild, []byte(cron), 0); err != nil {
		return err
	}
	if err := db.Put("clockin:notice:title:"+c.feat+":"+guild, []byte(title), 0); err != nil {
		return err
	}
	if err := db.Put("clockin:notice:content:"+c.feat+":"+guild, []byte(content), 0); err != nil {
		return err
	}
	return c.cronNotice()
}

func (c *channelClockIn) SetClock(guild, channel string) error {
	if err := db.SAdd("clockin:notice:guild", []byte(guild)); err != nil {
		return err
	}
	return db.Put("clockin:notice:clock:"+c.feat+":"+guild, []byte(channel), 0)
}
