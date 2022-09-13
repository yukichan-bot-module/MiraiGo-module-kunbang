package kunbang

import (
	"strings"
	"sync"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

var instance *kunbang
var logger = utils.GetModuleLogger("com.aimerneige.kunbang")

type kunbang struct {
}

func init() {
	instance = &kunbang{}
	bot.RegisterModule(instance)
}

func (k *kunbang) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "com.aimerneige.kunbang",
		Instance: instance,
	}
}

// Init 初始化过程
// 在此处可以进行 Module 的初始化配置
// 如配置读取
func (k *kunbang) Init() {
}

// PostInit 第二次初始化
// 再次过程中可以进行跨 Module 的动作
// 如通用数据库等等
func (k *kunbang) PostInit() {
}

// Serve 注册服务函数部分
func (k *kunbang) Serve(b *bot.Bot) {
	b.GroupMessageEvent.Subscribe(func(c *client.QQClient, msg *message.GroupMessage) {
		// 消息元素数量少于 2，直接退出
		if len(msg.Elements) < 2 {
			return
		}
		atTarget := int64(0)
		isAt := false
		isKunbang := false
		// 只关心前俩个消息元素
		for i := 0; i < 2; i++ {
			ele := msg.Elements[i]
			switch e := ele.(type) {
			case *message.AtElement:
				if !isAt {
					isAt = true
					atTarget = e.Target
				}
			case *message.TextElement:
				if !isKunbang {
					content := strings.TrimSpace(e.Content)
					if content == "捆绑" {
						isKunbang = true
					}
				}
			}
		}
		// 消息符合规则，尝试禁言并发送消息
		if isAt && isKunbang {
			targetMemberInfo, err := c.GetMemberInfo(msg.GroupCode, atTarget)
			if err != nil {
				logger.WithError(err).
					WithField("GroupCode", msg.GroupCode).
					WithField("GroupName", msg.GroupName).
					WithField("Target", atTarget).
					Error("Fail to get group member info.")
				return
			}
			if err := targetMemberInfo.Mute(uint32(16)); err == nil {
				sendingMsg := message.NewSendingMessage()
				sendingMsg.Append(message.NewAt(msg.Sender.Uin))
				sendingMsg.Append(message.NewText("把"))
				// 检查一下是不是在捆绑自己
				if atTarget != msg.Sender.Uin {
					sendingMsg.Append(message.NewAt(atTarget))
				} else {
					sendingMsg.Append(message.NewText("自己"))
				}
				sendingMsg.Append(message.NewText("绑起来了！"))
				c.SendGroupMessage(msg.GroupCode, sendingMsg)
			}
		}
	})
}

// Start 此函数会新开携程进行调用
// ```go
//
//	go exampleModule.Start()
//
// ```
// 可以利用此部分进行后台操作
// 如 http 服务器等等
func (k *kunbang) Start(b *bot.Bot) {
}

// Stop 结束部分
// 一般调用此函数时，程序接收到 os.Interrupt 信号
// 即将退出
// 在此处应该释放相应的资源或者对状态进行保存
func (k *kunbang) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	// 别忘了解锁
	defer wg.Done()
}
