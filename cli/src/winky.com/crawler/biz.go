package main

import (
	"fmt"
	"github.com/panthesingh/goson"
	"log"
	"time"
)

/*
 * 队列业务分发处理器
 */
func procBiz() {
	que := conf.Redis.Queue
	for {
		queStr, err := rpop(que)
		if err != nil {
			log.Println("queue pop msg failed:", err)
			time.Sleep(3 * time.Second)
			continue
		}

		jsonObj, err := goson.Parse([]byte(queStr))
		if err != nil {
			trace(fmt.Sprintf("ParseJson json %s failed: %v", queStr, err))
		}
		businessType := jsonObj.Get("bizType").Int()
		if businessType <= 0 {
			trace(fmt.Sprintf("ParseJson not exists bizType %s", businessType))
		}

		switch businessType {
		case BusinessArticle:
			articleId := jsonObj.Get("bizContent").Get("articleId").Int()
			article(articleId)
		default:
			trace("unserved business type")
		}
	}
}

func procCrawler() {

	//新车新闻 买车1	行情6
	//http://auto.sohu.com/xinchenews/index.shtml
	//http://auto.sohu.com/xinchenews/index_1473.shtml
	//http://auto.sohu.com/xinchenews/index_1375.shtml
	sohuHandle("http://auto.sohu.com/xinchenews/index.shtml", "gb2312", 1, 6)
	sohuBatchHandle("http://auto.sohu.com/xinchenews/index_%d.shtml", "gb2312", 1, 6, 1375, 1473)

	//对比评测 买车1	评测7
	//http://auto.sohu.com/tag/0590/000016590.shtml
	//http://auto.sohu.com/tag/0590/000016590_181.shtml
	//http://auto.sohu.com/tag/0590/000016590_83.shtml
	sohuHandle("http://auto.sohu.com/tag/0590/000016590.shtml", "gb2312", 1, 7)
	sohuBatchHandle("http://auto.sohu.com/tag/0590/000016590_%d.shtml", "gb2312", 1, 7, 83, 181)

	//试驾报告	买车1 评测7
	//http://auto.sohu.com/tag/0796/000016796.shtml
	//http://auto.sohu.com/tag/0796/000016796_145.shtml
	//http://auto.sohu.com/tag/0796/000016796_47.shtml

	sohuHandle("http://auto.sohu.com/tag/0796/000016796.shtml", "gb2312", 1, 7)
	sohuBatchHandle("http://auto.sohu.com/tag/0796/000016796_%d.shtml", "gb2312", 1, 7, 47, 145)

	//深度测试	http://auto.sohu.com/tag/0160/000025160.shtml	买车1	评测7
	//http://auto.sohu.com/tag/0160/000025160_57.shtml
	//http://auto.sohu.com/tag/0160/000025160_1.shtml

	sohuHandle("http://auto.sohu.com/tag/0160/000025160.shtml", "gb2312", 1, 7)
	sohuBatchHandle("http://auto.sohu.com/tag/0160/000025160_%d.shtml", "gb2312", 1, 7, 1, 57)

	//车展新车	http://auto.sohu.com/tag/0553/000015553.shtml	买车1	行情6
	//http://auto.sohu.com/tag/0553/000015553_236.shtml
	//http://auto.sohu.com/tag/0553/000015553_138.shtml

	sohuHandle("http://auto.sohu.com/tag/0553/000015553.shtml", "gb2312", 1, 6)
	sohuBatchHandle("http://auto.sohu.com/tag/0553/000015553_%d.shtml", "gb2312", 1, 6, 138, 236)


	//新车曝光	http://auto.sohu.com/tag/0478/000016478.shtml	买车1	行情6
	//http://auto.sohu.com/tag/0478/000016478_181.shtml
	//http://auto.sohu.com/tag/0478/000016478_83.shtml

	sohuHandle("http://auto.sohu.com/tag/0478/000016478.shtml", "gb2312", 1, 6)
	sohuBatchHandle("http://auto.sohu.com/tag/0478/000016478_%d.shtml", "gb2312", 1, 6, 83, 181)


	//新车到店	http://auto.sohu.com/tag/0696/000032696.shtml	买车1	行情6
	sohuHandle("http://auto.sohu.com/tag/0696/000032696.shtml", "gb2312", 1, 6)


	//维修保养	http://auto.sohu.com/baoyang/index.shtml	用车2	保养14
	//http://auto.sohu.com/baoyang/index_219.shtml
	//http://auto.sohu.com/baoyang/index_121.shtml
	sohuHandle("http://auto.sohu.com/baoyang/index.shtml", "gb2312", 2, 14)
	sohuBatchHandle("http://auto.sohu.com/baoyang/index_%d.shtml", "gb2312", 2, 14, 131, 219)


	//行车技巧	http://auto.sohu.com/jiqiao/index.shtml		用车2	经验11
	//http://auto.sohu.com/jiqiao/index_225.shtml
	//http://auto.sohu.com/jiqiao/index_127.shtml
	sohuHandle("http://auto.sohu.com/jiqiao/index.shtml", "gb2312", 2, 11)
	sohuBatchHandle("http://auto.sohu.com/jiqiao/index_%d.shtml", "gb2312", 2, 11, 127, 225)


	//日常养护	http://auto.sohu.com/tag/0761/000015761.shtml	用车2	保养14
	//http://auto.sohu.com/tag/0761/000015761_75.shtml
	//http://auto.sohu.com/tag/0761/000015761_1.shtml
	sohuHandle("http://auto.sohu.com/tag/0761/000015761.shtml", "gb2312", 2, 14)
	sohuBatchHandle("http://auto.sohu.com/tag/0761/000015761_%d.shtml", "gb2312", 2, 14, 1, 75)


	//新手用车	http://auto.sohu.com/tag/0540/000032540.shtml	用车2	经验11
	//http://auto.sohu.com/tag/0540/000032540_1.shtml
	sohuHandle("http://auto.sohu.com/tag/0540/000032540.shtml", "gb2312", 2, 11)
	sohuBatchHandle("http://auto.sohu.com/tag/0540/000032540_%d.shtml", "gb2312", 2, 11, 1, 1)


	//事故处理	http://auto.sohu.com/tag/0612/000016612.shtml	用车2	经验11
	//http://auto.sohu.com/tag/0612/000016612_20.shtml
	//http://auto.sohu.com/tag/0612/000016612_1.shtml
	sohuHandle("http://auto.sohu.com/tag/0612/000016612.shtml", "gb2312", 2, 11)
	sohuBatchHandle("http://auto.sohu.com/tag/0612/000016612_%d.shtml", "gb2312", 2, 11, 1, 20)


	//4S保养	http://auto.sohu.com/tag/0548/000015548.shtml	用车2	保养14
	//http://auto.sohu.com/tag/0548/000015548_27.shtml
	//http://auto.sohu.com/tag/0548/000015548_1.shtml
	sohuHandle("http://auto.sohu.com/tag/0548/000015548.shtml", "gb2312", 2, 14)
	sohuBatchHandle("http://auto.sohu.com/tag/0548/000015548_%d.shtml", "gb2312", 2, 14, 1, 27)


	//汽车改装	http://auto.sohu.com/transform/index.shtml	玩车3	改装18
	//http://auto.sohu.com/transform/index_92.shtml
	//http://auto.sohu.com/transform/index_1.shtml
	sohuHandle("http://auto.sohu.com/transform/index.shtml", "gb2312", 3, 18)
	sohuBatchHandle("http://auto.sohu.com/transform/index_%d.shtml", "gb2312", 3, 18, 1, 92)


	//历史文化	http://auto.sohu.com/culture/index.shtml	玩车3	趣闻15
	//http://auto.sohu.com/culture/index_138.shtml
	//http://auto.sohu.com/culture/index_40.shtml
	sohuHandle("http://auto.sohu.com/culture/index.shtml", "gb2312", 3, 15)
	sohuBatchHandle("http://auto.sohu.com/culture/index_%d.shtml", "gb2312", 3, 15, 40, 138)


	//自驾游记	http://auto.sohu.com/zijia/index.shtml		玩车3	自驾17
	//http://auto.sohu.com/zijia/index_36.shtml
	//http://auto.sohu.com/zijia/index_1.shtml
	sohuHandle("http://auto.sohu.com/zijia/index.shtml", "gb2312", 3, 17)
	sohuBatchHandle("http://auto.sohu.com/zijia/index_%d.shtml", "gb2312", 3, 17, 1, 36)


	//卖二手车	http://auto.sohu.com/tag/0997/000021997.shtml	卖车4	高价卖19
	//http://auto.sohu.com/tag/0997/000021997_2.shtml
	//http://auto.sohu.com/tag/0997/000021997_1.shtml
	sohuHandle("http://auto.sohu.com/tag/0997/000021997.shtml", "gb2312", 4, 19)
	sohuBatchHandle("http://auto.sohu.com/tag/0997/000021997_%d.shtml", "gb2312", 4, 19, 1, 2)

}