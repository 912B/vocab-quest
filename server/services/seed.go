package services

import (
	"database/sql"
	"log"
	"vocab-quest/server/models"
)

func SeedDatabase(db *sql.DB) error {
	log.Println("Checking database seeds...")

	// 1. Check & Seed Dictionary
	var dictID int
	dictName := "Grade 5 Vocabulary"
	err := db.QueryRow("SELECT id FROM dictionaries WHERE name = ?", dictName).Scan(&dictID)

	if err == sql.ErrNoRows {
		log.Println("Creating Grade 5 Dictionary...")
		res, err := db.Exec(`INSERT INTO dictionaries (name, description, is_active) VALUES (?, ?, ?)`,
			dictName, "小学五年级上/下册核心词汇表 (PEP标准)", true)
		if err != nil {
			return err
		}
		id, _ := res.LastInsertId()
		dictID = int(id)
	} else if err != nil {
		return err
	}

	// Add Initial Words (Standard Grade 5 School Exam Vocabulary - PEP Edition)
	initialWords := []models.Word{
		// Unit 1: My Day (我的作息)
		{Text: "exercise", Definition: "锻炼; 做运动", Difficulty: 1},
		{Text: "eat breakfast", Definition: "吃早饭", Difficulty: 1},
		{Text: "eat dinner", Definition: "吃晚饭", Difficulty: 1},
		{Text: "eat lunch", Definition: "吃午饭", Difficulty: 1},
		{Text: "do morning exercises", Definition: "做早操", Difficulty: 1},
		{Text: "have class", Definition: "上课", Difficulty: 1},
		{Text: "play sports", Definition: "进行体育运动", Difficulty: 1},
		{Text: "clean my room", Definition: "打扫我的房间", Difficulty: 1},
		{Text: "go for a walk", Definition: "散步", Difficulty: 1},
		{Text: "go shopping", Definition: "去购物", Difficulty: 1},
		{Text: "take a dancing class", Definition: "上舞蹈课", Difficulty: 1},
		{Text: "when", Definition: "什么时候", Difficulty: 1},
		{Text: "after", Definition: "在...之后", Difficulty: 1},
		{Text: "start", Definition: "开始", Difficulty: 1},
		{Text: "usually", Definition: "通常", Difficulty: 1},
		{Text: "Spanish", Definition: "西班牙语", Difficulty: 2},
		{Text: "late", Definition: "晚; 迟", Difficulty: 1},
		{Text: "a.m.", Definition: "上午", Difficulty: 1},
		{Text: "p.m.", Definition: "下午", Difficulty: 1},
		{Text: "work", Definition: "工作", Difficulty: 1},
		{Text: "island", Definition: "岛屿", Difficulty: 2},

		// Unit 2: My Favourite Season (我最喜欢的季节)
		{Text: "spring", Definition: "春天", Difficulty: 1},
		{Text: "summer", Definition: "夏天", Difficulty: 1},
		{Text: "autumn", Definition: "秋天", Difficulty: 1},
		{Text: "winter", Definition: "冬天", Difficulty: 1},
		{Text: "season", Definition: "季节", Difficulty: 1},
		{Text: "picnic", Definition: "野餐", Difficulty: 1},
		{Text: "go on a picnic", Definition: "去野餐", Difficulty: 1},
		{Text: "pick apples", Definition: "摘苹果", Difficulty: 1},
		{Text: "snowman", Definition: "雪人", Difficulty: 1},
		{Text: "make a snowman", Definition: "堆雪人", Difficulty: 1},
		{Text: "go swimming", Definition: "去游泳", Difficulty: 1},
		{Text: "which", Definition: "哪一个", Difficulty: 1},
		{Text: "best", Definition: "最好的", Difficulty: 1},
		{Text: "snow", Definition: "雪", Difficulty: 1},
		{Text: "good job", Definition: "做得好", Difficulty: 1},
		{Text: "because", Definition: "因为", Difficulty: 1},
		{Text: "vacation", Definition: "假期", Difficulty: 1},
		{Text: "all", Definition: "全; 完全", Difficulty: 1},
		{Text: "pink", Definition: "粉色", Difficulty: 1},
		{Text: "lovely", Definition: "可爱的", Difficulty: 1},
		{Text: "leaf", Definition: "叶子 (复数 leaves)", Difficulty: 1},
		{Text: "fall", Definition: "落下; 秋天", Difficulty: 1},
		{Text: "paint", Definition: "绘画", Difficulty: 1},

		// Unit 3: School Calendar & Months (校历与月份)
		{Text: "January", Definition: "一月", Difficulty: 1},
		{Text: "February", Definition: "二月", Difficulty: 1},
		{Text: "March", Definition: "三月", Difficulty: 1},
		{Text: "April", Definition: "四月", Difficulty: 1},
		{Text: "May", Definition: "五月", Difficulty: 1},
		{Text: "June", Definition: "六月", Difficulty: 1},
		{Text: "July", Definition: "七月", Difficulty: 1},
		{Text: "August", Definition: "八月", Difficulty: 1},
		{Text: "September", Definition: "九月", Difficulty: 1},
		{Text: "October", Definition: "十月", Difficulty: 1},
		{Text: "November", Definition: "十一月", Difficulty: 1},
		{Text: "December", Definition: "十二月", Difficulty: 1},
		{Text: "few", Definition: "不多; 很少", Difficulty: 1},
		{Text: "party", Definition: "聚会; 派对", Difficulty: 1},
		{Text: "trip", Definition: "旅行", Difficulty: 1},
		{Text: "school trip", Definition: "学校郊游", Difficulty: 1},
		{Text: "sports meet", Definition: "运动会", Difficulty: 1},
		{Text: "Easter", Definition: "复活节", Difficulty: 2},
		{Text: "contest", Definition: "比赛; 竞赛", Difficulty: 2},
		{Text: "Great Wall", Definition: "长城", Difficulty: 1},
		{Text: "RSVP", Definition: "请回复 (Ritpondez s'il vous plait)", Difficulty: 3},

		// Unit 4: When is the Art Show? (序数词与节日)
		{Text: "first", Definition: "第一 (1st)", Difficulty: 1},
		{Text: "second", Definition: "第二 (2nd)", Difficulty: 1},
		{Text: "third", Definition: "第三 (3rd)", Difficulty: 1},
		{Text: "fourth", Definition: "第四 (4th)", Difficulty: 1},
		{Text: "fifth", Definition: "第五 (5th)", Difficulty: 1},
		{Text: "twelfth", Definition: "第十二 (12th)", Difficulty: 2},
		{Text: "twentieth", Definition: "第二十 (20th)", Difficulty: 2},
		{Text: "twenty-first", Definition: "第二十一 (21st)", Difficulty: 2},
		{Text: "thirtieth", Definition: "第三十 (30th)", Difficulty: 2},
		{Text: "special", Definition: "特殊的", Difficulty: 1},
		{Text: "kitten", Definition: "小猫", Difficulty: 1},
		{Text: "diary", Definition: "日记", Difficulty: 2},
		{Text: "make a noise", Definition: "吵闹", Difficulty: 1},
		{Text: "walk", Definition: "走", Difficulty: 1},
		{Text: "fur", Definition: "毛皮", Difficulty: 2},
		{Text: "open", Definition: "开着的", Difficulty: 1},

		// Unit 5: Whose dog is it? (名词性物主代词)
		{Text: "mine", Definition: "我的", Difficulty: 1},
		{Text: "yours", Definition: "你的; 你们的", Difficulty: 1},
		{Text: "his", Definition: "他的", Difficulty: 1},
		{Text: "hers", Definition: "她的", Difficulty: 1},
		{Text: "theirs", Definition: "他们的", Difficulty: 1},
		{Text: "ours", Definition: "我们的", Difficulty: 1},
		{Text: "climbing", Definition: "正在爬", Difficulty: 1},
		{Text: "eating", Definition: "正在吃", Difficulty: 1},
		{Text: "playing", Definition: "正在玩", Difficulty: 1},
		{Text: "jumping", Definition: "正在跳", Difficulty: 1},
		{Text: "drinking", Definition: "正在喝", Difficulty: 1},
		{Text: "sleeping", Definition: "正在睡", Difficulty: 1},
		{Text: "each", Definition: "每一", Difficulty: 1},
		{Text: "other", Definition: "其他", Difficulty: 1},
		{Text: "each other", Definition: "互相", Difficulty: 1},
		{Text: "excited", Definition: "兴奋的", Difficulty: 1},
		{Text: "like", Definition: "像...一样", Difficulty: 1},

		// Unit 6: Work Quietly (进行时与指令)
		{Text: "doing morning exercises", Definition: "正在做早操", Difficulty: 1},
		{Text: "having... class", Definition: "正在上...课", Difficulty: 1},
		{Text: "reading a book", Definition: "正在看书", Difficulty: 1},
		{Text: "listening to music", Definition: "正在听音乐", Difficulty: 1},
		{Text: "keep to the right", Definition: "靠右行", Difficulty: 1},
		{Text: "keep your desk clean", Definition: "保持桌面整洁", Difficulty: 1},
		{Text: "talk quietly", Definition: "小声说话", Difficulty: 1},
		{Text: "take turns", Definition: "按顺序来", Difficulty: 1},
		{Text: "bamboo", Definition: "竹子", Difficulty: 1},
		{Text: "its", Definition: "它的 (指事物/动物)", Difficulty: 1},
		{Text: "show", Definition: "给...看; 展示", Difficulty: 1},
		{Text: "anything", Definition: "任何事物", Difficulty: 1},
		{Text: "else", Definition: "另外; 其他", Difficulty: 1},
		{Text: "exhibition", Definition: "展览", Difficulty: 2},
		{Text: "say", Definition: "说; 讲", Difficulty: 1},
		{Text: "sushi", Definition: "寿司", Difficulty: 1},
		{Text: "teach", Definition: "教", Difficulty: 1},
		{Text: "Canadian", Definition: "加拿大的", Difficulty: 2},
	}

	stmt, err := db.Prepare(`INSERT INTO words (dictionary_id, text, definition, difficulty) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	checkStmt, err := db.Prepare(`SELECT id FROM words WHERE dictionary_id = ? AND text = ?`)
	if err != nil {
		return err
	}
	defer checkStmt.Close()

	for _, w := range initialWords {
		var exists int
		err := checkStmt.QueryRow(dictID, w.Text).Scan(&exists)
		if err == sql.ErrNoRows {
			stmt.Exec(dictID, w.Text, w.Definition, w.Difficulty)
		}
	}

	// 2. Check & Seed Users
	var userCount int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if userCount == 0 {
		log.Println("Seeding Default Users...")
		users := []models.User{
			{Username: "admin", Password: "admin", Role: "admin", Avatar: "assets/icons/rocket.png"}, // Rocket, Planet
			{Username: "Pilot", Password: "🛸👽", Role: "user", Avatar: "assets/icons/ufo.png"},        // UFO, Alien
			{Username: "Engineer", Password: "🔧🔋", Role: "user", Avatar: "assets/icons/wrench.png"},  // Wrench, Battery
		}

		stmtUser, err := db.Prepare("INSERT INTO users (username, password, role, avatar) VALUES (?, ?, ?, ?)")
		if err == nil {
			defer stmtUser.Close()
			for _, u := range users {
				stmtUser.Exec(u.Username, u.Password, u.Role, u.Avatar)
			}
		}
	}

	return nil
}
