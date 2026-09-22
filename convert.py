import pandas as pd
import os

# data klasöründeki parquet dosyasının yolunu belirt (dosya adını kendi dosya adına göre kontrol edebilirsin)
parquet_path = os.path.join("data", "stajyer_tavsiye_sistem.parquet")
json_output_path = os.path.join("data", "products.json")

print(" Parquet dosyası okunuyor...")
df = pd.read_parquet(parquet_path)

# Veriyi JSON formatına çevir ve data klasörüne kaydet
df.to_json(json_output_path, orient="records", force_ascii=False, indent=4)

print(f" Başarıyla dönüştürüldü! Dosya konumu: {json_output_path}")
print(f"Toplam ürün/satır sayısı: {len(df)}")
print(f" Sütunlar: {list(df.columns)}")