# Kushena Backend API

Kushena is a restaurant employee and order management system. This backend API is built using Golang, Gin-Gonic, MongoDB, JWT authentication, and bcrypt for password hashing.

## Features
- Employee authentication (Login, JWT-based authentication)
- Employee attendance tracking (Check-in, Check-out, Attendance report)
- Order management (Create, Update, Delete, Retrieve orders)
- Food and Drink management
- Role-based access control (Manager role)
- Reports (Daily, Weekly, Monthly, Yearly)

## Technologies Used
- **Golang**: Backend development
- **Gin-Gonic**: HTTP framework for building RESTful APIs
- **MongoDB**: NoSQL database for data storage
- **JWT (JSON Web Token)**: Authentication and authorization
- **Bcrypt**: Password hashing
- **Email Services**: Sending emails (if implemented)

## Setup Instructions

### Prerequisites
- Install Go (>=1.18)
- Install MongoDB
- Set up environment variables

### Environment Variables
Create a `.env` file and define the following variables:
```env
PORT=8080
MONGO_URI=mongodb://localhost:27017/kushena
JWT_SECRET=your_secret_key
```

### Installation
Clone the repository and install dependencies:
```sh
git clone https://github.com/yesetoda/kushena.git
cd kushena
go mod tidy
```

### Run the Application
```sh
go run main.go
```

## API Endpoints

### Authentication
| Method | Endpoint            | Description |
|--------|--------------------|-------------|
| POST   | `/employee/login`  | Employee login |

### Attendance Management
| Method | Endpoint              | Description |
|--------|----------------------|-------------|
| POST   | `/checkin`           | Employee check-in |
| POST   | `/checkout`          | Employee check-out |
| GET    | `/attendance`        | Get attendance records |
| GET    | `/checkstatus`       | Get employee check-in status |
| GET    | `/todaysworkingtime` | Get today's working hours |

### Reports (Manager Only)
| Method | Endpoint    | Description |
|--------|------------|-------------|
| GET    | `/report/daily`   | Get daily report |
| GET    | `/report/weekly`  | Get weekly report |
| GET    | `/report/monthly` | Get monthly report |
| GET    | `/report/yearly`  | Get yearly report |

### Employee Management (Manager Only)
| Method | Endpoint                  | Description |
|--------|--------------------------|-------------|
| POST   | `/manage/employee`       | Create an employee |
| GET    | `/manage/employee/:id`   | Get employee by ID |
| PATCH  | `/manage/employee`       | Update employee details |
| DELETE | `/manage/employee/:id`   | Delete an employee |
| GET    | `/manage/employees`      | Get all employees |

### Order Management
| Method | Endpoint           | Description |
|--------|------------------|-------------|
| POST   | `/action/order`  | Create an order |
| PATCH  | `/action/order`  | Update an order |
| DELETE | `/action/order/:id` | Delete an order |
| GET    | `/action/order/:id` | Get order by ID |
| GET    | `/action/orders` | Get all orders |
| GET    | `/action/myorders` | Get all orders for logged-in user |

### Food & Drink Management
| Method | Endpoint          | Description |
|--------|-----------------|-------------|
| POST   | `/action/food`  | Create food item |
| PATCH  | `/action/food`  | Update food item |
| DELETE | `/action/food/:id` | Delete food item |
| GET    | `/action/food/:id` | Get food by ID |
| GET    | `/action/foods` | Get all foods |
| POST   | `/action/drink` | Create drink item |
| PATCH  | `/action/drink` | Update drink item |
| DELETE | `/action/drink/:id` | Delete drink item |
| GET    | `/action/drink/:id` | Get drink by ID |
| GET    | `/action/drinks` | Get all drinks |

## Authentication & Authorization
- JWT authentication is required for most endpoints.
- Manager endpoints require a valid JWT with `Manager` role.

## Error Handling
- Returns `404` for undefined routes.
- Returns `403` for unauthorized access.
- Returns `400` for bad requests.

## Contributing
1. Fork the repository.
2. Create a new branch (`feature/your-feature`).
3. Commit changes and push to the branch.
4. Open a pull request.

## License
MIT License

## Contact
For inquiries, reach out to `yesetoda@gmail.com`

