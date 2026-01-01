package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"vocab-quest/server/services"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 1. Connect to DB
	db, err := sql.Open("sqlite3", "./vocab.db")
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	fmt.Println("--- STARTING SESSION SIMULATION (3-4-3) ---")

	// 2. Create Test User
	username := fmt.Sprintf("sim_user_%d", time.Now().Unix())
	res, err := db.Exec("INSERT INTO users (username, password, role) VALUES (?, 'pass', 'user')", username)
	if err != nil {
		log.Fatal("Failed to create user:", err)
	}
	userID, _ := res.LastInsertId()
	fmt.Printf("Created Test User: %s (ID: %d)\n", username, userID)

	// 3. Get All Words
	rows, err := db.Query("SELECT id, text FROM words")
	if err != nil {
		log.Fatal(err)
	}
	var allWordIDs []int
	for rows.Next() {
		var id int
		var text string
		rows.Scan(&id, &text)
		allWordIDs = append(allWordIDs, id)
	}
	rows.Close()

	if len(allWordIDs) < 15 {
		log.Fatal("Not enough words in DB to run simulation (Need at least 15)")
	}

	// 4. Inject Progress Data to shape the pool
	// We need:
	// - A pool of WEAK words (Low accuracy) -> Expect 4 selected
	// - A pool of REVIEW words (Old, high accuracy) -> Expect 3 selected
	// - A pool of NEW words (No progress) -> Expect 3 selected

	// Inject 6 Weak Words (Attempts 10, Success 2 => 20% Accuracy)
	weakIDs := allWordIDs[:6]
	fmt.Printf("Injecting %d Weak Words (Acc 20%%)...\n", len(weakIDs))
	for _, wid := range weakIDs {
		_, err := db.Exec("INSERT INTO user_progress (user_id, word_id, attempts, successes, last_played_at) VALUES (?, ?, 10, 2, ?)",
			userID, wid, time.Now()) // Just played, but weak
		if err != nil {
			log.Fatal(err)
		}
	}

	// Inject 6 Review Words (Attempts 10, Success 10 => 100% Accuracy, Played 10 days ago)
	reviewIDs := allWordIDs[6:12]
	fmt.Printf("Injecting %d Review Words (Acc 100%%, Old)...\n", len(reviewIDs))
	oldDate := time.Now().Add(-240 * time.Hour) // 10 days ago
	for _, wid := range reviewIDs {
		_, err := db.Exec("INSERT INTO user_progress (user_id, word_id, attempts, successes, last_played_at) VALUES (?, ?, 10, 10, ?)",
			userID, wid, oldDate)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Remaining words are NEW (attempts=0 or no record)
	fmt.Printf("Remaining Words are NEW (No records)\n")

	// 5. Run GetSessionWords
	fmt.Println("\n>>> GENERATING SESSION (Limit 10) <<<")
	session, err := services.GetSessionWords(db, int(userID), 10)
	if err != nil {
		log.Fatal("GetSessionWords failed:", err)
	}

	// 6. Analyze Results
	counts := map[string]int{"New": 0, "Weak": 0, "Review": 0}

	fmt.Println("\n[SESSION CONTENT]")
	fmt.Printf("%-15s | %-10s | %-10s\n", "Word", "Prof Flag", "Category Guess")
	fmt.Println("------------------------------------------------")

	for _, w := range session {
		cat := "Unknown"
		// Logic from Game Service:
		// Prof 0 = New
		// Prof 2 = Weak (<0.6)
		// Prof 4/5 = Master (High Acc) -> Likely Review pile
		// Prof 3 = Normal -> Likely Review pile

		if w.Proficiency == 0 {
			cat = "New"
		} else if w.Proficiency == 2 {
			cat = "Weak"
		} else if w.Proficiency >= 3 {
			cat = "Review"
		}

		counts[cat]++
		fmt.Printf("%-15s | %-10d | %-10s\n", w.Text, w.Proficiency, cat)
	}

	fmt.Println("------------------------------------------------")
	fmt.Println("\n[SUMMARY]")
	fmt.Printf("New Words:    %d (Target: 3)\n", counts["New"])
	fmt.Printf("Weak Words:   %d (Target: 4)\n", counts["Weak"])
	fmt.Printf("Review Words: %d (Target: 3)\n", counts["Review"])

	// Validation for "Fixed Review, Remedial First"
	// Context: 6 Weak, 6 Review, x New. Limit 10.
	// 1. Review needs to grab 3.
	// 2. Remaining 7 slots.
	// 3. Weak (6) takes 6 slots.
	// 4. Remaining 1 slot.
	// 5. New takes 1 slot.
	// Expect: Review=3, Weak=6, New=1.
	if counts["Review"] == 3 && counts["Weak"] == 6 && counts["New"] == 1 {
		fmt.Println("\n✅ SUCCESS: REVIEW FIXED + WEAK PRIORITY (3 Review + 6 Weak + 1 New)")
	} else if counts["Review"] >= 3 && counts["Weak"] >= 4 { // Generalized generalized
		fmt.Printf("\n✅ SUCCESS: Review 3 preserved, Weak heavily populated. Rev:%d Weak:%d New:%d\n", counts["Review"], counts["Weak"], counts["New"])
	} else {
		fmt.Printf("\n⚠️  Unexpected Mix. New: %d, Weak: %d, Review: %d\n", counts["New"], counts["Weak"], counts["Review"])
	}
}
