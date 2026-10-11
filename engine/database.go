package engine

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "modernc.org/sqlite"
)

type TickLog struct {
	TickNo    int
	Category  string
	Price     float64
	Inflation float64
	Demand    int
	Supply    int
}

type AgentTransactionLog struct {
	TickNo     int
	AgentName  string
	ItemName   string
	Category   string
	Price      float64
	NewBalance float64
}

type DBManager struct {
	db          *sql.DB
	tickChan    chan TickLog
	txChan      chan AgentTransactionLog
	wg          sync.WaitGroup
	closeSignal chan struct{}
}

func NewDBManager(dbPath string) (*DBManager, error) {
	database, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("veritabanı açılamadı: %v", err)
	}

	mgr := &DBManager{
		db:          database,
		tickChan:    make(chan TickLog, 500),
		txChan:      make(chan AgentTransactionLog, 500),
		closeSignal: make(chan struct{}),
	}

	if err := mgr.initTables(); err != nil {
		return nil, err
	}

	mgr.wg.Add(1)
	go mgr.startWorker()

	return mgr, nil
}

func (m *DBManager) initTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS tick_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tick_no INTEGER,
			category TEXT,
			price REAL,
			inflation REAL,
			demand INTEGER,
			supply INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS agent_transactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tick_no INTEGER,
			agent_name TEXT,
			item_name TEXT,
			category TEXT,
			price REAL,
			new_balance REAL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, q := range queries {
		if _, err := m.db.Exec(q); err != nil {
			return fmt.Errorf("tablo oluşturma hatası: %v", err)
		}
	}
	return nil
}

func (m *DBManager) LogTick(logData TickLog) {
	m.tickChan <- logData
}

func (m *DBManager) LogTransaction(txData AgentTransactionLog) {
	m.txChan <- txData
}

func (m *DBManager) startWorker() {
	defer m.wg.Done()

	for {
		select {
		case tLog := <-m.tickChan:
			m.insertTickLog(tLog)
		case txLog := <-m.txChan:
			m.insertTransactionLog(txLog)
		case <-m.closeSignal:
			for len(m.tickChan) > 0 {
				m.insertTickLog(<-m.tickChan)
			}
			for len(m.txChan) > 0 {
				m.insertTransactionLog(<-m.txChan)
			}
			return
		}
	}
}

func (m *DBManager) insertTickLog(l TickLog) {
	query := `INSERT INTO tick_logs (tick_no, category, price, inflation, demand, supply) VALUES (?, ?, ?, ?, ?, ?)`
	if _, err := m.db.Exec(query, l.TickNo, l.Category, l.Price, l.Inflation, l.Demand, l.Supply); err != nil {
		log.Printf("[DB HATA] Tick kaydı atılamadı: %v", err)
	}
}

func (m *DBManager) insertTransactionLog(l AgentTransactionLog) {
	query := `INSERT INTO agent_transactions (tick_no, agent_name, item_name, category, price, new_balance) VALUES (?, ?, ?, ?, ?, ?)`
	if _, err := m.db.Exec(query, l.TickNo, l.AgentName, l.ItemName, l.Category, l.Price, l.NewBalance); err != nil {
		log.Printf("[DB HATA] İşlem kaydı atılamadı: %v", err)
	}
}

func (m *DBManager) Close() {
	close(m.closeSignal)
	m.wg.Wait()
	m.db.Close()
	fmt.Println("💾 [DB] Tüm veriler SQLite veritabanına kaydedildi ve bağlantı kapatıldı.")
}