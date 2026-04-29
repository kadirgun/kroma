### Document idle kapatma

onvisiblitychange devredışı bırakıldı. document.visibilityState her zaman visible dönecek.

### Canvas noise

—canvas-noise=12564565454546

### WebGL noise

—canvas-noise=12564565454546

### WebRTC spoofing

—visitor-ip=181.12.15.35

### Video Codecs

proprietary_codecs = true

### Navigator Permission

accelerometer|background-sync|camera|clipboard-read|clipboard-write|geolocation|gyroscope|magnetometer|microphone|midi|notifications|payment-handler|persistent-storage bunlar olmalı

### WebGL Vendor/Renderer

—webgl-vendor —webgl-renderer

### Speech Voices

GOOGLE_CHROME_BRANDING sorgusu kaldırıldı

### GPUAdapterInfo

—gpu-vendor —gpu-arch

### userAgentData Chrome ve Tam versiyon

args.gn branding_file_path ayarlandı.

### Timezone

playwright timezoneId

### Hardware Concurrency

--cpu-cores=32

### DNS over HTTPS

--doh-template=”https://cloudflare-dns.com/dns-query%7B?dns%7D”

### Accept-Language Header

--accept-lang=en,en-US,tr,tr-TR

### Inactive Window

--start-maximized verince tüm pencereler inactive başlayacak

### AppUserModelId

açılan pencereler aynı simge grubunda toplanacak

### Webgl Enthropy

artık webgl ise optypes kontrolü yapılmıyor. renk sayısı yüksekse gürültü ekleniyor.

### Useragent version

--user-agent-version=99.55.66.33

### DeviceMemory

yeni versiyonda 32 GB'a kadar gösteriyor. 8,16,32 şeklinde --device-memory-gb ile verilecek.

### Input.scheduleMouseEvent request

Input.scheduleMouseEvent ile event listesi gönderiliyor. Eventları gönderdikten sonra hemen cevap geliyor. süre hesabı yapıp beklemek gerek.

### Input.dispatchKeyEvent modifiers AltGraph ve Numlock

AltGraph: 16, NumLock: 32

### WebAuthn auto cancel

--cancel-webauthn bayrağı ile webathn register istekleri otomatik iptal ediliyor
