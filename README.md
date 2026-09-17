#  Autonomous AI Agent Economy Simulation

Go dilinin yüksek eşzamanlılık (concurrency) gücünü ve veri bilimi dünyasının gerçek e-ticaret verilerini birleştiren otonom ajan simülasyon motoru.

---

##  Projenin Amacı ve Mimarisi
Bu proje, bağımsız yapay zeka ajanlarının (Traders, Merchants) bir pazar yeri ortamında eşzamanlı olarak etkileşime girdiği, ürün incelediği ve gelecekte ekonomik kararlar alacağı modüler bir simülasyon altyapısıdır. 

* **Veri Katmanı (Python):** Gerçek dünya e-ticaret verileri (`.parquet` formatında, ~965k ürün) Pandas yardımıyla işlenir ve Go motorunun hızlıca tüketebileceği hafif bir yapıya dönüştürülür.
* **Simülasyon Motoru (Go):** Goroutines, Channels, Tickers ve `context` tabanlı **Graceful Shutdown** mimarisi kullanılarak sıfır veri kaybı ve yüksek performansla çalışacak şekilde tasarlanmıştır.

---

##  Kullanılan Teknolojiler

* **Dil:** Go (Golang), Python (ETL ve Veri Dönüşümü)
* **Veri Formatı:** Parquet, JSON (Yerel `.gitignore` korumalı)
* **Mimari Desenler:** Concurrency (Goroutine/Mutex), Ticker Loop, Interface-based Agent Design, Graceful Shutdown

---

##  Proje Yapısı

```text
ai-agent-economy/
├── data/                 # Büyük veri setleri ve JSON çıktıları (.gitignore ile korunur)
├── engine/               # Simülasyon motoru, ajan tanımları ve pazar yeri mantığı
│   ├── agent.go          # Ajan arayüzü ve davranışları (Act)
│   ├── marketplace.go    # 965k+ ürün havuzu yönetimi
│   └── tick.go           # Zaman (Tick) motoru ve eşzamanlı tetikleme
├── convert.py            # Parquet -> JSON veri dönüştürme betiği
├── main.go               # Uygulama giriş noktası ve sinyal yakalama (Graceful Shutdown)
├── go.mod                # Go modül bağımlılıkları
└── README.md             # Proje dokümantasyonu