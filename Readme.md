# API Documentation

## 1. **User Routes**

### **POST /register**
**Description**: Register a new user.  
**Request Body**:
- `email`: string (required)
- `password`: string (required)
- `gender`: string (required)
- `role`: string (optional, default: `user`)

**Response**:
- Status Code: 201 (Created)
  - Body: `{ "message": "User registered successfully" }`
- Status Code: 400 (Bad Request)
  - Body: `{ "error": "Bad request", "message": "Error description" }`

---

### **POST /login**
**Description**: Login an existing user and get a JWT token.  
**Request Body**:
- `email`: string (required)
- `password`: string (required)

**Response**:
- Status Code: 200 (OK)
  - Body: `{ "message": "Login successful", "token": "JWT_TOKEN_HERE" }`
- Status Code: 400 (Bad Request)
  - Body: `{ "error": "Bad request", "message": "Error description" }`
- Status Code: 401 (Unauthorized)
  - Body: `{ "error": "Unauthorized", "message": "Email or password is invalid" }`

---

## 2. **Movie Routes**

### **POST /movies**
**Description**: Admin can create a new movie.  
**Request Body**:
- `title`: string (required)
- `genre`: string (required)
- `director`: string (optional)
- `release_date`: string (required, format: YYYY-MM-DD)
- `description`: string (optional)

**Response**:
- Status Code: 201 (Created)
  - Body: `{ "message": "Movie created successfully" }`
- Status Code: 400 (Bad Request)
  - Body: `{ "error": "Bad request", "message": "Error description" }`
- Status Code: 403 (Forbidden)
  - Body: `{ "error": "Forbidden", "message": "Access denied" }`

---

### **PUT /movies/:movie_id**
**Description**: Admin can update an existing movie by its ID.  
**Request Body**:
- `title`: string (optional)
- `genre`: string (optional)
- `director`: string (optional)
- `release_date`: string (optional, format: YYYY-MM-DD)
- `description`: string (optional)

**Response**:
- Status Code: 200 (OK)
  - Body: `{ "message": "Movie updated successfully" }`
- Status Code: 400 (Bad Request)
  - Body: `{ "error": "Bad request", "message": "Error description" }`
- Status Code: 404 (Not Found)
  - Body: `{ "error": "Not Found", "message": "Movie not found" }`
- Status Code: 403 (Forbidden)
  - Body: `{ "error": "Forbidden", "message": "Access denied" }`

---

### **GET /movies**
**Description**: Public route to get a list of all movies with pagination.  
**Query Parameters**:
- `page`: int (optional, default: 1)
- `limit`: int (optional, default: 10)

**Response**:
- Status Code: 200 (OK)
  - Body: `{ "movies": [...], "page": 1, "total": 100 }`

---

### **GET /movies/search**
**Description**: Public route to search movies by query string.  
**Query Parameters**:
- `query`: string (required)

**Response**:
- Status Code: 200 (OK)
  - Body: `{ "movies": [...] }`

---

## 3. **Stats Routes**

### **POST /stats/:movie_id/view**
**Description**: Authenticated users can track views for a movie.  
**Request Body**:
- `user_id`: int (required)

**Response**:
- Status Code: 200 (OK)
  - Body: `{ "message": "View tracked successfully" }`
- Status Code: 401 (Unauthorized)
  - Body: `{ "error": "Unauthorized", "message": "User not authenticated" }`

---

### **POST /stats/:movie_id/vote**
**Description**: Authenticated users with the "user" role can vote for a movie.  
**Request Body**:
- `user_id`: int (required)
- `rating`: int (required, range: 1-5)

**Response**:
- Status Code: 200 (OK)
  - Body: `{ "message": "Vote recorded successfully" }`
- Status Code: 401 (Unauthorized)
  - Body: `{ "error": "Unauthorized", "message": "User not authenticated" }`
- Status Code: 403 (Forbidden)
  - Body: `{ "error": "Forbidden", "message": "Only users with role 'user' can vote" }`

---

### **POST /stats/:movie_id/unvote**
**Description**: Authenticated users with the "user" role can unvote for a movie.  
**Response**:
- Status Code: 200 (OK)
  - Body: `{ "message": "Vote removed successfully" }`
- Status Code: 401 (Unauthorized)
  - Body: `{ "error": "Unauthorized", "message": "User not authenticated" }`
- Status Code: 403 (Forbidden)
  - Body: `{ "error": "Forbidden", "message": "Only users with role 'user' can unvote" }`

---

### **GET /stats/most-viewed-genre-movie**
**Description**: Admin can view the most viewed genre of movies.  
**Response**:
- Status Code: 200 (OK)
  - Body: `{ "genre": "Action", "view_count": 1500 }`
- Status Code: 403 (Forbidden)
  - Body: `{ "error": "Forbidden", "message": "Access denied" }`

---

### **GET /stats/most-voted-genre-movie**
**Description**: Admin can view the most voted genre of movies.  
**Response**:
- Status Code: 200 (OK)
  - Body: `{ "genre": "Drama", "vote_count": 1200 }`
- Status Code: 403 (Forbidden)
  - Body: `{ "error": "Forbidden", "message": "Access denied" }`

---

### **POST /stats/:movie_id/trace**
**Description**: Authenticated users with the "user" role can trace viewership of a movie.  
**Request Body**:
- `trace_id`: string (required)

**Response**:
- Status Code: 200 (OK)
  - Body: `{ "message": "Viewership trace completed" }`
- Status Code: 401 (Unauthorized)
  - Body: `{ "error": "Unauthorized", "message": "User not authenticated" }`
- Status Code: 403 (Forbidden)
  - Body: `{ "error": "Forbidden", "message": "Only users with role 'user' can trace" }`

---

### **GET /stats/user/voted-movies**
**Description**: Authenticated users can view the movies they have voted for.  
**Response**:
- Status Code: 200 (OK)
  - Body: `{ "voted_movies": [{ "movie_id": 1, "title": "Movie Title", "vote": 5 }, ...] }`
- Status Code: 401 (Unauthorized)
  - Body: `{ "error": "Unauthorized", "message": "User not authenticated" }`

---

## Middleware

### **AuthMiddleware**
**Description**: Ensures that the user is authenticated before accessing protected routes.

**Response**:
- Status Code: 401 (Unauthorized)
  - Body: `{ "error": "Unauthorized", "message": "User not authenticated" }`

### **RoleMiddleware(role string)**
**Description**: Ensures that the user has the required role (e.g., "admin", "user") before accessing protected routes.

**Response**:
- Status Code: 403 (Forbidden)
  - Body: `{ "error": "Forbidden", "message": "Access denied" }`

---

## Notes
- All routes that require authentication will need a valid JWT token passed in the `Authorization` header as a Bearer token.
- Only users with the appropriate roles (e.g., "admin", "user") will be able to access certain routes.

---

Dokumentasi ini memberikan gambaran lengkap mengenai penggunaan API untuk aplikasi Anda, termasuk endpoint yang tersedia, format data, dan status code yang dapat dikembalikan.