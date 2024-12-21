# Treehole Image Hosting

---
## Upload Image

To upload an image, send a `POST` request to `{hostname}/api/uploadImage`. Include the photo in the request body with the form-data field named "source".

### Example:
    http://localhost:8000/api/uploadImage

## Get Image

To retrieve an image, use the `GET` method with the following endpoint: `{hostname}/api/i/:year/:month/:day/:identifier`.

### Example:
    http://localhost:8000/api/i/2024/12/06/6288772352016bf28f1a571d0.jpg

## System Environment Variables
- DB_URL
  - root:`PASSWORD`@tcp(localhost:`PORT`)/`YOUR_DATABASE_NAME`?parseTime=true&loc=Asia%2fShanghai
- HOST_NAME (change it to your own host, which directly exports to users)
  - http://localhost:8000
  - https://image.fduhole.com
---
