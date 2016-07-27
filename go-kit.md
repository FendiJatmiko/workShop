gokit adalah salah satu tool programing terdistribusi untuk membangun microservices, 
    yang berguna untuk menyelesaikan permasalahan yang biasanya muncul ketika membangun 
    sistem terdistribusi, sehingga programer dapat lebih fokus pada proses bisnis.
sedikit sejarah mengenai microservices => secara tradisional, aplikasi web dibangun menggunakan 
    pendekatan monolithic dimana seluruh aplikasi dibuat, didesain, dideploy dan di maintain secara keseluruhan
        contohnya adalah aplikasi wordpress. Wordpress merupakan contoh paling mudah untuk menggambarkan arsitektur 
        aplikasi yang bersifat monolitik, dimana dalam satu aplikasi kita dapat memiliki frontend sekaligus backend.
        semua Fitur security, performance, manajemen, konten, ad, statistik, semuanya dibangun dengan menggunakan
        PHP dan database MYSQL dalam source code yang sama.
    pada contoh aplikasi web diatas mau tidak mau akan timbul banyak masalah umum dari waktu ke waktu, misalnya
    cukup mudah untuk layer abstak yang bocor antar modul, kemudian bagian aplikasi yang mungkin memerlukan 
    power yang berbeda satu sama lain yang memaksa developer untuk men scale, membanung dan mendeploy ulang seluruh aplikasi
        Aplikasi berbasis microservices bertujuan untuk mengatasi masalah diatas.

      
        Ketika aplikasi monolitik memiliki satu source logs, satu source metrik, satu aplikasi untuk di deploy 
        dan satu API untuk di rate limit,sedangkan microservis memiliki beberapa sources, bebrapa hal yang menjadi
        fokus utama microservis antara lain :
        1. rate limiters yaitu memberi limit pada batas threshold atas
        2. Serialization adalah konversi antara bahasa struktur data ke byte stream untuk mempresentasikan ke sistem yang lain.
            sistem yang lain tersebut biasanya berupa browser(json/xml/html) atau database
        3. Logging adalah sebuah catatan output dari aplikasi yang terstruktur yang diurutkan berdasarkan waktu
        4. Metrics adalah catatan yang di instrumenkan dari bagian aplikasi yang termasuk jarak latency, perhitungan request, sistem health dan lain lain. 
        5. Circuit Breakers untuk mencegah resilency terhadap kerusakan intermittent
        6. Request Tracing untuk mendiagnosa masalah antar servis dan membuat ulang sistem secara keseluruhan
        7. Services Discovery memungkinkan servis yang berbeda - beda untuk menemukan satu sama lain

go-kit adalah salah satu tools untuk menangani permasalahan diatas, terdiri dari beberapa abstraksi dan di encode 
menjadi paket - paket yang menyediakan seperangkat interfaces yang dibutuhkan developer.
    Berikut adalah contoh program untuk menghitung total integer yang dimasukan dan akan direspon dengan
    total nya menggunakan go-kit. 
    Pertama kita membuat interface yang akan digunakan.
~~~    
    type Counter interface {
	Add(int) int
	}	
~~~
dengan implementasi sebagai berikut : 
~~~	
	type countService struct {
		v  int
		mu sync.Mutex
	}
	func (c *countService) Add(v int) int {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.v += v
		return c.v
	}
~~~	
kemudian endpoint pada go-kit merepresentasikan satu RPC dan menjadi dasar untuk pembuatan klien - server,
berikut contoh implementasi endpoint pada go-kit  :
~~~	
	type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
~~~	
kemudian adapter pada go-kit memungkinkan struct yang diimplementasikan kepada satu interfaces 
yang digunakan ketika membutuhkan interface lainnya, jadi pada kasus ini adapter akan digunakan sebagai endpoint
untuk menghandle request dan respon akan menggunakan fungsi bernama encoder/encoder.
~~~
	func decodeAddRequest(r *http.Request) (interface{}, error) {
		var req addRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return req, nil
	}

	func encodeResponse(w http.ResponseWriter, response interface{}) error {
		return json.NewEncoder(w).Encode(response)
	}
~~~	
dan untuk http transport go-kit :
~~~
	func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request)
~~~
yang selanjutnya menggunakan fungsi build-in
~~~ 
		decodeAddRequest,
		encodeResponse,
~~~
metode ini akan diencode oleh
~~~
		return json.NewEncoder(w).Encode(response)
~~~
dan dilanjutkan pada
~~~
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
~~~
dan didecode hasil dari JSON.


Kemudian menggunakan Middleware yang akan digunakan untuk menggabungkan endpoint dan juga menambah generik endpoint
	pada setiap requestnya, yaitu dengan cara membuat fungsi 
	yang mengambil endpoint dan memiliki return value yang berupa endpoint yang baru.
	Middleware akan digunakan untuk proses <i>logging</i>, <i>metrics</i> dan <i>rate limitting</i>, <i>rate limitting</i> akan
	menggunakan <i>juju's token bucket</i> dan <i>go-kit rate limitting middleware</i>
Kemudian akan menggunakan servis metrics middleware 
~~~
func metricsMiddleware(requestCount metrics.Counter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			requestCount.Add(1)
			return next(ctx, request)
		}
	}
}
~~~
yang berfungsi untuk menambahkan nilai satu pada setiap reques(request.count),
yang akan diekspos pada servis build-in "expvar" dan dipanggil menggunakan. 
~~~
requestCount := expvar.NewCounter("request.count")
~~~
pada main program
yang terakhir logging Middleware yang akan memberikan output reques path, reques id, request dan respon data
Logging middleware yang mengandalkan extrator function
~~~
func beforeIDExtractor(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, requestIDKey, r.Header.Get("X-Request-Id"))
}

func beforePATHExtractor(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, pathKey, r.URL.EscapedPath())
}
~~~
 yang berguna untuk mengextract path dan request-id dari 
*http-Request kemudian menambahkannya pada context yang telah disediakan oleh Middleware dan endpoint.
untuk melihat hasil program diatas dijalankan, dapat 
