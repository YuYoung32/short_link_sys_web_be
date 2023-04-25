/**
 * Created by YuYoung on 2023/4/25
 * Description: 生成测试短链
 */

package link

import (
	"errors"
	"github.com/bits-and-blooms/bloom/v3"
	"gorm.io/gorm"
	"math/rand"
	"short_link_sys_web_be/conf"
	"short_link_sys_web_be/database"
	"short_link_sys_web_be/link_gen"
	"short_link_sys_web_be/log"
	"strconv"
)

// GenerateLinkData 生成测试短链, 需要提前初始化conf和database
func GenerateLinkData(list []database.Link) error {
	var _linkGenAlgorithm link_gen.LinkGen
	var _shortLinkBF *bloom.BloomFilter
	var _longLinkBF *bloom.BloomFilter

	logger := log.GetLogger()
	db := database.GetDBInstance()
	mapAlgorithm := map[string]link_gen.LinkGen{
		"murmurHash":   link_gen.MurmurHash{},
		"xxHash":       link_gen.XXHash{},
		"fnvHash":      link_gen.FNVHash{},
		"simpleSeq":    link_gen.SimpleSequencer{},
		"snowflakeSeq": link_gen.SnowflakeSequencer{},
	}

	algorithmName := conf.GlobalConfig.GetString("handler.link.algorithm")
	if _linkGenAlgorithm = mapAlgorithm[algorithmName]; _linkGenAlgorithm == nil {
		return errors.New("bad config: handler.link.algorithm")
	} else {
		logger.Info("set link algorithm to: " + conf.GlobalConfig.GetString("handler.link.algorithm"))
	}

	// 布隆过滤器初始化, 用于检测是否已经存在该长链
	_longLinkBF = bloomFilterInit("long_link")
	// 使用hash算法时, 才会使用短链的布隆过滤器
	if _linkGenAlgorithm.GetType() == link_gen.HashType {
		_shortLinkBF = bloomFilterInit("short_link")
	}

	for i, item := range list {
		// 绝对链接检测与转换
		longLink := processRawLink(item.LongLink)
		// 布隆过滤器检测是否已经存在长链, 若存在无需再次添加
		if _longLinkBF.TestString(longLink) {
			// 排除假阳性
			err := db.Where(&database.Link{LongLink: longLink}).Take(&database.Link{}).Error
			if err == nil {
				// 已确认数据库中有该长链
				continue
			} else if err != gorm.ErrRecordNotFound {
				return err
			}
		}

		_genUniqueLink := func(longLink string) string {
			var shortLink string
			switch _linkGenAlgorithm.GetType() {
			case link_gen.HashType:
				for {
					shortLink = _linkGenAlgorithm.GenLink(longLink)
					if !_shortLinkBF.TestString(shortLink) {
						break
					}
					longLink = longLink + strconv.Itoa(rand.Int())
					shortLink = _linkGenAlgorithm.GenLink(longLink)
				}
			case link_gen.SeqType:
				shortLink = _linkGenAlgorithm.GenLink(longLink)
			}
			return shortLink
		}

		// 已确认数据库中无长链
		_longLinkBF.AddString(longLink)
		shortLink := _genUniqueLink(longLink)
		if _shortLinkBF != nil {
			_shortLinkBF.AddString(shortLink)
		}
		list[i].ShortLink = shortLink
	}

	err := db.Create(list).Error
	return err
}
