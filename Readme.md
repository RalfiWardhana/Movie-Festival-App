# Movie API Documentation

## Base URL
`http://localhost:9191/api/v1`

---

## Authentication

### Bearer Token
Beberapa endpoint memerlukan autentikasi dengan **Bearer Token**. Tambahkan header berikut ke setiap permintaan yang membutuhkan autentikasi:
```json
{
  "Authorization": "Bearer <your_token>"
}
```

---

## Endpoints

### **User**

#### Register User or Admin
**POST** `/user/register`

**Request Body:**
```json
{
  "email": "queen@mail.com",
  "password": "12345",
  "gender": "Perempuan",
  "role": "user"
}
```

#### Login
**POST** `/user/login`

**Request Body:**
```json
{
  "email": "ralfi@mail.com",
  "password": "ralfi789"
}
```

---

### **Movie**

#### Create and Upload Movie
**POST** `/movies`

**Authorization:** Required (Bearer Token)

**Request Form Data:**
- `title`: (string) Movie title
- `description`: (string) Movie description
- `duration`: (string) Duration (e.g., "2 jam")
- `artist`: (string) List of actors
- `genre_id`: (integer) Genre ID
- `file`: (file) Movie file

---

#### Update Movie
**PUT** `/movies/:id`

**Authorization:** Required (Bearer Token)

**Request Form Data:**
- `title`: (string) Updated title

---

#### List Movies with Pagination
**GET** `/movies`

**Query Parameters:**
- `page`: (integer, optional) Page number (default: 1)
- `limit`: (integer, optional) Items per page (default: 10)

---

#### Search Movies
**GET** `/movies/search`

**Query Parameters:**
- `title`: (string) Title search keyword
- `artist`: (string) Artist search keyword

---

### **Vote and View**

#### Track Movie Viewership
**POST** `/stats/:movieId/view`

**Authorization:** Required (Bearer Token)

---

#### Vote a Movie
**POST** `/stats/:movieId/vote`

**Authorization:** Required (Bearer Token)

---

#### Unvote a Movie
**POST** `/stats/:movieId/unvote`

**Authorization:** Required (Bearer Token)

---

#### Most Viewed Movie and Genre
**GET** `/stats/most-viewed-genre-movie`

**Authorization:** Required (Bearer Token)

---

#### Most Voted Movie and Genre (Admin Only)
**GET** `/stats/most-voted-genre-movie`

**Authorization:** Required (Bearer Token)

---
```

