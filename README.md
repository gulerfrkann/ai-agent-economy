# 🚀 Autonomous AI Agent Economy Simulation

Go dilinin yüksek eşzamanlılık (concurrency) gücünü ve gerçek e-ticaret verilerini birleştiren otonom ajan simülasyon motoru.

---

##  Projenin Amacı ve Mimarisi
Bu proje, bağımsız yapay zeka ajanlarının (Traders, Merchants) bir pazar yeri ortamında eşzamanlı olarak etkileşime girdiği, bütçe yönetimi, gelir döngüsü ve kategori uzmanlığına dayalı kararlar aldığı modüler bir simülasyon altyapısıdır. 

* **Veri Katmanı (Python):** Gerçek dünya e-ticaret verileri (`.parquet` formatında, ~965k+ ürün) Pandas yardımıyla işlenir ve Go motorunun hızlıca tüketebileceği optimize edilmiş bir yapıya dönüştürülür.
* **Simülasyon Motoru (Go):** Goroutines, Channels, Tickers, Mutex ve `context` tabanlı **Graceful Shutdown** mimarisi kullanılarak sıfır veri kaybı ve yüksek performansla çalışacak şekilde tasarlanmıştır.

---

##  Temel Özellikler ve Ekonomik Mekanizmalar
1. **Akıllı Ürün Seçimi (Niche Interests):** Ajanlar rastgele harcama yapmak yerine, kendi uzmanlık alanlarına ve ilgi duydukları kategorilere (`Interests`) öncelik vererek pazar yerinden ürün arar.
2. **Gelir Döngüsü (Income / Earn):** Her tick (tur) başında ajanlara taban gelir (örn. +150 TL) eklenerek sistemin tıkanması önlenir ve sürekli bir ekonomik çark oluşturulur.
3. **Bakiye & Envanter Yönetimi:** Ajanlar anlık bakiyelerini kontrol eder; bütçelerine uygun ürünleri envanterlerine ekler, yetersiz durumlarda ise akıllıca pas geçerler.
4. **Detaylı Final Raporlaması:** Simülasyon sonlandırıldığında (`CTRL + C`) ajanların kalan bakiyelerini, envanterlerindeki ürün sayılarını ve aldıkları başlıca ürünleri özetleyen şık bir rapor sunulur.

---

##  Kullanılan Teknolojiler

* **Dil:** Go (Golang), Python (ETL ve Veri Dönüşümü)
* **Veri Formatı:** Parquet, JSON (Yerel `.gitignore` korumalı)
* **Mimari Desenler:** Concurrency, Ticker Loop, Interface-based Agent Design, Graceful Shutdown, Ledger & Inventory Management

---

##  Proje Yapısı

```text
ai-agent-economy/
├── data/                 # Büyük veri setleri ve JSON çıktıları (.gitignore ile korunur)
├── engine/               # Simülasyon motoru, ajan tanımları ve pazar yeri mantığı
│   ├── agent.go          # Ajan arayüzü, bütçe, ilgi alanları ve Act davranışı
│   ├── marketplace.go    # 965k+ ürün havuzu yönetimi ve dinamik fiyatlandırma
│   └── tick.go           # Zaman (Tick) motoru ve eşzamanlı tetikleme
├── convert.py            # Parquet -> JSON veri dönüştürme betiği
├── main.go               # Uygulama giriş noktası, simülasyon yönetimi ve final raporlama
├── go.mod                # Go modül bağımlılıkları
└── README.md             # Proje dokümantasyonu