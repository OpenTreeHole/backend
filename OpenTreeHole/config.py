# 站点配置

import os
import uuid

SITE_NAME = 'Open Tree Hole'  # 网站名称
TZ = "Asia/Shanghai"  # 时区
LANGUAGE = "zh-Hans"
ALLOW_CONNECT_HOSTS = ['opentreehole.org']  # 允许连接的域名
EMAIL_WHITELIST = ["test.com"]  # 允许注册树洞的邮箱域名
MIN_PASSWORD_LENGTH = 8  # 允许的最短用户密码长度
VALIDATION_CODE_EXPIRE_TIME = 5  # 验证码失效时间（分钟）
MAX_TAGS = 5
MAX_TAG_LENGTH = 8
NAME_LIST = ['陌客', '安逸', '独照', '亡心', '梦巷', '逐风', '花憩', '余了', '在乎', '上隐', '璃茉', '擱淺', '双臂', '寰鸾', '不烟', '长野', '怪味', '知人', '沐夏', '旭明', '昔年', '穆梦', '自私', '冬渡', '栀瑜', '临风', '孤崖', '栀蓝',
             '看懂', '犊子', '顾挽', '旧夢', '煞尾', '安烨', '超标', '诗句', '冰尘', '洛夏', '灼夏', '振海', '欤你', '柠心', '野性', '鸢浅', '允你', '寄生', '忆夕', '屿卿', '十四', '流晞', '礼失', '沙痕', '淘汰', '忆情', '芥末', '情場',
             '烟沙', '绾起', '凶傲', '遇热', '思君', '凉牧', '仄言', '荒漠', '南征', '轻轨', '忆梦', '葵顷', '孤亡', '花漓', '离岸', '北宸', '烟水', '猫景', '得闲', '恻隐', '海心', '夙缘', '呆橘', '曦夏', '浔梦', '柠汐', '苏苏', '流萤',
             '贝兎', '成思', '风声', '初始', '天下', '墨汐', '温笙', '陪伴', '海星', '断忆', '麦兜', '嘟嘟', '尐脸', '哇塞', '浅陌', '陶醉', '饮湿', '执迷', '凉牧', '冰雨', '画生', '允你', '无畏', '热血', '清涩', '罪过', '栀蓝', '忆心',
             '歎之', '心裂', '逆流', '梦缘', '萌兽', '淇淋', '思春', '颜控', '呛心', '徒劳', '偏食', '栀蓝', '情殇', '抵年', '临渊', '屿风', '衬眉', '浅巷', '心抑', '清悦', '隔壁', '千帆', '酒泉', '幽殤', '俊民', '宇恒', '墨漓', '词庸',
             '俗人', '浮梦', '安若', '坠落', '離開', '风轻', '森迷', '过分', '不见', '独自', '祭陌', '暮昼', '慕芹', '挚爱', '红鸾', '尘缘', '森屿', '角度', '趣味', '追梦', '与你', '寻找', '寒冬', '陌颜', '浅酌', '溺爱', '良人', '痛心',
             '陪伴', '夏姬', '奇葩', '白头', '缠绕', '起伏', '闭眸', '猫余', '话梗', '夏冰', '玖随', '古巷', '鹤隐', '南桥', '太阳', '污女', '北街', '日记', '泠鸢', '寰鸾', '邀月', '致爱', '尐脸', '迷鹿', '碾压', '温书', '媚眼', '光澜',
             '染玖', '無言', '沐离', '锦瑟', '陌森', '一时', '子夏', '不闻', '北岸', '格调', '安生', '兜兜', '玩腻', '萌主', '灵魂', '伤心', '速递', '沉溺', '不顾', '供电', '南玖', '自渡', '宥沐', '诱夏', '浅意', '涼城', '負卿', '调皮',
             '轨車', '难遇', '稚气', '和蛙', '炎凉', '夜巷', '一休', '漓殇', '心癌', '初愈', '北殷', '安柒', '祥阔', '云荒', '出发', '浪徒', '清词', '湛蓝', '路途', '偏要', '梦回', '墨渊', '仙糖', '南简', '沉默', '暮迟', '女王', '玩家',
             '野慌', '孤屿', '夏姬', '烟沙', '松懈', '谈天', '夜巷', '军权', '橘子', '初墨', '秋摇', '逝爱', '夜染', '寥寥', '昔望', '甜心', '晚归', '顾挽', '听寒', '梦港', '窈窕', '俊迈', '德元', '傲骨', '别搅', '词庸', '夏利', '夏倩',
             '芯話', '轻叩', '天涯', '眉目', '生白', '偏食', '云深', '黑名', '任性', '绝色', '溺于', '璎珞', '温油', '峰弟', '一休', '元嘉', '文来', '氤氲', '摘星', '夏柒', '慵懒', '甜屁', '疯度', '温妤', '宠幼', '归州', '糜废', '未挽',
             '归遇', '百川', '决绝', '随便', '兜兜', '夏祭', '凉橙', '茶颜', '未脱', '告白', '精緻', '少女', '妄囍', '鱼溺', '素顏', '安安', '清涩', '浪迹', '尾戒', '酷刑', '调皮', '诺言', '俘获', '嘤咛', '余香', '韶华', '挑衅', '脱缰',
             '甘来', '绝色', '西门', '鱼干', '半枫', '舞琴', '清絮', '轻叩', '取决', '孤凉', '久安', '独霸', '枕书', '模仿', '菇凉', '鱼芗', '现实', '物类', '心系', '巡游', '浅眠', '祈风', '初妆', '锐翰', '萤火', '倾情', '璃火', '江东',
             '殇雪', '池华', '黎开', '青丝', '青词', '风霄', '舟尼', '北觅', '笑望', '沉溺', '勾画', '绝版', '夜旅', '七木', '峰弟', '一休', '含莲', '问薇', '鹤隐', '沉栋', '归遇', '旧巷', '决绝', '空荡', '原创', '陌上', '小鱼', '歡囍',
             '林鹊', '夏熙', '顺庆', '夕陌', '囍旧', '染玖', '冷妆', '崩溃', '蓝鸢', '荒凉', '夕颜', '拌你', '心上', '甜未', '逗逗', '柚柠', '梦罢', '冷月', '鬼迹', '懵懂', '梦醒', '传说', '逆流', '繁心', '心如', '千帆', '反叛', '浅忘',
             '稚初', '凉牧', '风渺', '於穆', '绝凌', '迷途', '浅笑', '蹉跎', '斑驳', '荆棘', '闹巷', '束缚', '初闻', '夏沫', '浅葬', '春水', '供电', '遇热', '南玖', '绍钧', '北丧', '寒冬', '北葵', '栀瑜', '锦谧', '難尋', '无味', '梦巷',
             '不堪', '讨好', '慌屿', '从蓉', '失心', '烟霸', '桃花', '如故', '爱我', '沙雕', '幽夏', '宏恺', '末年', '慵懒', '殇泪', '北栀', '沐白', '后来', '浅笙', '久别', '知遇', '鬼迹', '衡虑', '告别', '凉墨', '临渊', '西风', '探月',
             '凛然', '不负', '棒棒', '心誶', '折磨', '落兮', '风灵', '恶魔', '折扇', '甜岛', '多坎', '擱淺', '肆意', '一半', '无我', '寡欢', '如煙', '风度', '一缕', '草裙', '酒香', '私奔', '青瓷', '唇红', '难枕', '孤凉', '苦泪', '黛眉',
             '浅瞳', '一夜', '忌惮', '人世', '邂逅', '初晴', '双臂', '酒腻', '栀寒', '吧唧', '妙彤', '心海', '初识', '老弟', '闲云', '浅忘', '少年', '触觉', '晚吟', '未挽', '北诗', '桔栀', '陌颜', '故里', '曾将', '七味', '空蒙', '虐爱',
             '唯爱', '宿觞', '忆情', '服软', '橘香', '碾压', '刻薄', '傲气', '二囍', '孤寐', '愚剧', '堇夏', '边侣', '流年', '刺骨', '湮灭', '祭陌', '夜寐', '星夜', '童年', '小鱼', '单殺', '浅忘', '寒风', '獨淰', '夏桐', '暗涌', '虞生',
             '卷耳', '笙念', '茫然', '顾染', '酒笙', '累心', '邂逅', '走野', '辜予', '小孩', '趁早', '眸敛', '劣迹', '友人', '妄凝', '仙儿', '雨樱', '甜屿', '三年', '情愫', '细诉', '乞许', '淡网', '十雾', '语酌', '配角', '英雄', '奶盖',
             '逆天', '囚宠', '煙花', '忘年', '难拥', '一口', '凌云', '超标', '若梦', '闭眸', '清念', '德元', '娇气', '殇怹', '孤寂', '赴我', '故人', '随性', '拥抱', '之南', '布丁', '萌兽', '陌念', '眸敛', '青萍', '渔民', '夏晴', '北祀',
             '笠含', '振海', '星承', '归宿', '寒潮', '机遇', '枕梦', '微风', '难遇', '炼狱', '句点', '发光', '豆蔻', '眸海', '夏课', '云暖', '漂泊', '暖男', '染玖', '舍予', '噬骨', '体味', '腻歪', '倾寒', '迎天', '秋蝶', '瑾夏', '停步',
             '孤酒', '落花', '無歡', '槿栀', '堅強', '人散', '独钓', '相念', '无语', '犊子', '几欢', '渊鱼', '毒打', '棒棒', '冷颜', '谈花', '病娇', '岛徒', '聆回', '妄言', '黎晴', '轻轨', '顾染', '偏执', '无情', '疏离', '酒腻', '如煙',
             '逗比', '于倩', '思春', '初识', '缠绵', '私奔', '锐阵', '德元', '执傲', '蚀妆', '晴栀', '热血', '酒痕', '離開', '逐风', '暮辞', '屿暖', '风情', '多坎', '旧爱', '存在', '素歆', '思归', '东寻', '清悦', '演繹', '晨曦', '邮寄',
             '星辉', '泠崖', '沧浪', '晚吟', '边侣', '浪漫', '淡然', '模仿', '潦倒', '二二', '有爱', '害羞', '素衣', '如故', '乖兽', '颜控', '夏芷', '大傻', '浅殇', '异情', '悲歌', '独醉', '深岛', '残年', '甜未', '下沉', '晚期', '亡心',
             '理念', '热情', '兜兜', '炙雪', '孤猫', '夏浅', '香浅', '梦毁', '暗谣', '特别', '躁动', '绝版', '青睐', '慕言', '俗趣', '挽留', '无畏', '浪漫', '初始', '若颜', '忆泙', '痴怨', '花枝', '拿糖', '浅歌', '奢念', '心如', '浅羽',
             '友人', '鸢尾', '逗霸', '未初', '刺青', '云雾', '蓝天', '森眸', '信马', '沐白', '木舟', '心上', '满志', '沧桑', '流年', '囚宠', '相約', '虚妄', '煞尾', '玩世', '赖床', '物类', '海心', '恬淡', '思仙', '往事', '人渣', '自嘲',
             '苦泪', '软酱', '自然', '人散', '炙雪', '梓桑', '雾月', '文饰', '毒打', '沙雕', '熟吻', '离安', '宇恒', '森雾', '星战', '苟活', '北陌', '殇怹', '子夏', '北执', '临风', '相忘', '浅酌', '相約', '野慌', '尘埃', '鸽屿', '从蓉',
             '少女', '绵衫', '伸手', '诀别', '大傻', '风吟', '寒栀', '心癌', '蓝天', '長野', '淡网', '入戏', '残雪', '长歡', '末玖', '泡芙', '瑾凉', '明日', '不弃', '欢颜', '堇年', '剑圣', '蕾溪', '顽劣', '冬灼', '傲娇', '陶醉', '烟霸',
             '小川', '无妄', '若兮', '扰梦', '停泊', '风声', '擱淺', '野途', '静心', '方式', '相倚', '承欢', '卖萌', '师傅', '一碗', '诱夏', '锐翰', '浩博', '子真', '喵语', '怀桔', '陌森', '流徙', '笙南', '逊色', '孑然', '无痕', '情人',
             '衾曦', '归往', '失心', '逆流', '琴瑟', '春花', '北战', '耳畔', '静听', '浅眠', '谈花', '宥沐', '逆夏', '越界', '茶味', '取决', '浅离', '蜕变', '魅绪', '不知', '夙兴', '雨夜', '重来', '简白', '悠闲', '琉羽', '暗谣', '笑话',
             '方式', '相倚', '抚风', '夏芷', '默闻', '炽热', '空城', '柠凉', '忆海', '知南', '暧昧', '北丧', '予别', '厌世', '会傲', '秋水', '东渡', '忆梦', '挽与', '荡羕', '别离', '枕梦', '暮想', '妄图', '甜屁', '薄凉', '自私', '沙鹭',
             '若暮', '触觉', '安瑾', '倾心', '罪过', '不弃', '素颜', '寄风', '青尤', '侯乔', '告白', '执著', '温书', '顾妄', '夏芷', '博裕', '清引', '忆白', '风云', '空船', '看清', '烟冷', '甜甜', '阁楼', '浅唱', '琼花', '芷芹', '桃酥',
             '千鸢', '纵情', '心痛', '流走', '浪友', '柠澈', '华衣', '凌春', '柒安', '自嘲', '奈奈', '不败', '轻吟', '栀蓝', '局外', '绝恋', '山谷', '静谧', '讨好', '棱人', '萝莉', '北觅', '倚梅', '崆野', '浅岸', '甜酒', '宛栀', '旭明',
             '清风', '喵叽', '梦徒', '安辞', '难喻', '桔栀', '記忔', '北挽', '酒薄', '抵年', '浮笙', '丑角', '暮余', '此月', '酸奶', '贫僧', '酣眠', '明瑞', '阳夏', '拽爷']

# 数据库配置
DATABASE_HOST = "localhost"  # 数据库主机
DATABASE_PORT = 3306  # 数据库端口
DATABASE_NAME = "open_tree_hole"  # 数据库名称
DATABASE_USER = ""  # 数据库用户
DATABASE_PASSWORD = ""  # 数据库密码
REDIS_ADDRESS = 'redis://localhost:6379'  # redis 缓存地址

# 邮件配置
EMAIL_HOST = ''
EMAIL_PORT = 587
EMAIL_HOST_USER = ''
EMAIL_HOST_PASSWORD = ''
EMAIL_USE_TLS = True
EMAIL_USE_SSL = True
DEFAULT_FROM_EMAIL = ''  # 默认发件人地址

# 图片配置
MAX_IMAGE_SIZE = 20  # 最大上传图片大小（MB）
# 采用 Github 图床
GITHUB_OWENER = 'OpenTreeHole'
GITHUB_TOKEN = ''
GITHUB_REPO = ''
GITHUB_BRANCH = ''
SECRET_KEY = str(uuid.uuid1())  # 足够长的密码，供 Django 安全机制

# 用环境变量中的配置覆盖
envs = os.environ
local = locals().copy()  # 拷贝一份，否则运行时 locals() 会改变
for item in local:
    if item.startswith('_') or item == 'os':  # 内置变量名不考虑
        continue
    if item in envs:
        try:
            exec(f'{item} = eval(envs.get(item))')  # 非字符串类型使用 eval() 转换
        except NameError:
            exec(f'{item} = envs.get(item)')  # 否则直接为字符串
