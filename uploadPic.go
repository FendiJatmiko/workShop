package main
  import (
        "fmt"
        "image"
        "image/jpeg"
        "net/http"
  )

  var (
    images = make(map[string] image.Image)
  )
  //funsi untuk mengolah html syntax
  func HandleRoot(w http.ResponseWriter, r *http.Request){
      fmt.Fprintf(w, `
          <html><body>
              <form method="post" enctype="multipart/form-data" action="/upload"
              name="upload">
                  <label for="file">Choose image:</label>
                  <input type="file" name="image"/> </br>
                  <input type="submit" name="submit" value=Upload/>
                </form>
            </body></html>
        `)
  }
  //fungsi untuk menangani proses upload image
  func HandleUpload( w http.ResponseWriter, r *http.Request) {
      file, header, _:= r.FormFile("image ")
      image, _, _ := image.Decode(file)
      images[header.Filename] = image
      http.Redirect(w, r, "/image?name="+header.Filename, 303)
  }
//encode image dengan hanya jpeg imange
func HandleImange (w http.ResponseWriter, r *http.Request){
    imageName := r.FormValue("name")
    image     := images[imageName]
    jpeg.Encode(w, image, &jpeg.Options{Quality:
      jpeg.DefaultQuality })
}
  func main(){
    http.HandleFunc("/", HandleRoot)
    http.HandleFunc("/upload", HandleUpload)
    http.HandleFunc("/image", HandleImange)
    http.ListenAndServe(":8000", nil)
  }
